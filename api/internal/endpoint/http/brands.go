package endpoint_http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/IndominusByte/catalyst-be/api/internal/constant"
	brandsentity "github.com/IndominusByte/catalyst-be/api/internal/entity/brands"
	"github.com/creent-production/cdk-go/response"
)

type brandsUsecaseIface interface {
	Create(ctx context.Context, rw http.ResponseWriter, payload *brandsentity.JsonCreateSchema)
}

func AddBrands(r *http.ServeMux, uc brandsUsecaseIface) {
	const prefix = "/brand"

	r.HandleFunc(prefix, func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			var p brandsentity.JsonCreateSchema

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
}
