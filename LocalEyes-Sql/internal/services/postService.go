package services

import (
	"localEyes/internal/interfaces"
	"localEyes/internal/models"
	"time"
)

type PostService struct {
	repo interfaces.PostRepository
}

func NewPostService(repo interfaces.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) CreatePost(userId int, title, content, postType string) error {
	post := &models.Post{
		UId:       userId,
		Title:     title,
		Content:   content,
		Type:      postType,
		CreatedAt: time.Now(),
		Likes:     0,
	}
	err := s.repo.Create(post)
	return err
}

func (s *PostService) UpdateMyPost(postId, userId int, title, content string) error {
	err := s.repo.UpdateUserPost(postId, userId, title, content)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostService) GiveAllPosts() ([]*models.Post, error) {
	posts, err := s.repo.GetAllPosts()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostService) GiveMyPosts(UId int) ([]*models.Post, error) {
	posts, err := s.repo.GetPostsByUId(UId)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostService) DeleteMyPost(UId, PId int) error {
	err := s.repo.DeleteByUIdPId(UId, PId)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostService) Like(PId int) error {
	err := s.repo.UpdateLike(PId)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostService) GiveFilteredPosts(filterType string) ([]*models.Post, error) {
	posts, err := s.repo.GetPostsByFilter(filterType)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostService) PostIdExist(PId int) (bool, error) {
	posts, err := s.repo.GetPostsByPId(PId)
	if err != nil {
		return false, err
	}
	return len(posts) > 0, nil
}
