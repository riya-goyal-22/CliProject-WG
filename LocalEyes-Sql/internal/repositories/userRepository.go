package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"localEyes/config"
	"localEyes/internal/models"
	"localEyes/utils"
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
	columns := []string{"username", "password", "is_active", "city", "dwelling_age", "tag", "notification"}
	query := config.InsertQuery(config.UserTable, columns)
	//query := "INSERT INTO users (username, password, is_active, city, dwelling_age, tag, notification) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err = r.DB.Exec(query, user.Username, user.Password, user.IsActive, user.City, user.DwellingAge, user.Tag, notification)
	return err
}

func (r *MySQLUserRepository) FindByUId(UId int) (*models.User, error) {
	var user models.User
	columns := []string{"id", "username", "password", "is_active", "city", "dwelling_age", "tag", "notification"}
	condition := "id"
	query := config.SelectQuery(config.UserTable, condition, "", columns)
	//query := "SELECT id, username, password, is_active, city, dwelling_age, tag, notification FROM users WHERE id = ?"
	var notification []byte
	err := r.DB.QueryRow(query, UId).Scan(&user.UId, &user.Username, &user.Password, &user.IsActive, &user.City, &user.DwellingAge, &user.Tag, &notification)
	err = json.Unmarshal(notification, &user.Notification)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *MySQLUserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	columns := []string{"id", "username", "password", "is_active", "city", "dwelling_age", "tag", "notification"}
	condition1 := "username"
	query := config.SelectQuery(config.UserTable, condition1, "", columns)
	//query := "SELECT id, username, password, is_active, city, dwelling_age, tag, notification FROM users WHERE username = ?"
	var notification []byte
	err := r.DB.QueryRow(query, username).Scan(&user.UId, &user.Username, &user.Password, &user.IsActive, &user.City, &user.DwellingAge, &user.Tag, &notification)
	err = json.Unmarshal(notification, &user.Notification)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *MySQLUserRepository) FindByUsernamePassword(username, password string) (*models.User, error) {
	var user models.User
	columns := []string{"id", "username", "password", "is_active", "city", "dwelling_age", "tag", "notification"}
	condition1 := "username"
	condition2 := "password"
	query := config.SelectQuery(config.UserTable, condition1, condition2, columns)
	//query := "SELECT id, username, password, is_active, city, dwelling_age, tag, notification FROM users WHERE username = ? AND password = ?"
	var notification []byte
	err := r.DB.QueryRow(query, username, password).Scan(&user.UId, &user.Username, &user.Password, &user.IsActive, &user.City, &user.DwellingAge, &user.Tag, &notification)
	err = json.Unmarshal(notification, &user.Notification)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *MySQLUserRepository) FindAdminByUsernamePassword(username, password string) (*models.Admin, error) {
	var admin models.Admin
	columns := []string{"id", "username", "password"}
	condition1 := "username"
	condition2 := "password"
	query := config.SelectQuery(config.UserTable, condition1, condition2, columns)
	//query := "SELECT id, username, password FROM users WHERE username = ? AND password = ?"
	row := r.DB.QueryRow(query, username, password)
	if row != nil {
		err := row.Scan(&admin.User.UId, &admin.User.Username, &admin.User.Password)
		return &admin, err
	}
	return nil, errors.New("not found")
}

func (r *MySQLUserRepository) GetAllUsers() ([]*models.User, error) {
	columns := []string{"id", "username", "password", "is_active", "city", "dwelling_age", "tag", "notification"}
	query := config.SelectQuery(config.UserTable, "", "", columns)
	//query := "SELECT id, username, password, is_active, city, dwelling_age, tag, notification FROM users"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			utils.Logger.Println("ERROR: Error closing rows:", err)
		}
	}(rows)

	var users []*models.User
	for rows.Next() {
		var user models.User
		var notification []byte
		if err := rows.Scan(&user.UId, &user.Username, &user.Password, &user.IsActive, &user.City, &user.DwellingAge, &user.Tag, &notification); err != nil {
			return nil, err
		}
		err = json.Unmarshal(notification, &user.Notification)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *MySQLUserRepository) DeleteByUId(UId int) error {
	condition1 := "id"
	query := config.DeleteQuery(config.UserTable, condition1, "")
	//query := "DELETE FROM users WHERE id = ?"
	result, err := r.DB.Exec(query, UId)
	if result != nil {
		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New(config.Red + "No user exist with this id" + config.Reset)
		}
	}
	return err
}

func (r *MySQLUserRepository) UpdateActiveStatus(UId int, status bool) error {
	condition1 := "id"
	columns := []string{"is_active"}
	query := config.UpdateQuery(config.UserTable, condition1, "", columns)
	//query := "UPDATE users SET is_active = ? WHERE id = ?"
	result, err := r.DB.Exec(query, status, UId)
	if result != nil {
		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New(config.Red + "No inActive user exist with this id" + config.Reset)
		}
	}
	return err
}

func (r *MySQLUserRepository) PushNotification(UId int, title string) error {
	columns := "notification= JSON_ARRAY_APPEND(notification, '$' ,?)"
	condition1 := "id!=?"
	condition2 := "username!=?"
	query := config.UpdateQueryWithValue(config.UserTable, condition1, condition2, columns)
	//query := "UPDATE users SET notification= JSON_ARRAY_APPEND(notification, '$' ,?) WHERE id != ?"
	notification := "New post: " + title + "\n"
	_, err := r.DB.Exec(query, notification, UId, "admin")
	return err
}

func (r *MySQLUserRepository) ClearNotification(UId int) error {
	columns := []string{"notification"}
	condition1 := "id"
	query := config.UpdateQuery(config.UserTable, condition1, "", columns)
	//query := "UPDATE users SET notification =?  WHERE id = ?"
	_, err := r.DB.Exec(query, "[]", UId)
	return err
}
