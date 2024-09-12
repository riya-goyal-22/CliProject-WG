package repositories_test

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"localEyes/config"
	"localEyes/internal/models"
	"localEyes/internal/repositories"
	"testing"
	"time"
)

func TestMySQLPostRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock instance: %v", err)
	}
	defer db.Close()

	repo := repositories.NewMySQLPostRepository(db)

	post := &models.Post{
		UId:       1,
		Title:     "Test Post",
		Type:      "food",
		Content:   "This is a test post",
		Likes:     0,
		CreatedAt: time.Now(),
	}

	mock.ExpectExec("INSERT INTO posts").WithArgs(post.UId, post.Title, post.Type, post.Content, post.Likes, post.CreatedAt).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(post)
	assert.NoError(t, err)
}

func TestMySQLPostRepository_GetAllPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock instance: %v", err)
	}
	defer db.Close()

	repo := repositories.NewMySQLPostRepository(db)

	rows := sqlmock.NewRows([]string{"post_id", "user_id", "title", "type", "content", "likes", "created_at"}).
		AddRow(1, 1, "Test Post", "food", "This is a test post", 0, "2006-01-02T15:04:05Z")

	mock.ExpectQuery(`^SELECT post_id, user_id, title, type, content, likes, created_at FROM posts$`).WillReturnRows(rows)

	posts, err := repo.GetAllPosts()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(posts) == 0 {
		t.Fatal("expected at least one post, but got none")
	}
	assert.Len(t, posts, 1)
	assert.Equal(t, "Test Post", posts[0].Title)
}

func TestMySQLPostRepository_DeleteByPId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock instance: %v", err)
	}
	defer db.Close()

	repo := repositories.NewMySQLPostRepository(db)

	mock.ExpectExec("DELETE FROM posts WHERE post_id = ?").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteByPId(1)
	assert.NoError(t, err)
}

func TestMySQLPostRepository_DeleteByPId_NoRowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock instance: %v", err)
	}
	defer db.Close()

	repo := repositories.NewMySQLPostRepository(db)

	mock.ExpectExec("DELETE FROM posts WHERE post_id = ?").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.DeleteByPId(1)
	assert.EqualError(t, err, config.Red+"No Post exist with this id"+config.Reset)
}

func TestMySQLPostRepository_GetPostsByFilter(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock instance: %v", err)
	}
	defer db.Close()

	repo := repositories.NewMySQLPostRepository(db)

	rows := sqlmock.NewRows([]string{"post_id", "user_id", "title", "type", "content", "likes", "created_at"}).
		AddRow(1, 1, "Test Post", "food", "This is a test post", 0, "2024-09-08 00:00:00")

	// Ensure the expected query matches exactly with the actual query
	mock.ExpectQuery(`^SELECT post_id, user_id, title,type, content, likes,created_at FROM posts WHERE type = \?$`).
		WithArgs("food").
		WillReturnRows(rows)

	posts, err := repo.GetPostsByFilter("food")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(posts) == 0 {
		t.Fatal("expected at least one post, but got none")
	}
	assert.Len(t, posts, 1)
	assert.Equal(t, "Test Post", posts[0].Title)
}

func TestMySQLPostRepository_UpdateUserPost(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock instance: %v", err)
	}
	defer db.Close()

	repo := repositories.NewMySQLPostRepository(db)

	mock.ExpectExec("UPDATE posts SET title = \\?, content = \\? WHERE post_id = \\? AND user_id = \\?").
		WithArgs("Updated Title", "Updated Content", 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateUserPost(1, 1, "Updated Title", "Updated Content")
	assert.NoError(t, err)
}

func TestMySQLPostRepository_UpdateUserPost_NoRowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock instance: %v", err)
	}
	defer db.Close()

	repo := repositories.NewMySQLPostRepository(db)

	mock.ExpectExec("UPDATE posts SET title = \\?, content = \\? WHERE post_id = \\? AND user_id = \\?").
		WithArgs("Updated Title", "Updated Content", 1, 1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.UpdateUserPost(1, 1, "Updated Title", "Updated Content")
	assert.EqualError(t, err, config.Red+"You can only update your post"+config.Reset)
}

func TestMySQLPostRepository_UpdateLike(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock instance: %v", err)
	}
	defer db.Close()

	repo := repositories.NewMySQLPostRepository(db)

	mock.ExpectExec("^UPDATE posts SET likes = likes \\+ 1 WHERE post_id = \\?$").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateLike(1)
	assert.NoError(t, err)
}

func TestMySQLPostRepository_UpdateLike_NoRowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock instance: %v", err)
	}
	defer db.Close()

	repo := repositories.NewMySQLPostRepository(db)

	// Mock the expectation for the query
	mock.ExpectExec(`^UPDATE posts SET likes = likes \+ 1 WHERE post_id = \?$`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected

	// Call the method to be tested
	err = repo.UpdateLike(1)
	if err == nil {
		t.Fatal("expected error, but got nil")
	}
	expectedErr := config.Red + "No post exist with this id" + config.Reset
	if err.Error() != expectedErr {
		t.Errorf("Error message not equal:\n\texpected: %v\n\tactual  : %v", expectedErr, err.Error())
	}
}

func TestDeleteByUIdPId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	UId := 1
	PId := 10

	mock.ExpectExec("^DELETE FROM posts WHERE post_id = \\? AND user_id=\\?$").
		WithArgs(PId, UId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := repositories.NewMySQLPostRepository(db)
	err = repo.DeleteByUIdPId(UId, PId)

	assert.NoError(t, err)
}

func TestDeleteByUIdPId_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	UId := 1
	PId := 999

	mock.ExpectExec("^DELETE FROM posts WHERE post_id = \\? AND user_id=\\?$").
		WithArgs(PId, UId).
		WillReturnResult(sqlmock.NewResult(0, 0))

	repo := repositories.NewMySQLPostRepository(db)
	err = repo.DeleteByUIdPId(UId, PId)

	assert.Error(t, err)
	assert.EqualError(t, err, config.Red+"No Post exist with this id"+config.Reset)
}

func TestDeleteByUId_Successful(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	UId := 1

	mock.ExpectExec("^DELETE FROM posts WHERE user_id = \\?$").
		WithArgs(UId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := repositories.NewMySQLPostRepository(db)
	err = repo.DeleteByUId(UId)

	assert.NoError(t, err)
}

func TestDeleteByUId_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	UId := 1

	mock.ExpectExec("^DELETE FROM posts WHERE user_id = \\?$").
		WithArgs(UId).
		WillReturnError(errors.New("delete error"))

	repo := repositories.NewMySQLPostRepository(db)
	err = repo.DeleteByUId(UId)

	assert.Error(t, err)
	assert.EqualError(t, err, "delete error")
}

func TestGetPostsByUId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	UId := 1
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{"post_id", "user_id", "title", "type", "content", "likes", "created_at"}).
		AddRow(1, UId, "Title 1", "food", "Content 1", 10, createdAt).
		AddRow(2, UId, "Title 2", "travel", "Content 2", 15, createdAt)

	mock.ExpectQuery("^SELECT post_id, user_id, title,type, content, likes,created_at FROM posts WHERE user_id = \\?$").
		WithArgs(UId).
		WillReturnRows(rows)

	repo := repositories.NewMySQLPostRepository(db)
	posts, err := repo.GetPostsByUId(UId)

	assert.NoError(t, err)
	assert.Len(t, posts, 2)
	assert.Equal(t, 1, posts[0].PostId)
	assert.Equal(t, "Title 1", posts[0].Title)
	assert.Equal(t, 10, posts[0].Likes)
}

func TestGetPostsByUId_NoPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	UId := 1

	rows := sqlmock.NewRows([]string{"post_id", "user_id", "title", "type", "content", "likes", "created_at"})

	mock.ExpectQuery("^SELECT post_id, user_id, title,type, content, likes,created_at FROM posts WHERE user_id = \\?$").
		WithArgs(UId).
		WillReturnRows(rows)

	repo := repositories.NewMySQLPostRepository(db)
	posts, err := repo.GetPostsByUId(UId)

	assert.NoError(t, err)
	assert.Len(t, posts, 0)
}

func TestGetPostsByUId_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	UId := 1

	mock.ExpectQuery("^SELECT post_id, user_id, title,type, content, likes,created_at FROM posts WHERE user_id = \\?$").
		WithArgs(UId).
		WillReturnError(errors.New("query error"))

	repo := repositories.NewMySQLPostRepository(db)
	posts, err := repo.GetPostsByUId(UId)

	assert.Error(t, err)
	assert.Nil(t, posts)
	assert.EqualError(t, err, "query error")
}

func TestGetPostsByPId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	PId := 1
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{"post_id", "user_id", "title", "type", "content", "likes", "created_at"}).
		AddRow(PId, 1, "Title 1", "food", "Content 1", 10, createdAt)

	mock.ExpectQuery("^SELECT post_id, user_id, title,type, content, likes,created_at FROM posts WHERE post_id = \\?$").
		WithArgs(PId).
		WillReturnRows(rows)

	repo := repositories.NewMySQLPostRepository(db)
	posts, err := repo.GetPostsByPId(PId)

	assert.NoError(t, err)
	assert.Len(t, posts, 1)
	assert.Equal(t, PId, posts[0].PostId)
	assert.Equal(t, "Title 1", posts[0].Title)
	assert.Equal(t, 10, posts[0].Likes)
}

func TestGetPostsByPId_NoPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	PId := 1

	rows := sqlmock.NewRows([]string{"post_id", "user_id", "title", "type", "content", "likes", "created_at"})

	mock.ExpectQuery("^SELECT post_id, user_id, title,type, content, likes,created_at FROM posts WHERE post_id = \\?$").
		WithArgs(PId).
		WillReturnRows(rows)

	repo := repositories.NewMySQLPostRepository(db)
	posts, err := repo.GetPostsByPId(PId)

	assert.NoError(t, err)
	assert.Len(t, posts, 0)
}

func TestGetPostsByPId_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	PId := 1

	mock.ExpectQuery("^SELECT post_id, user_id, title,type, content, likes,created_at FROM posts WHERE post_id = \\?").
		WithArgs(PId).
		WillReturnError(errors.New("query error"))

	repo := repositories.NewMySQLPostRepository(db)
	posts, err := repo.GetPostsByPId(PId)

	assert.Error(t, err)
	assert.Nil(t, posts)
	assert.EqualError(t, err, "query error")
}
