package scope

import (
	"log/slog"

	postgresdb "github.com/Gen1usBruh/warehouse-api/internal/storage/postgres/sqlc"
)

type Dependencies struct {
	Sl *slog.Logger
	Db *postgresdb.Queries
}
