package users

import (
	"context"

	usersentity "github.com/IndominusByte/catalyst-be/api/internal/entity/users"
	"github.com/jmoiron/sqlx"
)

type RepoUsers struct {
	db      *sqlx.DB
	queries map[string]string
	execs   map[string]string
}

var queries = map[string]string{
	"getUserByDynamic": `SELECT id, email, name, password, created_at, updated_at FROM account.users`,
}
var execs = map[string]string{}

func New(db *sqlx.DB) (*RepoUsers, error) {
	rp := &RepoUsers{
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
func (r *RepoUsers) Validate() error {
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

func (r *RepoUsers) GetUserById(ctx context.Context, id int) (*usersentity.User, error) {
	var t usersentity.User
	stmt, _ := r.db.PrepareNamedContext(ctx, r.queries["getUserByDynamic"]+" WHERE id = :id")

	return &t, stmt.GetContext(ctx, &t, usersentity.User{Id: id})
}
