package middlewares_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"localEyes/internal/middlewares"
	"localEyes/internal/models"
	"localEyes/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticationMiddleware(t *testing.T) {
	tests := []struct {
		name            string
		authHeader      string
		expectedCode    int
		expectedMessage string
		mockSetup       func()
	}{
		{
			name:            "Missing Authorization Header",
			authHeader:      "",
			expectedCode:    http.StatusUnauthorized,
			expectedMessage: "Missing authentication token",
		},
		{
			name:            "Invalid Token",
			authHeader:      "Bearer invalidtoken",
			expectedCode:    http.StatusUnauthorized,
			expectedMessage: "Invalid token",
		},
		//{
		//	name:            "Valid Token",
		//	authHeader:      "Bearer validtoken",
		//	expectedCode:    http.StatusOK,
		//	expectedMessage: "",
		//	mockSetup: func() {
		//		utils.ValidateTokenFunc = func(token string) bool {
		//			return true
		//		}
		//		defer func() { utils.ValidateTokenFunc = utils.ValidateToken }()
		//	},
		//},
	}

	//utils.NewUnauthorizedErrorFunc = func(message string) map[string]string {
	//	return map[string]string{"error": message}
	//}
	//defer func(){ utils.NewUnauthorizedErrorFunc=utils.NewUnauthorizedError}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}
			req := httptest.NewRequest("GET", "/some-endpoint", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			middleware := middlewares.AuthenticationMiddleware(handler)
			middleware.ServeHTTP(rr, req)

			var actualRes models.Response
			if err := json.Unmarshal(rr.Body.Bytes(), &actualRes); err != nil {
				t.Fatalf("failed to unmarshal actual response: %v", err)
			}

			assert.Equal(t, tt.expectedCode, actualRes.Code)
			assert.Equal(t, tt.expectedMessage, actualRes.Message)
		})
	}
}

func TestAdminAuthMiddleware(t *testing.T) {
	tests := []struct {
		name            string
		authHeader      string
		expectedCode    int
		expectedMessage string
	}{
		{
			name:            "Missing Authorization Header",
			authHeader:      "",
			expectedCode:    http.StatusUnauthorized,
			expectedMessage: "Missing authentication token",
		},
		{
			name:            "Invalid Admin Token",
			authHeader:      "Bearer invalidadmintoken",
			expectedCode:    http.StatusUnauthorized,
			expectedMessage: "Invalid token",
		},
	}

	// Mock the utils functions for testing
	utils.ValidateAdminTokenFunc = func(token string) bool {
		return token == "Bearer validadmintoken"
	}
	defer func() { utils.ValidateAdminTokenFunc = utils.ValidateAdminToken }()
	//utils.NewUnauthorizedError = func(message string) map[string]string {
	//	return map[string]string{"error": message}
	//}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/admin-endpoint", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			middleware := middlewares.AdminAuthMiddleware(handler)
			middleware.ServeHTTP(rr, req)

			var actualRes models.Response
			err := json.Unmarshal(rr.Body.Bytes(), &actualRes)
			if err != nil {
				t.Fatalf("failed to unmarshal actual response: %v", err)
			}
			assert.Equal(t, tt.expectedCode, actualRes.Code)
			assert.Equal(t, tt.expectedMessage, actualRes.Message)
		})
	}
}
