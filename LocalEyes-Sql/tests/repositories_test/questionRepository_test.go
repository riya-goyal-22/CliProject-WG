package repositories_test

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"localEyes/internal/models"
	"localEyes/internal/repositories"
	"testing"
	"time"
)

func TestCreateQuestion(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLQuestionRepository(db)

	question := &models.Question{
		PostId:    "post-1",
		UserId:    "user-1",
		QId:       "q-1",
		Text:      "What is this?",
		Replies:   []string{"Reply 1", "Reply 2"},
		CreatedAt: time.Now(),
	}

	replies, _ := json.Marshal(question.Replies)
	mock.ExpectExec("INSERT INTO questions").
		WithArgs(question.PostId, question.UserId, question.QId, question.Text, replies, question.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(question)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllQuestions(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLQuestionRepository(db)

	rows := sqlmock.NewRows([]string{"q_id", "post_id", "uuid", "text", "replies", "created_at"}).
		AddRow("q-1", "post-1", "user-1", "What is this?", `["Reply 1"]`, time.Now().Format("2006-01-02T15:04:05Z")).
		AddRow("q-2", "post-1", "user-1", "Another question?", `["Reply 2"]`, time.Now().Format("2006-01-02T15:04:05Z"))

	mock.ExpectQuery("SELECT q_id, post_id, uuid, text, replies, created_at FROM questions").
		WillReturnRows(rows)

	questions, err := repo.GetAllQuestions()
	assert.NoError(t, err)
	assert.Len(t, questions, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteByQIdUId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLQuestionRepository(db)

	mock.ExpectExec("DELETE FROM questions WHERE q_id= \\? AND uuid= \\?").
		WithArgs("q-1", "user-1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteByQIdUId("q-1", "user-1")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteByPId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLQuestionRepository(db)

	mock.ExpectExec("DELETE FROM questions WHERE post_id = \\?").
		WithArgs("post-1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteByPId("post-1")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteByQId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLQuestionRepository(db)

	mock.ExpectExec("DELETE FROM questions WHERE q_id = \\?").
		WithArgs("q-1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteByQId("q-1")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetQuestionsByPId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLQuestionRepository(db)

	rows := sqlmock.NewRows([]string{"q_id", "post_id", "uuid", "text", "replies", "created_at"}).
		AddRow("q-1", "post-1", "user-1", "What is this?", `["Reply 1"]`, time.Now().Format("2006-01-02T15:04:05Z"))

	mock.ExpectQuery("SELECT q_id, post_id, uuid, text, replies, created_at FROM questions WHERE post_id = \\?").
		WithArgs("post-1").
		WillReturnRows(rows)

	questions, err := repo.GetQuestionsByPId("post-1")
	assert.NoError(t, err)
	assert.Len(t, questions, 1)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateQuestion(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLQuestionRepository(db)

	// Update the regex to avoid complex escaping
	mock.ExpectExec("^UPDATE questions SET replies= JSON_ARRAY_APPEND\\(replies, '\\$' ,\\?\\) WHERE q_id=\\?$").
		WithArgs("New reply", "q-1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateQuestion("q-1", "New reply")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
