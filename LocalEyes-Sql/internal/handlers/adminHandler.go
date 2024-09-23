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

type AdminHandler struct {
	service interfaces.AdminServiceInterface
}

func NewAdminHandler(service interfaces.AdminServiceInterface) *AdminHandler {
	return &AdminHandler{service: service}
}

func (handler *AdminHandler) DisplayUsers(w http.ResponseWriter, r *http.Request) {
	users, err := handler.service.GetAllUsers()
	queryParams := r.URL.Query()
	limitString := queryParams.Get("limit")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Error getting all users")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error in DisplayUsers:" + err.Error())
		}
		return
	}
	var responseData []models.ResponseUser
	for _, user := range users {
		responseData = append(responseData, models.ResponseUser{
			UId:         user.UId,
			Username:    user.Username,
			City:        user.City,
			LivingSince: user.DwellingAge,
			Tag:         user.Tag,
		})
	}
	if limitString != "" {
		limit, _ := strconv.Atoi(limitString)
		if limit < len(responseData) {
			responseData = responseData[:limit]
		}
	}
	response := &models.Response{
		Data:    responseData,
		Code:    http.StatusOK,
		Message: "Success",
	}
	w.WriteHeader(http.StatusOK)
	utils.Logger.Info("Users displayed successfully")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("Error in DisplayUsers encoding response:" + err.Error())
	}
}

func (handler *AdminHandler) DisplayPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := handler.service.GetAllPosts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Error getting all posts")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response in DisplayPosts:" + err.Error())
		}
		return
	}
	var responseData []models.ResponsePost
	for _, post := range posts {
		responseData = append(responseData, models.ResponsePost{
			PostId:    post.PostId,
			UId:       post.UId,
			Title:     post.Title,
			Type:      post.Type,
			Content:   post.Content,
			Likes:     post.Likes,
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	response := models.Response{
		Data:    responseData,
		Code:    http.StatusOK,
		Message: "Success",
	}
	w.WriteHeader(http.StatusOK)
	utils.Logger.Info("Posts displayed successfully")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("Error encoding response in DisplayPosts:" + err.Error())
	}
}

func (handler *AdminHandler) DisplayQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := handler.service.GetAllQuestions()
	queryParams := r.URL.Query()
	limitString := queryParams.Get("limit")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Error getting all questions")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response in DisplayQuestions:" + err.Error())
		}
		return
	}
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
	response := &models.Response{
		Data:    responseData,
		Code:    http.StatusOK,
		Message: "Success",
	}
	w.WriteHeader(http.StatusOK)
	utils.Logger.Info("Questions displayed successfully")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("Error encoding response in DisplayQuestions:" + err.Error())
	}
}

func (handler *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["user_id"]
	err := handler.service.DeleteUser(userId)
	if err != nil {
		if errors.Is(err, utils.NoUser) {
			w.WriteHeader(http.StatusNotFound)
			response := utils.NewNotFoundError("User not found")
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Logger.Error("Error encoding response in DeleteUser:" + err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Error deleting user")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response in DeleteUser:" + err.Error())
		}
		return
	}
	response := &models.Response{
		Code:    http.StatusOK,
		Message: "Successfully deleted user",
	}
	w.WriteHeader(http.StatusOK)
	utils.Logger.Info("User deleted successfully")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("Error encoding response in DeleteUser:" + err.Error())
	}
}

func (handler *AdminHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["post_id"]
	err := handler.service.DeletePost(postId)
	if err != nil {
		if errors.Is(err, utils.NoPost) {
			w.WriteHeader(http.StatusNotFound)
			response := utils.NewNotFoundError("Post not found")
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Logger.Error("Error in encoding response DeletePost:" + err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Error deleting post")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response in DeletePost:" + err.Error())
		}
		return
	}
	response := &models.Response{
		Code:    http.StatusOK,
		Message: "Successfully deleted post",
	}
	w.WriteHeader(http.StatusOK)
	utils.Logger.Info("Post deleted successfully")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("Error encoding response in DeletePost:" + err.Error())
	}
}

func (handler *AdminHandler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	questionId := mux.Vars(r)["ques_id"]
	err := handler.service.DeleteQuestion(questionId)
	if err != nil {
		if errors.Is(err, utils.NoQuestion) {
			w.WriteHeader(http.StatusNotFound)
			response := utils.NewNotFoundError("Question not found")
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Logger.Error("Error encoding response in DeleteQuestion:" + err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Error deleting question")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response in DeleteQuestion:" + err.Error())
		}
		return
	}
	response := &models.Response{
		Code:    http.StatusOK,
		Message: "Successfully deleted question",
	}
	w.WriteHeader(http.StatusOK)
	utils.Logger.Info("Question deleted successfully")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("Error encoding response in DeleteQuestion:" + err.Error())
	}
}

func (handler *AdminHandler) ReactivateUser(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["user_id"]
	err := handler.service.ReActivate(userId)
	if err != nil {
		if errors.Is(err, utils.NoUser) {
			w.WriteHeader(http.StatusNotFound)
			response := utils.NewNotFoundError("No Inactive user found")
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Logger.Error("Error encoding response in ReactivateUser:" + err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Error reactivating user")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response in ReactivateUser:" + err.Error())
		}
		return
	}
	response := &models.Response{
		Code:    http.StatusOK,
		Message: "Successfully activated user",
		Data:    nil,
	}
	w.WriteHeader(http.StatusOK)
	utils.Logger.Info("User reactivated successfully")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("Error encoding response in ReactivateUser:" + err.Error())
	}
}
