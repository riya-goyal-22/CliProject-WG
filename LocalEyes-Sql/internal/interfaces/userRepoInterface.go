package interfaces

import "localEyes/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	FindByUId(uId string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByUsernamePassword(username string, password string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	DeleteByUId(uId string) error
	UpdateActiveStatus(uId string, status bool) error
	PushNotification(uId string, title string) error
	ClearNotification(uId string) error
}
