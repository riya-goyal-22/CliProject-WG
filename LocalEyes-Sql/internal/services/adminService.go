package services

import (
	"localEyes/internal/interfaces"
	"localEyes/internal/models"
)

type AdminService struct {
	UserRepo interfaces.UserRepository
	PostRepo interfaces.PostRepository
	QuesRepo interfaces.QuestionRepository
}

func NewAdminService(userRepo interfaces.UserRepository, postRepo interfaces.PostRepository, quesRepo interfaces.QuestionRepository) *AdminService {
	return &AdminService{UserRepo: userRepo, PostRepo: postRepo, QuesRepo: quesRepo}
}

func (s *AdminService) GetAllUsers() ([]*models.User, error) {
	users, err := s.UserRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *AdminService) GetAllPosts() ([]*models.Post, error) {
	posts, err := s.PostRepo.GetAllPosts()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *AdminService) GetAllQuestions() ([]*models.Question, error) {
	questions, err := s.QuesRepo.GetAllQuestions()
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (s *AdminService) DeleteUser(uId string) error {
	err1 := s.UserRepo.DeleteByUId(uId)
	//err2 := s.PostRepo.DeleteByUId(UId)  //post not deleted when a user is deleted
	if err1 != nil {
		return err1
	}
	return nil
}

func (s *AdminService) DeletePost(pId string) error {
	err1 := s.PostRepo.DeleteByPId(pId)
	err2 := s.QuesRepo.DeleteByPId(pId)
	if err1 != nil {
		return err1
	} else if err2 != nil {
		return err2
	}
	return nil
}

func (s *AdminService) DeleteQuestion(qId string) error {
	err := s.QuesRepo.DeleteByQId(qId)
	if err != nil {
		return err
	}
	return nil
}

func (s *AdminService) ReActivate(uId string) error {
	err := s.UserRepo.UpdateActiveStatus(uId, true)
	if err != nil {
		return err
	}
	return nil
}
