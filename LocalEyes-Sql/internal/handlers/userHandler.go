package handlers

import (
	_ "database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"localEyes/internal/interfaces"
	"localEyes/internal/models"
	"localEyes/utils"
	"log"
	"net/http"
)

type UserHandler struct {
	service interfaces.UserServiceInterface
}

func NewUserHandler(service interfaces.UserServiceInterface) *UserHandler {
	return &UserHandler{service}
}

func (handler *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var client models.Client
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.NewBadRequestError("Invalid JSON body")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}
	if !handler.service.ValidateUsername(client.Username) {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.NewBadRequestError("Username not available")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}
	if !handler.service.ValidateEmail(client.Email) {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.NewBadRequestError("Email already registered")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}
	if !utils.ValidatePassword(client.Password) {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.NewBadRequestError("Password not strong")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}
	var livingSinceInYears = client.LivingSince.Days/365.0 + client.LivingSince.Months/12.0 + client.LivingSince.Years

	err = handler.service.Signup(client.Username, client.Password, client.Email, livingSinceInYears)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Error signing up")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
	utils.Logger.Info("User signed up successfully")
	response := &models.Response{
		Message: "User created successfully",
		Code:    http.StatusOK,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("ERROR: Error encoding response")
	}
}

func (handler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var client models.Client
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.NewBadRequestError("Invalid JSON body")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}

	user, err := handler.service.Login(client.Username, client.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := utils.NewUnauthorizedError(err.Error())
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}
	generatedToken, err := utils.GenerateTokenFunc(user.Username, user.UId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Error generating token ")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Logger.Info("User logged in successfully")
	response := models.Response{
		Data:    generatedToken,
		Code:    http.StatusOK,
		Message: "User logged in successfully",
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("ERROR: Error encoding response")
	}
}

func (handler *UserHandler) DeActivate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bearerToken := r.Header.Get("Authorization")
	claims, err := utils.ExtractClaimsFunc(bearerToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Invalid token")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}
	id := claims["id"].(string)
	//intId := int(id)
	err = handler.service.DeActivate(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Error deactivating user")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.Logger.Info("User deactivated successfully")
	response := &models.Response{
		Message: "User Deactivated successfully",
		Code:    http.StatusOK,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("ERROR: Error encoding response")
	}
	return
}

func (handler *UserHandler) ViewProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bearerToken := r.Header.Get("Authorization")
	claims, err := utils.ExtractClaimsFunc(bearerToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewUnauthorizedError("Error extracting claims")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}
	id := claims["id"].(string)
	user, err := handler.service.GetUserById(id)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := utils.NewUnauthorizedError("Invalid token")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}
	responseUser := &models.ResponseUser{
		UId:          user.UId,
		Email:        user.Email,
		Username:     user.Username,
		City:         user.City,
		LivingSince:  user.DwellingAge,
		Tag:          user.Tag,
		ActiveStatus: user.IsActive,
	}
	response := models.Response{
		Data:    responseUser,
		Code:    http.StatusOK,
		Message: "User viewed successfully",
	}
	w.WriteHeader(http.StatusOK)
	utils.Logger.Info("User viewed successfully")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("ERROR: Error encoding response")
	}
	return
}

func (handler *UserHandler) ViewNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bearerToken := r.Header.Get("Authorization")
	claims, err := utils.ExtractClaimsFunc(bearerToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := utils.NewUnauthorizedError("Invalid token")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}
	id := claims["id"].(string)

	notifications, err := handler.service.GetNotifications(id)
	if err != nil {
		if errors.Is(err, utils.NoUser) {
			w.WriteHeader(http.StatusNotFound)
			response := utils.NewNotFoundError("No user found with that id")
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				utils.Logger.Error("ERROR: Error encoding response")
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Error getting notifications")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}
	response := &models.Response{
		Message: "Success",
		Code:    http.StatusOK,
		Data:    notifications,
	}
	w.WriteHeader(http.StatusOK)
	utils.Logger.Info("User viewed successfully")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("ERROR: Error encoding response")
	}

}

func (handler *UserHandler) SendOtp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userEmail models.UserEmail
	err := json.NewDecoder(r.Body).Decode(&userEmail)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		response := utils.NewBadRequestError("Invalid JSON body")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}
	err = handler.service.SendOtp(userEmail.Email)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Internal server error")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	response := &models.Response{
		Message: "Success",
		Code:    http.StatusOK,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("ERROR: Error encoding response")
	}
}

func (handler *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var resetUser models.ResetPasswordUser
	err := json.NewDecoder(r.Body).Decode(&resetUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.NewBadRequestError("Invalid JSON body")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}
	err = handler.service.PasswordReset(resetUser)
	if err != nil {
		fmt.Println("reset ", err)
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Internal server error")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	response := &models.Response{
		Message: "Success",
		Code:    http.StatusOK,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("ERROR: Error encoding response")
	}
}

func (handler *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := mux.Vars(r)["user_id"]
	user, err := handler.service.GetUserById(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Error getting user")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
	}
	responseUser := models.ResponseUser{
		UId:         user.UId,
		Email:       user.Email,
		Username:    user.Username,
		City:        user.City,
		LivingSince: user.DwellingAge,
		Tag:         user.Tag,
	}
	response := models.Response{
		Data:    responseUser,
		Code:    http.StatusOK,
		Message: "Success",
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("ERROR: Error encoding response")
	}
}

func (handler *UserHandler) UpdateUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	bearerToken := r.Header.Get("Authorization")
	claims, err := utils.ExtractClaimsFunc(bearerToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		response := utils.NewUnauthorizedError("Invalid token")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}
	userId := claims["id"].(string)
	var newUser models.Client
	err = json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := utils.NewBadRequestError("Invalid JSON body")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
		return
	}
	err = handler.service.UpdateUser(userId, &newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := utils.NewInternalServerError("Error updating user")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			utils.Logger.Error("ERROR: Error encoding response")
		}
	}
	w.WriteHeader(http.StatusOK)
	response := &models.Response{
		Message: "Success",
		Code:    http.StatusOK,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		utils.Logger.Error("ERROR: Error encoding response")
	}
}
