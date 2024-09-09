package repositories_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"localEyes/internal/models"
	"localEyes/internal/repositories"
)

func TestMySQLUserRepository_Create(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create a repository with the mock DB
	repo := repositories.NewMySQLUserRepository(db)

	// Create a mock user
	user := &models.User{
		Username:    "test_user",
		Password:    "password123",
		IsActive:    true,
		City:        "New York",
		DwellingAge: 5,
		Tag:         "tag1",
		Notification: []string{
			"Welcome to LocalEyes",
		},
	}

	// Marshal the notification to JSON
	notification, _ := json.Marshal(user.Notification)

	// Expect the insert query
	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Username, user.Password, user.IsActive, user.City, user.DwellingAge, user.Tag, notification).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the Create method
	err = repo.Create(user)

	// Assert no error was returned
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_FindByUId(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create a repository with the mock DB
	repo := repositories.NewMySQLUserRepository(db)

	// Mock the user that will be returned
	user := &models.User{
		UId:         1,
		Username:    "test_user",
		Password:    "password123",
		IsActive:    true,
		City:        "New York",
		DwellingAge: 5,
		Tag:         "tag1",
		Notification: []string{
			"Welcome to LocalEyes",
		},
	}
	notification, _ := json.Marshal(user.Notification)

	// Expect the select query
	mock.ExpectQuery("SELECT id, username, password, is_active, city, dwelling_age, tag, notification FROM users WHERE id = ?").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "is_active", "city", "dwelling_age", "tag", "notification"}).
			AddRow(user.UId, user.Username, user.Password, user.IsActive, user.City, user.DwellingAge, user.Tag, notification))

	// Call the FindByUId method
	result, err := repo.FindByUId(1)

	// Assert no error was returned
	assert.NoError(t, err)

	// Validate the returned user
	assert.Equal(t, user.UId, result.UId)
	assert.Equal(t, user.Username, result.Username)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_FindByUsername(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create a repository with the mock DB
	repo := repositories.NewMySQLUserRepository(db)

	// Mock the user that will be returned
	user := &models.User{
		UId:         1,
		Username:    "test_user",
		Password:    "password123",
		IsActive:    true,
		City:        "New York",
		DwellingAge: 5,
		Tag:         "tag1",
		Notification: []string{
			"Welcome to LocalEyes",
		},
	}
	notification, _ := json.Marshal(user.Notification)

	// Expect the select query
	mock.ExpectQuery("SELECT id, username, password, is_active, city, dwelling_age, tag, notification FROM users WHERE username = ?").
		WithArgs("test_user").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "is_active", "city", "dwelling_age", "tag", "notification"}).
			AddRow(user.UId, user.Username, user.Password, user.IsActive, user.City, user.DwellingAge, user.Tag, notification))

	// Call the FindByUsername method
	result, err := repo.FindByUsername("test_user")

	// Assert no error was returned
	assert.NoError(t, err)

	// Validate the returned user
	assert.Equal(t, user.Username, result.Username)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_UpdateActiveStatus(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database connection: %v", err)
	}
	defer db.Close()

	// Create a repository with the mock DB
	repo := repositories.NewMySQLUserRepository(db)

	// Define the expected query and arguments
	expectedQuery := "UPDATE users SET is_active = \\? WHERE id = \\?"

	// Expect the update query
	mock.ExpectExec(expectedQuery).
		WithArgs(true, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the UpdateActiveStatus method
	err = repo.UpdateActiveStatus(1, true)

	// Assert no error was returned
	assert.NoError(t, err)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expected SQL queries were not executed: %v", err)
	}
}

func TestMySQLUserRepository_PushNotification(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create a repository with the mock DB
	repo := repositories.NewMySQLUserRepository(db)

	// Expect the update query
	mock.ExpectExec("UPDATE users SET notification= JSON_ARRAY_APPEND").
		WithArgs("New post: Test Post\n", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the PushNotification method
	err = repo.PushNotification(1, "Test Post")

	// Assert no error was returned
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMySQLUserRepository_ClearNotification(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create a repository with the mock DB
	repo := repositories.NewMySQLUserRepository(db)

	// Expect the update query to clear notifications (set notification to an empty array)
	mock.ExpectExec("UPDATE users SET notification = '\\[\\]' WHERE id = \\?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call the ClearNotification method
	err = repo.ClearNotification(1)

	// Assert no error was returned
	assert.NoError(t, err)

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByUsernamePassword_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock instance: %v", err)
	}
	defer db.Close()

	repo := repositories.NewMySQLUserRepository(db)

	// Define the expected results
	notification := json.RawMessage(`{"email":"example@example.com"}`)
	rows := sqlmock.NewRows([]string{"id", "username", "password", "is_active", "city", "dwelling_age", "tag", "notification"}).
		AddRow(1, "testuser", "testpass", true, "New York", 5, "Admin", notification)

	// Set the expectation for the query
	mock.ExpectQuery("SELECT id, username, password, is_active, city, dwelling_age, tag, notification FROM users WHERE username = \\? AND password = \\?").
		WithArgs("testuser", "testpass").
		WillReturnRows(rows)

	// Call the method to be tested
	user, err := repo.FindByUsernamePassword("testuser", "testpass")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.UId)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "testpass", user.Password)
	assert.True(t, user.IsActive)
	assert.Equal(t, "New York", user.City)
	assert.Equal(t, 5, user.DwellingAge)
	assert.Equal(t, "Admin", user.Tag)
	//assert.Equal(t, `{"email":"example@example.com"}`,(user.Notification))
}

func TestFindByUsernamePassword_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock instance: %v", err)
	}
	defer db.Close()

	repo := repositories.NewMySQLUserRepository(db)

	// Set the expectation for the query
	mock.ExpectQuery("SELECT id, username, password, is_active, city, dwelling_age, tag, notification FROM users WHERE username = \\? AND password = \\?").
		WithArgs("nonexistentuser", "wrongpass").
		WillReturnError(sql.ErrNoRows)

	// Call the method to be tested
	_, err = repo.FindByUsernamePassword("nonexistentuser", "wrongpass")
	assert.Error(t, err)
	//assert.Nil(t, user)
}

func TestFindByUsernamePassword_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock instance: %v", err)
	}
	defer db.Close()

	repo := repositories.NewMySQLUserRepository(db)

	// Set the expectation for the query
	mock.ExpectQuery("SELECT id, username, password, is_active, city, dwelling_age, tag, notification FROM users WHERE username = \\? AND password = \\?").
		WithArgs("testuser", "testpass").
		WillReturnError(sql.ErrConnDone) // Simulate a connection error

	// Call the method to be tested
	_, err = repo.FindByUsernamePassword("testuser", "testpass")
	assert.Error(t, err)
	//assert.Nil(t, user)
}

func TestFindAdminByUsernamePassword_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	username := "admin"
	password := "password123"

	rows := sqlmock.NewRows([]string{"id", "username", "password"}).
		AddRow(1, username, password)

	mock.ExpectQuery("^SELECT id, username, password FROM users WHERE username = \\? AND password = \\?$").
		WithArgs(username, password).
		WillReturnRows(rows)

	repo := repositories.NewMySQLUserRepository(db)
	admin, err := repo.FindAdminByUsernamePassword(username, password)

	assert.NoError(t, err)
	assert.NotNil(t, admin)
	assert.Equal(t, 1, admin.User.UId)
	assert.Equal(t, username, admin.User.Username)
	assert.Equal(t, password, admin.User.Password)
}

func TestFindAdminByUsernamePassword_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	username := "admin"
	password := "wrongpassword"

	mock.ExpectQuery("^SELECT id, username, password FROM users WHERE username = \\? AND password = \\?$").
		WithArgs(username, password).
		WillReturnError(sql.ErrNoRows)

	repo := repositories.NewMySQLUserRepository(db)
	_, err = repo.FindAdminByUsernamePassword(username, password)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestGetAllUsers_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "username", "password", "is_active", "city", "dwelling_age", "tag", "notification"}).
		AddRow(1, "user1", "pass1", true, "CityA", 5, "Tag1", `["notif1"]`).
		AddRow(2, "user2", "pass2", false, "CityB", 10, "Tag2", `["notif2"]`)

	mock.ExpectQuery("^SELECT id, username, password, is_active, city, dwelling_age, tag, notification FROM users$").
		WillReturnRows(rows)

	repo := repositories.NewMySQLUserRepository(db)
	users, err := repo.GetAllUsers()

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, 1, users[0].UId)
	assert.Equal(t, "user1", users[0].Username)
	assert.Equal(t, "CityA", users[0].City)
}

func TestGetAllUsers_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery("^SELECT id, username, password, is_active, city, dwelling_age, tag, notification FROM users$").
		WillReturnError(errors.New("query error"))

	repo := repositories.NewMySQLUserRepository(db)
	users, err := repo.GetAllUsers()

	assert.Error(t, err)
	assert.Nil(t, users)
	assert.EqualError(t, err, "query error")
}

func TestDeleteByUId_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	UId := 1
	mock.ExpectExec("^DELETE FROM users WHERE id = \\?$").
		WithArgs(UId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := repositories.NewMySQLUserRepository(db)
	err = repo.DeleteByUId(UId)

	assert.NoError(t, err)
}
