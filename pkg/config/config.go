package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Client ClientConfig `mapstructure:"client"`
}

type ServerConfig struct {
	Network  string           `mapstructure:"network"`
	Address  string           `mapstructure:"address"`
	Timeouts SrvTimeoutConfig `mapstructure:"timeouts"`
	Limits   SrvLimitConfig   `mapstructure:"limits"`
}

type SrvTimeoutConfig struct {
	Read     time.Duration `mapstructure:"read"`
	Write    time.Duration `mapstructure:"write"`
	Idle     time.Duration `mapstructure:"idle"`
	Shutdown time.Duration `mapstructure:"shutdown"`
}

type SrvLimitConfig struct {
	MaxConnectionsSize int   `mapstructure:"max_connections_size"`
	MaxMessageSize     int64 `mapstructure:"max_message_size"`
}

type ClientConfig struct {
	Network  string              `mapstructure:"network"`
	Address  string              `mapstructure:"address"`
	Timeouts ClinetTimeoutConfig `mapstructure:"timeouts"`
	Limits   ClientLimitConfig   `mapstructure:"limits"`
}

type ClinetTimeoutConfig struct {
	Read            time.Duration `mapstructure:"read"`
	Write           time.Duration `mapstructure:"write"`
	Connect         time.Duration `mapstructure:"connect"`
	KeepAlivePeriod time.Duration `mapstructure:"keep_alive_period"`
}

type ClientLimitConfig struct {
	MaxRetries int           `mapstructure:"max_retries"`
	RetryDelay time.Duration `mapstructure:"retry_delay"`
	KeepAlive  bool          `mapstructure:"keep_alive"`
}

func Load(path string) (*Config, error) {
	viper.SetConfigFile(path)

	//  viper.SetDefault("server.network", "tcp")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// viper.AutomaticEnv()
	// viper.BindEnv("server.address", "SERVER_ADDRESS")

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
