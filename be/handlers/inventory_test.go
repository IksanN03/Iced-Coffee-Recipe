package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInventoryEndpoints(t *testing.T) {
	r := setupTestRouter()

	// Test data
	inventoryItem := map[string]interface{}{
		"item_name":     "Mineral Water",
		"quantity":      1,
		"uom":           "Liter",
		"price_per_qty": 5000,
	}

	t.Run("Add Inventory", func(t *testing.T) {
		jsonData, _ := json.Marshal(inventoryItem)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/inventory", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", TestToken)
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("Get Inventory", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/inventory?search=Mineral Water", nil)
		req.Header.Set("Authorization", TestToken)
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("Update Inventory", func(t *testing.T) {
		jsonData, _ := json.Marshal(inventoryItem)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/inventory/1", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", TestToken)
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("Delete Inventory", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/inventory/1", nil)
		req.Header.Set("Authorization", TestToken)
		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})
}
