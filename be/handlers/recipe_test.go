package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipeEndpoints(t *testing.T) {
	r := setupTestRouter()

	recipeData := map[string]interface{}{
		"number_of_cups": 1,
		"ingredients": map[string]interface{}{
			"Aren Sugar": map[string]interface{}{
				"amount": 15,
				"unit":   "g",
			},
			"Milk": map[string]interface{}{
				"amount": 150,
				"unit":   "ml",
			},
			"Ice Cube": map[string]interface{}{
				"amount": 20,
				"unit":   "g",
			},
			"Plastic Cup": map[string]interface{}{
				"amount": 1,
				"unit":   "pcs",
			},
			"Coffee Bean": map[string]interface{}{
				"amount": 20,
				"unit":   "g",
			},
			"Mineral Water": map[string]interface{}{
				"amount": 50,
				"unit":   "ml",
			},
		},
	}

	t.Run("Add Recipe", func(t *testing.T) {
		jsonData, _ := json.Marshal(recipeData)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/recipe", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", TestToken)
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})
		assert.NotNil(t, data["sku"])
		assert.Equal(t, float64(13250), data["cogs"])
	})

	t.Run("Get Recipe", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/recipe?search=IC-20250131-001", nil)
		req.Header.Set("Authorization", TestToken)
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("Update Recipe", func(t *testing.T) {
		recipeData["number_of_cups"] = 2
		jsonData, _ := json.Marshal(recipeData)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/recipe/1", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", TestToken)
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].(map[string]interface{})
		assert.Equal(t, float64(26500), data["cogs"])
	})
}
