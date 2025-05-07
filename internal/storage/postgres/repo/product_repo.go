package repo

import (
	"context"

	"github.com/Gen1usBruh/warehouse-api/internal/domain/product"
	db "github.com/Gen1usBruh/warehouse-api/internal/storage/postgres/sqlc"
)

type ProductRepo struct {
	q *db.Queries
}

func NewProductRepo(q *db.Queries) *ProductRepo {
	return &ProductRepo{q: q}
}

func (r *ProductRepo) Create(ctx context.Context, p product.Product) (int32, error) {
	return r.q.CreateProduct(ctx, db.CreateProductParams{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Quantity:    p.Quantity,
	})
}

func (r *ProductRepo) GetByID(ctx context.Context, id int32) (product.Product, error) {
	row, err := r.q.GetProductByID(ctx, id)
	if err != nil {
		return product.Product{}, err
	}
	return product.Product(row), nil
}

func (r *ProductRepo) Update(ctx context.Context, p product.Product) error {
	return r.q.UpdateProduct(ctx, db.UpdateProductParams(p))
}

func (r *ProductRepo) Delete(ctx context.Context, id int32) error {
	return r.q.DeleteProduct(ctx, id)
}

func (r *ProductRepo) List(ctx context.Context) ([]product.Product, error) {
	rows, err := r.q.ListProducts(ctx)
	if err != nil {
		return nil, err
	}
	var result []product.Product
	for _, row := range rows {
		result = append(result, product.Product(row))
	}
	return result, nil
}
