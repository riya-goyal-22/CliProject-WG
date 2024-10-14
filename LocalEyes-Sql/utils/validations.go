package utils

import (
	"localEyes/internal/interfaces"
	"localEyes/internal/models"
	"strings"
)

func ValidateUsername(username string, userRepo interfaces.UserRepository) bool {
	if username == "admin" || username == "Admin" {
		return false
	}
	user, err := userRepo.FindByUsername(username)
	if user == nil && err != nil {
		return true
	}
	return false
}

func ValidatePassword(password string) bool {
	if len(password) > 5 {
		if strings.Contains(password, "@") || strings.Contains(password, "#") || strings.Contains(password, "$") || strings.Contains(password, "%") || strings.Contains(password, "^") || strings.Contains(password, "*") {
			if strings.Contains(password, "1") || strings.Contains(password, "2") || strings.Contains(password, "3") || strings.Contains(password, "4") || strings.Contains(password, "5") || strings.Contains(password, "6") || strings.Contains(password, "7") || strings.Contains(password, "8") || strings.Contains(password, "9") || strings.Contains(password, "0") {
				return true
			}
		}
	}
	return false
}

func ValidateFilter(filter string) bool {
	return filter == "food" || filter == "travel" || filter == "shopping" || filter == "other" || filter == ""
}

func SetTag(value float64) string {
	if value > 1.0 {
		return "resident"
	}
	return "newbie"
}

func ValidatePostRequest(post models.RequestPost) (bool, error) {
	if post.Title == "" {
		return false, TitleMissing
	}
	if post.Content == "" {
		return false, ContentMissing
	}
	if post.Type == "" {
		return false, TypeMissing
	}
	if post.Type != "food" && post.Type != "shopping" && post.Type != "travel" && post.Type != "other" {
		return false, InvalidPost
	}
	return true, nil
}
