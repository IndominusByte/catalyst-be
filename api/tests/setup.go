package tests

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/IndominusByte/catalyst-be/api/internal/config"
	handler_http "github.com/IndominusByte/catalyst-be/api/internal/endpoint/http/handler"

	brandsrepo "github.com/IndominusByte/catalyst-be/api/internal/repo/brands"
	ordersrepo "github.com/IndominusByte/catalyst-be/api/internal/repo/orders"
	productsrepo "github.com/IndominusByte/catalyst-be/api/internal/repo/products"
	usersrepo "github.com/IndominusByte/catalyst-be/api/internal/repo/users"
)

type setupRepo struct {
	brandsRepo   brandsrepo.RepoBrands
	productsRepo productsrepo.RepoProducts
	ordersRepo   ordersrepo.RepoOrders
	usersRepo    usersrepo.RepoUsers
}

func setupEnvironment() (*setupRepo, *handler_http.Server) {
	// init config
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	// connect the db
	db, err := config.DBConnect(cfg)
	if err != nil {
		panic(err)
	}
	// mount router
	r := handler_http.CreateNewServer(db, cfg)
	if err := r.MountHandlers(); err != nil {
		panic(err)
	}

	// you can insert your behaviors here
	brandsRepo, _ := brandsrepo.New(db)
	productsRepo, _ := productsrepo.New(db)
	ordersRepo, _ := ordersrepo.New(db)
	usersRepo, _ := usersrepo.New(db)

	setuprepo := setupRepo{
		brandsRepo:   *brandsRepo,
		productsRepo: *productsRepo,
		ordersRepo:   *ordersRepo,
		usersRepo:    *usersRepo,
	}

	return &setuprepo, r
}

// executeRequest, creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(req *http.Request, s *handler_http.Server) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, req)

	return rr
}

func createForm(form map[string]string) (string, io.Reader, error) {
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)
	defer mp.Close()
	for key, val := range form {
		if strings.HasPrefix(val, "@") {
			val = val[1:]
			if len(val) < 1 {
				mp.CreateFormFile(key, "")
				continue
			}
			file, err := os.Open(val)
			if err != nil {
				return "", nil, err
			}
			defer file.Close()
			filename := strings.Split(val, "/")
			part, err := mp.CreateFormFile(key, filename[len(filename)-1])
			if err != nil {
				return "", nil, err
			}
			io.Copy(part, file)
		} else {
			mp.WriteField(key, val)
		}
	}
	return mp.FormDataContentType(), body, nil
}

func createMaximum(length int) string {
	word := ""
	for i := 0; i < length; i++ {
		word += "a"
	}
	return word
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
