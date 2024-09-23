package interfaces

import (
	"localEyes/internal/models"
)

type PostRepository interface {
	Create(post *models.Post) error
	GetAllPosts() ([]*models.Post, error)
	DeleteByPId(pId string) error
	DeleteByUIdPId(uId string, pId string) error
	GetPostsByFilter(filter string) ([]*models.Post, error)
	GetPostsByUId(uId string) ([]*models.Post, error)
	GetPostByPId(pId string) (*models.PostWithQuestions, error)
	UpdateUserPost(pId string, uId string, title string, content string) error
	UpdateLike(pId string) error
}
