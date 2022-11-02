package orders

import (
	"context"

	"github.com/jmoiron/sqlx"

	ordersentity "github.com/IndominusByte/catalyst-be/api/internal/entity/orders"
)

type RepoOrders struct {
	db      *sqlx.DB
	queries map[string]string
	execs   map[string]string
}

var queries = map[string]string{
	"getOrderByDynamic": `SELECT transaction.orders.id, transaction.orders.buyer_id, account.users.name AS buyer_name, transaction.orders.product_id, transaction.products.name AS product_name, transaction.orders.qty, transaction.orders.price, transaction.orders.total_price, transaction.orders.created_at, transaction.orders.updated_at FROM transaction.orders LEFT JOIN transaction.products ON transaction.products.id = transaction.orders.product_id LEFT JOIN account.users ON account.users.id = transaction.orders.buyer_id`,
}
var execs = map[string]string{
	"insertOrder": `INSERT INTO transaction.orders (buyer_id, product_id, qty, price, total_price) VALUES (:buyer_id, :product_id, :qty, :price, :total_price) RETURNING id`,
}

func New(db *sqlx.DB) (*RepoOrders, error) {
	rp := &RepoOrders{
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
func (r *RepoOrders) Validate() error {
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

func (r *RepoOrders) Insert(ctx context.Context, payload *ordersentity.JsonCreateSchema) int {
	var id int
	stmt, _ := r.db.PrepareNamedContext(ctx, r.execs["insertOrder"])
	stmt.QueryRowxContext(ctx, payload).Scan(&id)

	return id
}

func (r *RepoOrders) GetOrderById(ctx context.Context, id int) (*ordersentity.Order, error) {
	var t ordersentity.Order
	stmt, _ := r.db.PrepareNamedContext(ctx, r.queries["getOrderByDynamic"]+" WHERE transaction.orders.id = :id")

	return &t, stmt.GetContext(ctx, &t, ordersentity.Order{Id: id})
}

func (r *RepoOrders) GetOrderByBuyerId(ctx context.Context, buyerId int) (*ordersentity.Order, error) {
	var t ordersentity.Order
	stmt, _ := r.db.PrepareNamedContext(ctx, r.queries["getOrderByDynamic"]+" WHERE transaction.orders.buyer_id = :buyer_id")

	return &t, stmt.GetContext(ctx, &t, ordersentity.Order{BuyerId: buyerId})
}
