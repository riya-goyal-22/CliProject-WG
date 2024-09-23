package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"localEyes/internal/handlers"
	"localEyes/internal/models"
	"localEyes/tests/mocks"
	"localEyes/utils"
)

func TestAdminHandler_DisplayUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockAdminServiceInterface(ctrl)
	handler := handlers.NewAdminHandler(mockService)

	// Successful response
	mockService.EXPECT().GetAllUsers().Return([]*models.User{
		{UId: "1", Username: "user1", City: "CityA", DwellingAge: 5, Tag: "tag1"},
	}, nil)

	req, _ := http.NewRequest("GET", "/admin/users", nil)
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/admin/users", handler.DisplayUsers)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.Response
	json.NewDecoder(rr.Body).Decode(&response)
	if response.Code != http.StatusOK || response.Message != "Success" {
		t.Errorf("unexpected response: got %v", response)
	}

	// Error case
	mockService.EXPECT().GetAllUsers().Return(nil, utils.NoUser)

	req, _ = http.NewRequest("GET", "/admin/users", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestAdminHandler_DisplayPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockAdminServiceInterface(ctrl)
	handler := handlers.NewAdminHandler(mockService)

	// Successful response
	mockService.EXPECT().GetAllPosts().Return([]*models.Post{
		{PostId: "1", UId: "1", Title: "Post 1", Type: "TypeA", Content: "Content 1", Likes: 10, CreatedAt: time.Now()},
	}, nil)

	req, _ := http.NewRequest("GET", "/admin/posts", nil)
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/admin/posts", handler.DisplayPosts)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.Response
	json.NewDecoder(rr.Body).Decode(&response)
	if response.Code != http.StatusOK || response.Message != "Success" {
		t.Errorf("unexpected response: got %v", response)
	}

	// Error case
	mockService.EXPECT().GetAllPosts().Return(nil, utils.NoPost)

	req, _ = http.NewRequest("GET", "/admin/posts", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestAdminHandler_DisplayQuestions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockAdminServiceInterface(ctrl)
	handler := handlers.NewAdminHandler(mockService)

	replies := []string{"reply"}
	// Successful response
	mockService.EXPECT().GetAllQuestions().Return([]*models.Question{
		{QId: "1", PostId: "1", UserId: "1", Text: "Question 1", Replies: replies, CreatedAt: time.Now()},
	}, nil)

	req, _ := http.NewRequest("GET", "/admin/questions", nil)
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/admin/questions", handler.DisplayQuestions)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.Response
	json.NewDecoder(rr.Body).Decode(&response)
	if response.Code != http.StatusOK || response.Message != "Success" {
		t.Errorf("unexpected response: got %v", response)
	}

	// Error case
	mockService.EXPECT().GetAllQuestions().Return(nil, utils.NoQuestion)

	req, _ = http.NewRequest("GET", "/admin/questions", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
}

func TestAdminHandler_DeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockAdminServiceInterface(ctrl)
	handler := handlers.NewAdminHandler(mockService)

	userId := "1"

	// Successful deletion
	mockService.EXPECT().DeleteUser(userId).Return(nil)

	req, _ := http.NewRequest("DELETE", "/admin/users/"+userId, nil)
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/admin/users/{user_id}", handler.DeleteUser)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.Response
	json.NewDecoder(rr.Body).Decode(&response)
	if response.Code != http.StatusOK || response.Message != "Successfully deleted user" {
		t.Errorf("unexpected response: got %v", response)
	}

	// User not found case
	mockService.EXPECT().DeleteUser(userId).Return(utils.NoUser)

	req, _ = http.NewRequest("DELETE", "/admin/users/"+userId, nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestAdminHandler_DeletePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockAdminServiceInterface(ctrl)
	handler := handlers.NewAdminHandler(mockService)

	postId := "1"

	// Successful deletion
	mockService.EXPECT().DeletePost(postId).Return(nil)

	req, _ := http.NewRequest("DELETE", "/admin/posts/"+postId, nil)
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/admin/posts/{post_id}", handler.DeletePost)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.Response
	json.NewDecoder(rr.Body).Decode(&response)
	if response.Code != http.StatusOK || response.Message != "Successfully deleted post" {
		t.Errorf("unexpected response: got %v", response)
	}

	// Post not found case
	mockService.EXPECT().DeletePost(postId).Return(utils.NoPost)

	req, _ = http.NewRequest("DELETE", "/admin/posts/"+postId, nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestAdminHandler_DeleteQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockAdminServiceInterface(ctrl)
	handler := handlers.NewAdminHandler(mockService)

	questionId := "1"

	// Successful deletion
	mockService.EXPECT().DeleteQuestion(questionId).Return(nil)

	req, _ := http.NewRequest("DELETE", "/admin/questions/"+questionId, nil)
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/admin/questions/{ques_id}", handler.DeleteQuestion)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.Response
	json.NewDecoder(rr.Body).Decode(&response)
	if response.Code != http.StatusOK || response.Message != "Successfully deleted question" {
		t.Errorf("unexpected response: got %v", response)
	}

	// Question not found case
	mockService.EXPECT().DeleteQuestion(questionId).Return(utils.NoQuestion)

	req, _ = http.NewRequest("DELETE", "/admin/questions/"+questionId, nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestAdminHandler_ReActivateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockAdminServiceInterface(ctrl)
	handler := handlers.NewAdminHandler(mockService)

	userId := "1"

	// Successful reactivation
	mockService.EXPECT().ReActivate(userId).Return(nil)

	req, _ := http.NewRequest("PATCH", "/admin/users/"+userId+"/reactivate", nil)
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/admin/users/{user_id}/reactivate", handler.ReactivateUser)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response models.Response
	json.NewDecoder(rr.Body).Decode(&response)
	if response.Code != http.StatusOK || response.Message != "Successfully activated user" || response.Data != nil {
		t.Errorf("unexpected response: got %v", response)
	}

	// User not found case
	mockService.EXPECT().ReActivate(userId).Return(utils.NoUser)

	req, _ = http.NewRequest("PATCH", "/admin/users/"+userId+"/reactivate", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
