package utils

import "errors"

var NotYourPost = errors.New("no post of yours exist with this id")
var NotYourQuestion = errors.New("no question of yours exist with this id")
var NoPost = errors.New("no post exist with this id")
var NoQuestion = errors.New("no question exist with this id")
var NoUser = errors.New("no user exist with this id")
var TitleMissing = errors.New("required field 'title' is missing")
var ContentMissing = errors.New("required field 'content' is missing")
var TypeMissing = errors.New("required field 'type' is missing")
var InvalidPost = errors.New("invalid post type")
