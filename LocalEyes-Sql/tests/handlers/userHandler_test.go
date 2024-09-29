package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"localEyes/internal/handlers"
	"localEyes/internal/models"
	"localEyes/tests/mocks"
	"localEyes/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_SignUp_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockUserServiceInterface(ctrl)
	userHandler := handlers.NewUserHandler(mockService)

	client := models.Client{
		Username: "testuser",
		Password: "StrongPassword@123",
		LivingSince: models.LivingSince{
			Days:   1,
			Months: 0,
			Years:  1,
		},
	}

	mockService.EXPECT().Signup(client.Username, client.Password, 1, gomock.Any()).Return(nil)
	mockService.EXPECT().ValidateUsername(client.Username).Return(true)
	body, _ := json.Marshal(client)
	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	userHandler.SignUp(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201, got %d", res.StatusCode)
	}

	var response models.Response
	json.NewDecoder(res.Body).Decode(&response)
	if response.Message != "User created successfully" {
		t.Errorf("expected message 'User created successfully', got '%s'", response.Message)
	}
}

func TestUserHandler_SignUp_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockUserServiceInterface(ctrl)
	userHandler := handlers.NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString("{invalid-json}"))
	w := httptest.NewRecorder()

	userHandler.SignUp(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", res.StatusCode)
	}
}

func TestUserHandler_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockUserServiceInterface(ctrl)
	userHandler := handlers.NewUserHandler(mockService)

	client := models.Client{
		Username: "testuser",
		Password: "StrongPassword123",
	}

	mockService.EXPECT().Login(client.Username, client.Password).Return(&models.User{Username: client.Username}, nil)
	mockToken := "generated.token.here"
	utils.GenerateTokenFunc = func(username string, uid string) (string, error) {
		return mockToken, nil
	}
	defer func() { utils.GenerateTokenFunc = utils.GenerateToken }()
	body, _ := json.Marshal(client)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	userHandler.Login(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	var response models.Response
	json.NewDecoder(res.Body).Decode(&response)
	if response.Data != mockToken {
		t.Errorf("expected token '%s', got '%s'", mockToken, response.Data)
	}
}

func TestUserHandler_DeActivate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockUserServiceInterface(ctrl)
	userHandler := handlers.NewUserHandler(mockService)

	// Assuming a valid bearer token is passed
	req := httptest.NewRequest(http.MethodPost, "/deactivate", nil)
	req.Header.Set("Authorization", "Bearer valid.token.here")
	w := httptest.NewRecorder()
	utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
		return jwt.MapClaims{"id": "userId"}, nil
	}
	defer func() { utils.ExtractClaimsFunc = utils.ExtractClaims }() // Reset after test
	mockService.EXPECT().DeActivate("userId").Return(nil)

	userHandler.DeActivate(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}
}

func TestUserHandler_ViewProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockUserServiceInterface(ctrl)
	userHandler := handlers.NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req.Header.Set("Authorization", "Bearer valid.token.here")
	w := httptest.NewRecorder()
	utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
		return jwt.MapClaims{"id": "userId"}, nil
	}
	defer func() { utils.ExtractClaimsFunc = utils.ExtractClaims }() // Reset after test

	mockService.EXPECT().GetUserById("userId").Return(&models.User{Username: "testuser"}, nil)

	userHandler.ViewProfile(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}
}

func TestUserHandler_ViewNotifications_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockUserServiceInterface(ctrl)
	userHandler := handlers.NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/notifications", nil)
	req.Header.Set("Authorization", "Bearer valid.token.here")
	w := httptest.NewRecorder()
	result := []string{}
	utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
		return jwt.MapClaims{"id": "userId"}, nil
	}
	defer func() { utils.ExtractClaimsFunc = utils.ExtractClaims }() // Reset after test

	mockService.EXPECT().GetNotifications("userId").Return(&result, nil)

	userHandler.ViewNotifications(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}
}

func TestUserHandler_SignUp_ErrorCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockUserServiceInterface(ctrl)
	handler := handlers.NewUserHandler(mockService)

	tests := []struct {
		name         string
		body         string
		mockSetup    func()
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "Invalid JSON body",
			body:         "invalid json",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid JSON body",
		},
		{
			name: "Username not available",
			body: `{"username": "existingUser", "password": "Password123", "living_since": {}}`,
			mockSetup: func() {
				mockService.EXPECT().ValidateUsername("existingUser").Return(false)
			},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Username not available",
		},
		{
			name: "Weak password",
			body: `{"username": "newUser", "password": "weak", "living_since": {}}`,
			mockSetup: func() {
				mockService.EXPECT().ValidateUsername("newUser").Return(true)
			},
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Password not strong",
		},
		{
			name: "Error signing up",
			body: `{"username": "newUser", "password": "Password@123", "living_since": {}}`,
			mockSetup: func() {
				mockService.EXPECT().ValidateUsername("newUser").Return(true)
				mockService.EXPECT().Signup("newUser", "Password@123", 0, "newbie").Return(errors.New("signup error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "Error signing up",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer([]byte(tt.body)))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.SignUp(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)

			var actualResponse models.Response
			err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMsg, actualResponse.Message)
		})
	}
}

func TestUserHandler_Login_ErrorCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockUserServiceInterface(ctrl)
	handler := handlers.NewUserHandler(mockService)

	tests := []struct {
		name         string
		body         string
		mockSetup    func()
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "Invalid JSON body",
			body:         "invalid json",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid JSON body",
		},
		{
			name: "Unauthorized login",
			body: `{"username": "user", "password": "wrong"}`,
			mockSetup: func() {
				mockService.EXPECT().Login("user", "wrong").Return(nil, errors.New("invalid credentials"))
			},
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "invalid credentials",
		},
		{
			name: "Error generating token",
			body: `{"username": "user", "password": "Password@123"}`,
			mockSetup: func() {
				mockService.EXPECT().Login("user", "Password@123").Return(&models.User{Username: "user"}, nil)
				utils.GenerateTokenFunc = func(username string, uid string) (string, error) {
					return "", errors.New("token error")
				}
			},
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "Error generating token ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req, err := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(tt.body)))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.Login(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)

			var actualResponse struct {
				Token   string `json:"token"`
				Code    int    `json:"code"`
				Message string `json:"message"`
			}
			err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMsg, actualResponse.Message)
		})
	}
}

func TestUserHandler_DeActivate_ErrorCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockUserServiceInterface(ctrl)
	handler := handlers.NewUserHandler(mockService)

	tests := []struct {
		name         string
		bearerToken  string
		mockSetup    func()
		expectedCode int
		expectedMsg  string
	}{
		{
			name:        "Invalid token",
			bearerToken: "Bearer invalid",
			mockSetup: func() {
				utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
					return nil, errors.New("invalid token")
				}
			},
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "Invalid token",
		},
		{
			name:        "Error deactivating user",
			bearerToken: "Bearer valid",
			mockSetup: func() {
				utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
					return jwt.MapClaims{"id": "userId"}, nil
				}
				mockService.EXPECT().DeActivate("userId").Return(errors.New("deactivation error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "Error deactivating user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req, err := http.NewRequest("DELETE", "/deactivate", nil)
			assert.NoError(t, err)
			req.Header.Set("Authorization", tt.bearerToken)

			rr := httptest.NewRecorder()
			handler.DeActivate(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)

			var actualResponse models.Response
			err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMsg, actualResponse.Message)
		})
	}
}

func TestUserHandler_ViewProfile_ErrorCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockUserServiceInterface(ctrl)
	handler := handlers.NewUserHandler(mockService)

	tests := []struct {
		name         string
		bearerToken  string
		mockSetup    func()
		expectedCode int
		expectedMsg  string
	}{
		{
			name:        "Invalid token",
			bearerToken: "Bearer invalid",
			mockSetup: func() {
				utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
					return nil, errors.New("invalid token")
				}
			},
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "Error extracting claims",
		},
		{
			name:        "Unauthorized user",
			bearerToken: "Bearer valid",
			mockSetup: func() {
				utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
					return jwt.MapClaims{"id": "userId"}, nil
				}
				mockService.EXPECT().GetUserById("userId").Return(nil, errors.New("unauthorized"))
			},
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "Unauthorized user",
		},
		{
			name:        "Error retrieving user",
			bearerToken: "Bearer valid",
			mockSetup: func() {
				utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
					return jwt.MapClaims{"id": "userId"}, nil
				}
				mockService.EXPECT().GetUserById("userId").Return(nil, errors.New("user not found"))
			},
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "Unauthorized user",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req, err := http.NewRequest("GET", "/profile", nil)
			assert.NoError(t, err)
			req.Header.Set("Authorization", tt.bearerToken)

			rr := httptest.NewRecorder()
			handler.ViewProfile(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)

			var actualResponse models.Response
			err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMsg, actualResponse.Message)
		})
	}
}

func TestUserHandler_ViewNotifications_ErrorCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockUserServiceInterface(ctrl)
	handler := handlers.NewUserHandler(mockService)

	tests := []struct {
		name         string
		bearerToken  string
		mockSetup    func()
		expectedCode int
		expectedMsg  string
	}{
		{
			name:        "Invalid token",
			bearerToken: "Bearer invalid",
			mockSetup: func() {
				utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
					return nil, errors.New("invalid token")
				}
			},
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "Invalid token",
		},
		{
			name:        "No user found with that id",
			bearerToken: "Bearer valid",
			mockSetup: func() {
				utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
					return jwt.MapClaims{"id": "userId"}, nil
				}
				mockService.EXPECT().GetNotifications("userId").Return(nil, utils.NoUser)
			},
			expectedCode: http.StatusNotFound,
			expectedMsg:  "No user found with that id",
		},
		{
			name:        "Error getting notifications",
			bearerToken: "Bearer valid",
			mockSetup: func() {
				utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
					return jwt.MapClaims{"id": "userId"}, nil
				}
				mockService.EXPECT().GetNotifications("userId").Return(nil, errors.New("error getting notifications"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "Error getting notifications",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			req, err := http.NewRequest("GET", "/notifications", nil)
			assert.NoError(t, err)
			req.Header.Set("Authorization", tt.bearerToken)

			rr := httptest.NewRecorder()
			handler.ViewNotifications(rr, req)

			assert.Equal(t, tt.expectedCode, rr.Code)

			var actualResponse models.Response
			err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedMsg, actualResponse.Message)
		})
	}
}
