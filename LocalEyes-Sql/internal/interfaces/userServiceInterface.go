package interfaces

import "localEyes/internal/models"

type UserServiceInterface interface {
	Signup(username string, password string, email string, dwellingAge float64) error
	Login(username string, password string) (*models.User, error)
	DeActivate(uid string) error
	NotifyUsers(uid string, title string) error
	GetNotifications(uid string) ([]string, error)
	GetUserById(uid string) (*models.User, error)
	ValidateUsername(username string) bool
	ValidateEmail(email string) bool
	SendOtp(email string) error
	PasswordReset(resetUser models.ResetPasswordUser) error
	UpdateUser(uId string, requestUser *models.Client) error
}
