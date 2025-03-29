package middlewares_test

import (
	"encoding/json"
	"localEyes/internal/middlewares"
	"localEyes/internal/models"
	"localEyes/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticationMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		authHeader   string
		expectedCode int
		expectedBody *models.Response // Assuming this is the type for your error response
	}{
		{
			name:         "Missing token",
			authHeader:   "",
			expectedCode: http.StatusUnauthorized,
			expectedBody: utils.NewUnauthorizedError("Missing authentication token"),
		},
		{
			name:         "Invalid token",
			authHeader:   "Bearer invalid_token",
			expectedCode: http.StatusUnauthorized,
			expectedBody: utils.NewUnauthorizedError("Invalid token"),
		},
		{
			name:         "Valid token",
			authHeader:   "Bearer valid_token",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the ValidateTokenFunc
			utils.ValidateTokenFunc = func(token string) bool {
				if token == "Bearer valid_token" {
					return true
				}
				return false
			}

			// Create a request
			req, err := http.NewRequest("GET", "/test", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Create a ResponseRecorder to capture the response
			rr := httptest.NewRecorder()

			// Create a test handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK) // For valid token case
			})

			// Call the middleware
			middleware := middlewares.AuthenticationMiddleware(handler)
			middleware.ServeHTTP(rr, req)

			// Check the status code
			if rr.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, rr.Code)
			}

			// If we expect an error response, check the body
			if tt.expectedCode == http.StatusUnauthorized {
				var response models.Response
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				if response.Message != tt.expectedBody.Message {
					t.Errorf("Expected message '%s', got '%s'", tt.expectedBody.Message, response.Message)
				}
			}
		})
	}
}

func TestAdminAuthMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		authHeader   string
		expectedCode int
		expectedBody *models.Response // Assuming this is the type for your error response
	}{
		{
			name:         "Missing token",
			authHeader:   "",
			expectedCode: http.StatusUnauthorized,
			expectedBody: utils.NewUnauthorizedError("Missing authentication token"),
		},
		{
			name:         "Invalid token",
			authHeader:   "Bearer invalid_admin_token",
			expectedCode: http.StatusUnauthorized,
			expectedBody: utils.NewUnauthorizedError("Invalid token"),
		},
		{
			name:         "Valid admin token",
			authHeader:   "Bearer valid_admin_token",
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock the ValidateAdminTokenFunc
			utils.ValidateAdminTokenFunc = func(token string) bool {
				if token == "Bearer valid_admin_token" {
					return true
				}
				return false
			}

			// Create a request
			req, err := http.NewRequest("GET", "/admin/test", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Create a ResponseRecorder to capture the response
			rr := httptest.NewRecorder()

			// Create a test handler
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK) // For valid token case
			})

			// Call the middleware
			middleware := middlewares.AdminAuthMiddleware(handler)
			middleware.ServeHTTP(rr, req)

			// Check the status code
			if rr.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, rr.Code)
			}

			// If we expect an error response, check the body
			if tt.expectedCode == http.StatusUnauthorized {
				var response models.Response
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				if response.Message != tt.expectedBody.Message {
					t.Errorf("Expected message '%s', got '%s'", tt.expectedBody.Message, response.Message)
				}
			}
		})
	}
}
