package endpoint_http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/IndominusByte/catalyst-be/api/internal/constant"
	productsentity "github.com/IndominusByte/catalyst-be/api/internal/entity/products"
	"github.com/creent-production/cdk-go/parser"
	"github.com/creent-production/cdk-go/response"
	"github.com/creent-production/cdk-go/validation"
)

type productsUsecaseIface interface {
	Create(ctx context.Context, rw http.ResponseWriter, payload *productsentity.JsonCreateSchema)
	GetById(ctx context.Context, rw http.ResponseWriter, productId int)
	GetByBrandId(ctx context.Context, rw http.ResponseWriter, payload *productsentity.QueryParamAllProductSchema)
}

func AddProducts(r *http.ServeMux, uc productsUsecaseIface) {
	const prefix = "/product"

	r.HandleFunc(prefix, func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			var p productsentity.JsonCreateSchema

			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
					constant.Body: constant.FailedParseBody,
				})
				return
			}

			uc.Create(r.Context(), rw, &p)
		default:
			response.WriteJSONResponse(rw, 404, nil, map[string]interface{}{
				constant.Body: constant.ResourceNotFound,
			})
		}
	})

	r.HandleFunc(prefix+"/", func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			productId, _ := parser.ParsePathToInt("/product/(.*)", r.URL.Path)

			uc.GetById(r.Context(), rw, productId)
		default:
			response.WriteJSONResponse(rw, 404, nil, map[string]interface{}{
				constant.Body: constant.ResourceNotFound,
			})
		}
	})

	r.HandleFunc(prefix+"/brand", func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			var p productsentity.QueryParamAllProductSchema

			if err := validation.ParseRequest(&p, r.URL.Query()); err != nil {
				response.WriteJSONResponse(rw, 422, nil, map[string]interface{}{
					constant.Body: constant.FailedParseBody,
				})
				return
			}

			uc.GetByBrandId(r.Context(), rw, &p)
		default:
			response.WriteJSONResponse(rw, 404, nil, map[string]interface{}{
				constant.Body: constant.ResourceNotFound,
			})
		}
	})
}
