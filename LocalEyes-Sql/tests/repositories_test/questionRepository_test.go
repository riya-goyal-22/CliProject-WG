package repositories_test

import (
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"localEyes/constants"
	"localEyes/internal/models"
	"localEyes/internal/repositories"
	"testing"
	"time"
)

func TestMySQLQuestionRepository_Create(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create a repository with the mock DB
	repo := repositories.NewMySQLQuestionRepository(db)

	// Create a mock question
	question := &models.Question{
		PostId:    1,
		UserId:    1,
		Text:      "What's your favorite food?",
		Replies:   []string{"Pizza", "Burger"},
		CreatedAt: time.Now(),
	}

	// Marshal the replies to JSON
	replies, _ := json.Marshal(question.Replies)

	// Expect the insert query
	mock.ExpectExec("INSERT INTO questions").
		WithArgs(question.PostId, question.UserId, question.Text, replies, question.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the Create method
	err = repo.Create(question)

	// Assert no error was returned
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLQuestionRepository_GetAllQuestions(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database connection: %v", err)
	}
	defer db.Close()

	// Create a repository with the mock DB
	repo := repositories.NewMySQLQuestionRepository(db)

	// Mock the rows that will be returned by the query
	rows := sqlmock.NewRows([]string{"q_id", "post_id", "user_id", "text", "replies", "created_at"}).
		AddRow(1, 1, 1, "What is your favorite food?", `["Pizza", "Burger"]`, "2006-01-02T15:04:05Z")

	mock.ExpectQuery("^SELECT q_id, post_id, user_id, text, replies, created_at FROM questions$").
		WillReturnRows(rows)

	// Call the GetAllQuestions method
	questions, err := repo.GetAllQuestions()

	// Assert no error was returned
	assert.NoError(t, err)

	// Check if we got the correct number of questions
	assert.Len(t, questions, 1)

	// Validate the returned question
	if len(questions) > 0 {
		assert.Equal(t, 1, questions[0].QId)
		assert.Equal(t, "What is your favorite food?", questions[0].Text)
		assert.ElementsMatch(t, []string{"Pizza", "Burger"}, questions[0].Replies)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expected SQL queries were not executed: %v", err)
	}
}

//func TestMySQLQuestionRepository_DeleteByQIdUId(t *testing.T) {
//	// Create a mock database connection
//	db, mock, err := sqlmock.New()
//	assert.NoError(t, err)
//	defer db.Close()
//
//	// Create a repository with the mock DB
//	repo := repositories.NewMySQLQuestionRepository(db)
//
//	// Expect the delete query
//	mock.ExpectExec("^DELETE FROM questions WHERE q_id = \\? AND user_id = \\?$").
//		WithArgs(1, 1).
//		WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate successful deletion
//
//	// Call the DeleteByQIdUId method
//	err = repo.DeleteByQIdUId(1, 1)
//
//	// Assert no error was returned
//	assert.NoError(t, err)
//
//	// Ensure all expectations were met
//	assert.NoError(t, mock.ExpectationsWereMet())
//}

func TestDeleteByQIdUId(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to initialize sqlmock: %v", err)
	}
	defer db.Close()

	// Create an instance of MySQLQuestionRepository
	repo := &repositories.MySQLQuestionRepository{DB: db}

	// Define test cases
	tests := []struct {
		name            string
		qId, uId        int
		expectedError   error
		affectedRows    int64
		mockExpectation func()
	}{
		{
			name:          "successful delete",
			qId:           1,
			uId:           1,
			expectedError: nil,
			affectedRows:  1,
			mockExpectation: func() {
				mock.ExpectExec("DELETE FROM questions WHERE q_id = \\? AND user_id = \\?").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name:          "no rows affected",
			qId:           2,
			uId:           2,
			expectedError: errors.New(constants.Red + "No Question exist with this id" + constants.Reset),
			affectedRows:  0,
			mockExpectation: func() {
				mock.ExpectExec("DELETE FROM questions WHERE q_id = \\? AND user_id = \\?").
					WithArgs(2, 2).
					WillReturnResult(sqlmock.NewResult(1, 0))
			},
		},
		{
			name:          "error executing query",
			qId:           3,
			uId:           3,
			expectedError: errors.New("query execution error"),
			affectedRows:  0,
			mockExpectation: func() {
				mock.ExpectExec("DELETE FROM questions WHERE q_id = \\? AND user_id = \\?").
					WithArgs(3, 3).
					WillReturnError(errors.New("query execution error"))
			},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the mock expectations
			tt.mockExpectation()

			// Call the method under test
			err := repo.DeleteByQIdUId(tt.qId, tt.uId)

			// Assert the results
			assert.Equal(t, tt.expectedError, err)
			if tt.expectedError == nil {
				mock.ExpectationsWereMet()
			}
		})
	}
}

func TestUpdateQuestion_Success(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	QId := 1
	answer := "New answer"

	// Set up mock expectations
	mock.ExpectExec("UPDATE questions SET replies= JSON_ARRAY_APPEND").
		WithArgs(answer, QId).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate a successful update with 1 affected row

	// Initialize repository
	repo := repositories.NewMySQLQuestionRepository(db)

	// Call the method
	err = repo.UpdateQuestion(QId, answer)

	// Assert
	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}

func TestUpdateQuestion_NoRowsAffected(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	QId := 1
	answer := "New answer"

	// Set up mock expectations
	mock.ExpectExec("UPDATE questions SET replies= JSON_ARRAY_APPEND").
		WithArgs(answer, QId).
		WillReturnResult(sqlmock.NewResult(1, 0)) // Simulate no rows affected

	// Initialize repositories
	repo := repositories.NewMySQLQuestionRepository(db)

	// Call the method
	err = repo.UpdateQuestion(QId, answer)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, constants.Red+"No Question exist with this id"+constants.Reset, err.Error())
	mock.ExpectationsWereMet()
}

func TestUpdateQuestion_QueryError(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	QId := 1
	answer := "New answer"

	// Set up mock expectations
	mock.ExpectExec("UPDATE questions SET replies= JSON_ARRAY_APPEND").
		WithArgs(answer, QId).
		WillReturnError(errors.New("query error")) // Simulate a query error

	// Initialize repositories
	repo := repositories.NewMySQLQuestionRepository(db)

	// Call the method
	err = repo.UpdateQuestion(QId, answer)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "query error", err.Error())
	mock.ExpectationsWereMet()
}

func TestDeleteByPId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLQuestionRepository(db)

	// Set up mock expectations
	PId := 1
	mock.ExpectExec("^DELETE FROM questions WHERE post_id = \\?$").
		WithArgs(PId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the method
	err = repo.DeleteByPId(PId)

	// Assert results
	assert.NoError(t, err)
}

func TestDeleteByPId_NoRowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLQuestionRepository(db)

	// Set up mock expectations
	PId := 1
	mock.ExpectExec("DELETE FROM questions WHERE post_id = ?").
		WithArgs(PId).
		WillReturnResult(sqlmock.NewResult(1, 0))

	// Call the method
	err = repo.DeleteByPId(PId)

	// Assert results
	assert.EqualError(t, err, constants.Red+"No Question exist with this id"+constants.Reset)
}

func TestDeleteByPId_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLQuestionRepository(db)

	// Set up mock expectations
	PId := 1
	mock.ExpectExec("DELETE FROM questions WHERE post_id = ?").
		WithArgs(PId).
		WillReturnError(errors.New("some error"))

	// Call the method
	err = repo.DeleteByPId(PId)

	// Assert results
	assert.EqualError(t, err, "some error")
}

func TestDeleteByQId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLQuestionRepository(db)

	// Set up mock expectations
	QId := 1
	mock.ExpectExec("^DELETE FROM questions WHERE q_id = \\?$").
		WithArgs(QId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the method
	err = repo.DeleteByQId(QId)

	// Assert results
	assert.NoError(t, err)
}

func TestDeleteByQId_NoRowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLQuestionRepository(db)

	// Set up mock expectations
	QId := 1
	mock.ExpectExec("^DELETE FROM questions WHERE q_id = \\?$").
		WithArgs(QId).
		WillReturnResult(sqlmock.NewResult(1, 0))

	// Call the method
	err = repo.DeleteByQId(QId)

	// Assert results
	assert.EqualError(t, err, constants.Red+"No Question exist with this id"+constants.Reset)
}

func TestDeleteByQId_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLQuestionRepository(db)

	// Set up mock expectations
	QId := 1
	mock.ExpectExec("^DELETE FROM questions WHERE q_id = \\?$").
		WithArgs(QId).
		WillReturnError(errors.New("some error"))

	// Call the method
	err = repo.DeleteByQId(QId)

	// Assert results
	assert.EqualError(t, err, "some error")
}

func TestGetQuestionsByPId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLQuestionRepository(db)

	PId := 1
	createdAt := time.Now()
	replies := `[{"text":"First reply"}]`

	// Set up mock expectations
	rows := sqlmock.NewRows([]string{"q_id", "post_id", "user_id", "text", "replies", "created_at"}).
		AddRow(1, PId, 1, "Test question", []byte(replies), createdAt)
	mock.ExpectQuery("^SELECT q_id, post_id,user_id, text, replies ,created_at FROM questions WHERE post_id = \\?$").
		WithArgs(PId).
		WillReturnRows(rows)

	// Call the method
	result, err := repo.GetQuestionsByPId(PId)

	// Assert results
	assert.NoError(t, err)
	//assert.Len(t, result, 1)
	assert.Equal(t, "Test question", result[0].Text)
	assert.Equal(t, createdAt, result[0].CreatedAt)
	assert.Equal(t, 1, result[0].QId)

}

func TestGetQuestionsByPId_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLQuestionRepository(db)

	PId := 1

	// Set up mock expectations
	mock.ExpectQuery("^SELECT q_id, post_id,user_id, text, replies ,created_at FROM questions WHERE post_id = \\?$").
		WithArgs(PId).
		WillReturnError(errors.New("some error"))

	// Call the method
	_, err = repo.GetQuestionsByPId(PId)

	// Assert results
	assert.EqualError(t, err, "some error")
}
