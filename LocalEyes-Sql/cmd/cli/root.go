//go:build !test
// +build !test

package cli

import (
	"fmt"
	"localEyes/constants"
	"localEyes/internal/services"
	"localEyes/utils"
)

func RootCli(userService *services.UserService, postService *services.PostService, questionService *services.QuestionService, adminService *services.AdminService) {
	for {
		fmt.Println(constants.Magenta + "\n=====================================================")
		fmt.Println("Welcome to Local Eyes!")
		fmt.Println("=====================================================" + constants.Reset)
		fmt.Println(constants.Blue + "1. Sign Up")
		fmt.Println("2. Log In")
		fmt.Println("3. Admin login")
		fmt.Println("4. Exit" + constants.Reset)

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
			fmt.Println(constants.Red + "Invalid choice, please try again." + constants.Reset)
		}
	}
}
