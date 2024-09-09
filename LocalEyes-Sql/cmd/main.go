//go:build !test
// +build !test

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"localEyes/cmd/cli"
	"localEyes/internal/repositories"
	"localEyes/internal/services"
	"log"
	"os"
	"sync"
	"time"
)

var dbClient *sql.DB
var once sync.Once

func init() {
	err := godotenv.Load("../config.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	GetSQLClient()
	dbClient.SetConnMaxLifetime(time.Hour * 1)
}

func main() {
	userRepo := repositories.NewMySQLUserRepository(dbClient)
	userService := services.NewUserService(userRepo)

	postRepo := repositories.NewMySQLPostRepository(dbClient)
	postService := services.NewPostService(postRepo)

	questionRepo := repositories.NewMySQLQuestionRepository(dbClient)
	questionService := services.NewQuestionService(questionRepo)

	adminService := services.NewAdminService(userRepo, postRepo, questionRepo)

	cli.RootCli(userService, postService, questionService, adminService)
}

func GetSQLClient() {
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			//constants.DBUser,
			//constants.DBPassword,
			//constants.DBHost,
			//constants.DBPort,
			//constants.DBName,
			os.Getenv("DBUser"),
			os.Getenv("DBPassword"),
			os.Getenv("DBHost"),
			os.Getenv("DBPort"),
			os.Getenv("DBName"),
		)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Fatal(err)
		}

		if err := db.Ping(); err != nil {
			log.Fatal(err)
		}

		dbClient = db
	})
}
