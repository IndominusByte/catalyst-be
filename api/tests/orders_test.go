package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"

	productsentity "github.com/IndominusByte/catalyst-be/api/internal/entity/products"
	"github.com/stretchr/testify/assert"
)

const (
	prefixOrder      = "/order"
	nameProductOrder = "asdasdasdasd"
)

func TestUpOrder(t *testing.T) {
	repo, _ := setupEnvironment()

	brand, _ := repo.brandsRepo.GetBrandByName(context.Background(), nameBrand)
	// create product
	repo.productsRepo.Insert(context.Background(), &productsentity.JsonCreateSchema{
		Name:        nameProductOrder,
		Description: "asdasd",
		Price:       1,
		BrandId:     brand.Id,
	})
}
func TestValidationCreateOrder(t *testing.T) {
	_, s := setupEnvironment()

	var data map[string]interface{}

	tests := [...]struct {
		name    string
		payload map[string]int
	}{
		{
			name:    "required",
			payload: map[string]int{"buyer_id": 0, "product_id": 0, "qty": 0},
		},
		{
			name:    "minimum",
			payload: map[string]int{"buyer_id": -1, "product_id": -1, "qty": -1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.payload)
			if err != nil {
				panic(err)
			}

			req, _ := http.NewRequest(http.MethodPost, prefixOrder, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			response := executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "required":
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["buyer_id"].(string))
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["product_id"].(string))
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["qty"].(string))
			case "minimum":
				assert.Equal(t, "Shorter than minimum length 1.", data["detail_message"].(map[string]interface{})["buyer_id"].(string))
				assert.Equal(t, "Shorter than minimum length 1.", data["detail_message"].(map[string]interface{})["product_id"].(string))
				assert.Equal(t, "Shorter than minimum length 1.", data["detail_message"].(map[string]interface{})["qty"].(string))
			}
			assert.Equal(t, 422, response.Result().StatusCode)
		})
	}

}

func TestCreateOrder(t *testing.T) {
	repo, s := setupEnvironment()
	// get id
	product, _ := repo.productsRepo.GetProductByName(context.Background(), nameProductOrder)

	var data map[string]interface{}

	tests := [...]struct {
		name       string
		payload    map[string]int
		expected   string
		statusCode int
	}{
		{
			name:       "user not found",
			payload:    map[string]int{"buyer_id": 999, "product_id": 999, "qty": 1},
			expected:   "User not found.",
			statusCode: 404,
		},
		{
			name:       "product not found",
			payload:    map[string]int{"buyer_id": 1, "product_id": 999, "qty": 1},
			expected:   "Product not found.",
			statusCode: 404,
		},

		{
			name:       "success",
			payload:    map[string]int{"buyer_id": 1, "product_id": product.Id, "qty": 1},
			expected:   "Successfully add a new order.",
			statusCode: 201,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.payload)
			if err != nil {
				panic(err)
			}

			req, _ := http.NewRequest(http.MethodPost, prefixOrder, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			response := executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			assert.Equal(t, test.expected, data["detail_message"].(map[string]interface{})["_app"].(string))
			assert.Equal(t, test.statusCode, response.Result().StatusCode)
		})
	}
}

func TestGetOrderById(t *testing.T) {
	repo, s := setupEnvironment()

	order, _ := repo.ordersRepo.GetOrderByBuyerId(context.Background(), 1)

	var data map[string]interface{}

	tests := [...]struct {
		name string
		url  string
	}{
		{
			name: "not found",
			url:  prefixOrder + "/99999999",
		},
		{
			name: "success",
			url:  prefixOrder + "/" + strconv.Itoa(order.Id),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, test.url, nil)

			response := executeRequest(req, s)

			body, _ := io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "not found":
				assert.Equal(t, "Order not found.", data["detail_message"].(map[string]interface{})["_app"].(string))
				assert.Equal(t, 404, response.Result().StatusCode)
			case "success":
				assert.NotNil(t, data["results"])
				assert.Equal(t, 200, response.Result().StatusCode)
			}
		})
	}
}
