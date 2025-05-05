package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Gen1usBruh/warehouse-api/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(dbConfig *config.Database) (*pgxpool.Pool, error) {

	if dbConfig == nil {
		return nil, errors.New("missing config")
	}

	const op = "storage.postgres.ConnectDB"

	dbUrl := dsn(dbConfig)

	pgxConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return nil, fmt.Errorf("%s | %w", op, err)
	}

	if dbConfig.MaxConns != 0 {
		pgxConfig.MaxConns = dbConfig.MaxConns
	} else {
		pgxConfig.MaxConns = 15
	}
	if dbConfig.MaxIdleConnections != 0 {
		pgxConfig.MinConns = dbConfig.MaxIdleConnections
	} else {
		pgxConfig.MinConns = 10
	}
	if dbConfig.MaxConnLifeTimeInSeconds != 0 {
		pgxConfig.MaxConnLifetime = time.Second * time.Duration(dbConfig.MaxConnLifeTimeInSeconds)
	} else {
		pgxConfig.MaxConnLifetime = 25 * time.Minute
	}
	if dbConfig.MaxConnIdleTimeInSeconds != 0 {
		pgxConfig.MaxConnIdleTime = time.Second * time.Duration(dbConfig.MaxConnIdleTimeInSeconds)
	} else {
		pgxConfig.MaxConnLifetime = 5 * time.Minute
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("%s | %w", op, err)
	}

	if err = dbPool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("%s | %w", op, err)
	}

	return dbPool, nil
}

func dsn(d *config.Database) string {
	return fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v?sslmode=%v",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.DBName,
		d.SSLMode,
	)
}
