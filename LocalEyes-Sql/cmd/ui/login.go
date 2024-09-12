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

func login(userService *services.UserService, questionService *services.QuestionService, postService *services.PostService) {
	fmt.Println(config.Blue + "==============================")
	fmt.Println("LOGIN")
	fmt.Println("=============================" + config.Reset)
	username := utils.PromptInput("Enter your username:")
	prompt := &promptui.Prompt{
		Label:     config.Cyan + "Enter your password" + config.Reset,
		Mask:      '*',
		IsConfirm: false,
	}
	password := utils.PromptPassword(prompt)
	user, err := userService.Login(username, password)
	if err != nil {
		fmt.Println(err)
		return
	}
	user.NotifyChannel = make(chan string)
	go func() {
		for _, s := range user.Notification {
			user.NotifyChannel <- s
		}
	}()
	//for _, s := range user.Notification {
	//	user.NotifyChannel <- s
	//}
	fmt.Println(config.Green + "\nUser logged in successfully 😊" + config.Reset)
	for i := 0; i < len(user.Notification); i++ {
		select {
		case msg := <-user.NotifyChannel:
			fmt.Print(config.Gray + msg + config.Reset)
		default:
			break
		}
	}
	err = userService.UnNotifyUsers(user.UId)
	if err != nil {
		fmt.Println(err)
	}

	for {
		fmt.Println(config.Blue + "\n1.View my Profile")
		fmt.Println("2.Manage posts")
		fmt.Println("3.Deactivate account")
		fmt.Println("4.Return" + config.Reset)
		choice := utils.GetChoice()
		switch choice {
		case 1:
			fmt.Println(config.Magenta + "\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			fmt.Println("Welcome ", user.Username)
			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~" + config.Reset)
			fmt.Println("City:", user.City)
			fmt.Println("Type of user:", user.Tag)
			fmt.Printf("Living in City for:%v years\n", user.DwellingAge)
		case 2:
			managePost(postService, questionService, userService, user.UId)
		case 3:
			err := userService.DeActivate(user.UId)
			if err != nil {
				fmt.Println(config.Red + "Error Deactivating user:" + err.Error() + config.Reset)
				utils.Logger.Println("ERROR: Error Deactivating user :" + err.Error())
			} else {
				fmt.Println(config.Green + "User Deactivated successfully" + config.Reset)
				utils.Logger.Println("INFO: User Deactivated with id-", user.UId)
				return
			}
		case 4:
			return
		default:
			fmt.Println(config.Red + "Invalid Choice,Try Again" + config.Reset)
		}
	}
}
