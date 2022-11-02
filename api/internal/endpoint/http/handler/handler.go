package handler_http

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/IndominusByte/catalyst-be/api/internal/config"
	"github.com/jmoiron/sqlx"

	endpoint_http "github.com/IndominusByte/catalyst-be/api/internal/endpoint/http"
	brandsrepo "github.com/IndominusByte/catalyst-be/api/internal/repo/brands"
	ordersrepo "github.com/IndominusByte/catalyst-be/api/internal/repo/orders"
	productsrepo "github.com/IndominusByte/catalyst-be/api/internal/repo/products"
	usersrepo "github.com/IndominusByte/catalyst-be/api/internal/repo/users"
	brandsusecase "github.com/IndominusByte/catalyst-be/api/internal/usecase/brands"
	ordersusecase "github.com/IndominusByte/catalyst-be/api/internal/usecase/orders"
	productsusecase "github.com/IndominusByte/catalyst-be/api/internal/usecase/products"
)

type Server struct {
	Router *http.ServeMux
	// Db config can be added here
	db  *sqlx.DB
	cfg *config.Config
}

func CreateNewServer(db *sqlx.DB, cfg *config.Config) *Server {
	s := &Server{db: db, cfg: cfg}
	s.Router = http.NewServeMux()
	return s
}

func (s *Server) MountHandlers() error {
	s.Router.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("doc.json"), //The url pointing to API definition
	))

	// you can insert your behaviors here
	brandsRepo, err := brandsrepo.New(s.db)
	if err != nil {
		return err
	}
	brandsUsecase := brandsusecase.NewBrandsUsecase(brandsRepo)
	endpoint_http.AddBrands(s.Router, brandsUsecase)

	productsRepo, err := productsrepo.New(s.db)
	if err != nil {
		return err
	}
	productsUsecase := productsusecase.NewProductsUsecase(productsRepo, brandsRepo)
	endpoint_http.AddProducts(s.Router, productsUsecase)

	ordersRepo, err := ordersrepo.New(s.db)
	if err != nil {
		return err
	}

	usersRepo, err := usersrepo.New(s.db)
	if err != nil {
		return err
	}
	ordersUsecase := ordersusecase.NewOrdersUsecase(ordersRepo, productsRepo, usersRepo)
	endpoint_http.AddOrders(s.Router, ordersUsecase)

	return nil
}
