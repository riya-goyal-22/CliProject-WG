package models

type Admin struct {
	User User `bson:"user"`
}
