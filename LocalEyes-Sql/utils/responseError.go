package utils

import (
	"localEyes/internal/models"
	"net/http"
)

func NewNotFoundError(message string) *models.Response {
	return &models.Response{
		Message: message,
		Code:    http.StatusNotFound,
	}
}

func NewInternalServerError(message string) *models.Response {
	return &models.Response{
		Message: message,
		Code:    http.StatusInternalServerError,
	}
}

func NewBadRequestError(message string) *models.Response {
	return &models.Response{
		Message: message,
		Code:    http.StatusBadRequest,
	}
}

func NewUnauthorizedError(message string) *models.Response {
	return &models.Response{
		Message: message,
		Code:    http.StatusUnauthorized,
	}
}
