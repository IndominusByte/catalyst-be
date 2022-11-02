package orders

import (
	"context"
	"fmt"
	"net/http"

	"github.com/IndominusByte/catalyst-be/api/internal/constant"
	ordersentity "github.com/IndominusByte/catalyst-be/api/internal/entity/orders"
	"github.com/creent-production/cdk-go/response"
	"github.com/creent-production/cdk-go/validation"
)

type OrdersUsecase struct {
	ordersRepo   ordersRepo
	productsRepo productsRepo
	usersRepo    usersRepo
}

func NewOrdersUsecase(orderRepo ordersRepo, productRepo productsRepo, userRepo usersRepo) *OrdersUsecase {
	return &OrdersUsecase{
		ordersRepo:   orderRepo,
		productsRepo: productRepo,
		usersRepo:    userRepo,
	}
}

func (uc *OrdersUsecase) Create(ctx context.Context, rw http.ResponseWriter, payload *ordersentity.JsonCreateSchema) {
	if err := validation.StructValidate(payload); err != nil {
		response.WriteJSONResponse(rw, 422, nil, err)
		return
	}

	// check user id not found
	if _, err := uc.usersRepo.GetUserById(ctx, payload.BuyerId); err != nil {
		response.WriteJSONResponse(rw, 404, nil, map[string]interface{}{
			constant.App: "User not found.",
		})
		return
	}

	product, err := uc.productsRepo.GetProductById(ctx, payload.ProductId)
	if err != nil {
		response.WriteJSONResponse(rw, 404, nil, map[string]interface{}{
			constant.App: "Product not found.",
		})
		return
	}

	// insert into db
	payload.Price = product.Price
	payload.TotalPrice = payload.Qty * payload.Price
	uc.ordersRepo.Insert(ctx, payload)

	response.WriteJSONResponse(rw, 201, nil, map[string]interface{}{
		constant.App: "Successfully add a new order.",
	})
}

func (uc *OrdersUsecase) GetById(ctx context.Context, rw http.ResponseWriter, orderId int) {
	t, err := uc.ordersRepo.GetOrderById(ctx, orderId)
	fmt.Println(err)
	if err != nil {
		response.WriteJSONResponse(rw, 404, nil, map[string]interface{}{
			constant.App: "Order not found.",
		})
		return
	}
	response.WriteJSONResponse(rw, 200, t, nil)
}
