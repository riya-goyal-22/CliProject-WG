package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"localEyes/internal/handlers"
	"localEyes/internal/models"
	"localEyes/tests/mocks"
	"localEyes/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPostHandler_CreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockPostServiceInterface(ctrl)
	handler := handlers.NewPostHandler(mockService)

	requestBody := `{"title": "Test Post", "content": "This is a test.", "type": "travel"}`
	req, err := http.NewRequest("POST", "/posts", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer valid-token")

	// Mock the ExtractClaims function
	utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
		return jwt.MapClaims{"id": "userId"}, nil
	}
	defer func() { utils.ExtractClaimsFunc = utils.ExtractClaims }() // Reset after test

	mockService.EXPECT().CreatePost("userId", "Test Post", "This is a test.", "travel").Return(nil)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/posts", handler.CreatePost).Methods("POST")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
	expectedResponse := models.Response{
		Code:    http.StatusOK,
		Message: "Post created successfully",
		Data:    nil,
	}
	var actualRes models.Response
	if err := json.Unmarshal(rr.Body.Bytes(), &actualRes); err != nil {
		t.Fatalf("failed to unmarshal actual response: %v", err)
	}
	if actualRes.Code != expectedResponse.Code {
		t.Errorf("handler returned wrong status code: got %v want %v", actualRes.Code, expectedResponse.Code)
	}
	if actualRes.Message != expectedResponse.Message {
		t.Errorf("handler returned wrong message: got %v want %v", actualRes.Message, expectedResponse.Message)
	}
}

func TestPostHandler_CreatePost_ErrorCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockPostServiceInterface(ctrl)
	handler := handlers.NewPostHandler(mockService)

	tests := []struct {
		name         string
		body         string
		mockSetup    func()
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "Missing request body",
			body:         "",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Missing Request body",
		},
		{
			name: "Invalid token",
			body: `{"title": "Test Post", "content": "This is a test.", "type": "travel"}`,
			mockSetup: func() {
				utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
					return nil, errors.New("invalid token")
				}
			},
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "Invalid token",
		},
		{
			name: "Error creating post",
			body: `{"title": "Test Post", "content": "This is a test.", "type": "travel"}`,
			mockSetup: func() {
				utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
					return jwt.MapClaims{"id": "userId"}, nil
				}
				mockService.EXPECT().CreatePost("userId", "Test Post", "This is a test.", "travel").Return(errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "Error creating post",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/posts", bytes.NewBuffer([]byte(tt.body)))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", "Bearer valid-token")

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/posts", handler.CreatePost).Methods("POST")

			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}

			var actualRes models.Response
			if err := json.Unmarshal(rr.Body.Bytes(), &actualRes); err != nil {
				t.Fatalf("failed to unmarshal actual response: %v", err)
			}

			if actualRes.Code != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", actualRes.Code, tt.expectedCode)
			}
			if actualRes.Message != tt.expectedMsg {
				t.Errorf("handler returned wrong message: got %v want %v", actualRes.Message, tt.expectedMsg)
			}
		})
	}
}

func TestPostHandler_DisplayPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockPostServiceInterface(ctrl)
	handler := handlers.NewPostHandler(mockService)

	posts := []*models.Post{
		{PostId: "1", UId: "userId", Title: "Post 1", Content: "Content 1", Type: "travel", Likes: 0},
	}
	mockService.EXPECT().GiveAllPosts().Return(posts, nil)

	req, err := http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/posts", handler.DisplayPosts).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedResponse := models.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    posts,
	}
	var actualRes models.Response
	if err := json.Unmarshal(rr.Body.Bytes(), &actualRes); err != nil {
		t.Fatalf("failed to unmarshal actual response: %v", err)
	}
	if actualRes.Code != expectedResponse.Code {
		t.Errorf("handler returned wrong status code: got %v want %v", actualRes.Code, expectedResponse.Code)
	}
	if actualRes.Message != expectedResponse.Message {
		t.Errorf("handler returned wrong message: got %v want %v", actualRes.Message, expectedResponse.Message)
	}
}

func TestPostHandler_UpdatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockPostServiceInterface(ctrl)
	handler := handlers.NewPostHandler(mockService)

	postId := "1"
	requestBody := `{"title": "Updated Post", "content": "This is an updated test."}`
	req, err := http.NewRequest("PUT", "/posts/"+postId, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer valid-token")
	// Mock the ExtractClaims function
	utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
		return jwt.MapClaims{"id": "userId"}, nil
	}
	defer func() { utils.ExtractClaimsFunc = utils.ExtractClaims }() // Reset after test

	mockService.EXPECT().UpdateMyPost(postId, "userId", "Updated Post", "This is an updated test.").Return(nil)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/posts/{post_id}", handler.UpdatePost).Methods("PUT")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expectedResponse := models.Response{
		Code:    http.StatusOK,
		Message: "Post updated successfully",
		Data:    nil,
	}
	var actualRes models.Response
	if err := json.Unmarshal(rr.Body.Bytes(), &actualRes); err != nil {
		t.Fatalf("failed to unmarshal actual response: %v", err)
	}
	if actualRes.Code != expectedResponse.Code {
		t.Errorf("handler returned wrong status code: got %v want %v", actualRes.Code, expectedResponse.Code)
	}
	if actualRes.Message != expectedResponse.Message {
		t.Errorf("handler returned wrong message: got %v want %v", actualRes.Message, expectedResponse.Message)
	}
}

func TestPostHandler_UpdatePost_ErrorCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockPostServiceInterface(ctrl)
	handler := handlers.NewPostHandler(mockService)

	tests := []struct {
		name         string
		body         string
		mockSetup    func()
		expectedCode int
		expectedMsg  string
	}{
		{
			name:         "Invalid/Missing request body",
			body:         "",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid/Missing Request body",
		},
		{
			name:         "Missing required fields",
			body:         `{"title": "", "content": ""}`,
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Missing Required fields",
		},
		{
			name: "Invalid token",
			body: `{"title": "Updated Post", "content": "Updated content"}`,
			mockSetup: func() {
				utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
					return nil, errors.New("invalid token")
				}
			},
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "Invalid token",
		},
		{
			name: "Error updating post",
			body: `{"title": "Updated Post", "content": "Updated content"}`,
			mockSetup: func() {
				mockService.EXPECT().UpdateMyPost(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("database error"))
				utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
					return jwt.MapClaims{"id": "userId"}, nil
				}
			},
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "Error updating post",
		},
		{
			name: "Not your post",
			body: `{"title": "Updated Post", "content": "Updated content"}`,
			mockSetup: func() {
				mockService.EXPECT().UpdateMyPost(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(utils.NotYourPost)
				utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
					return jwt.MapClaims{"id": "userId"}, nil
				}
			},
			expectedCode: http.StatusNotFound,
			expectedMsg:  "no post of yours exist with this id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("PUT", "/posts/post_id", bytes.NewBuffer([]byte(tt.body)))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", "Bearer valid-token")

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/posts/{post_id}", handler.UpdatePost).Methods("PUT")

			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}

			var actualRes models.Response
			if err := json.Unmarshal(rr.Body.Bytes(), &actualRes); err != nil {
				t.Fatalf("failed to unmarshal actual response: %v", err)
			}

			if actualRes.Code != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", actualRes.Code, tt.expectedCode)
			}
			if actualRes.Message != tt.expectedMsg {
				t.Errorf("handler returned wrong message: got %v want %v", actualRes.Message, tt.expectedMsg)
			}
		})
	}
}

func TestPostHandler_DeletePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockPostServiceInterface(ctrl)
	handler := handlers.NewPostHandler(mockService)

	postId := "1"
	req, err := http.NewRequest("DELETE", "/posts/"+postId, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer valid-token")
	// Mock the ExtractClaims function
	utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
		return jwt.MapClaims{"id": "userId"}, nil
	}
	defer func() { utils.ExtractClaimsFunc = utils.ExtractClaims }() // Reset after test

	mockService.EXPECT().DeleteMyPost("userId", postId).Return(nil)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/posts/{post_id}", handler.DeletePost).Methods("DELETE")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expectedResponse := models.Response{
		Code:    http.StatusOK,
		Message: "Post deleted successfully",
		Data:    nil,
	}
	var actualRes models.Response
	if err := json.Unmarshal(rr.Body.Bytes(), &actualRes); err != nil {
		t.Fatalf("failed to unmarshal actual response: %v", err)
	}
	if actualRes.Code != expectedResponse.Code {
		t.Errorf("handler returned wrong status code: got %v want %v", actualRes.Code, expectedResponse.Code)
	}
	if actualRes.Message != expectedResponse.Message {
		t.Errorf("handler returned wrong message: got %v want %v", actualRes.Message, expectedResponse.Message)
	}
}

func TestPostHandler_DisplayPosts_ErrorCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockPostServiceInterface(ctrl)
	handler := handlers.NewPostHandler(mockService)

	tests := []struct {
		name         string
		filter       string
		limit        string
		mockSetup    func()
		expectedCode int
		expectedMsg  string
	}{
		{
			name:   "Error fetching all posts",
			filter: "",
			mockSetup: func() {
				mockService.EXPECT().GiveAllPosts().Return(nil, errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "Error displaying posts",
		},
		{
			name:   "Error fetching filtered posts",
			filter: "food",
			mockSetup: func() {
				mockService.EXPECT().GiveFilteredPosts("food").Return(nil, errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "Error displaying posts in filterdatabase error",
		},
		{
			name:         "Invalid filter",
			filter:       "invalid",
			expectedCode: http.StatusBadRequest,
			expectedMsg:  "Invalid filter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/posts?filter="+tt.filter+"&limit="+tt.limit, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/posts", handler.DisplayPosts).Methods("GET")

			if tt.mockSetup != nil {
				tt.mockSetup()
			}

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}

			var actualRes models.Response
			if err := json.Unmarshal(rr.Body.Bytes(), &actualRes); err != nil {
				t.Fatalf("failed to unmarshal actual response: %v", err)
			}

			if actualRes.Code != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", actualRes.Code, tt.expectedCode)
			}
			if actualRes.Message != tt.expectedMsg {
				t.Errorf("handler returned wrong message: got %v want %v", actualRes.Message, tt.expectedMsg)
			}
		})
	}
}

func TestPostHandler_LikePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockPostServiceInterface(ctrl)
	handler := handlers.NewPostHandler(mockService)

	postId := "1"
	req, err := http.NewRequest("POST", "/posts/"+postId+"/like", nil)
	if err != nil {
		t.Fatal(err)
	}

	mockService.EXPECT().Like(postId).Return(nil)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/posts/{post_id}/like", handler.LikePost).Methods("POST")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expectedResponse := models.Response{
		Code:    http.StatusOK,
		Message: "Post liked successfully",
		Data:    nil,
	}
	var actualRes models.Response
	if err := json.Unmarshal(rr.Body.Bytes(), &actualRes); err != nil {
		t.Fatalf("failed to unmarshal actual response: %v", err)
	}
	if actualRes.Code != expectedResponse.Code {
		t.Errorf("handler returned wrong status code: got %v want %v", actualRes.Code, expectedResponse.Code)
	}
	if actualRes.Message != expectedResponse.Message {
		t.Errorf("handler returned wrong message: got %v want %v", actualRes.Message, expectedResponse.Message)
	}
}

func TestPostHandler_DisplayPostById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockPostServiceInterface(ctrl)
	handler := handlers.NewPostHandler(mockService)

	postId := "1"
	posts := &models.PostWithQuestions{
		PostId: "1", UId: "userId", Title: "Post 1", Content: "Content 1", Type: "travel", Likes: 0,
	}
	mockService.EXPECT().GivePostById(postId).Return(posts, nil)

	req, err := http.NewRequest("GET", "/posts/"+postId, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/posts/{post_id}", handler.DisplayPostById).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expectedResponse := models.Response{
		Code:    http.StatusOK,
		Message: "Successfully displayed post",
		Data:    posts,
	}
	var actualRes models.Response
	if err := json.Unmarshal(rr.Body.Bytes(), &actualRes); err != nil {
		t.Fatalf("failed to unmarshal actual response: %v", err)
	}
	if actualRes.Code != expectedResponse.Code {
		t.Errorf("handler returned wrong status code: got %v want %v", actualRes.Code, expectedResponse.Code)
	}
	if actualRes.Message != expectedResponse.Message {
		t.Errorf("handler returned wrong message: got %v want %v", actualRes.Message, expectedResponse.Message)
	}

}

func TestPostHandler_DisplayUserPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockPostServiceInterface(ctrl)
	handler := handlers.NewPostHandler(mockService)

	// Mock the ExtractClaims function
	utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
		return jwt.MapClaims{"id": "userId"}, nil
	}
	defer func() { utils.ExtractClaimsFunc = utils.ExtractClaims }() // Reset after test

	req, err := http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer valid-token")

	// Mocked posts
	mockPosts := []*models.Post{
		{PostId: "1", UId: "userId", Title: "Post 1", Type: "travel", Content: "Content 1", Likes: 0, CreatedAt: time.Now()},
		{PostId: "2", UId: "userId", Title: "Post 2", Type: "food", Content: "Content 2", Likes: 5, CreatedAt: time.Now()},
	}

	mockService.EXPECT().GiveMyPosts("userId").Return(mockPosts, nil)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/posts", handler.DisplayUserPosts).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedResponse := models.Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data: []models.ResponsePost{
			{PostId: "1", UId: "userId", Title: "Post 1", Type: "travel", Content: "Content 1", Likes: 0, CreatedAt: "2006-01-02 15:04:05"},
			{PostId: "2", UId: "userId", Title: "Post 2", Type: "food", Content: "Content 2", Likes: 5, CreatedAt: "2006-01-02 15:04:05"},
		},
	}

	var actualRes models.Response
	if err := json.Unmarshal(rr.Body.Bytes(), &actualRes); err != nil {
		t.Fatalf("failed to unmarshal actual response: %v", err)
	}

	if actualRes.Code != expectedResponse.Code {
		t.Errorf("handler returned wrong status code: got %v want %v", actualRes.Code, expectedResponse.Code)
	}
	if actualRes.Message != expectedResponse.Message {
		t.Errorf("handler returned wrong message: got %v want %v", actualRes.Message, expectedResponse.Message)
	}

}

func TestPostHandler_DisplayUserPosts_ErrorCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockPostServiceInterface(ctrl)
	handler := handlers.NewPostHandler(mockService)

	tests := []struct {
		name         string
		mockSetup    func()
		expectedCode int
		expectedMsg  string
	}{
		{
			name: "Unauthorized token",
			mockSetup: func() {
				utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
					return nil, errors.New("invalid token")
				}
			},
			expectedCode: http.StatusUnauthorized,
			expectedMsg:  "Invalid token",
		},
		{
			name: "Error fetching user posts",
			mockSetup: func() {
				utils.ExtractClaimsFunc = func(token string) (jwt.MapClaims, error) {
					return jwt.MapClaims{"id": "userId"}, nil
				}
				mockService.EXPECT().GiveMyPosts("userId").Return(nil, errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedMsg:  "Error displaying posts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/posts", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", "Bearer valid-token")

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/posts", handler.DisplayUserPosts).Methods("GET")

			tt.mockSetup()
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedCode)
			}

			var actualRes models.Response
			if err := json.Unmarshal(rr.Body.Bytes(), &actualRes); err != nil {
				t.Fatalf("failed to unmarshal actual response: %v", err)
			}

			if actualRes.Code != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v", actualRes.Code, tt.expectedCode)
			}
			if actualRes.Message != tt.expectedMsg {
				t.Errorf("handler returned wrong message: got %v want %v", actualRes.Message, tt.expectedMsg)
			}
		})
	}
}
