package models

import (
	"time"
)

type Post struct {
	PostId    int       `bson:"id"`
	UId       int       `bson:"userId"`
	Title     string    `bson:"title"`
	Type      string    `bson:"type"`
	Content   string    `bson:"content"`
	Likes     int       `bson:"likes"`
	CreatedAt time.Time `bson:"created_at"`
}
