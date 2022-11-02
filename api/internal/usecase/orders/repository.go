package orders

import (
	"context"

	ordersentity "github.com/IndominusByte/catalyst-be/api/internal/entity/orders"
	productsentity "github.com/IndominusByte/catalyst-be/api/internal/entity/products"
	usersentity "github.com/IndominusByte/catalyst-be/api/internal/entity/users"
)

type ordersRepo interface {
	Insert(ctx context.Context, payload *ordersentity.JsonCreateSchema) int
	GetOrderById(ctx context.Context, id int) (*ordersentity.Order, error)
}
type productsRepo interface {
	GetProductById(ctx context.Context, id int) (*productsentity.Product, error)
}
type usersRepo interface {
	GetUserById(ctx context.Context, id int) (*usersentity.User, error)
}
