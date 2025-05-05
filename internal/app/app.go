package app

import (
	"net/http"

	"github.com/Gen1usBruh/warehouse-api/internal/config"
)

type App struct {
	*http.Server
}

func NewApp(server config.Server, restServer http.Handler) (*App, error) {
	httpServer := http.Server{
		Addr:         server.Address,
		Handler:      restServer,
		ReadTimeout:  server.Timeout,
		WriteTimeout: server.Timeout,
		IdleTimeout:  server.IdleTimeout,
	}

	return &App{
		Server: &httpServer,
	}, nil
}
