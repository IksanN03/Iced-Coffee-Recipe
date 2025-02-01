package route

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSetupRoutes(t *testing.T) {
	router := gin.New()
	SetupRoutes(router)

	tests := []struct {
		name     string
		path     string
		method   string
		wantCode int
	}{
		{"Auth Submit Email", "/auth/submit-email", "POST", 400},
		{"Auth Magic Link", "/auth/magic-link", "GET", 400},
		{"Get Inventory", "/inventory", "GET", 401},
		{"Add Inventory", "/inventory", "POST", 401},
		{"Update Inventory", "/inventory/1", "PUT", 401},
		{"Delete Inventory", "/inventory/1", "DELETE", 401},
		{"Add Recipe", "/recipe", "POST", 401},
		// Add these test cases in the tests slice
		{"Get Recipe", "/recipe", "GET", 401},
		{"Update Recipe", "/recipe/1", "PUT", 401},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}
