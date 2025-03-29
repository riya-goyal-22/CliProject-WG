package services_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"localEyes/internal/models"
	"localEyes/internal/services"
	"localEyes/tests/mocks"
	"testing"
)

func TestCreatePost_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mock.NewMockPostRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)

	postService := services.NewPostService(mockPostRepo, mockUserRepo)

	userId := "user-1"
	title := "Test Title"
	content := "Test Content"
	postType := "general"

	mockPostRepo.EXPECT().Create(gomock.Any()).Return(nil)
	mockUserRepo.EXPECT().PushNotification(userId, title).Return(nil)

	err := postService.CreatePost(userId, title, content, postType)

	assert.NoError(t, err)
}

func TestCreatePost_ErrorOnCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mock.NewMockPostRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)

	postService := services.NewPostService(mockPostRepo, mockUserRepo)

	userId := "user-1"
	title := "Test Title"
	content := "Test Content"
	postType := "general"

	mockPostRepo.EXPECT().Create(gomock.Any()).Return(errors.New("create error"))

	err := postService.CreatePost(userId, title, content, postType)

	assert.Error(t, err)
	assert.Equal(t, "create error", err.Error())
}

func TestCreatePost_ErrorOnNotification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPostRepo := mock.NewMockPostRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)

	postService := services.NewPostService(mockPostRepo, mockUserRepo)

	userId := "user-1"
	title := "Test Title"
	content := "Test Content"
	postType := "general"

	mockPostRepo.EXPECT().Create(gomock.Any()).Return(nil)
	mockUserRepo.EXPECT().PushNotification(userId, title).Return(errors.New("notification error"))

	err := postService.CreatePost(userId, title, content, postType)

	assert.Error(t, err)
	assert.Equal(t, "notification error", err.Error())
}

func TestUpdateMyPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockPostRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	service := services.NewPostService(mockRepo, mockUserRepo)

	postId := "1"
	userId := "1"
	title := "Updated Title"
	content := "Updated Content"

	// Set up expectations
	mockRepo.EXPECT().UpdateUserPost(postId, userId, title, content).Return(nil)

	// Call the method
	err := service.UpdateMyPost(postId, userId, title, content)

	// Assert results
	assert.NoError(t, err)
}

func TestGiveAllPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockPostRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	service := services.NewPostService(mockRepo, mockUserRepo)

	posts := []*models.Post{
		{UId: "1", Title: "Post 1"},
		{UId: "2", Title: "Post 2"},
	}

	// Set up expectations
	mockRepo.EXPECT().GetAllPosts().Return(posts, nil)

	// Call the method
	result, err := service.GiveAllPosts()

	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, posts, result)
}

func TestGiveMyPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockPostRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	service := services.NewPostService(mockRepo, mockUserRepo)

	userId := "1"
	posts := []*models.Post{
		{UId: userId, Title: "My Post 1"},
		{UId: userId, Title: "My Post 2"},
	}

	// Set up expectations
	mockRepo.EXPECT().GetPostsByUId(userId).Return(posts, nil)

	// Call the method
	result, err := service.GiveMyPosts(userId)

	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, posts, result)
}

func TestDeleteMyPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockPostRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	service := services.NewPostService(mockRepo, mockUserRepo)

	userId := "1"
	postId := "1"

	// Set up expectations
	mockRepo.EXPECT().DeleteByUIdPId(userId, postId).Return(nil)

	// Call the method
	err := service.DeleteMyPost(userId, postId)

	// Assert results
	assert.NoError(t, err)
}

func TestLike(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockPostRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	service := services.NewPostService(mockRepo, mockUserRepo)

	postId := "1"

	// Set up expectations
	mockRepo.EXPECT().UpdateLike(postId).Return(nil)

	// Call the method
	err := service.Like(postId)

	// Assert results
	assert.NoError(t, err)
}

func TestGiveFilteredPosts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockPostRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	service := services.NewPostService(mockRepo, mockUserRepo)

	filterType := "Food"
	posts := []*models.Post{
		{Title: "Food Post 1"},
		{Title: "Food Post 2"},
	}

	// Set up expectations
	mockRepo.EXPECT().GetPostsByFilter(filterType).Return(posts, nil)

	// Call the method
	result, err := service.GiveFilteredPosts(filterType)

	// Assert results
	assert.NoError(t, err)
	assert.Equal(t, posts, result)
}
