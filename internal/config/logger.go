package config

type Logger struct {
	Level string `env:"LOG_LEVEL" envDefault:"info"`
}