package repositories

import (
	"database/sql"
	"encoding/json"
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

func (r *MySQLUserRepository) FindByUId(uId string) (*models.User, error) {
	var user models.User
	columns := []string{"uuid", "username", "password", "is_active", "city", "dwelling_age", "tag", "notification"}
	condition := "uuid"
	query := config.SelectQuery(config.UserTable, condition, "", columns)
	//query := "SELECT uuid, username, password, is_active, city, dwelling_age, tag, notification FROM users WHERE uuid = ?"
	var notification []byte
	err := r.DB.QueryRow(query, uId).Scan(&user.UId, &user.Username, &user.Password, &user.IsActive, &user.City, &user.DwellingAge, &user.Tag, &notification)
	err = json.Unmarshal(notification, &user.Notification)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *MySQLUserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	columns := []string{"uuid", "username", "password", "is_active", "city", "dwelling_age", "tag", "notification"}
	condition1 := "username"
	query := config.SelectQuery(config.UserTable, condition1, "", columns)
	//query := "SELECT uuid, username, password, is_active, city, dwelling_age, tag, notification FROM users WHERE username = ?"
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
	columns := []string{"uuid", "username", "password", "is_active", "city", "dwelling_age", "tag", "notification"}
	condition1 := "username"
	condition2 := "password"
	query := config.SelectQuery(config.UserTable, condition1, condition2, columns)
	//query := "SELECT uuid, username, password, is_active, city, dwelling_age, tag, notification FROM users WHERE username = ? AND password = ?"
	var notification []byte
	err := r.DB.QueryRow(query, username, password).Scan(&user.UId, &user.Username, &user.Password, &user.IsActive, &user.City, &user.DwellingAge, &user.Tag, &notification)
	err = json.Unmarshal(notification, &user.Notification)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *MySQLUserRepository) GetAllUsers() ([]*models.User, error) {
	columns := []string{"uuid", "username", "password", "is_active", "city", "dwelling_age", "tag", "notification"}
	query := config.SelectQuery(config.UserTable, "", "", columns)
	//query := "SELECT uuid, username, password, is_active, city, dwelling_age, tag, notification FROM users"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			utils.Logger.Error("ERROR: Error closing rows:" + err.Error())
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

func (r *MySQLUserRepository) DeleteByUId(uId string) error {
	condition1 := "uuid"
	query := config.DeleteQuery(config.UserTable, condition1, "")
	//query := "DELETE FROM users WHERE uuid = ?"
	result, err := r.DB.Exec(query, uId)
	if result != nil {
		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return utils.NoUser
		}
	}
	return err
}

func (r *MySQLUserRepository) UpdateActiveStatus(uId string, status bool) error {
	condition1 := "uuid"
	columns := []string{"is_active"}
	query := config.UpdateQuery(config.UserTable, condition1, "", columns)
	//query := "UPDATE users SET is_active = ? WHERE uuid = ?"
	result, err := r.DB.Exec(query, status, uId)
	if result != nil {
		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return utils.NoUser
		}
	}
	return err
}

func (r *MySQLUserRepository) PushNotification(uId, title string) error {
	columns := "notification= JSON_ARRAY_APPEND(notification, '$' ,?)"
	condition1 := "uuid!"
	condition2 := "username!"
	query := config.UpdateQueryWithValue(config.UserTable, condition1, condition2, columns)
	//query := "UPDATE users SET notification= JSON_ARRAY_APPEND(notification, '$' ,?) WHERE uuid != ?"
	notification := "New post: " + title + "\n"
	_, err := r.DB.Exec(query, notification, uId, "admin")
	return err
}

func (r *MySQLUserRepository) ClearNotification(uId string) error {
	columns := []string{"notification"}
	condition1 := "uuid"
	query := config.UpdateQuery(config.UserTable, condition1, "", columns)
	//query := "UPDATE users SET notification =?  WHERE uuid = ?"
	_, err := r.DB.Exec(query, "[]", uId)
	return err
}
