package config

import "time"

type Server struct {
	Address string `env:"SERVER_ADDRESS"`
	Timeout time.Duration `env:"SERVER_TIMEOUT"`
	IdleTimeout time.Duration `env:"SERVER_IDLE_TIMEOUT"`
}