package services_test

import (
	"errors"
	"localEyes/internal/models"
	"localEyes/internal/services"
	"localEyes/tests/mocks"
	"testing"
	_ "time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestQuestionService_AskQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockQuestionRepository(ctrl)
	questionService := services.NewQuestionService(mockRepo)

	tests := []struct {
		name    string
		userId  int
		postId  int
		content string
		mockErr error
		wantErr bool
	}{
		{
			name:    "successful question creation",
			userId:  1,
			postId:  1,
			content: "Is this place good?",
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "repo returns error on create",
			userId:  1,
			postId:  1,
			content: "Is this place good?",
			mockErr: errors.New("DB error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().Create(gomock.Any()).Return(tt.mockErr)

			err := questionService.AskQuestion(tt.userId, tt.postId, tt.content)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestQuestionService_DeleteQuesByPId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockQuestionRepository(ctrl)
	questionService := services.NewQuestionService(mockRepo)

	tests := []struct {
		name    string
		postId  int
		mockErr error
		wantErr bool
	}{
		{
			name:    "successful deletion",
			postId:  1,
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "repo returns error on delete",
			postId:  1,
			mockErr: errors.New("DB error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().DeleteByPId(tt.postId).Return(tt.mockErr)

			err := questionService.DeleteQuesByPId(tt.postId)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestQuestionService_GetPostQuestions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockQuestionRepository(ctrl)
	questionService := services.NewQuestionService(mockRepo)

	tests := []struct {
		name       string
		postId     int
		mockResult []*models.Question
		mockErr    error
		wantErr    bool
	}{
		{
			name:       "successful retrieval",
			postId:     1,
			mockResult: []*models.Question{{PostId: 1, Text: "Question 1"}},
			mockErr:    nil,
			wantErr:    false,
		},
		{
			name:       "repo returns error",
			postId:     1,
			mockResult: nil,
			mockErr:    errors.New("DB error"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().GetQuestionsByPId(tt.postId).Return(tt.mockResult, tt.mockErr)

			result, err := questionService.GetPostQuestions(tt.postId)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResult, result)
			}
		})
	}
}

func TestQuestionService_AddAnswer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockQuestionRepository(ctrl)
	questionService := services.NewQuestionService(mockRepo)

	tests := []struct {
		name    string
		qId     int
		answer  string
		mockErr error
		wantErr bool
	}{
		{
			name:    "successful update",
			qId:     1,
			answer:  "Yes, it's great!",
			mockErr: nil,
			wantErr: false,
		},
		{
			name:    "repo returns error on update",
			qId:     1,
			answer:  "Yes, it's great!",
			mockErr: errors.New("DB error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().UpdateQuestion(tt.qId, tt.answer).Return(tt.mockErr)

			err := questionService.AddAnswer(tt.qId, tt.answer)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteUserQues_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockQuestionRepository(ctrl)
	service := services.NewQuestionService(mockRepo)

	QId := 1
	UId := 123

	// Setup expectations
	mockRepo.EXPECT().DeleteByQIdUId(QId, UId).Return(nil)

	// Call the service method
	err := service.DeleteUserQues(UId, QId)

	// Assert that there was no error
	assert.NoError(t, err)
}

func TestDeleteUserQues_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockQuestionRepository(ctrl)
	service := services.NewQuestionService(mockRepo)

	QId := 1
	UId := 123
	expectedError := errors.New("error deleting question")

	// Setup expectations
	mockRepo.EXPECT().DeleteByQIdUId(QId, UId).Return(expectedError)

	// Call the service method
	err := service.DeleteUserQues(UId, QId)

	// Assert that the expected error is returned
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}
