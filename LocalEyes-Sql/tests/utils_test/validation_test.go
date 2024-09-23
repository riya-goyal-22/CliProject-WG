package utils_test

import (
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"localEyes/tests/mocks"
	"localEyes/utils"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		expected bool
	}{
		{"password1@", true},    // Valid password
		{"short1@", true},       // valid
		{"longpassword", false}, // No special character
		{"Password@", false},    // No number
	}

	for _, test := range tests {
		result := utils.ValidatePassword(test.password)
		assert.Equal(t, test.expected, result)
	}
}

func TestValidateFilter(t *testing.T) {
	tests := []struct {
		filter   string
		expected bool
	}{
		{"food", true},
		{"travel", true},
		{"shopping", true},
		{"other", true},
		{"invalid", false},
	}

	for _, test := range tests {
		result := utils.ValidateFilter(test.filter)
		assert.Equal(t, test.expected, result)
	}
}

func TestValidateUsername_UsernameExists(t *testing.T) {
	// Initialize gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a new MockUserRepository
	mockRepo := mock.NewMockUserRepository(ctrl)

	// Simulate FindByUsername returning an existing user (no error)
	mockRepo.EXPECT().FindByUsername("existinguser").Return(nil, nil)

	// Test for existing username in the repository
	result := utils.ValidateUsername("existinguser", mockRepo)
	assert.False(t, result, "Existing username should not be valid")
}

func TestValidateUsername_UsernameNotFound(t *testing.T) {
	// Initialize gomock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a new MockUserRepository
	mockRepo := mock.NewMockUserRepository(ctrl)

	// Simulate FindByUsername returning sql.ErrNoRows for non-existing username
	mockRepo.EXPECT().FindByUsername("newuser").Return(nil, sql.ErrNoRows)

	// Test for valid username when not found in the repository
	result := utils.ValidateUsername("newuser", mockRepo)
	assert.True(t, result, "Username not found in the repository should be valid")
}
