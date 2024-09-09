package interfaces

import (
	"localEyes/internal/models"
)

type QuestionRepository interface {
	Create(question *models.Question) error
	GetAllQuestions() ([]*models.Question, error)
	DeleteByQIdUId(QId, UId int) error
	DeleteByPId(PId int) error
	GetQuestionsByPId(PId int) ([]*models.Question, error)
	UpdateQuestion(QId int, answer string) error
	DeleteByQId(QId int) error
}
