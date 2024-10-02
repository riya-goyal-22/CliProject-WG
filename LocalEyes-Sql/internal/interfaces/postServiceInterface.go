package interfaces

import "localEyes/internal/models"

type PostServiceInterface interface {
	CreatePost(userId string, title string, content string, postType string) error
	UpdateMyPost(postId string, userId string, title string, content string) error
	GiveAllPosts() ([]*models.Post, error)
	GiveMyPosts(uId string) ([]*models.Post, error)
	DeleteMyPost(uId string, pId string) error
	Like(pId string) error
	GiveFilteredPosts(filterType string) ([]*models.Post, error)
	GivePostById(pId string) (*models.PostWithQuestions, error)
}
