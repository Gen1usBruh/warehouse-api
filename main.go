package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gen1usBruh/warehouse-api/internal/app"
	"github.com/Gen1usBruh/warehouse-api/internal/config"
	"github.com/Gen1usBruh/warehouse-api/internal/logger/sl"
	"github.com/Gen1usBruh/warehouse-api/internal/rest"
	"github.com/Gen1usBruh/warehouse-api/internal/scope"
	"github.com/Gen1usBruh/warehouse-api/internal/storage/postgres"
	postgresdb "github.com/Gen1usBruh/warehouse-api/internal/storage/postgres/sqlc"
	"github.com/Gen1usBruh/warehouse-api/internal/storage/postgres/repo"
	"github.com/Gen1usBruh/warehouse-api/internal/usecase"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
        log.Fatalf("No .env file found: %v\n", err)
    }
	conf, err := config.New()
	if err != nil {
		log.Fatalf("Could not create config: %v\n", err)
	}

	conn, err := postgres.ConnectDB(&conf.Database)
	if err != nil {
		log.Fatalf("Could not connect to postgres: %v\n", err)
	}

	qConn := postgresdb.New(conn)
	productRepo := repo.NewProductRepo(qConn)
	productUC := usecase.NewProductUseCase(productRepo)

	restServer := rest.NewHandler(rest.HandlerConfig{
		Dep: &scope.Dependencies{
			Sl: sl.SetupLogger(&conf.Logger),
			Product: productUC,
		},
	})

	server, err := app.NewApp(conf.Server, restServer)
	if err != nil {
		log.Fatalf("Unable to start server: %v\n", err)
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
