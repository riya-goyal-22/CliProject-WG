package utils

import (
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"strings"
)

var NotYourPost = errors.New("no post of yours exist with this id")
var NotYourQuestion = errors.New("no question of yours exist with this id")
var NoPost = errors.New("no post exist with this id")
var NoQuestion = errors.New("no question exist with this id")
var NoUser = errors.New("no user exist with this id")

func GenerateRandomId() string {
	id := uuid.New()
	idString := base64.RawStdEncoding.EncodeToString(id[:])[:4]
	if strings.Contains(idString, "/") {
		idString = strings.Replace(idString, "/", "A", -1)
		return idString
	}
	return idString
}
