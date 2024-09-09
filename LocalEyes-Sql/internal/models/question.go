package models

import (
	"time"
)

type Question struct {
	QId       int       `bson:"q_id"`
	PostId    int       `bson:"post_id"`
	UserId    int       `bson:"user_id"`
	Text      string    `bson:"text"`
	Replies   []string  `bson:"replies"`
	CreatedAt time.Time `bson:"created_at"`
}
