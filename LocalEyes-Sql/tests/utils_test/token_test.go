package utils_test

import (
	"github.com/stretchr/testify/assert"
	"localEyes/utils"
	"os"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	os.Setenv("Secret", "mysecret")
	os.Setenv("AdminUsername", "admin")

	username := "testuser"
	uid := "12345"
	token, err := utils.GenerateToken(username, uid)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate token
	claims, err := utils.ExtractClaims("Bearer " + token)
	assert.NoError(t, err)
	assert.Equal(t, username, claims["sub"])
	assert.Equal(t, uid, claims["id"])
}

func TestValidateToken_Valid(t *testing.T) {
	os.Setenv("Secret", "mysecret")
	token, _ := utils.GenerateToken("testuser", "12345")

	valid := utils.ValidateToken("Bearer " + token)
	assert.True(t, valid)
}

func TestValidateToken_Invalid(t *testing.T) {
	valid := utils.ValidateToken("Bearer invalidtoken")
	assert.False(t, valid)
}

func TestValidateAdminToken_Valid(t *testing.T) {
	os.Setenv("Secret", "mysecret")
	os.Setenv("AdminUsername", "admin")
	token, _ := utils.GenerateToken("admin", "12345")

	valid := utils.ValidateAdminToken("Bearer " + token)
	assert.True(t, valid)
}

func TestValidateAdminToken_Invalid(t *testing.T) {
	os.Setenv("Secret", "mysecret")
	os.Setenv("AdminUsername", "admin")
	token, _ := utils.GenerateToken("testuser", "12345")

	valid := utils.ValidateAdminToken("Bearer " + token)
	assert.False(t, valid)
}

func TestExtractClaims_Valid(t *testing.T) {
	os.Setenv("Secret", "mysecret")
	token, _ := utils.GenerateToken("testuser", "12345")

	claims, err := utils.ExtractClaims("Bearer " + token)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", claims["sub"])
	assert.Equal(t, "12345", claims["id"])
}

func TestExtractClaims_Invalid(t *testing.T) {
	claims, err := utils.ExtractClaims("Bearer invalidtoken")
	assert.Error(t, err)
	assert.Nil(t, claims)
}
