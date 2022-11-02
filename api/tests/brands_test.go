package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	prefixBrand = "/brand"
	nameBrand   = "asdasd"
	nameBrand2  = "asdasd2"
)

func TestValidationCreateBrand(t *testing.T) {
	_, s := setupEnvironment()

	var data map[string]interface{}

	tests := [...]struct {
		name    string
		payload map[string]string
	}{
		{
			name:    "required",
			payload: map[string]string{"name": ""},
		},
		{
			name:    "minimum",
			payload: map[string]string{"name": "a"},
		},
		{
			name:    "maximum",
			payload: map[string]string{"name": createMaximum(200)},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.payload)
			if err != nil {
				panic(err)
			}

			req, _ := http.NewRequest(http.MethodPost, prefixBrand, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			response := executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			switch test.name {
			case "required":
				assert.Equal(t, "Missing data for required field.", data["detail_message"].(map[string]interface{})["name"].(string))
			case "minimum":
				assert.Equal(t, "Shorter than minimum length 3.", data["detail_message"].(map[string]interface{})["name"].(string))
			case "maximum":
				assert.Equal(t, "Longer than maximum length 100.", data["detail_message"].(map[string]interface{})["name"].(string))
			}
			assert.Equal(t, 422, response.Result().StatusCode)
		})
	}
}

func TestCreateBrand(t *testing.T) {
	_, s := setupEnvironment()

	var data map[string]interface{}

	tests := [...]struct {
		name       string
		payload    map[string]string
		expected   string
		statusCode int
	}{
		{
			name:       "success",
			payload:    map[string]string{"name": nameBrand},
			expected:   "Successfully add a new brand.",
			statusCode: 201,
		},
		{
			name:       "success",
			payload:    map[string]string{"name": nameBrand2},
			expected:   "Successfully add a new brand.",
			statusCode: 201,
		},
		{
			name:       "name already taken",
			payload:    map[string]string{"name": nameBrand},
			expected:   "The name has already been taken.",
			statusCode: 400,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.payload)
			if err != nil {
				panic(err)
			}

			req, _ := http.NewRequest(http.MethodPost, prefixBrand, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			response := executeRequest(req, s)

			body, _ = io.ReadAll(response.Result().Body)
			json.Unmarshal(body, &data)

			assert.Equal(t, test.expected, data["detail_message"].(map[string]interface{})["_app"].(string))
			assert.Equal(t, test.statusCode, response.Result().StatusCode)
		})
	}
}
