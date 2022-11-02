package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	prefixProduct = "/product"
	nameProduct   = "asdasd"
	nameProduct2  = "asdasd2"
)

func TestValidationCreateProduct(t *testing.T) {
	_, s := setupEnvironment()

	var data map[string]interface{}

	tests := [...]struct {
		name    string
		payload map[string]interface{}
	}{
		{
			name:    "minimum",
			payload: map[string]interface{}{"name": "a", "description": "a", "price": -1, "brand_id": -1},
		},
		{
			name:    "maximum",
			payload: map[string]interface{}{"name": createMaximum(200), "description": createMaximum(200)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.payload)
			if err != nil {
				panic(err)
			}

			req, _ := http.NewRequest(http.MethodPost, prefixProduct, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			response := executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "minimum":
				assert.Equal(t, "Shorter than minimum length 3.", data["detail_message"].(map[string]interface{})["name"].(string))
				assert.Equal(t, "Shorter than minimum length 5.", data["detail_message"].(map[string]interface{})["description"].(string))
				assert.Equal(t, "Shorter than minimum length 1.", data["detail_message"].(map[string]interface{})["price"].(string))
				assert.Equal(t, "Shorter than minimum length 1.", data["detail_message"].(map[string]interface{})["brand_id"].(string))
			case "maximum":
				assert.Equal(t, "Longer than maximum length 100.", data["detail_message"].(map[string]interface{})["name"].(string))
			}
			assert.Equal(t, 422, response.Result().StatusCode)
		})
	}

}

func TestCreateProduct(t *testing.T) {
	repo, s := setupEnvironment()

	var data map[string]interface{}

	brand, _ := repo.brandsRepo.GetBrandByName(context.Background(), nameBrand)

	tests := [...]struct {
		name       string
		payload    map[string]interface{}
		expected   string
		statusCode int
	}{
		{
			name:       "brand not found",
			payload:    map[string]interface{}{"name": nameProduct, "description": "asdasd", "price": 1, "brand_id": 999999999},
			expected:   "Brand not found.",
			statusCode: 404,
		},
		{
			name:       "success",
			payload:    map[string]interface{}{"name": nameProduct, "description": "asdasd", "price": 1, "brand_id": brand.Id},
			expected:   "Successfully add a new product.",
			statusCode: 201,
		},
		{
			name:       "name duplicate",
			payload:    map[string]interface{}{"name": nameProduct, "description": "asdasd", "price": 1, "brand_id": brand.Id},
			expected:   "The name has already been taken.",
			statusCode: 400,
		},
		{
			name:       "success 2",
			payload:    map[string]interface{}{"name": nameProduct2, "description": "asdasd", "price": 1, "brand_id": brand.Id},
			expected:   "Successfully add a new product.",
			statusCode: 201,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.payload)
			if err != nil {
				panic(err)
			}

			req, _ := http.NewRequest(http.MethodPost, prefixProduct, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			response := executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			assert.Equal(t, test.expected, data["detail_message"].(map[string]interface{})["_app"].(string))
			assert.Equal(t, test.statusCode, response.Result().StatusCode)
		})
	}
}

func TestValidationGetByIdProduct(t *testing.T) {
	_, s := setupEnvironment()

	var data map[string]interface{}

	tests := [...]struct {
		name string
		url  string
	}{
		{
			name: "type data",
			url:  prefixProduct + "/abc",
		},
		{
			name: "minimum",
			url:  prefixProduct + "/-1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, test.url, nil)

			response := executeRequest(req, s)

			body, _ := io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			assert.Equal(t, "Product not found.", data["detail_message"].(map[string]interface{})["_app"].(string))
			assert.Equal(t, 404, response.Result().StatusCode)
		})
	}

}

func TestGetByIdProduct(t *testing.T) {
	repo, s := setupEnvironment()

	var data map[string]interface{}

	// get id
	product, _ := repo.productsRepo.GetProductByName(context.Background(), nameProduct)

	tests := [...]struct {
		name string
		url  string
	}{
		{
			name: "not found",
			url:  prefixProduct + "/99999999",
		},
		{
			name: "success",
			url:  prefixProduct + "/" + strconv.Itoa(product.Id),
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
				assert.Equal(t, "Product not found.", data["detail_message"].(map[string]interface{})["_app"].(string))
				assert.Equal(t, 404, response.Result().StatusCode)
			case "success":
				assert.NotNil(t, data["results"])
				assert.Equal(t, 200, response.Result().StatusCode)
			}
		})
	}
}

func TestValidationGetByBrandId(t *testing.T) {
	_, s := setupEnvironment()

	var data map[string]interface{}

	tests := [...]struct {
		name string
		url  string
	}{
		{
			name: "empty",
			url:  prefixProduct + "/brand",
		},
		{
			name: "type data",
			url:  prefixProduct + "/brand" + "?page=a&per_page=a&brand_id=a",
		},
		{
			name: "minimum",
			url:  prefixProduct + "/brand" + "?page=-1&per_page=-1&brand_id=-1",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, test.url, nil)

			response := executeRequest(req, s)

			body, _ := io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "empty":
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["page"].(string))
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["per_page"].(string))
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["brand_id"].(string))
			case "type data":
				assert.Equal(t, "Invalid input type.", data["detail_message"].(map[string]interface{})["_body"].(string))
			case "minimum":
				assert.Equal(t, "Must be greater than or equal to 1.", data["detail_message"].(map[string]interface{})["page"].(string))
				assert.Equal(t, "Must be greater than or equal to 1.", data["detail_message"].(map[string]interface{})["per_page"].(string))
				assert.Equal(t, "Must be greater than or equal to 1.", data["detail_message"].(map[string]interface{})["brand_id"].(string))
			}
			assert.Equal(t, 422, response.Result().StatusCode)
		})
	}

}

func TestGetByBrandId(t *testing.T) {
	repo, s := setupEnvironment()

	brand, _ := repo.brandsRepo.GetBrandByName(context.Background(), nameBrand)
	var data map[string]interface{}

	req, _ := http.NewRequest(http.MethodGet, prefixProduct+"/brand?page=1&per_page=1&brand_id="+strconv.Itoa(brand.Id), nil)

	response := executeRequest(req, s)

	body, _ := io.ReadAll(response.Result().Body)
	json.Unmarshal(body, &data)

	assert.NotNil(t, data["results"].(map[string]interface{})["data"])
	assert.Equal(t, 200, response.Result().StatusCode)
}
