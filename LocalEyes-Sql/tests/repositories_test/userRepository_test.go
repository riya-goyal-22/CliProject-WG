package repositories_test

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"localEyes/internal/models"
	"localEyes/internal/repositories"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLUserRepository(db)

	user := &models.User{
		Username:     "testuser",
		Password:     "password123",
		IsActive:     true,
		City:         "Test City",
		DwellingAge:  5,
		Tag:          "test",
		Notification: []string{"Notification 1"},
	}

	notification, _ := json.Marshal(user.Notification)
	mock.ExpectExec("^INSERT INTO users \\(username, password, is_active, city, dwelling_age, tag, notification\\) VALUES \\(\\?, \\?, \\?, \\?, \\?, \\?, \\?\\)$").
		WithArgs(user.Username, user.Password, user.IsActive, user.City, user.DwellingAge, user.Tag, notification).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByUId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLUserRepository(db)

	notification := []byte(`["Notification 1"]`)
	mock.ExpectQuery("^SELECT uuid, username, password, is_active, city, dwelling_age, tag, notification FROM users WHERE uuid = \\?$").
		WithArgs("user-uuid").
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "username", "password", "is_active", "city", "dwelling_age", "tag", "notification"}).
			AddRow("user-uuid", "testuser", "password123", true, "Test City", 5, "test", notification))

	user, err := repo.FindByUId("user-uuid")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLUserRepository(db)

	notification := []byte(`["Notification 1"]`)
	mock.ExpectQuery("^SELECT uuid, username, password, is_active, city, dwelling_age, tag, notification FROM users WHERE username = \\?$").
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "username", "password", "is_active", "city", "dwelling_age", "tag", "notification"}).
			AddRow("user-uuid", "testuser", "password123", true, "Test City", 5, "test", notification))

	user, err := repo.FindByUsername("testuser")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user-uuid", user.UId)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByUsernamePassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLUserRepository(db)

	notification := []byte(`["Notification 1"]`)
	mock.ExpectQuery("^SELECT uuid, username, password, is_active, city, dwelling_age, tag, notification FROM users WHERE username = \\? AND password = \\?$").
		WithArgs("testuser", "password123").
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "username", "password", "is_active", "city", "dwelling_age", "tag", "notification"}).
			AddRow("user-uuid", "testuser", "password123", true, "Test City", 5, "test", notification))

	user, err := repo.FindByUsernamePassword("testuser", "password123")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "user-uuid", user.UId)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLUserRepository(db)

	notification := []byte(`["Notification 1"]`)
	mock.ExpectQuery("^SELECT uuid, username, password, is_active, city, dwelling_age, tag, notification FROM users$").
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "username", "password", "is_active", "city", "dwelling_age", "tag", "notification"}).
			AddRow("user-uuid", "testuser", "password123", true, "Test City", 5, "test", notification))

	users, err := repo.GetAllUsers()
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "testuser", users[0].Username)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteByUId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLUserRepository(db)

	mock.ExpectExec("^DELETE FROM users WHERE uuid= \\? AND username!= \\?$").
		WithArgs("user-uuid").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteByUId("user-uuid")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateActiveStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLUserRepository(db)

	mock.ExpectExec("^UPDATE users SET is_active = \\? WHERE uuid = \\?$").
		WithArgs(true, "user-uuid").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.UpdateActiveStatus("user-uuid", true)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPushNotification(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLUserRepository(db)

	mock.ExpectExec("^UPDATE users SET notification= JSON_ARRAY_APPEND\\(notification, '\\$' ,\\?\\) WHERE uuid!=\\? AND username!=\\?$").
		WithArgs("New post: Test Title\n", "user-uuid", "admin").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.PushNotification("user-uuid", "Test Title")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestClearNotification(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.NewMySQLUserRepository(db)

	mock.ExpectExec("^UPDATE users SET notification = \\?  WHERE uuid = \\?$").
		WithArgs("[]", "user-uuid").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.ClearNotification("user-uuid")
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
