package interfaces

import "localEyes/internal/models"

type AdminServiceInterface interface {
	GetAllUsers(limit, offset int, search string) ([]*models.User, error)
	GetAllPosts(limit, offset int, search, filter string) ([]*models.Post, error)
	GetAllQuestions() ([]*models.Question, error)
	DeleteUser(uId string) error
	DeletePost(pId string) error
	DeleteQuestion(qId string) error
	ReActivate(uId string) error
}
