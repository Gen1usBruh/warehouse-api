package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/Gen1usBruh/warehouse-api/internal/domain/product"
)

type ProductUseCase struct {
	repo product.Repository
}

func NewProductUseCase(r product.Repository) *ProductUseCase {
	return &ProductUseCase{repo: r}
}

var (
	ErrPriceLimit     = errors.New("price exceeds maximum allowed value of $10,000")
	ErrNameIsReserved = errors.New("product name is reserved")
	ErrQuantityLimit  = errors.New("quantity exceeds maximum allowed value of 1000 units")
)

func IsBusinessError(err error) bool {
	return errors.Is(err, ErrPriceLimit) ||
		errors.Is(err, ErrNameIsReserved) ||
		errors.Is(err, ErrQuantityLimit)
}

func validateProduct(p product.Product) error {
	if strings.EqualFold(p.Name, "Sarkor") || strings.EqualFold(p.Name, "Sochnaya Dolina") {
		return ErrNameIsReserved
	}

	if p.Price > 10000 {
		return ErrPriceLimit
	}

	if p.Quantity > 1000 {
		return ErrQuantityLimit
	}

	return nil
}

func (u *ProductUseCase) Create(ctx context.Context, p product.Product) (int32, error) {
	if err := validateProduct(p); err != nil {
		return 0, err
	}

	return u.repo.Create(ctx, p)
}

func (u *ProductUseCase) GetByID(ctx context.Context, id int32) (product.Product, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *ProductUseCase) Update(ctx context.Context, p product.Product) error {
	return u.repo.Update(ctx, p)
}

func (u *ProductUseCase) Delete(ctx context.Context, id int32) error {
	return u.repo.Delete(ctx, id)
}

func (u *ProductUseCase) List(ctx context.Context) ([]product.Product, error) {
	return u.repo.List(ctx)
}
