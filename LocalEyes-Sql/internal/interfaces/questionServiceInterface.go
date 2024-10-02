package interfaces

import "localEyes/internal/models"

type QuestionServiceInterface interface {
	AskQuestion(userId string, postId string, content string) error
	DeleteQuesByPId(postId string) error
	DeleteUserQues(uId string, qId string) error
	GetPostQuestions(pId string) ([]*models.Question, error)
	AddAnswer(qId string, answer string) error
}
