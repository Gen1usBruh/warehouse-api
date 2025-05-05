package config

import "github.com/caarlos0/env/v11"

func New() (*Config, error) {
	var conf Config
	err := env.Parse(&conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

type Config struct {
	Database Database
	Server   Server
	Logger   Logger
}