package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"localEyes/cmd/ui"
	"localEyes/config"
	"localEyes/internal/repositories"
	"localEyes/internal/services"
	"localEyes/utils"
	"log"
)

var dbClient *sql.DB

func init() {
	err := godotenv.Load("../config.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	dbClient = config.GetSQLClient()
	utils.InitLoggerFile()
}

func main() {
	defer config.CloseDBClient()
	defer utils.CloseLoggerFile()
	userService := services.NewUserService(repositories.NewMySQLUserRepository(dbClient))

	postService := services.NewPostService(repositories.NewMySQLPostRepository(dbClient))

	questionService := services.NewQuestionService(repositories.NewMySQLQuestionRepository(dbClient))

	adminService := services.NewAdminService(repositories.NewMySQLUserRepository(dbClient),
		repositories.NewMySQLPostRepository(dbClient),
		repositories.NewMySQLQuestionRepository(dbClient))

	ui.RootCli(userService, postService, questionService, adminService)

	fmt.Println(config.Magenta + "Thank you 😊, Visit Again" + config.Reset)
}
