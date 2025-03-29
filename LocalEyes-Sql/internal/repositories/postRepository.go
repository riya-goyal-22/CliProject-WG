package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"localEyes/config"
	"localEyes/internal/models"
	"localEyes/utils"
	"strings"
	"time"
)

type MySQLPostRepository struct {
	DB *sql.DB
}

func NewMySQLPostRepository(Db *sql.DB) *MySQLPostRepository {
	return &MySQLPostRepository{
		DB: Db,
	}
}

func (r *MySQLPostRepository) Create(post *models.Post) error {
	columns := []string{"uuid", "post_id", "title", "type", "content", "likes", "created_at", "users_who_liked"}
	query := config.InsertQuery(config.PostTable, columns)
	emptyArray := []byte("[]")
	//query := "INSERT INTO posts (user_id, title,type, content, likes,created_at) VALUES (?, ?, ?, ?, ?,?)"
	_, err := r.DB.Exec(query, post.UId, post.PostId, post.Title, post.Type, post.Content, post.Likes, post.CreatedAt, emptyArray)
	return err
}

func (r *MySQLPostRepository) GetAllPosts(limit, offset int, search string, filter string) ([]*models.Post, error) {
	columns := []string{"post_id", "uuid", "title", "type", "content", "likes", "created_at", "users_who_liked"}
	query := config.SelectQuery(config.PostTable, "", "", columns)
	var conditions []string
	var params []interface{}
	if search != "" {
		conditions = append(conditions, " (title LIKE CONCAT('%', ?, '%') OR content LIKE CONCAT('%', ?, '%'))")
		params = append(params, search, search)
	}
	if filter != "" {
		conditions = append(conditions, " type = ?")
		params = append(params, filter)
	}
	if len(conditions) > 0 {
		query = query + " WHERE" + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	params = append(params, limit, offset)
	fmt.Println(query)
	fmt.Println(params)
	rows, err := r.DB.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			utils.Logger.Error("ERROR : Error closing rows:" + err.Error())
		}
	}(rows)

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		var createdAt string
		var usersWhoLiked []byte
		if err := rows.Scan(&post.PostId, &post.UId, &post.Title, &post.Type, &post.Content, &post.Likes, &createdAt, &usersWhoLiked); err != nil {
			return nil, err
		}
		if createdAt != "" {
			parsedTime, err := time.Parse("2006-01-02T15:04:05Z", createdAt)
			if err != nil {
				return nil, err
			}
			post.CreatedAt = parsedTime
		}
		if usersWhoLiked != nil {
			err := json.Unmarshal(usersWhoLiked, &post.Users)
			if err != nil {
				return nil, err
			}
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *MySQLPostRepository) DeleteByPId(pId string) error {
	condition1 := "post_id"
	query := config.DeleteQuery(config.PostTable, condition1, "")
	//query := "DELETE FROM posts WHERE post_id = ?"
	result, err := r.DB.Exec(query, pId)
	if result != nil {
		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return utils.NoPost
		}
	}
	return err
}

func (r *MySQLPostRepository) DeleteByUIdPId(uId, pId string) error {
	condition1 := "post_id"
	condition2 := "uuid"
	query := config.DeleteQuery(config.PostTable, condition1, condition2)
	//query := "DELETE FROM posts WHERE post_id = ? AND uuid=?"
	result, err := r.DB.Exec(query, pId, uId)
	if result != nil {
		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return utils.NotYourPost
		}
	}
	return err
}

func (r *MySQLPostRepository) GetPostsByFilter(filter string) ([]*models.Post, error) {
	columns := []string{"post_id", "uuid", "title", "type", "content", "likes", "created_at", "users_who_liked"}
	condition1 := "type"
	query := config.SelectQuery(config.PostTable, condition1, "", columns)
	//query := "SELECT post_id, uuid, title,type, content, likes,created_at FROM posts WHERE type = ?"
	rows, err := r.DB.Query(query, filter)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			utils.Logger.Error("ERROR: Error closing rows:" + err.Error())
		}
	}(rows)

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		var createdAt string
		var usersWhoLiked []byte
		if err := rows.Scan(&post.PostId, &post.UId, &post.Title, &post.Type, &post.Content, &post.Likes, &createdAt, &usersWhoLiked); err != nil {
			return nil, err
		}
		if createdAt != "" {
			parsedTime, err := time.Parse("2006-01-02T15:04:05Z", createdAt)
			if err != nil {
				return nil, err
			}
			post.CreatedAt = parsedTime
		}
		if usersWhoLiked != nil {
			err := json.Unmarshal(usersWhoLiked, &post.Users)
			if err != nil {
				return nil, err
			}
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *MySQLPostRepository) GetPostsByUId(uId string) ([]*models.Post, error) {
	condition1 := "uuid"
	columns := []string{"post_id", "uuid", "title", "type", "content", "likes", "created_at", "users_who_liked"}
	query := config.SelectQuery(config.PostTable, condition1, "", columns)
	//query := "SELECT post_id, uuid, title,type, content, likes,created_at FROM posts WHERE uuid = ?"
	rows, err := r.DB.Query(query, uId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			utils.Logger.Error("ERROR : Error closing rows:" + err.Error())
		}
	}(rows)

	var posts []*models.Post
	var usersWhoLiked []byte
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.PostId, &post.UId, &post.Title, &post.Type, &post.Content, &post.Likes, &post.CreatedAt, &usersWhoLiked); err != nil {
			return nil, err
		}
		if usersWhoLiked != nil {
			err := json.Unmarshal(usersWhoLiked, &post.Users)
			if err != nil {
				return nil, err
			}
		}
		posts = append(posts, &post)
	}
	return posts, nil
}

func (r *MySQLPostRepository) GetPostByPId(pId string) (*models.Post, error) {
	columns := []string{"post_id", "uuid", "title", "type", "content", "likes", "created_at", "users_who_liked"}
	condition1 := "post_id"
	query := config.SelectQuery(config.PostTable, condition1, "", columns)
	//query := "SELECT post_id, uuid, title,type, content, likes,created_at FROM posts WHERE post_id = ?"
	var post models.Post
	var createdAt string
	var usersWhoLiked []byte
	rows, err := r.DB.Query(query, pId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&post.PostId, &post.UId, &post.Title, &post.Type, &post.Content, &post.Likes, &createdAt, &usersWhoLiked)
		if err != nil {
			return nil, err
		}
	}
	if createdAt != "" {
		parsedTime, err := time.Parse("2006-01-02T15:04:05Z", createdAt)
		if err != nil {
			return nil, err
		}
		post.CreatedAt = parsedTime
	}
	if usersWhoLiked != nil {
		err := json.Unmarshal(usersWhoLiked, &post.Users)
		if err != nil {
			return nil, err
		}
	}
	return &post, nil
}

func (r *MySQLPostRepository) UpdateUserPost(pId, uId, title, content string) error {
	columns := []string{"title", "content"}
	condition1 := "post_id"
	condition2 := "uuid"
	query := config.UpdateQuery(config.PostTable, condition1, condition2, columns)
	//query := "UPDATE posts SET title = ?, content = ? WHERE post_id = ? AND uuid = ?"
	result, err := r.DB.Exec(query, title, content, pId, uId)
	if result != nil {
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return utils.NotYourPost
		}
	}
	return err
}

func (r *MySQLPostRepository) Like(uId, pId string) error {
	columns := "likes = likes+1,users_who_liked = JSON_ARRAY_APPEND(users_who_liked, '$' , ?)"
	condition1 := "post_id"
	query := config.UpdateQueryWithValue(config.PostTable, condition1, "", columns)
	//query := "UPDATE posts SET likes = likes+1 WHERE post_id = ?"
	result, err := r.DB.Exec(query, uId, pId)
	if result != nil {
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return utils.NoPost
		}
	}
	return err
}

func (r *MySQLPostRepository) Dislike(uId, pId string) error {
	columns := "likes = likes-1,users_who_liked = JSON_REMOVE(users_who_liked, JSON_UNQUOTE(JSON_SEARCH(users_who_liked, 'one', ?)))"
	condition1 := "post_id"
	query := config.UpdateQueryWithValue(config.PostTable, condition1, "", columns)
	//query := "UPDATE posts SET likes = likes+1 WHERE post_id = ?"
	result, err := r.DB.Exec(query, uId, pId)
	if result != nil {
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return utils.NoPost
		}
	}
	return err
}
