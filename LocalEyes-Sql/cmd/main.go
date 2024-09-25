package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"localEyes/config"
	"localEyes/internal/handlers"
	"localEyes/internal/middlewares"
	"localEyes/internal/repositories"
	"localEyes/internal/services"
	"log"
	"net/http"
)

var dbClient *sql.DB

func init() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	dbClient = config.GetSQLClient()
}

func main() {
	defer config.CloseDBClient()

	router := mux.NewRouter()
	userService := services.NewUserService(repositories.NewMySQLUserRepository(dbClient))
	postService := services.NewPostService(repositories.NewMySQLPostRepository(dbClient), repositories.NewMySQLUserRepository(dbClient))
	questionService := services.NewQuestionService(repositories.NewMySQLQuestionRepository(dbClient))

	adminService := services.NewAdminService(repositories.NewMySQLUserRepository(dbClient),
		repositories.NewMySQLPostRepository(dbClient),
		repositories.NewMySQLQuestionRepository(dbClient))

	userHandler := handlers.NewUserHandler(userService)
	postHandler := handlers.NewPostHandler(postService)
	questionHandler := handlers.NewQuestionHandler(questionService)
	adminHandler := handlers.NewAdminHandler(adminService)

	router.HandleFunc("/signup", userHandler.SignUp).Methods("POST")
	router.HandleFunc("/login", userHandler.Login).Methods("POST")

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(middlewares.AuthenticationMiddleware)
	apiRouter.HandleFunc("/user/deactivate", userHandler.DeActivate).Methods("POST")
	apiRouter.HandleFunc("/user/profile", userHandler.ViewProfile).Methods("GET")
	apiRouter.HandleFunc("/user/notification", userHandler.ViewNotifications).Methods("GET")
	apiRouter.HandleFunc("/posts/all", postHandler.DisplayPosts).Methods("GET")
	apiRouter.HandleFunc("/post", postHandler.CreatePost).Methods("POST")
	apiRouter.HandleFunc("/post/{post_id}", postHandler.DisplayPostById).Methods("GET")
	apiRouter.HandleFunc("/post/{post_id}/like", postHandler.LikePost).Methods("POST")
	apiRouter.HandleFunc("/user/posts/all", postHandler.DisplayUserPosts).Methods("GET")
	apiRouter.HandleFunc("/user/post/{post_id}", postHandler.UpdatePost).Methods("PUT")
	apiRouter.HandleFunc("/user/post/{post_id}", postHandler.DeletePost).Methods("DELETE")
	apiRouter.HandleFunc("/post/{post_id}/questions/all", questionHandler.GetQuestions).Methods("GET")
	apiRouter.HandleFunc("/post/{post_id}/question", questionHandler.CreateQuestion).Methods("POST")
	apiRouter.HandleFunc("/post/{post_id}/question/{ques_id}", questionHandler.AddAnswer).Methods("PUT")
	apiRouter.HandleFunc("/post/{post_id}/question/{ques_id}", questionHandler.DeleteQuestion).Methods("DELETE")

	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middlewares.AdminAuthMiddleware)
	adminRouter.HandleFunc("/users", adminHandler.DisplayUsers).Methods("GET")
	adminRouter.HandleFunc("/questions", adminHandler.DisplayQuestions).Methods("GET")
	adminRouter.HandleFunc("/user/{user_id}", adminHandler.DeleteUser).Methods("DELETE")
	adminRouter.HandleFunc("/post/{post_id}", adminHandler.DeletePost).Methods("DELETE")
	adminRouter.HandleFunc("/question/{ques_id}", adminHandler.DeleteQuestion).Methods("DELETE")
	adminRouter.HandleFunc("reactivate/user/{user_id}", adminHandler.ReactivateUser).Methods("POST")

	//ui.RootCli(userService, postService, questionService, adminService)
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(config.Magenta + "Thank you 😊, Visit Again" + config.Reset)
}
