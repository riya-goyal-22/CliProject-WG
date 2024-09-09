package services_test

import (
	"localEyes/internal/models"
	"localEyes/internal/services"
	"localEyes/tests/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)
	service := services.NewPostService(mockRepo)

	post := &models.Post{
		UId:       1,
		Title:     "Test Post",
		Content:   "Test Content",
		Type:      "Travel",
		CreatedAt: time.Now(),
		Likes:     0,
	}

	// Set up expectations
	mockRepo.EXPECT().Create(post).Return(nil)

	// Call the method
	err := service.CreatePost(post.UId, post.Title, post.Content, post.Type)

	// Assert results
	assert.NoError(t, err)
}

func TestUpdateMyPost(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)
	service := services.NewPostService(mockRepo)

	postId := 1
	userId := 1
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

	mockRepo := mocks.NewMockPostRepository(ctrl)
	service := services.NewPostService(mockRepo)

	posts := []*models.Post{
		{UId: 1, Title: "Post 1"},
		{UId: 2, Title: "Post 2"},
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

	mockRepo := mocks.NewMockPostRepository(ctrl)
	service := services.NewPostService(mockRepo)

	userId := 1
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

	mockRepo := mocks.NewMockPostRepository(ctrl)
	service := services.NewPostService(mockRepo)

	userId := 1
	postId := 1

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

	mockRepo := mocks.NewMockPostRepository(ctrl)
	service := services.NewPostService(mockRepo)

	postId := 1

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

	mockRepo := mocks.NewMockPostRepository(ctrl)
	service := services.NewPostService(mockRepo)

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

func TestPostIdExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockPostRepository(ctrl)
	service := services.NewPostService(mockRepo)

	postId := 1
	posts := []*models.Post{
		{UId: 1, Title: "Existing Post"},
	}

	// Set up expectations
	mockRepo.EXPECT().GetPostsByPId(postId).Return(posts, nil)

	// Call the method
	exists, err := service.PostIdExist(postId)

	// Assert results
	assert.NoError(t, err)
	assert.True(t, exists)
}
