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

func TestInventoryOperations(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/inventory", GetInventory)
	r.POST("/inventory", AddInventory)

	t.Run("Get Inventory", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/inventory?page=1&limit=10", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	})

	t.Run("Add Inventory", func(t *testing.T) {
		item := map[string]interface{}{
			"item_name":     "Coffee Bean",
			"quantity":      1.0,
			"uom":           "kg",
			"price_per_qty": 100000,
		}
		jsonData, _ := json.Marshal(item)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/inventory", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	})
}
