package endpoint_http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/IndominusByte/catalyst-be/api/internal/constant"
	ordersentity "github.com/IndominusByte/catalyst-be/api/internal/entity/orders"
	"github.com/creent-production/cdk-go/parser"
	"github.com/creent-production/cdk-go/response"
)

type ordersUsecaseIface interface {
	Create(ctx context.Context, rw http.ResponseWriter, payload *ordersentity.JsonCreateSchema)
	GetById(ctx context.Context, rw http.ResponseWriter, orderId int)
}

func AddOrders(r *http.ServeMux, uc ordersUsecaseIface) {
	const prefix = "/order"

	r.HandleFunc(prefix, func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			var p ordersentity.JsonCreateSchema

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
			orderId, _ := parser.ParsePathToInt("/order/(.*)", r.URL.Path)

			uc.GetById(r.Context(), rw, orderId)
		default:
			response.WriteJSONResponse(rw, 404, nil, map[string]interface{}{
				constant.Body: constant.ResourceNotFound,
			})
		}
	})

}
