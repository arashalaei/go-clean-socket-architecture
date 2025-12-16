package tcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/arashalaei/go-clean-socket-architecture/internal/delivery/tcp/dto"
	"github.com/arashalaei/go-clean-socket-architecture/internal/usecase/school"
)

type IServer interface {
	Start(ctx context.Context) error
	Shutdown() error
	RegisterHandler(reqType RequestType, handler RequestHandler)
	ServerHandlers
}

type ServerHandlers interface {
	CreateSchoolHandler(ctx context.Context, payload json.RawMessage) (interface{}, error)
	CreatePersonHandler(ctx context.Context, payload json.RawMessage) (interface{}, error)
	CreateClassHandler(ctx context.Context, payload json.RawMessage) (interface{}, error)
	WhoAmIHandler(ctx context.Context, payload json.RawMessage) (interface{}, error)
}

type SrvCfg struct {
	Network         string
	Address         string
	MaxConnections  int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	MaxMessageSize  int64
	ShutdownTimeout time.Duration
}

type RequestType string

const (
	CreateSchool RequestType = "creat_school"
	CreatePerson RequestType = "creat_person"
	CreateClass  RequestType = "creat_class"
	WhoAmI       RequestType = "who_am_i"
)

type server struct {
	// optional
	cfg    *SrvCfg
	logger *log.Logger

	listener    net.Listener
	connections sync.Map
	connCount   chan struct{}
	handlers    map[RequestType]RequestHandler
	wg          sync.WaitGroup
	mu          sync.RWMutex

	// use cases
	schoolUsecases *school.SchoolUsecases
}

type RequestHandler func(ctx context.Context, payload json.RawMessage) (interface{}, error)

func NewServer(ops ...srvops) IServer {
	s := &server{
		// default cfg
		cfg: &SrvCfg{
			Network:         "tcp",
			Address:         ":8080",
			MaxConnections:  1000,
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			IdleTimeout:     5 * time.Minute,
			MaxMessageSize:  1024 * 1024, // 1MB
			ShutdownTimeout: 30 * time.Second,
		},
		// default logger
		logger:    log.New(os.Stdout, "[TCP Server]", log.LstdFlags),
		connCount: make(chan struct{}, 1000),
		handlers:  make(map[RequestType]RequestHandler),
		wg:        sync.WaitGroup{},
		mu:        sync.RWMutex{},
	}

	// set options
	for _, o := range ops {
		o(s)
	}
	return s
}

type srvops func(*server)

func WithCfg(cfg SrvCfg) srvops {
	return func(s *server) {
		s.cfg = &cfg
		s.connCount = make(chan struct{}, cfg.MaxConnections)
	}
}

func WithLogger(l log.Logger) srvops {
	return func(s *server) {
		s.logger = &l
	}
}

func WithSchoolUsecases(su school.SchoolUsecases) srvops {
	return func(s *server) {
		s.schoolUsecases = &su
	}
}

func (s *server) Start(ctx context.Context) error {
	l, err := net.Listen(s.cfg.Network, s.cfg.Address)
	if err != nil {
		return fmt.Errorf("failed to start listener: %w", err)
	}
	s.listener = l
	s.logger.Printf("Server is listening on: %s\n", s.cfg.Address)

	go s.acceptConn(ctx)

	return nil
}

func (s *server) RegisterHandler(reqType RequestType, handler RequestHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[reqType] = handler
}

func (s *server) acceptConn(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				return
			}

			select {
			case s.connCount <- struct{}{}:
				s.wg.Add(1)
				go s.handleConn(ctx, conn)
			default:
				s.logger.Println("Connection limit reached, rejecting connection")
				conn.Close()
			}
		}
	}
}

func (s *server) handleConn(ctx context.Context, conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()
	defer func() { <-s.connCount }()

	clienAddr := conn.RemoteAddr().String()
	s.connections.Store(clienAddr, conn)
	defer s.connections.Delete(clienAddr)

	s.logger.Printf("New connection from %s\n", clienAddr)
	defer s.logger.Printf("Connection closed: %s\n", clienAddr)

	reader := bufio.NewReader(conn)
	conn.SetDeadline(time.Now().Add(s.cfg.IdleTimeout))

	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn.SetReadDeadline(time.Now().Add(s.cfg.ReadTimeout))
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					return
				}
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					s.logger.Printf("Read timeout for %s\n", clienAddr)
					return
				}
				s.logger.Printf("Read error for %s: %v\n", clienAddr, err)
				return
			}

			if int64(len(line)) > s.cfg.MaxMessageSize {
				s.json(conn, dto.Response{
					Status:  false,
					Message: "message too large",
				})
				continue
			}

			conn.SetDeadline(time.Now().Add(s.cfg.IdleTimeout))
			s.processRequest(ctx, conn, line)
		}
	}
}

func (s *server) processRequest(ctx context.Context, conn net.Conn, data []byte) {
	var req dto.Request

	if err := json.Unmarshal(data, &req); err != nil {
		s.json(conn, dto.Response{
			Status:  false,
			Message: "invalid JSON format",
		})
		return
	}

	s.mu.RLock()
	handler, ok := s.handlers[RequestType(req.Type)]
	s.mu.RUnlock()

	if !ok {
		s.json(conn, dto.Response{
			Status:  false,
			Message: fmt.Sprintf("unknown request type: %s", req.Type),
		})
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	result, err := handler(ctx, req.Payload)
	if err != nil {
		s.json(conn, dto.Response{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	s.json(conn, dto.Response{
		Status: true,
		Data:   result,
	})
}

func (s *server) json(conn net.Conn, res dto.Response) {
	conn.SetWriteDeadline(time.Now().Add(s.cfg.WriteTimeout))

	data, err := json.Marshal(res)
	if err != nil {
		s.logger.Printf("Failed to marshal response: %v\n", err)
		return
	}

	data = append(data, '\n')

	if _, err = conn.Write(data); err != nil {
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			s.logger.Printf("Wrire timeout")
			return
		}
		s.logger.Printf("Failed to write response: %v\n", err)
	}
}

func (s *server) Shutdown() error {
	s.logger.Println("Starting graceful shutdown...")
	if err := s.listener.Close(); err != nil {
		s.logger.Printf("Error closing listener: %v\n", err)
	}
	return nil
}
