package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
}

type ServerConfig struct {
	Network  string        `mapstructure:"network"`
	Address  string        `mapstructure:"address"`
	Timeouts TimeoutConfig `mapstructure:"timeouts"`
	Limits   LimitConfig   `mapstructure:"limits"`
}

type TimeoutConfig struct {
	Read     time.Duration `mapstructure:"read"`
	Write    time.Duration `mapstructure:"write"`
	Idle     time.Duration `mapstructure:"idle"`
	Shutdown time.Duration `mapstructure:"shutdown"`
}

type LimitConfig struct {
	MaxConnectionsSize int   `mapstructure:"max_connections_size"`
	MaxMessageSize     int64 `mapstructure:"max_message_size"`
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
