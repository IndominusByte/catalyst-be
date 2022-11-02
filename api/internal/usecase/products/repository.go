package products

import (
	"context"

	brandsentity "github.com/IndominusByte/catalyst-be/api/internal/entity/brands"
	productsentity "github.com/IndominusByte/catalyst-be/api/internal/entity/products"
)

type productsRepo interface {
	GetProductByName(ctx context.Context, name string) (*productsentity.Product, error)
	GetProductById(ctx context.Context, id int) (*productsentity.Product, error)
	GetAllProductByBrandPaginate(ctx context.Context,
		payload *productsentity.QueryParamAllProductSchema) (*productsentity.ProductPaginate, error)
	Insert(ctx context.Context, payload *productsentity.JsonCreateSchema) int
}

type brandsRepo interface {
	GetBrandById(ctx context.Context, id int) (*brandsentity.Brand, error)
}
