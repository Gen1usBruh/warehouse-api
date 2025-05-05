package app

import (
	"net/http"

	"github.com/Gen1usBruh/warehouse-api/internal/config"
)

type App struct {
	*http.Server
}

func NewApp(server config.Server, restServer http.Handler) (*App, error) {
	sm := http.NewServeMux()

	httpServer := http.Server{
		Addr:         server.Address,
		Handler:      sm,
		ReadTimeout:  server.Timeout,
		WriteTimeout: server.Timeout,
		IdleTimeout:  server.IdleTimeout,
	}

	return &App{
		Server: &httpServer,
	}, nil
}