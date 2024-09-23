package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"localEyes/internal/interfaces"
	"localEyes/internal/models"
	"localEyes/utils"
	"net/http"
	"strconv"
)

type QuestionHandler struct {
	service interfaces.QuestionServiceInterface
}

func NewQuestionHandler(service interfaces.QuestionServiceInterface) *QuestionHandler {
	return &QuestionHandler{service: service}
}

func (handler *QuestionHandler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var question models.RequestQuestion
	err := json.NewDecoder(r.Body).Decode(&question)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.NewBadRequestError("Invalid JSON payload")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: error encoding response")
		}
		return
	}
	postId := mux.Vars(r)["post_id"]
	bearerToken := r.Header.Get("Authorization")
	claims, err := utils.ExtractClaims(bearerToken)
	id := claims["id"].(string)
	//intId := int(id)
	err = handler.service.AskQuestion(id, postId, question.Question)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("error while asking question")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: error encoding response")
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	utils.Logger.Info("Question created")
	response := &models.Response{
		Code:    http.StatusOK,
		Message: "Question Created",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("ERROR: error encoding response")
	}
}

func (handler *QuestionHandler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["post_id"]
	questions, err := handler.service.GetPostQuestions(postId)
	queryParams := r.URL.Query()
	limitString := queryParams.Get("limit")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("error while getting questions")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: error encoding response")
		}
		return
	}
	utils.Logger.Info("Questions retrieved")
	var responseData []models.ResponseQuestion
	for _, question := range questions {
		responseData = append(responseData, models.ResponseQuestion{
			QId:       question.QId,
			PostId:    question.PostId,
			UserId:    question.UserId,
			Text:      question.Text,
			Replies:   question.Replies,
			CreatedAt: question.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	if limitString != "" {
		limit, _ := strconv.Atoi(limitString)
		if limit < len(responseData) {
			responseData = responseData[:limit]
		}
	}
	response := models.Response{
		Data:    responseData,
		Code:    http.StatusOK,
		Message: "Success",
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("ERROR: error encoding response")
	}
}

func (handler *QuestionHandler) AddAnswer(w http.ResponseWriter, r *http.Request) {
	quesId := mux.Vars(r)["ques_id"]
	var request models.RequestAnswer
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.NewBadRequestError("Invalid JSON payload")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: error encoding response")
		}
		return
	}
	err = handler.service.AddAnswer(quesId, request.Answer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("error while adding answer " + err.Error())
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: error encoding response")
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Logger.Info("Answer added")
	response := &models.Response{
		Code:    http.StatusOK,
		Message: "Answer Added",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("ERROR: error encoding response")
	}
}

func (handler *QuestionHandler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	quesId := mux.Vars(r)["ques_id"]
	bearerToken := r.Header.Get("Authorization")
	claims, err := utils.ExtractClaims(bearerToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := utils.NewUnauthorizedError("Invalid Token")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: error encoding response")
		}
		return
	}
	id := claims["id"].(string)
	//intId := int(id)
	err = handler.service.DeleteUserQues(id, quesId)
	if err != nil {
		if errors.Is(err, utils.NotYourQuestion) {
			w.WriteHeader(http.StatusNotFound)
			response := utils.NewNotFoundError("Question Not Found")
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Logger.Error("ERROR: error encoding response")
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("error while deleting question")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: error encoding response")
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Logger.Info("Question deleted")
	response := &models.Response{
		Code:    http.StatusOK,
		Message: "Question Deleted",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("ERROR: error encoding response")
	}
}
