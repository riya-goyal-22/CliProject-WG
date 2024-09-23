package interfaces

import "localEyes/internal/models"

type UserServiceInterface interface {
	Signup(username string, password string, dwellingAge int, tag string) error
	Login(username string, password string) (*models.User, error)
	DeActivate(uid string) error
	NotifyUsers(uid string, title string) error
	GetNotifications(uid string) (*[]string, error)
	GetUserById(uid string) (*models.User, error)
	ValidateUsername(username string) bool
}
