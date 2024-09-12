package utils

import (
	"localEyes/internal/interfaces"
	"strings"
)

func ValidateUsername(username string, userRepo interfaces.UserRepository) bool {
	if username == "admin" || username == "Admin" {
		return false
	}
	user, err := userRepo.FindByUsername(username)
	if user == nil || err != nil {
		return true
	}
	//if errors.Is(err, sql.ErrNoRows) {
	//	return true
	//}
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
