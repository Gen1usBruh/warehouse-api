package product

import "context"

type Repository interface {
	Create(ctx context.Context, p Product) (int32, error)
	GetByID(ctx context.Context, id int32) (Product, error)
	Update(ctx context.Context, p Product) error
	Delete(ctx context.Context, id int32) error
	List(ctx context.Context) ([]Product, error)
}
