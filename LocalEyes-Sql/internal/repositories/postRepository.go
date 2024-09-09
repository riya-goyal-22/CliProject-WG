package repositories

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"localEyes/constants"
	"localEyes/internal/models"
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
	query := "INSERT INTO posts (user_id, title,type, content, likes,created_at) VALUES (?, ?, ?, ?, ?,?)"
	_, err := r.DB.Exec(query, post.UId, post.Title, post.Type, post.Content, post.Likes, post.CreatedAt)
	return err
}

func (r *MySQLPostRepository) GetAllPosts() ([]*models.Post, error) {
	query := "SELECT post_id, user_id, title, type, content, likes, created_at FROM posts"
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		var createdAt string
		if err := rows.Scan(&post.PostId, &post.UId, &post.Title, &post.Type, &post.Content, &post.Likes, &createdAt); err != nil {
			return nil, err
		}
		if createdAt != "" {
			parsedTime, err := time.Parse("2006-01-02T15:04:05Z", createdAt)
			if err != nil {
				return nil, err
			}
			post.CreatedAt = parsedTime
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *MySQLPostRepository) DeleteByPId(PId int) error {
	query := "DELETE FROM posts WHERE post_id = ?"
	result, err := r.DB.Exec(query, PId)
	if result != nil {
		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New(constants.Red + "No Post exist with this id" + constants.Reset)
		}
	}
	return err
}

func (r *MySQLPostRepository) DeleteByUIdPId(UId, PId int) error {
	query := "DELETE FROM posts WHERE post_id = ? AND user_id=?"
	result, err := r.DB.Exec(query, PId, UId)
	if result != nil {
		affectedRows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New(constants.Red + "No Post exist with this id" + constants.Reset)
		}
	}
	return err
}

func (r *MySQLPostRepository) DeleteByUId(UId int) error {
	query := "DELETE FROM posts WHERE user_id = ?"
	_, err := r.DB.Exec(query, UId)
	return err
}

func (r *MySQLPostRepository) GetPostsByFilter(filter string) ([]*models.Post, error) {
	query := "SELECT post_id, user_id, title,type, content, likes,created_at FROM posts WHERE type = ?"
	rows, err := r.DB.Query(query, filter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		var createdAt string
		if err := rows.Scan(&post.PostId, &post.UId, &post.Title, &post.Type, &post.Content, &post.Likes, &createdAt); err != nil {
			return nil, err
		}
		if createdAt != "" {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAt)
			if err != nil {
				return nil, err
			}
			post.CreatedAt = parsedTime
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *MySQLPostRepository) GetPostsByUId(UId int) ([]*models.Post, error) {
	query := "SELECT post_id, user_id, title,type, content, likes,created_at FROM posts WHERE user_id = ?"
	rows, err := r.DB.Query(query, UId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.PostId, &post.UId, &post.Title, &post.Type, &post.Content, &post.Likes, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *MySQLPostRepository) GetPostsByPId(PId int) ([]*models.Post, error) {
	query := "SELECT post_id, user_id, title,type, content, likes,created_at FROM posts WHERE post_id = ?"
	rows, err := r.DB.Query(query, PId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.PostId, &post.UId, &post.Title, &post.Type, &post.Content, &post.Likes, &post.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *MySQLPostRepository) UpdateUserPost(PId, UId int, title, content string) error {
	query := "UPDATE posts SET title = ?, content = ? WHERE post_id = ? AND user_id = ?"
	result, err := r.DB.Exec(query, title, content, PId, UId)
	if result != nil {
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return errors.New(constants.Red + "You can only update your post" + constants.Reset)
		}
	}
	return err
}

func (r *MySQLPostRepository) UpdateLike(PId int) error {
	query := "UPDATE posts SET likes = likes + 1 WHERE post_id = ?"
	result, err := r.DB.Exec(query, PId)
	if result != nil {
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return errors.New(constants.Red + "No post exist with this id" + constants.Reset)
		}
	}
	return err
}
