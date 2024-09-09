package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"localEyes/constants"
	"localEyes/internal/models"
)

type MySQLUserRepository struct {
	DB *sql.DB
}

func NewMySQLUserRepository(Db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{
		DB: Db,
	}
}

func (r *MySQLUserRepository) Create(user *models.User) error {
	notification, err := json.Marshal(user.Notification)
	query := "INSERT INTO users (username, password, is_active, city, dwelling_age, tag, notification) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err = r.DB.Exec(query, user.Username, user.Password, user.IsActive, user.City, user.DwellingAge, user.Tag, notification)
	return err
}

func (r *MySQLUserRepository) FindByUId(UId int) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, password, is_active, city, dwelling_age, tag, notification FROM users WHERE id = ?"
	var notification []byte
	err := r.DB.QueryRow(query, UId).Scan(&user.UId, &user.Username, &user.Password, &user.IsActive, &user.City, &user.DwellingAge, &user.Tag, &notification)
	json.Unmarshal(notification, &user.Notification)
	return &user, err
}

func (r *MySQLUserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, password, is_active, city, dwelling_age, tag, notification FROM users WHERE username = ?"
	var notification []byte
	err := r.DB.QueryRow(query, username).Scan(&user.UId, &user.Username, &user.Password, &user.IsActive, &user.City, &user.DwellingAge, &user.Tag, &notification)
	json.Unmarshal(notification, &user.Notification)
	return &user, err
}

func (r *MySQLUserRepository) FindByUsernamePassword(username, password string) (*models.User, error) {
	var user models.User
	query := "SELECT id, username, password, is_active, city, dwelling_age, tag, notification FROM users WHERE username = ? AND password = ?"
	var notification []byte
	err := r.DB.QueryRow(query, username, password).Scan(&user.UId, &user.Username, &user.Password, &user.IsActive, &user.City, &user.DwellingAge, &user.Tag, &notification)
	json.Unmarshal(notification, &user.Notification)
	return &user, err
}

func (r *MySQLUserRepository) FindAdminByUsernamePassword(username, password string) (*models.Admin, error) {
	var admin models.Admin
	query := "SELECT id, username, password FROM users WHERE username = ? AND password = ?"
	err := r.DB.QueryRow(query, username, password).Scan(&admin.User.UId, &admin.User.Username, &admin.User.Password)
	return &admin, err
}

func (r *MySQLUserRepository) GetAllUsers() ([]*models.User, error) {
	query := "SELECT id, username, password, is_active, city, dwelling_age, tag, notification FROM users"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		var notification []byte
		if err := rows.Scan(&user.UId, &user.Username, &user.Password, &user.IsActive, &user.City, &user.DwellingAge, &user.Tag, &notification); err != nil {
			return nil, err
		}
		json.Unmarshal(notification, &user.Notification)
		users = append(users, &user)
	}
	return users, nil
}

func (r *MySQLUserRepository) DeleteByUId(UId int) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := r.DB.Exec(query, UId)
	if result != nil {
		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New(constants.Red + "No user exist with this id" + constants.Reset)
		}
	}
	return err
}

func (r *MySQLUserRepository) UpdateActiveStatus(UId int, status bool) error {
	query := "UPDATE users SET is_active = ? WHERE id = ?"
	result, err := r.DB.Exec(query, status, UId)
	if result != nil {
		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New(constants.Red + "No user exist with this id" + constants.Reset)
		}
	}
	return err
}

func (r *MySQLUserRepository) PushNotification(UId int, title string) error {
	//query := "UPDATE users SET notification = CONCAT(COALESCE(notification, ''), ?) WHERE id != ?"
	query := "UPDATE users SET notification= JSON_ARRAY_APPEND(notification, '$' ,?) WHERE id != ?"
	notification := "New post: " + title + "\n"
	_, err := r.DB.Exec(query, notification, UId)
	return err
}

func (r *MySQLUserRepository) ClearNotification(UId int) error {
	query := "UPDATE users SET notification = '[]' WHERE id = ?"
	_, err := r.DB.Exec(query, UId)
	return err
}
