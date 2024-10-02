package models

import (
	"time"
)

type Post struct {
	PostId    string    `json:"post_id"`
	UId       string    `json:"user_id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	Likes     int       `json:"likes"`
	CreatedAt time.Time `json:"created_at"`
}

type PostQuestion struct {
	QId          string   `json:"q_id"`
	QuestionText string   `json:"question_text"`
	Replies      []string `json:"replies"`
}

type PostWithQuestions struct {
	PostId    string         `json:"post_id"`
	UId       string         `json:"uuid"`
	Title     string         `json:"title"`
	Type      string         `json:"type"`
	Content   string         `json:"content"`
	Likes     int            `json:"likes"`
	CreatedAt time.Time      `json:"created_at"`
	Questions []PostQuestion `json:"questions"`
}
