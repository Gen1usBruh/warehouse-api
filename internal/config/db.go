package config

type Database struct {
	Host                     string `env:"POSTGRES_HOST,notEmpty"`
	Port                     string `env:"POSTGRES_PORT,notEmpty"`
	User                     string `env:"POSTGRES_USER,notEmpty"`
	Password                 string `env:"POSTGRES_PASSWORD,notEmpty"`
	DBName                   string `env:"POSTGRES_DATABASE,notEmpty"`
	SSLMode                  string `env:"POSTGRES_SSL_MODE,notEmpty"             envDefault:"disable"`
	MaxConns                 int32  `env:"POSTGRES_MAX_CONNS"                     envDefault:"15"`
	MaxIdleConnections       int32  `env:"POSTGRES_MAX_IDLE_CONNS"                envDefault:"10"`
	MaxConnIdleTimeInSeconds int32  `env:"POSTGRES_MAX_CONN_IDLE_TIME_IN_SECONDS" envDefault:"300"`  // 5 minutes
	MaxConnLifeTimeInSeconds int32  `env:"POSTGRES_MAX_CONN_LIFETIME_IN_SECONDS"  envDefault:"1500"` // 25 minutes
}