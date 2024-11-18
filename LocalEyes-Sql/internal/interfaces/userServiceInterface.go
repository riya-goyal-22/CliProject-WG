package interfaces

import "localEyes/internal/models"

type UserServiceInterface interface {
	Signup(username string, password string, dwellingAge float64, answer string) error
	Login(username string, password string) (*models.User, error)
	DeActivate(uid string) error
	NotifyUsers(uid string, title string) error
	GetNotifications(uid string) ([]string, error)
	GetUserById(uid string) (*models.User, error)
	ValidateUsername(username string) bool
	PasswordReset(resetUser models.ResetPasswordUser) error
	UpdateUser(uId string, requestUser *models.Client) error
}
