//go:build !test
// +build !test

package cli

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"localEyes/constants"
	"localEyes/internal/services"
	"localEyes/utils"
)

func adminLogin(adminService *services.AdminService) {
	fmt.Println(constants.Blue + "\n==============================")
	fmt.Println("ADMIN LOGIN")
	fmt.Println("=============================" + constants.Reset)
	//username := utils.PromptInput("Enter your username:")
	prompt := &promptui.Prompt{
		Label:     constants.Cyan + "Enter your password" + constants.Reset,
		Mask:      '*',
		IsConfirm: false,
	}
	password := utils.PromptPassword(prompt)
	_, err := adminService.Login(password)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(constants.Green + "\nAdmin logged in successfully\n" + constants.Reset)

	for {
		fmt.Println(constants.Blue + "\n1.View Users")
		fmt.Println("2.View Questions")
		fmt.Println("3.View Posts")
		fmt.Println("4.Delete a user")
		fmt.Println("5.Delete a question")
		fmt.Println("6.Delete a post")
		fmt.Println("7.ReActivate User")
		fmt.Println("8.Return" + constants.Reset)
		choice := utils.GetChoice()
		switch choice {
		case 1:
			users, err := adminService.GetAllUsers()
			if err != nil {
				fmt.Println(err)
			} else {
				displayUsers(users)
			}
		case 2:
			questions, err := adminService.GetAllQuestions()
			if err != nil {
				fmt.Println(err)
			} else {
				displayQuestions(questions)
			}
		case 3:
			posts, err := adminService.GetAllPosts()
			if err != nil {
				fmt.Println(err)
			} else {
				displayPosts(posts)
			}
		case 4:
			UId, err := utils.PromptIntInput("Enter User Id to delete user:")
			err = adminService.DeleteUser(UId)
			if err != nil {
				fmt.Println(constants.Red + "Error deleting user:" + err.Error() + constants.Reset)
			} else {
				fmt.Println(constants.Green + "User deleted" + constants.Reset)
			}
		case 5:
			QId, err := utils.PromptIntInput("Enter Question Id to delete question:")
			err = adminService.DeleteQuestion(QId)
			if err != nil {
				fmt.Println(constants.Red + "Error deleting question:" + err.Error() + constants.Reset)
			} else {
				fmt.Println(constants.Green + "Question deleted" + constants.Reset)
			}
		case 6:
			PId, err := utils.PromptIntInput("Enter Post Id to delete post:")
			err = adminService.DeletePost(PId)
			if err != nil {
				fmt.Println(constants.Red + "Error deleting post:" + err.Error() + constants.Reset)
			} else {
				fmt.Println(constants.Green + "Post deleted" + constants.Reset)
			}
		case 7:
			UId, err := utils.PromptIntInput("Enter User Id to Activate user:")
			err = adminService.ReActivate(UId)
			if err != nil {
				fmt.Println(constants.Red + "Error activating user:" + err.Error() + constants.Reset)
			} else {
				fmt.Println(constants.Green + "User activated" + constants.Reset)
			}
		case 8:
			return
		default:
			fmt.Println(constants.Red + "Invalid choice" + constants.Reset)
		}

	}
}
