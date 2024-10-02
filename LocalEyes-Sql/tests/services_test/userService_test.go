package services_test

import (
	"database/sql"
	"errors"
	"localEyes/internal/models"
	"localEyes/internal/services"
	"localEyes/tests/mocks"
	"localEyes/utils"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepository(ctrl)
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

	mockRepo := mock.NewMockUserRepository(ctrl)
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
			"invalid Account credentials",
		},
		{
			"Login Inactive Account",
			"testuser", "password",
			&models.User{Username: "testuser", Password: services.HashPassword("password"), IsActive: false},
			nil,
			"inActive Account",
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

	mockRepo := mock.NewMockUserRepository(ctrl)
	userService := services.NewUserService(mockRepo)

	tests := []struct {
		name          string
		UId           string
		mockError     error
		expectedError string
	}{
		{"DeActivate Success", "1", nil, ""},
		{"DeActivate Error", "1", errors.New("update error"), "update error"},
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

	mockRepo := mock.NewMockUserRepository(ctrl)
	userService := services.NewUserService(mockRepo)

	tests := []struct {
		name          string
		UId           string
		title         string
		mockError     error
		expectedError string
	}{
		{"NotifyUsers Success", "1", "Test Notification", nil, ""},
		{"NotifyUsers Error", "1", "Test Notification", errors.New("notification error"), "notification error"},
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

func TestGetNotifications_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockUserRepository(ctrl)

	userService := services.NewUserService(mockUserRepo)

	uid := "user-1"
	expectedNotifications := []string{"Notification 1", "Notification 2"}
	user := &models.User{
		UId:          uid,
		Username:     "testuser",
		Notification: expectedNotifications,
	}

	mockUserRepo.EXPECT().FindByUId(uid).Return(user, nil)
	mockUserRepo.EXPECT().ClearNotification(uid).Return(nil)

	notifications, err := userService.GetNotifications(uid)

	assert.NoError(t, err)
	assert.Equal(t, &expectedNotifications, notifications)
}

func TestGetNotifications_NoUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockUserRepository(ctrl)

	userService := services.NewUserService(mockUserRepo)

	uid := "non-existent-user"

	mockUserRepo.EXPECT().FindByUId(uid).Return(nil, sql.ErrNoRows)

	notifications, err := userService.GetNotifications(uid)

	assert.Error(t, err)
	assert.Equal(t, utils.NoUser, err)
	assert.Nil(t, notifications)
}

func TestGetNotifications_ErrorOnClear(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockUserRepository(ctrl)

	userService := services.NewUserService(mockUserRepo)

	uid := "user-1"
	user := &models.User{
		UId:          uid,
		Username:     "testuser",
		Notification: []string{"Notification 1"},
	}

	mockUserRepo.EXPECT().FindByUId(uid).Return(user, nil)
	mockUserRepo.EXPECT().ClearNotification(uid).Return(errors.New("clear error"))

	notifications, err := userService.GetNotifications(uid)

	assert.Error(t, err)
	assert.Equal(t, "clear error", err.Error())
	assert.Equal(t, (*[]string)(nil), notifications)
}

func TestValidateUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepository(ctrl)
	userService := services.UserService{Repo: mockRepo}

	tests := []struct {
		username   string
		mockReturn *models.User
		mockErr    error
		expected   bool
	}{
		{"admin", nil, nil, false},                   // Should be invalid
		{"Admin", nil, nil, false},                   // Should be invalid
		{"user123", nil, nil, true},                  // Valid username, not found
		{"existingUser", &models.User{}, nil, false}, // Already exists
		{"newUser", nil, nil, true},                  // Valid username, not found
		{"", nil, nil, true},                         // Valid username (empty)
	}

	for _, tt := range tests {
		if tt.username != "admin" && tt.username != "Admin" {
			if tt.mockReturn != nil {
				mockRepo.EXPECT().FindByUsername(tt.username).Return(tt.mockReturn, tt.mockErr)
			} else {
				mockRepo.EXPECT().FindByUsername(tt.username).Return(nil, tt.mockErr)
			}
		}

		t.Run(tt.username, func(t *testing.T) {
			result := userService.ValidateUsername(tt.username)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
