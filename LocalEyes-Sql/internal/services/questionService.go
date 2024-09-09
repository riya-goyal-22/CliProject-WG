package services

import (
	"localEyes/internal/interfaces"
	"localEyes/internal/models"
	"time"
)

type QuestionService struct {
	repo interfaces.QuestionRepository
}

func NewQuestionService(repo interfaces.QuestionRepository) *QuestionService {
	return &QuestionService{repo: repo}
}

func (s *QuestionService) AskQuestion(userId, postId int, content string) error {
	question := &models.Question{
		PostId:    postId,
		UserId:    userId,
		Text:      content,
		Replies:   make([]string, 0),
		CreatedAt: time.Now(),
	}
	return s.repo.Create(question)
}

func (s *QuestionService) DeleteQuesByPId(postId int) error {
	err := s.repo.DeleteByPId(postId)
	if err != nil {
		return err
	}
	return nil
}

func (s *QuestionService) DeleteUserQues(UId, QId int) error {
	err := s.repo.DeleteByQIdUId(QId, UId)
	if err != nil {
		return err
	}
	return nil
}

func (s *QuestionService) GetPostQuestions(PId int) ([]*models.Question, error) {
	questions, err := s.repo.GetQuestionsByPId(PId)
	if err != nil {
		return nil, err
	}
	return questions, nil
}

func (s *QuestionService) AddAnswer(QId int, answer string) error {
	err := s.repo.UpdateQuestion(QId, answer)
	if err != nil {
		return err
	}
	return nil
}
