package brands

import (
	"context"

	brandsentity "github.com/IndominusByte/catalyst-be/api/internal/entity/brands"
	"github.com/jmoiron/sqlx"
)

type RepoBrands struct {
	db      *sqlx.DB
	queries map[string]string
	execs   map[string]string
}

var queries = map[string]string{
	"getBrandByDynamic": `SELECT id, name, created_at, updated_at FROM transaction.brands`,
}
var execs = map[string]string{
	"insertBrand": `INSERT INTO transaction.brands (name) VALUES (:name) RETURNING id`,
}

func New(db *sqlx.DB) (*RepoBrands, error) {
	rp := &RepoBrands{
		db:      db,
		queries: queries,
		execs:   execs,
	}

	err := rp.Validate()
	if err != nil {
		return nil, err
	}

	return rp, nil
}

// Validate will validate sql query to db
func (r *RepoBrands) Validate() error {
	for _, q := range r.queries {
		_, err := r.db.PrepareNamedContext(context.Background(), q)
		if err != nil {
			return err
		}
	}

	for _, e := range r.execs {
		_, err := r.db.PrepareNamedContext(context.Background(), e)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RepoBrands) GetBrandByName(ctx context.Context, name string) (*brandsentity.Brand, error) {
	var t brandsentity.Brand
	stmt, _ := r.db.PrepareNamedContext(ctx, r.queries["getBrandByDynamic"]+" WHERE name = :name")

	return &t, stmt.GetContext(ctx, &t, brandsentity.Brand{Name: name})
}

func (r *RepoBrands) GetBrandById(ctx context.Context, id int) (*brandsentity.Brand, error) {
	var t brandsentity.Brand
	stmt, _ := r.db.PrepareNamedContext(ctx, r.queries["getBrandByDynamic"]+" WHERE id = :id")

	return &t, stmt.GetContext(ctx, &t, brandsentity.Brand{Id: id})
}

func (r *RepoBrands) Insert(ctx context.Context, payload *brandsentity.JsonCreateSchema) int {
	var id int
	stmt, _ := r.db.PrepareNamedContext(ctx, r.execs["insertBrand"])
	stmt.QueryRowxContext(ctx, payload).Scan(&id)

	return id
}
