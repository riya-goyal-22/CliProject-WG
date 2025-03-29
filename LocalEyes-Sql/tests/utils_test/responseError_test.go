package utils_test

import (
	"github.com/stretchr/testify/assert"
	"localEyes/utils"
	"net/http"
	"testing"
)

func TestNewNotFoundError(t *testing.T) {
	message := "Resource not found"
	response := utils.NewNotFoundError(message)

	assert.Equal(t, message, response.Message)
	assert.Equal(t, http.StatusNotFound, response.Code)
}

func TestNewInternalServerError(t *testing.T) {
	message := "Internal server error"
	response := utils.NewInternalServerError(message)

	assert.Equal(t, message, response.Message)
	assert.Equal(t, http.StatusInternalServerError, response.Code)
}

func TestNewBadRequestError(t *testing.T) {
	message := "Bad request"
	response := utils.NewBadRequestError(message)

	assert.Equal(t, message, response.Message)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

func TestNewUnauthorizedError(t *testing.T) {
	message := "Unauthorized access"
	response := utils.NewUnauthorizedError(message)

	assert.Equal(t, message, response.Message)
	assert.Equal(t, http.StatusUnauthorized, response.Code)
}
