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

func TestSubmitEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/auth/submit-email", SubmitEmail)

	tests := []struct {
		name     string
		email    string
		wantCode int
	}{
		{"Valid Email", "test@example.com", 200},
		{"Invalid Email", "invalid-email", 400},
		{"Empty Email", "", 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := map[string]string{"email": tt.email}
			jsonData, _ := json.Marshal(body)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/auth/submit-email", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}
