package services_test

import (
	"errors"
	"localEyes/constants"
	"localEyes/internal/models"
	"localEyes/internal/services"
	"localEyes/tests/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAdminService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	adminService := services.NewAdminService(mockUserRepo, nil, nil)

	password := "admin123"
	hashedPassword := services.HashPassword(password)

	mockAdmin := &models.Admin{
		models.User{Username: "admin",
			Password: hashedPassword},
	}

	mockUserRepo.EXPECT().
		FindAdminByUsernamePassword("admin", hashedPassword).
		Return(mockAdmin, nil)

	admin, err := adminService.Login(password)
	assert.NoError(t, err)
	assert.NotNil(t, admin)
	assert.Equal(t, "admin", admin.User.Username)
}

func TestAdminService_Login_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	adminService := services.NewAdminService(mockUserRepo, nil, nil)

	password := "wrongpassword"
	hashedPassword := services.HashPassword(password)

	mockUserRepo.EXPECT().
		FindAdminByUsernamePassword("admin", hashedPassword).
		Return(nil, errors.New("Invalid username or password"))

	admin, err := adminService.Login(password)
	assert.Error(t, err)
	assert.Equal(t, constants.Red+"Invalid username or password"+constants.Reset, err.Error())
	assert.Nil(t, admin)
}

func TestAdminService_GetAllUsers_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	adminService := services.NewAdminService(mockUserRepo, nil, nil)

	mockUsers := []*models.User{
		{Username: "user1"},
		{Username: "user2"},
	}

	mockUserRepo.EXPECT().
		GetAllUsers().
		Return(mockUsers, nil)

	users, err := adminService.GetAllUsers()
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, len(mockUsers), len(users))
}

func TestAdminService_GetAllUsers_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	adminService := services.NewAdminService(mockUserRepo, nil, nil)

	mockUserRepo.EXPECT().
		GetAllUsers().
		Return(nil, errors.New("database error"))

	users, err := adminService.GetAllUsers()
	assert.Error(t, err)
	assert.Nil(t, users)
}

func TestAdminService_GetAllPosts_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mocks.NewMockPostRepository(ctrl)
	adminService := services.NewAdminService(nil, mockPostRepo, nil)

	mockPosts := []*models.Post{
		{Title: "Post 1"},
		{Title: "Post 2"},
	}

	mockPostRepo.EXPECT().
		GetAllPosts().
		Return(mockPosts, nil)

	posts, err := adminService.GetAllPosts()
	assert.NoError(t, err)
	assert.NotNil(t, posts)
	assert.Equal(t, len(mockPosts), len(posts))
}

func TestAdminService_GetAllPosts_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mocks.NewMockPostRepository(ctrl)
	adminService := services.NewAdminService(nil, mockPostRepo, nil)

	mockPostRepo.EXPECT().
		GetAllPosts().
		Return(nil, errors.New("database error"))

	posts, err := adminService.GetAllPosts()
	assert.Error(t, err)
	assert.Nil(t, posts)
}

func TestAdminService_GetAllQuestions_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuesRepo := mocks.NewMockQuestionRepository(ctrl)
	adminService := services.NewAdminService(nil, nil, mockQuesRepo)

	mockQuestions := []*models.Question{
		{Text: "Question 1"},
		{Text: "Question 2"},
	}

	mockQuesRepo.EXPECT().
		GetAllQuestions().
		Return(mockQuestions, nil)

	questions, err := adminService.GetAllQuestions()
	assert.NoError(t, err)
	assert.NotNil(t, questions)
	assert.Equal(t, len(mockQuestions), len(questions))
}

func TestAdminService_GetAllQuestions_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuesRepo := mocks.NewMockQuestionRepository(ctrl)
	adminService := services.NewAdminService(nil, nil, mockQuesRepo)

	mockQuesRepo.EXPECT().
		GetAllQuestions().
		Return(nil, errors.New("database error"))

	questions, err := adminService.GetAllQuestions()
	assert.Error(t, err)
	assert.Nil(t, questions)
}

func TestAdminService_DeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPostRepo := mocks.NewMockPostRepository(ctrl)
	adminService := services.NewAdminService(mockUserRepo, mockPostRepo, nil)

	mockUserRepo.EXPECT().
		DeleteByUId(1).
		Return(nil)

	mockPostRepo.EXPECT().
		DeleteByUId(1).
		Return(nil)

	err := adminService.DeleteUser(1)
	assert.NoError(t, err)
}

func TestAdminService_DeleteUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockPostRepo := mocks.NewMockPostRepository(ctrl)
	adminService := services.NewAdminService(mockUserRepo, mockPostRepo, nil)

	mockUserRepo.EXPECT().
		DeleteByUId(1).
		Return(errors.New("delete user error"))
	mockPostRepo.EXPECT().
		DeleteByUId(1).
		Return(nil)
	err := adminService.DeleteUser(1)
	assert.Error(t, err)
}

func TestAdminService_DeletePost_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mocks.NewMockPostRepository(ctrl)
	mockQuesRepo := mocks.NewMockQuestionRepository(ctrl)
	adminService := services.NewAdminService(nil, mockPostRepo, mockQuesRepo)

	mockPostRepo.EXPECT().
		DeleteByPId(1).
		Return(nil)

	mockQuesRepo.EXPECT().
		DeleteByPId(1).
		Return(nil)

	err := adminService.DeletePost(1)
	assert.NoError(t, err)
}

func TestAdminService_DeletePost_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mocks.NewMockPostRepository(ctrl)
	mockQuesRepo := mocks.NewMockQuestionRepository(ctrl)
	adminService := services.NewAdminService(nil, mockPostRepo, mockQuesRepo)

	mockPostRepo.EXPECT().
		DeleteByPId(1).
		Return(errors.New("delete post error"))
	mockQuesRepo.EXPECT().
		DeleteByPId(1).
		Return(nil)
	err := adminService.DeletePost(1)
	assert.Error(t, err)
}

func TestAdminService_DeleteQuestion_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuesRepo := mocks.NewMockQuestionRepository(ctrl)
	adminService := services.NewAdminService(nil, nil, mockQuesRepo)

	mockQuesRepo.EXPECT().
		DeleteByQId(1).
		Return(nil)

	err := adminService.DeleteQuestion(1)
	assert.NoError(t, err)
}

func TestAdminService_DeleteQuestion_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuesRepo := mocks.NewMockQuestionRepository(ctrl)
	adminService := services.NewAdminService(nil, nil, mockQuesRepo)

	mockQuesRepo.EXPECT().
		DeleteByQId(1).
		Return(errors.New("delete question error"))

	err := adminService.DeleteQuestion(1)
	assert.Error(t, err)
}

func TestAdminService_ReActivate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	adminService := services.NewAdminService(mockUserRepo, nil, nil)

	mockUserRepo.EXPECT().
		UpdateActiveStatus(1, true).
		Return(nil)

	err := adminService.ReActivate(1)
	assert.NoError(t, err)
}

func TestAdminService_ReActivate_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	adminService := services.NewAdminService(mockUserRepo, nil, nil)

	mockUserRepo.EXPECT().
		UpdateActiveStatus(1, true).
		Return(errors.New("reactivate error"))

	err := adminService.ReActivate(1)
	assert.Error(t, err)
}
