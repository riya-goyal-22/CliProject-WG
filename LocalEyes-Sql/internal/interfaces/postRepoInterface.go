package interfaces

import (
	"localEyes/internal/models"
)

type PostRepository interface {
	Create(post *models.Post) error
	GetAllPosts() ([]*models.Post, error)
	DeleteByPId(PId int) error
	DeleteByUId(UId int) error
	GetPostsByFilter(filter string) ([]*models.Post, error)
	GetPostsByUId(UId int) ([]*models.Post, error)
	UpdateUserPost(PId int, UId int, title string, content string) error
	UpdateLike(PId int) error
	DeleteByUIdPId(UId, PId int) error
	GetPostsByPId(PId int) ([]*models.Post, error)
}
