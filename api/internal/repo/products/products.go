package products

import (
	"context"
	"fmt"

	"github.com/creent-production/cdk-go/pagination"
	"github.com/jmoiron/sqlx"

	productsentity "github.com/IndominusByte/catalyst-be/api/internal/entity/products"
)

type RepoProducts struct {
	db      *sqlx.DB
	queries map[string]string
	execs   map[string]string
}

var queries = map[string]string{
	"getProductByDynamic": `SELECT id, name, description, price, brand_id, created_at, updated_at FROM transaction.products`,
}
var execs = map[string]string{
	"insertProduct": `INSERT INTO transaction.products (name, description, price, brand_id) VALUES (:name, :description, :price, :brand_id) RETURNING id`,
}

func New(db *sqlx.DB) (*RepoProducts, error) {
	rp := &RepoProducts{
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
func (r *RepoProducts) Validate() error {
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

func (r *RepoProducts) GetProductByName(ctx context.Context, name string) (*productsentity.Product, error) {
	var t productsentity.Product
	stmt, _ := r.db.PrepareNamedContext(ctx, r.queries["getProductByDynamic"]+" WHERE name = :name")

	return &t, stmt.GetContext(ctx, &t, productsentity.Product{Name: name})
}

func (r *RepoProducts) GetProductById(ctx context.Context, id int) (*productsentity.Product, error) {
	var t productsentity.Product
	stmt, _ := r.db.PrepareNamedContext(ctx, r.queries["getProductByDynamic"]+" WHERE id = :id")

	return &t, stmt.GetContext(ctx, &t, productsentity.Product{Id: id})
}

func (r *RepoProducts) Insert(ctx context.Context, payload *productsentity.JsonCreateSchema) int {
	var id int
	stmt, _ := r.db.PrepareNamedContext(ctx, r.execs["insertProduct"])
	stmt.QueryRowxContext(ctx, payload).Scan(&id)

	return id
}

func (r *RepoProducts) GetAllProductByBrandPaginate(ctx context.Context,
	payload *productsentity.QueryParamAllProductSchema) (*productsentity.ProductPaginate, error) {

	var results productsentity.ProductPaginate

	query := r.queries["getProductByDynamic"] + ` WHERE brand_id = :brand_id`
	query += ` ORDER BY id DESC`

	// pagination
	var count struct{ Total int }
	stmt_count, _ := r.db.PrepareNamedContext(ctx, fmt.Sprintf("SELECT count(*) AS total FROM (%s) AS anon_1", query))
	err := stmt_count.GetContext(ctx, &count, payload)
	if err != nil {
		return &results, err
	}
	payload.Offset = (payload.Page - 1) * payload.PerPage

	// results
	query += ` LIMIT :per_page OFFSET :offset`
	stmt, _ := r.db.PrepareNamedContext(ctx, query)
	err = stmt.SelectContext(ctx, &results.Data, payload)
	if err != nil {
		return &results, err
	}

	paginate := pagination.Paginate{Page: payload.Page, PerPage: payload.PerPage, Total: count.Total}
	results.Total = paginate.Total
	results.NextNum = paginate.NextNum()
	results.PrevNum = paginate.PrevNum()
	results.Page = paginate.Page
	results.IterPages = paginate.IterPages()

	return &results, nil
}
