package utils_test

import (
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"localEyes/internal/models"
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

func TestValidateUsername_ReservedUsername(t *testing.T) {
	assert.False(t, utils.ValidateUsername("admin", nil), "Username 'admin' should not be valid")
	assert.False(t, utils.ValidateUsername("Admin", nil), "Username 'Admin' should not be valid")
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

// Test for SetTag function
func TestSetTag(t *testing.T) {
	tests := []struct {
		value    float64
		expected string
	}{
		{1.1, "resident"},
		{1.0, "newbie"},
		{0.5, "newbie"},
	}

	for _, test := range tests {
		result := utils.SetTag(test.value)
		if result != test.expected {
			t.Errorf("SetTag(%f) = %s; expected %s", test.value, result, test.expected)
		}
	}
}

// Test for ValidatePostRequest function
func TestValidatePostRequest(t *testing.T) {
	tests := []struct {
		post     models.RequestPost
		expected bool
		errMsg   string
	}{
		{models.RequestPost{"Title", "Content", "food"}, true, ""},
		{models.RequestPost{"", "Content", "food"}, false, "required field 'title' is missing"},
		{models.RequestPost{"Title", "", "food"}, false, "required field 'content' is missing"},
		{models.RequestPost{"Title", "Content", ""}, false, "required field 'type' is missing"},
		{models.RequestPost{"Title", "Content", "invalid"}, false, "invalid post type"},
	}

	for _, test := range tests {
		result, err := utils.ValidatePostRequest(test.post)
		if result != test.expected {
			t.Errorf("ValidatePostRequest(%v) = %v; expected %v", test.post, result, test.expected)
		}
		if test.expected == false && err != nil {
			if err.Error() != test.errMsg {
				t.Errorf("expected error message '%s', got '%s'", test.errMsg, err.Error())
			}
		} else if test.expected == true && err != nil {
			t.Errorf("unexpected error for input %v: %v", test.post, err)
		}
	}
}
