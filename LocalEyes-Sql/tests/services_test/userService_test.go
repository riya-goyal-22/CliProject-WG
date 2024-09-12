package services_test

import (
	"database/sql"
	"errors"
	"localEyes/config"
	"localEyes/internal/models"
	"localEyes/internal/services"
	"localEyes/tests/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := services.NewUserService(mockRepo)

	tests := []struct {
		name          string
		username      string
		password      string
		dwellingAge   int
		tag           string
		mockError     error
		expectedError string
	}{
		{"Signup Success", "testuser", "password", 5, "tourist", nil, ""},
		{"Signup Error", "testuser", "password", 5, "tourist", errors.New("creation error"), "creation error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().Create(gomock.Any()).Return(tt.mockError)

			err := userService.Signup(tt.username, tt.password, tt.dwellingAge, tt.tag)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := services.NewUserService(mockRepo)

	tests := []struct {
		name          string
		username      string
		password      string
		mockUser      *models.User
		mockError     error
		expectedError string
	}{
		{
			"Login Success",
			"testuser", "password",
			&models.User{Username: "testuser", Password: services.HashPassword("password"), IsActive: true},
			nil,
			"",
		},
		{
			"Login Invalid Credentials",
			"testuser", "password",
			nil,
			errors.New("invalid credentials"),
			config.Red + "Invalid Account credentials" + config.Reset,
		},
		{
			"Login Inactive Account",
			"testuser", "password",
			&models.User{Username: "testuser", Password: services.HashPassword("password"), IsActive: false},
			nil,
			config.Red + "InActive Account" + config.Reset,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().FindByUsernamePassword(tt.username, services.HashPassword(tt.password)).Return(tt.mockUser, tt.mockError)

			user, err := userService.Login(tt.username, tt.password)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.username, user.Username)
			}
		})
	}
}

func TestUserService_DeActivate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := services.NewUserService(mockRepo)

	tests := []struct {
		name          string
		UId           int
		mockError     error
		expectedError string
	}{
		{"DeActivate Success", 1, nil, ""},
		{"DeActivate Error", 1, errors.New("update error"), "update error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().UpdateActiveStatus(tt.UId, false).Return(tt.mockError)

			err := userService.DeActivate(tt.UId)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_NotifyUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := services.NewUserService(mockRepo)

	tests := []struct {
		name          string
		UId           int
		title         string
		mockError     error
		expectedError string
	}{
		{"NotifyUsers Success", 1, "Test Notification", nil, ""},
		{"NotifyUsers Error", 1, "Test Notification", errors.New("notification error"), "notification error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().PushNotification(tt.UId, tt.title).Return(tt.mockError)

			err := userService.NotifyUsers(tt.UId, tt.title)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_UnNotifyUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	userService := services.NewUserService(mockRepo)

	tests := []struct {
		name          string
		UId           int
		mockError     error
		expectedError string
	}{
		{"UnNotifyUsers Success", 1, nil, ""},
		{"UnNotifyUsers NoRows Error", 1, sql.ErrNoRows, ""},
		{"UnNotifyUsers Error", 1, errors.New("clear notification error"), "clear notification error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().ClearNotification(tt.UId).Return(tt.mockError)

			err := userService.UnNotifyUsers(tt.UId)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
