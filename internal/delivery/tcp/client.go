package tcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/arashalaei/go-clean-socket-architecture/internal/delivery/tcp/dto"
)

type ClientConfig struct {
	Network         string
	Address         string
	ConnectTimeout  time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxRetries      int
	RetryDelay      time.Duration
	KeepAlive       bool
	KeepAlivePeriod time.Duration
}

type Client struct {
	config *ClientConfig
	conn   net.Conn
	mu     sync.RWMutex
	reader *bufio.Reader
	writer *bufio.Writer
	closed bool
}

type clientOps func(*Client)

func WithClientCfg(cfg ClientConfig) clientOps {
	return func(c *Client) {
		c.config = &cfg
	}
}

func NewClient(ops ...clientOps) *Client {
	c := &Client{
		config: &ClientConfig{
			Network:         "tcp",
			Address:         "localhost:8080",
			ConnectTimeout:  10 * time.Second,
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			MaxRetries:      3,
			RetryDelay:      time.Second,
			KeepAlive:       true,
			KeepAlivePeriod: 30 * time.Second,
		},
	}

	for _, o := range ops {
		o(c)
	}

	return c
}

func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil {
		return fmt.Errorf("already connected")
	}

	dialer := &net.Dialer{
		Timeout:   c.config.ConnectTimeout,
		KeepAlive: c.config.KeepAlivePeriod,
	}

	conn, err := dialer.Dial(c.config.Network, c.config.Address)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	if c.config.KeepAlive {
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(c.config.KeepAlivePeriod)
		}
	}

	c.conn = conn
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)
	c.closed = false

	return nil
}

func (c *Client) Send(
	ctx context.Context,
	requestType RequestType,
	payload interface{},
) (*dto.Response, error) {
	c.mu.RLock()
	if c.conn == nil || c.closed {
		c.mu.RUnlock()
		return nil, fmt.Errorf("not connected")
	}
	c.mu.RUnlock()

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req := dto.Request{
		Type:    string(requestType),
		Payload: payloadBytes,
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	var resp *dto.Response
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.config.RetryDelay):
			}
		}

		resp, err = c.sendWithTimeout(ctx, reqBytes)
		if err == nil {
			return resp, nil
		}

		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	return nil, err
}

func (c *Client) sendWithTimeout(ctx context.Context, reqBytes []byte) (*dto.Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(c.config.WriteTimeout)
	}
	c.conn.SetWriteDeadline(deadline)

	reqBytes = append(reqBytes, '\n')
	if _, err := c.writer.Write(reqBytes); err != nil {
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			return nil, fmt.Errorf("timeout for write")
		}
		return nil, fmt.Errorf("failed to write request: %w", err)
	}

	if err := c.writer.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush: %w", err)
	}

	if !ok {
		deadline = time.Now().Add(c.config.ReadTimeout)
	}
	c.conn.SetReadDeadline(deadline)

	respBytes, err := c.reader.ReadBytes('\n')
	if err != nil {
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			return nil, fmt.Errorf("timeout for read")
		}
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var resp dto.Response
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &resp, nil
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil || c.closed {
		return nil
	}

	c.closed = true
	return c.conn.Close()
}

func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn != nil && !c.closed
}
