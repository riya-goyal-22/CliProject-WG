package interfaces

import "localEyes/internal/models"

type AdminServiceInterface interface {
	GetAllUsers() ([]*models.User, error)
	GetAllPosts() ([]*models.Post, error)
	GetAllQuestions() ([]*models.Question, error)
	DeleteUser(uId string) error
	DeletePost(pId string) error
	DeleteQuestion(qId string) error
	ReActivate(uId string) error
}
