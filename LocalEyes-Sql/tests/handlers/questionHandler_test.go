package handlers

import (
	"bytes"
	"encoding/json"
	"localEyes/internal/handlers"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"localEyes/internal/models"
	"localEyes/tests/mocks"
)

func TestQuestionHandler_CreateQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockQuestionServiceInterface(ctrl)
	handler := handlers.NewQuestionHandler(mockService)

	// Set environment variable for secret
	os.Setenv("Secret", "mysecret")

	question := models.RequestQuestion{Question: "What is your favorite color?"}
	body, _ := json.Marshal(question)
	req, err := http.NewRequest("POST", "/posts/1/questions", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"id": "userId"})
	tokenString, _ := token.SignedString([]byte(os.Getenv("Secret")))
	req.Header.Set("Authorization", "Bearer "+tokenString)

	mockService.EXPECT().AskQuestion("userId", "1", question.Question).Return(nil)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/posts/{post_id}/questions", handler.CreateQuestion).Methods("POST")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	expectedResponse := models.Response{
		Code:    http.StatusOK,
		Message: "Question Created",
		Data:    nil,
	}

	// Unmarshal the actual response to a struct for structured comparison
	var actualRes models.Response
	if err := json.Unmarshal(rr.Body.Bytes(), &actualRes); err != nil {
		t.Fatalf("failed to unmarshal actual response: %v", err)
	}

	// Check Code
	if actualRes.Code != expectedResponse.Code {
		t.Errorf("handler returned wrong code: got %v want %v", actualRes.Code, expectedResponse.Code)
	}

	// Check Message
	if actualRes.Message != expectedResponse.Message {
		t.Errorf("handler returned wrong message: got %v want %v", actualRes.Message, expectedResponse.Message)
	}

	// Optionally, check Data if needed
	if actualRes.Data != expectedResponse.Data {
		t.Errorf("handler returned unexpected data: got %v want %v", actualRes.Data, expectedResponse.Data)
	}
}

func TestQuestionHandler_GetQuestions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockQuestionServiceInterface(ctrl)
	handler := handlers.NewQuestionHandler(mockService)

	// Set environment variable for secret
	os.Setenv("Secret", "mysecret")

	req, err := http.NewRequest("GET", "/posts/1/questions", nil)
	if err != nil {
		t.Fatal(err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"id": "userId"})
	tokenString, _ := token.SignedString([]byte(os.Getenv("Secret")))
	req.Header.Set("Authorization", "Bearer "+tokenString)

	questions := []*models.Question{
		{QId: "1", PostId: "1", UserId: "userId", Text: "What is your favorite color?", Replies: []string{}, CreatedAt: time.Now()},
	}

	mockService.EXPECT().GetPostQuestions("1").Return(questions, nil)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/posts/{post_id}/questions", handler.GetQuestions).Methods("GET")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	//var response models.Response
	//json.Unmarshal(rr.Body.Bytes(), &response)
	//if len(response.Data) != len(questions) {
	//	t.Errorf("handler returned incorrect number of questions: got %v want %v", len(response.Data), len(questions))
	//}
}

func TestQuestionHandler_AddAnswer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockQuestionServiceInterface(ctrl)
	handler := handlers.NewQuestionHandler(mockService)

	// Set environment variable for secret
	os.Setenv("Secret", "mysecret")

	requestBody := `{"answer": "Blue"}`
	req, err := http.NewRequest("POST", "/questions/1/answers", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		t.Fatal(err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"id": "userId"})
	tokenString, _ := token.SignedString([]byte(os.Getenv("Secret")))
	req.Header.Set("Authorization", "Bearer "+tokenString)

	mockService.EXPECT().AddAnswer("1", "Blue").Return(nil)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/questions/{ques_id}/answers", handler.AddAnswer).Methods("POST")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expectedResponse := models.Response{
		Code:    http.StatusOK,
		Message: "Answer Added",
		Data:    nil,
	}
	var actualRes models.Response
	if err := json.Unmarshal(rr.Body.Bytes(), &actualRes); err != nil {
		t.Fatalf("failed to unmarshal actual response: %v", err)
	}
	if actualRes.Code != expectedResponse.Code {
		t.Errorf("handler returned wrong code: got %v want %v", actualRes.Code, expectedResponse.Code)
	}
	if actualRes.Message != expectedResponse.Message {
		t.Errorf("handler returned wrong message: got %v want %v", actualRes.Message, expectedResponse.Message)
	}
}

func TestQuestionHandler_DeleteQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock.NewMockQuestionServiceInterface(ctrl)
	handler := handlers.NewQuestionHandler(mockService)

	// Set environment variable for secret
	os.Setenv("Secret", "mysecret")

	req, err := http.NewRequest("DELETE", "/questions/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"id": "userId"})
	tokenString, _ := token.SignedString([]byte(os.Getenv("Secret")))
	req.Header.Set("Authorization", "Bearer "+tokenString)

	mockService.EXPECT().DeleteUserQues("userId", "1").Return(nil)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/questions/{ques_id}", handler.DeleteQuestion).Methods("DELETE")

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	expectedResponse := models.Response{
		Code:    http.StatusOK,
		Message: "Question Deleted",
		Data:    nil,
	}
	var actualRes models.Response
	if err := json.Unmarshal(rr.Body.Bytes(), &actualRes); err != nil {
		t.Fatalf("failed to unmarshal actual response: %v", err)
	}
	if actualRes.Code != expectedResponse.Code {
		t.Errorf("handler returned wrong code: got %v want %v", actualRes.Code, expectedResponse.Code)
	}
	if actualRes.Message != expectedResponse.Message {
		t.Errorf("handler returned wrong message: got %v want %v", actualRes.Message, expectedResponse.Message)
	}
}
