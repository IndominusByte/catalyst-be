package brands

import (
	"context"

	brandsentity "github.com/IndominusByte/catalyst-be/api/internal/entity/brands"
)

type brandsRepo interface {
	GetBrandByName(ctx context.Context, name string) (*brandsentity.Brand, error)
	Insert(ctx context.Context, payload *brandsentity.JsonCreateSchema) int
}
