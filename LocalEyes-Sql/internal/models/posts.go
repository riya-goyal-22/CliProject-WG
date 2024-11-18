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
	Users     []string  `json:"users"`
}

type PostWithQuestions struct {
	PostId    string             `json:"post_id"`
	UId       string             `json:"uuid"`
	Title     string             `json:"title"`
	Type      string             `json:"type"`
	Content   string             `json:"content"`
	Likes     int                `json:"likes"`
	CreatedAt time.Time          `json:"created_at"`
	Users     []string           `json:"users"`
	Questions []ResponseQuestion `json:"questions"`
}
