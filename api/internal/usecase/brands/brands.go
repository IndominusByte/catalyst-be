package brands

import (
	"context"
	"fmt"
	"net/http"

	"github.com/IndominusByte/catalyst-be/api/internal/constant"
	brandsentity "github.com/IndominusByte/catalyst-be/api/internal/entity/brands"
	"github.com/creent-production/cdk-go/response"
	"github.com/creent-production/cdk-go/validation"
)

type BrandsUsecase struct {
	brandsRepo brandsRepo
}

func NewBrandsUsecase(brandRepo brandsRepo) *BrandsUsecase {
	return &BrandsUsecase{
		brandsRepo: brandRepo,
	}
}

func (uc *BrandsUsecase) Create(ctx context.Context, rw http.ResponseWriter, payload *brandsentity.JsonCreateSchema) {
	if err := validation.StructValidate(payload); err != nil {
		response.WriteJSONResponse(rw, 422, nil, err)
		return
	}

	if _, err := uc.brandsRepo.GetBrandByName(ctx, payload.Name); err == nil {
		response.WriteJSONResponse(rw, 400, nil, map[string]interface{}{
			constant.App: fmt.Sprintf(constant.AlreadyTaken, "name"),
		})
		return
	}

	// save into database
	uc.brandsRepo.Insert(ctx, payload)

	response.WriteJSONResponse(rw, 201, nil, map[string]interface{}{
		constant.App: "Successfully add a new brand.",
	})
}
