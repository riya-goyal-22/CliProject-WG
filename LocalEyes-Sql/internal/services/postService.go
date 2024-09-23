package services

import (
	"localEyes/internal/interfaces"
	"localEyes/internal/models"
	"localEyes/utils"
	"time"
)

type PostService struct {
	repo     interfaces.PostRepository
	userRepo interfaces.UserRepository
}

func NewPostService(repo interfaces.PostRepository, userRepo interfaces.UserRepository) *PostService {
	return &PostService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *PostService) CreatePost(userId string, title, content, postType string) error {
	post := &models.Post{
		UId:       userId,
		PostId:    utils.GenerateRandomId(),
		Title:     title,
		Content:   content,
		Type:      postType,
		CreatedAt: time.Now(),
		Likes:     0,
	}
	err := s.repo.Create(post)
	if err != nil {
		return err
	}
	err = s.userRepo.PushNotification(userId, title)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostService) UpdateMyPost(postId, userId, title, content string) error {
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

func (s *PostService) GiveMyPosts(uId string) ([]*models.Post, error) {
	posts, err := s.repo.GetPostsByUId(uId)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *PostService) DeleteMyPost(uId, pId string) error {
	err := s.repo.DeleteByUIdPId(uId, pId)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostService) Like(pId string) error {
	err := s.repo.UpdateLike(pId)
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

//func (s *PostService) PostIdExist(pId string) (bool, error) {
//	posts, err := s.repo.GetPostByPId(pId)
//	if err != nil {
//		return false, err
//	}
//	return len(posts) > 0, nil
//}

func (s *PostService) GivePostById(pId string) (*models.PostWithQuestions, error) {
	post, err := s.repo.GetPostByPId(pId)
	if err != nil {
		return nil, err
	}
	return post, nil
}
