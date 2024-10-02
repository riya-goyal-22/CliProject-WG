package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"localEyes/internal/interfaces"
	"localEyes/internal/models"
	"localEyes/utils"
	"net/http"
	"strconv"
)

type PostHandler struct {
	service interfaces.PostServiceInterface
}

func NewPostHandler(service interfaces.PostServiceInterface) *PostHandler {
	return &PostHandler{
		service: service,
	}
}

func (handler *PostHandler) DisplayPosts(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	filter := queryParams.Get("filter")
	limitString := queryParams.Get("limit")
	if filter == "" {
		posts, err := handler.service.GiveAllPosts()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := utils.NewInternalServerError("Error displaying posts")
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Logger.Error("Error encoding response:" + err.Error())
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
		utils.Logger.Info("Successfully displayed posts")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response:" + err.Error())
		}
		return
	} else if filter == "food" || filter == "shopping" || filter == "other" || filter == "travel" {
		posts, err := handler.service.GiveFilteredPosts(filter)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := utils.NewInternalServerError("Error displaying posts in filter" + err.Error())
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Logger.Error("Error encoding response:" + err.Error())
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
		utils.Logger.Info("Successfully displayed filtered posts")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response:" + err.Error())
		}
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.NewBadRequestError("Invalid filter")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response:" + err.Error())
		}
	}

}

func (handler *PostHandler) DisplayUserPosts(w http.ResponseWriter, r *http.Request) {
	bearerToken := r.Header.Get("Authorization")
	claims, err := utils.ExtractClaimsFunc(bearerToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := utils.NewUnauthorizedError("Invalid token")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response:" + err.Error())
		}
		return
	}
	id := claims["id"].(string)
	//intId := (id)
	posts, err := handler.service.GiveMyPosts(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Error displaying posts")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response:" + err.Error())
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
	utils.Logger.Info("Successfully displayed user posts")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("Error encoding response:" + err.Error())
	}
	return
}

func (handler *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	bearerToken := r.Header.Get("Authorization")
	var requestPost models.RequestPost
	err := json.NewDecoder(r.Body).Decode(&requestPost)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.NewBadRequestError("Missing Request body")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response:" + err.Error())
		}
		return
	}
	isValid, err := utils.ValidatePostRequest(requestPost)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.NewBadRequestError(err.Error())
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response:" + err.Error())
		}
		return
	}
	if isValid {
		claims, err := utils.ExtractClaimsFunc(bearerToken)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			response := utils.NewUnauthorizedError("Invalid token")
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Logger.Error("Error encoding response:" + err.Error())
			}
			return
		}
		id := claims["id"].(string)
		//intId := int(id)
		err = handler.service.CreatePost(id, requestPost.Title, requestPost.Content, requestPost.Type)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := utils.NewInternalServerError("Error creating post")
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Logger.Error("Error encoding response:" + err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusCreated)
		utils.Logger.Info("Successfully created post")
		response := &models.Response{
			Message: "Post created successfully",
			Code:    http.StatusOK,
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response:" + err.Error())
		}
	}
}
func (handler *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	bearerToken := r.Header.Get("Authorization")
	postId := mux.Vars(r)["post_id"]

	post := struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.NewBadRequestError("Invalid/Missing Request body")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response:" + err.Error())
		}
		return
	}
	if post.Title == "" || post.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.NewBadRequestError("Missing Required fields")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response:" + err.Error())
		}
		return
	}
	claims, err := utils.ExtractClaimsFunc(bearerToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := utils.NewUnauthorizedError("Invalid token")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response:" + err.Error())
		}
		return
	}
	id := claims["id"].(string)
	//intId := int(id)
	err = handler.service.UpdateMyPost(postId, id, post.Title, post.Content)
	if err != nil {
		if errors.Is(err, utils.NotYourPost) {
			w.WriteHeader(http.StatusNotFound)
			response := utils.NewNotFoundError(err.Error())
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Logger.Error("Error encoding response:" + err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Error updating post")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response:" + err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Logger.Info("Successfully updated post")
	response := &models.Response{
		Message: "Post updated successfully",
		Code:    http.StatusOK,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("Error encoding response:" + err.Error())
	}
}

func (handler *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	bearerToken := r.Header.Get("Authorization")
	postId := mux.Vars(r)["post_id"]
	claims, err := utils.ExtractClaimsFunc(bearerToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := utils.NewUnauthorizedError("Invalid token")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response:" + err.Error())
		}
		return
	}
	id := claims["id"].(string)
	//intId := int(id)
	err = handler.service.DeleteMyPost(id, postId)
	if err != nil {
		if errors.Is(err, utils.NotYourPost) {
			w.WriteHeader(http.StatusNotFound)
			response := utils.NewNotFoundError(err.Error())
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Logger.Error("Error encoding response:" + err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Error deleting post")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response:" + err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Logger.Info("Successfully deleted post")
	response := &models.Response{
		Message: "Post deleted successfully",
		Code:    http.StatusOK,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("Error encoding response:" + err.Error())
	}
}

func (handler *PostHandler) LikePost(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["post_id"]
	err := handler.service.Like(postId)
	if err != nil {
		if errors.Is(err, utils.NoPost) {
			w.WriteHeader(http.StatusNotFound)
			response := utils.NewNotFoundError(err.Error())
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Logger.Error("Error encoding response:" + err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Error liking post")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response:" + err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Logger.Info("Successfully liked post")
	response := &models.Response{
		Message: "Post liked successfully",
		Code:    http.StatusOK,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("Error encoding response:" + err.Error())
	}
}

func (handler *PostHandler) DisplayPostById(w http.ResponseWriter, r *http.Request) {
	postId := mux.Vars(r)["post_id"]
	post, err := handler.service.GivePostById(postId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			response := utils.NewNotFoundError("No Post found with that id")
			err := json.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Logger.Error("Error encoding response:" + err.Error())
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Error getting post " + err.Error())
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("Error encoding response:" + err.Error())
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Logger.Info("Successfully displayed post")
	response := &models.Response{
		Message: "Successfully displayed post",
		Code:    http.StatusOK,
		Data:    post,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("Error encoding response:" + err.Error())
	}
}
