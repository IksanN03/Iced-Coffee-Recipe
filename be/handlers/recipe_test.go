package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddRecipe(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/recipe", AddRecipe)

	tests := []struct {
		name     string
		recipe   map[string]interface{}
		wantCode int
	}{
		{
			name: "Valid Recipe",
			recipe: map[string]interface{}{
				"number_of_cups": 1,
				"ingredients": map[string]interface{}{
					"Aren Sugar": map[string]interface{}{
						"amount": 15,
						"unit":   "g",
					},
				},
			},
			wantCode: 200,
		},
		{
			name: "Invalid Units",
			recipe: map[string]interface{}{
				"number_of_cups": 1,
				"ingredients": map[string]interface{}{
					"Aren Sugar": map[string]interface{}{
						"amount": 15,
						"unit":   "invalid",
					},
				},
			},
			wantCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.recipe)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/recipe", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func TestGetRecipe(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.Default()
    r.GET("/recipe", GetRecipe)

    t.Run("Get Recipes", func(t *testing.T) {
        w := httptest.NewRecorder()
        req, _ := http.NewRequest("GET", "/recipe?page=1&limit=10", nil)
        r.ServeHTTP(w, req)
        assert.Equal(t, 200, w.Code)
    })
}

func TestUpdateRecipe(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.Default()
    r.PUT("/recipe/:id", UpdateRecipe)

    recipe := map[string]interface{}{
        "number_of_cups": 2,
        "ingredients": map[string]interface{}{
            "Aren Sugar": map[string]interface{}{
                "amount": 30,
                "unit":   "g",
            },
        },
    }
    jsonData, _ := json.Marshal(recipe)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("PUT", "/recipe/1", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    r.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)
}
