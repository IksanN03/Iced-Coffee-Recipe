package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var TestToken = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Imlrc2FubnVyc2FsaW0xMjM0NTZAZ21haWwuY29tIiwiZXhwIjoxNzM5MzMzMTA0fQ.6eKh4Whm0xb4eUJ_zaOOPAJujIrOeEaocO_kmL-Ci58"

func TestAuthEndpoints(t *testing.T) {
	r := setupTestRouter()

	t.Run("Submit Email", func(t *testing.T) {
		tests := []struct {
			name     string
			email    string
			wantCode int
		}{
			{"Valid Email", "iksannursalim123456@gmail.com", 200},
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
	})

	t.Run("Verify Magic Link", func(t *testing.T) {
		tests := []struct {
			name     string
			token    string
			wantCode int
		}{
			{"Valid Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Imlrc2FubnVyc2FsaW0xMjM0NTZAZ21haWwuY29tIiwiZXhwIjoxNzM5MjQ2OTYyfQ.BbpzHRu_wLMm8e2oEgzkwj_O33jCL2QVi32Vca2Bdfo", 200},
			{"Invalid Token", "invalid-token", 401},
			{"Empty Token", "", 400},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/auth/magic-link?token="+tt.token, nil)
				r.ServeHTTP(w, req)

				assert.Equal(t, tt.wantCode, w.Code)
				if tt.wantCode == 200 {
					var response map[string]interface{}
					json.Unmarshal(w.Body.Bytes(), &response)
					data := response["data"].(map[string]interface{})
					TestToken = data["access_token"].(string)
					fmt.Printf("Using TestToken: %s\n", TestToken)

				}
			})
		}
	})
}
