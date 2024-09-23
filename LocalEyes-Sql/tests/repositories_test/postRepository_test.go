package repositories_test

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"localEyes/internal/models"
	"localEyes/internal/repositories"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLPostRepository(db)

	post := &models.Post{
		UId:       "user-uuid",
		PostId:    "post-uuid",
		Title:     "Test Post",
		Type:      "Type A",
		Content:   "This is a test post.",
		Likes:     0,
		CreatedAt: time.Now(),
	}

	mock.ExpectExec("INSERT INTO posts").
		WithArgs(post.UId, post.PostId, post.Title, post.Type, post.Content, post.Likes, post.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(post)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLPostRepository(db)

	rows := sqlmock.NewRows([]string{"post_id", "uuid", "title", "type", "content", "likes", "created_at"}).
		AddRow("post-uuid-1", "user-uuid", "Post 1", "Type A", "Content 1", 10, "2024-09-22T16:52:56Z").
		AddRow("post-uuid-2", "user-uuid", "Post 2", "Type B", "Content 2", 5, "2024-09-22T16:52:56Z")

	mock.ExpectQuery("SELECT post_id, uuid, title, type, content, likes, created_at FROM posts").
		WillReturnRows(rows)

	posts, err := repo.GetAllPosts()
	assert.NoError(t, err)
	assert.Len(t, posts, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteByPId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLPostRepository(db)

	postId := "post-uuid"

	mock.ExpectExec("DELETE FROM posts WHERE post_id = ?").
		WithArgs(postId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteByPId(postId)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLPostRepository_DeleteByPId_NoRowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock instance: %v", err)
	}
	defer db.Close()
	repo := repositories.NewMySQLPostRepository(db)
	postId := "post-uuid"
	mock.ExpectExec("DELETE FROM posts WHERE post_id = ?").
		WithArgs(postId).
		WillReturnResult(sqlmock.NewResult(0, 0))
	err = repo.DeleteByPId(postId)
	assert.Error(t, err)
}

func TestDeleteByUIdPId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLPostRepository(db)

	postId := "post-uuid"
	uId := "user-uuid"

	mock.ExpectExec("^DELETE FROM posts WHERE post_id = \\? AND uuid = \\?$").
		WithArgs(postId, uId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteByUIdPId(uId, postId)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPostsByFilter(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLPostRepository(db)

	rows := sqlmock.NewRows([]string{"post_id", "uuid", "title", "type", "content", "likes", "created_at"}).
		AddRow("post-uuid-1", "user-uuid", "Post 1", "Type A", "Content 1", 10, "2024-09-22T16:52:56Z").
		AddRow("post-uuid-2", "user-uuid", "Post 2", "Type A", "Content 2", 5, "2024-09-22T16:52:56Z")

	filter := "Type A"
	mock.ExpectQuery("SELECT post_id, uuid, title, type, content, likes, created_at FROM posts WHERE type = ?").
		WithArgs(filter).
		WillReturnRows(rows)

	posts, err := repo.GetPostsByFilter(filter)
	assert.NoError(t, err)
	assert.Len(t, posts, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPostsByUId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLPostRepository(db)

	rows := sqlmock.NewRows([]string{"post_id", "uuid", "title", "type", "content", "likes", "created_at"}).
		AddRow("post-uuid-1", "user-uuid", "Post 1", "Type A", "Content 1", 10, time.Now()).
		AddRow("post-uuid-2", "user-uuid", "Post 2", "Type A", "Content 2", 5, time.Now())

	uId := "user-uuid"
	mock.ExpectQuery("SELECT post_id, uuid, title, type, content, likes, created_at FROM posts WHERE uuid = ?").
		WithArgs(uId).
		WillReturnRows(rows)

	posts, err := repo.GetPostsByUId(uId)
	assert.NoError(t, err)
	assert.Len(t, posts, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPostByPId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLPostRepository(db)

	columns := []string{"p.post_id", "p.uuid", "p.title", "p.type", "p.content", "p.likes", "p.created_at", "q.q_id", "q.text", "q.replies"}
	rows := sqlmock.NewRows(columns).
		AddRow("post-uuid", "user-uuid", "Test Post", "Type A", "Content", 10, "2024-09-22T16:52:56Z", "q1", "What is this?", nil)

	mock.ExpectQuery("SELECT p.post_id, p.uuid, p.title, p.type, p.content, p.likes, p.created_at, q.q_id, q.text, q.replies FROM posts p LEFT JOIN questions q ON p.post_id = q.post_id WHERE p.post_id = ?").
		WithArgs("post-uuid").
		WillReturnRows(rows)

	post, err := repo.GetPostByPId("post-uuid")
	assert.NoError(t, err)
	assert.Equal(t, "Test Post", post.Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPostByPId_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	pId := "1"

	mock.ExpectQuery("SELECT p.post_id, p.uuid, p.title, p.type, p.content, p.likes, p.created_at, q.q_id, q.text, q.replies FROM posts p LEFT JOIN questions q ON p.post_id = q.post_id WHERE p.post_id = ?").
		WithArgs(pId).
		WillReturnError(errors.New("query error"))

	repo := repositories.NewMySQLPostRepository(db)
	posts, err := repo.GetPostByPId(pId)

	assert.Error(t, err)
	assert.Nil(t, posts)
	assert.EqualError(t, err, "query error")
}

func TestGetPostByPId_NoPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	pId := "1"

	rows := sqlmock.NewRows([]string{"post_id", "user_id", "title", "type", "content", "likes", "created_at"})

	mock.ExpectQuery("SELECT p.post_id, p.uuid, p.title, p.type, p.content, p.likes, p.created_at, q.q_id, q.text, q.replies FROM posts p LEFT JOIN questions q ON p.post_id = q.post_id WHERE p.post_id = ?").
		WithArgs(pId).
		WillReturnRows(rows)

	repo := repositories.NewMySQLPostRepository(db)
	_, err = repo.GetPostByPId(pId)

	assert.NoError(t, err)
}

func TestUpdateUserPost(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLPostRepository(db)

	postId := "post-uuid"
	uId := "user-uuid"
	title := "Updated Title"
	content := "Updated Content"

	mock.ExpectExec("^UPDATE posts SET title = \\?, content = \\? WHERE post_id = \\? AND uuid = \\?$").
		WithArgs(title, content, postId, uId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateUserPost(postId, uId, title, content)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateLike(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLPostRepository(db)

	postId := "post-uuid"

	mock.ExpectExec("^UPDATE posts SET likes = likes\\+1 WHERE post_id=\\?$").
		WithArgs(postId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateLike(postId)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
