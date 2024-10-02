package services

import (
	"localEyes/internal/interfaces"
	"localEyes/internal/models"
	"localEyes/utils"
	"time"
)

type QuestionService struct {
	repo interfaces.QuestionRepository
}

func NewQuestionService(repo interfaces.QuestionRepository) *QuestionService {
	return &QuestionService{repo: repo}
}

func (s *QuestionService) AskQuestion(userId, postId string, content string) error {
	question := &models.Question{
		PostId:    postId,
		UserId:    userId,
		QId:       utils.GenerateRandomId(),
		Text:      content,
		Replies:   make([]string, 0),
		CreatedAt: time.Now(),
	}
	return s.repo.Create(question)
}

func (s *QuestionService) DeleteQuesByPId(postId string) error {
	err := s.repo.DeleteByPId(postId)
	if err != nil {
		return err
	}
	return nil
}

func (s *QuestionService) DeleteUserQues(uId, qId string) error {
	err := s.repo.DeleteByQIdUId(qId, uId)
	if err != nil {
		return err
	}
	return nil
}

func (s *QuestionService) GetPostQuestions(pId string) ([]*models.Question, error) {
	questions, err := s.repo.GetQuestionsByPId(pId)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (s *QuestionService) AddAnswer(qId, answer string) error {
	err := s.repo.UpdateQuestion(qId, answer)
	if err != nil {
		return err
	}
	return nil
}
