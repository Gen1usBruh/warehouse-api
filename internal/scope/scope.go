package scope

import (
	"log/slog"

	usecase "github.com/Gen1usBruh/warehouse-api/internal/usecase"
)

type Dependencies struct {
	Sl      *slog.Logger
	Product *usecase.ProductUseCase
}
