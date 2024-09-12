package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"localEyes/config"
	"localEyes/internal/models"
	"localEyes/utils"
	"time"
)

type MySQLQuestionRepository struct {
	DB *sql.DB
}

func NewMySQLQuestionRepository(Db *sql.DB) *MySQLQuestionRepository {
	return &MySQLQuestionRepository{
		DB: Db,
	}
}

func (r *MySQLQuestionRepository) Create(question *models.Question) error {
	replies, err := json.Marshal(question.Replies)
	columns:=[]string{"post_id","user_id", "text", "replies","created_at"}
	query:=config.InsertQuery(config.QuestionTable,columns)
	//query := "INSERT INTO questions (post_id,user_id, text, replies,created_at) VALUES (?, ?, ?, ?,?)"
	_, err = r.DB.Exec(query, question.PostId, question.UserId, question.Text, replies, question.CreatedAt)
	return err
}

func (r *MySQLQuestionRepository) GetAllQuestions() ([]*models.Question, error) {
	columns:=[]string{"q_id","post_id","user_id","text","replies","created_at"}
	query:=config.SelectQuery(config.QuestionTable,"","",columns)
	//query := "SELECT q_id, post_id, user_id, text, replies, created_at FROM questions"
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

	var questions []*models.Question
	for rows.Next() {
		var question models.Question
		var replies string
		var createdAt string // Use string for raw scan

		// Scan the row into struct fields, using string for replies and createdAt
		if err := rows.Scan(&question.QId, &question.PostId, &question.UserId, &question.Text, &replies, &createdAt); err != nil {
			return nil, err
		}
		if replies != "" {
			if err := json.Unmarshal([]byte(replies), &question.Replies); err != nil {
				return nil, err
			}
		}

		if createdAt != "" {
			parsedTime, err := time.Parse("2006-01-02T15:04:05Z", createdAt)
			if err != nil {
				return nil, err
			}
			question.CreatedAt = parsedTime
		}

		questions = append(questions, &question)
	}
	return questions, nil
}
func (r *MySQLQuestionRepository) DeleteByQIdUId(QId, UId int) error {
	condition1:="q_id"
	condition2:="user_id"
	query:=config.DeleteQuery(config.QuestionTable,condition1,condition2)
	//query := "DELETE FROM questions WHERE q_id = ? AND user_id = ?"
	result, err := r.DB.Exec(query, QId, UId)
	if result != nil {
		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New(config.Red + "No Question exist with this id" + config.Reset)
		}
	}
	return err
}
func (r *MySQLQuestionRepository) DeleteByPId(PId int) error {
	condition1:="post_id"
	query:=config.DeleteQuery(config.QuestionTable,condition1,"")
	//query := "DELETE FROM questions WHERE post_id = ?"
	result, err := r.DB.Exec(query, PId)
	if result != nil {
		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New(config.Red + "No Question exist with this id" + config.Reset)
		}
	}
	return err
}
func (r *MySQLQuestionRepository) DeleteByQId(QId int) error {
	condition1:="q_id"
	query:=config.DeleteQuery(config.QuestionTable,condition1,"")
	//query := "DELETE FROM questions WHERE q_id = ?"
	result, err := r.DB.Exec(query, QId)
	if result != nil {
		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New(config.Red + "No Question exist with this id" + config.Reset)
		}
	}
	return err
}
func (r *MySQLQuestionRepository) GetQuestionsByPId(PId int) ([]*models.Question, error) {
	columns:=[]string{"q_id","post_id","user_id","text","replies","created_at"}
	condition1:="post_id"
	query:=config.SelectQuery(config.QuestionTable,condition1,"",columns)
	//query := "SELECT q_id, post_id,user_id, text, replies ,created_at FROM questions WHERE post_id = ?"
	rows, err := r.DB.Query(query, PId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*models.Question
	for rows.Next() {
		var question models.Question
		var replies []byte
		if err := rows.Scan(&question.QId, &question.PostId, &question.UserId, &question.Text, &replies, &question.CreatedAt); err != nil {
			return nil, err
		}
		err = json.Unmarshal(replies, &question.Replies)
		if err != nil {
			return nil, err
		}
		questions = append(questions, &question)
	}
	return questions, nil
}
func (r *MySQLQuestionRepository) UpdateQuestion(QId int, answer string) error {
	columns:=[]string{"replies"}
	condition1:="q_id"
	query:=config.UpdateQuery(config.QuestionTable,condition1,"",columns)
	//query := "UPDATE questions SET replies= JSON_ARRAY_APPEND(replies, '$' ,?) WHERE q_id = ?"
	result, err := r.DB.Exec(query, fmt.Sprintf("JSON_ARRAY_APPEND(replies, '$' ,%s)",answer), QId)
	if result != nil {
		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New(config.Red + "No Question exist with this id" + config.Reset)
		}
	}
	return err
}
