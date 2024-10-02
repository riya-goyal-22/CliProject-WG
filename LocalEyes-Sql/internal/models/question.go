package models

import (
	"time"
)

type Question struct {
	QId       string    `json:"question_id"`
	PostId    string    `json:"post_id"`
	UserId    string    `json:"user_id"`
	Text      string    `json:"text"`
	Replies   []string  `json:"replies"`
	CreatedAt time.Time `json:"created_at"`
}
