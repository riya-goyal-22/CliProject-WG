package interfaces

import (
	"localEyes/internal/models"
)

type QuestionRepository interface {
	Create(question *models.Question) error
	GetAllQuestions() ([]*models.Question, error)
	DeleteByQIdUId(qId string, uId string) error
	DeleteByPId(pId string) error
	DeleteByQId(qId string) error
	GetQuestionsByPId(pId string) ([]*models.Question, error)
	UpdateQuestion(qId string, answer string) error
}
