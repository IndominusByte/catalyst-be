package products

import (
	"context"
	"fmt"
	"net/http"

	"github.com/IndominusByte/catalyst-be/api/internal/constant"
	productsentity "github.com/IndominusByte/catalyst-be/api/internal/entity/products"
	"github.com/creent-production/cdk-go/response"
	"github.com/creent-production/cdk-go/validation"
)

type ProductsUsecase struct {
	productsRepo productsRepo
	brandsRepo   brandsRepo
}

func NewProductsUsecase(productRepo productsRepo, brandRepo brandsRepo) *ProductsUsecase {
	return &ProductsUsecase{
		productsRepo: productRepo,
		brandsRepo:   brandRepo,
	}
}

func (uc *ProductsUsecase) Create(ctx context.Context, rw http.ResponseWriter, payload *productsentity.JsonCreateSchema) {
	if err := validation.StructValidate(payload); err != nil {
		response.WriteJSONResponse(rw, 422, nil, err)
		return
	}

	// check name duplicate
	if _, err := uc.productsRepo.GetProductByName(ctx, payload.Name); err == nil {
		response.WriteJSONResponse(rw, 400, nil, map[string]interface{}{
			constant.App: fmt.Sprintf(constant.AlreadyTaken, "name"),
		})
		return
	}
	// check brand id not found
	if _, err := uc.brandsRepo.GetBrandById(ctx, payload.BrandId); err != nil {
		response.WriteJSONResponse(rw, 404, nil, map[string]interface{}{
			constant.App: "Brand not found.",
		})
		return
	}

	// insert into db
	uc.productsRepo.Insert(ctx, payload)

	response.WriteJSONResponse(rw, 201, nil, map[string]interface{}{
		constant.App: "Successfully add a new product.",
	})
}

func (uc *ProductsUsecase) GetById(ctx context.Context, rw http.ResponseWriter, productId int) {
	t, err := uc.productsRepo.GetProductById(ctx, productId)
	if err != nil {
		response.WriteJSONResponse(rw, 404, nil, map[string]interface{}{
			constant.App: "Product not found.",
		})
		return
	}
	response.WriteJSONResponse(rw, 200, t, nil)
}

func (uc *ProductsUsecase) GetByBrandId(ctx context.Context, rw http.ResponseWriter, payload *productsentity.QueryParamAllProductSchema) {

	if err := validation.StructValidate(payload); err != nil {
		response.WriteJSONResponse(rw, 422, nil, err)
		return
	}

	if _, err := uc.brandsRepo.GetBrandById(ctx, payload.BrandId); err != nil {
		response.WriteJSONResponse(rw, 404, nil, map[string]interface{}{
			constant.App: "Brand not found.",
		})
		return
	}

	results, _ := uc.productsRepo.GetAllProductByBrandPaginate(ctx, payload)

	response.WriteJSONResponse(rw, 200, results, nil)
}
