//go:build !test
// +build !test

package ui

import (
	"fmt"
	"localEyes/config"
	"localEyes/internal/services"
	"localEyes/utils"
)

func RootCli(userService *services.UserService, postService *services.PostService, questionService *services.QuestionService, adminService *services.AdminService) {
	for {
		fmt.Println(config.Magenta + "\n=====================================================")
		fmt.Println("Welcome to Local Eyes!")
		fmt.Println("=====================================================" + config.Reset)
		fmt.Println(config.Blue + "1. Sign Up")
		fmt.Println("2. Log In")
		fmt.Println("3. Admin login")
		fmt.Println("4. Exit" + config.Reset)

		choice := utils.GetChoice()
		switch choice {
		case 1:
			signUp(userService)
		case 2:
			login(userService, questionService, postService)
		case 3:
			adminLogin(adminService)
		case 4:
			return
		default:
			fmt.Println(config.Red + "Invalid choice, please try again." + config.Reset)
		}
	}
}
