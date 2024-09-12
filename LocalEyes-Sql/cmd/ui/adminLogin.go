//go:build !test
// +build !test

package ui

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"localEyes/config"
	"localEyes/internal/services"
	"localEyes/utils"
)

func adminLogin(adminService *services.AdminService) {
	fmt.Println(config.Blue + "\n==============================")
	fmt.Println("ADMIN LOGIN")
	fmt.Println("=============================" + config.Reset)
	//username := utils.PromptInput("Enter your username:")
	prompt := &promptui.Prompt{
		Label:     config.Cyan + "Enter your password" + config.Reset,
		Mask:      '*',
		IsConfirm: false,
	}
	password := utils.PromptPassword(prompt)
	_, err := adminService.Login(password)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(config.Green + "\nAdmin logged in successfully" + config.Reset)

	for {
		fmt.Println(config.Blue + "\n1.View Users")
		fmt.Println("2.View Questions")
		fmt.Println("3.View Posts")
		fmt.Println("4.Delete a user")
		fmt.Println("5.Delete a question")
		fmt.Println("6.Delete a post")
		fmt.Println("7.ReActivate User")
		fmt.Println("8.Return" + config.Reset)
		choice := utils.GetChoice()
		switch choice {
		case 1:
			users, err := adminService.GetAllUsers()
			if err != nil {
				fmt.Println(err)
			} else {
				displayUsers(users)
				utils.Logger.Println("INFO:Admin viewed all users")
			}
		case 2:
			questions, err := adminService.GetAllQuestions()
			if err != nil {
				fmt.Println(err)
			} else {
				displayQuestions(questions)
				utils.Logger.Println("INFO:Admin viewed all questions")
			}
		case 3:
			posts, err := adminService.GetAllPosts()
			if err != nil {
				fmt.Println(err)
			} else {
				displayPosts(posts)
				utils.Logger.Println("INFO:Admin viewed all posts")
			}
		case 4:
			uId, err := utils.PromptIntInput("Enter User Id to delete user:")
			err = adminService.DeleteUser(uId)
			if err != nil {
				fmt.Println(config.Red + "Error deleting user:" + err.Error() + config.Reset)
			} else {
				fmt.Println(config.Green + "User deleted" + config.Reset)
				utils.Logger.Println("INFO:Admin deleted user with id-", uId)
			}
		case 5:
			qId, err := utils.PromptIntInput("Enter Question Id to delete question:")
			err = adminService.DeleteQuestion(qId)
			if err != nil {
				fmt.Println(config.Red + "Error deleting question:" + err.Error() + config.Reset)
			} else {
				fmt.Println(config.Green + "Question deleted" + config.Reset)
				utils.Logger.Println("INFO:Admin deleted question with id-", qId)
			}
		case 6:
			pId, err := utils.PromptIntInput("Enter Post Id to delete post:")
			err = adminService.DeletePost(pId)
			if err != nil {
				fmt.Println(config.Red + "Error deleting post:" + err.Error() + config.Reset)
			} else {
				fmt.Println(config.Green + "Post deleted" + config.Reset)
				utils.Logger.Println("INFO:Admin deleted post with id-", pId)
			}
		case 7:
			uId, err := utils.PromptIntInput("Enter User Id to Activate user:")
			err = adminService.ReActivate(uId)
			if err != nil {
				fmt.Println(config.Red + "Error activating user:" + err.Error() + config.Reset)
			} else {
				fmt.Println(config.Green + "User activated" + config.Reset)
				utils.Logger.Println("INFO:Admin activated user with id-", uId)
			}
		case 8:
			return
		default:
			fmt.Println(config.Red + "Invalid choice" + config.Reset)
		}

	}
}
