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

func login(userService *services.UserService, questionService *services.QuestionService, postService *services.PostService) {
	fmt.Println(constants.Blue + "==============================")
	fmt.Println("LOGIN")
	fmt.Println("=============================" + constants.Reset)
	username := utils.PromptInput("Enter your username:")
	prompt := &promptui.Prompt{
		Label:     constants.Cyan + "Enter your password" + constants.Reset,
		Mask:      '*',
		IsConfirm: false,
	}
	password := utils.PromptPassword(prompt)
	user, err := userService.Login(username, password)
	if err != nil {
		fmt.Println(err)
		return
	}
	user.NotifyChannel = make(chan string, 5)
	for _, s := range user.Notification {
		user.NotifyChannel <- s
	}
	fmt.Println(constants.Green + "\nUser logged in successfully 😊" + constants.Reset)
	select {
	case msg := <-user.NotifyChannel:
		fmt.Println(constants.Gray + msg + "\n" + constants.Reset)
	default:
		break
	}
	err = userService.UnNotifyUsers(user.UId)
	if err != nil {
		fmt.Println(err)
	}

	for {
		fmt.Println(constants.Blue + "\n1.View my Profile")
		fmt.Println("2.Create post")
		fmt.Println("3.Update Post")
		fmt.Println("4 View Posts")
		fmt.Println("5.Filter Posts")
		fmt.Println("6.Open Post")
		fmt.Println("7.Like Post")
		fmt.Println("8.Delete Post")
		fmt.Println("9.Deactivate account")
		fmt.Println("10.Return" + constants.Reset)
		choice := utils.GetChoice()
		switch choice {
		case 1:
			fmt.Println(constants.Magenta + "\n~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
			fmt.Println("Welcome ", user.Username)
			fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~" + constants.Reset)
			fmt.Println("City:", user.City)
			fmt.Println("Type of user:", user.Tag)
			fmt.Printf("Living in City for:%v years\n", user.DwellingAge)
		case 2:
			postCreate(postService, userService, user.UId)
		case 3:
			myPosts, err := postService.GiveMyPosts(user.UId)
			if err != nil {
				fmt.Println(constants.Red + "Error loading posts:" + err.Error() + constants.Reset)
			} else {
				displayPosts(myPosts)
			}
			PId, err := utils.PromptIntInput("Enter post id to update:")
			if err != nil {
				fmt.Println(constants.Red + err.Error() + constants.Reset)
			}
			title := utils.PromptInput("Enter new post title:")
			content := utils.PromptInput("Enter new post content:")
			err = postService.UpdateMyPost(PId, user.UId, title, content)
			if err != nil {
				fmt.Println(constants.Red + "Error updating post:" + err.Error() + constants.Reset)
			} else {
				fmt.Println(constants.Green+"Post updated:", title)
			}
		case 4:
			posts, err := postService.GiveAllPosts()
			if err != nil {
				fmt.Println(constants.Red + "Error loading posts:" + err.Error() + constants.Reset)
			} else {
				displayPosts(posts)
			}

		case 5:
			var filterType string
			for {
				filterType = utils.PromptInput("Enter filter [food/travel/shopping/other]:")
				if utils.ValidateFilter(filterType) {
					break
				} else {
					fmt.Println("Invalid filter type:", filterType)
				}
			}
			posts, err := postService.GiveFilteredPosts(filterType)
			if err != nil {
				fmt.Println(constants.Red + "Error loading posts:" + err.Error() + constants.Reset)
			} else {
				displayPosts(posts)
			}

		case 6:
			PId, err := utils.PromptIntInput("Enter post id to open:")
			if err != nil {
				fmt.Println(constants.Red + err.Error() + constants.Reset)
				break
			}
			openPost(questionService, postService, PId, user.UId)

		case 7:
			PId, err := utils.PromptIntInput("Enter post id to like:")
			err = postService.Like(PId)
			if err != nil {
				fmt.Println(constants.Red + "Error liking post:" + err.Error() + constants.Reset)
			} else {
				fmt.Println(constants.Green + "Post Liked" + constants.Reset)
			}

		case 8:
			myPosts, err := postService.GiveMyPosts(user.UId)
			if err != nil {
				fmt.Println(constants.Red + "Error loading posts:" + err.Error() + constants.Reset)
			} else {
				displayPosts(myPosts)
			}
			PId, err := utils.PromptIntInput("Enter post id to delete:")
			if err != nil {
				fmt.Println(constants.Red + "Error taking postId input:" + err.Error() + constants.Reset)
			}
			err = postService.DeleteMyPost(user.UId, PId)
			if err != nil {
				fmt.Println(constants.Red + "Error deleting post:" + err.Error() + constants.Reset)
			} else {
				err = questionService.DeleteQuesByPId(PId)
				if err != nil {
					fmt.Println(constants.Red + "Error deleting question:" + err.Error() + constants.Reset)
				} else {
					fmt.Println(constants.Green + "Post deleted successfully" + constants.Reset)
				}
			}
		case 9:
			err := userService.DeActivate(user.UId)
			if err != nil {
				fmt.Println(constants.Red + "Error Deactivating user:" + err.Error() + constants.Reset)
			} else {
				fmt.Println(constants.Green + "User Deactivated successfully" + constants.Reset)
				return
			}
		case 10:
			return
		default:
			fmt.Println(constants.Red + "Invalid Choice,Try Again" + constants.Reset)
		}
	}
}

func postCreate(postService *services.PostService, userService *services.UserService, UId int) {
	fmt.Println(constants.Blue + "\n1.Create Food post")
	fmt.Println("2.Create Travel post")
	fmt.Println("3.Create Shopping post")
	fmt.Println("4.Create Other post" + constants.Reset)
	choice := utils.GetChoice()
	switch choice {
	case 1:
		title := utils.PromptInput("Enter post title:")
		content := utils.PromptInput("Enter post content:")
		err := postService.CreatePost(UId, title, content, "food")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(constants.Green+"Post created:", title)
			err := userService.NotifyUsers(UId, title)
			if err != nil {
				fmt.Println(err)
			}
		}
	case 2:
		title := utils.PromptInput("Enter post title:")
		content := utils.PromptInput("Enter post content:")
		err := postService.CreatePost(UId, title, content, "travel")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(constants.Green+"Post created:", title)
			err := userService.NotifyUsers(UId, title)
			if err != nil {
				fmt.Println(err)
			}
		}
	case 3:
		title := utils.PromptInput("Enter post title:")
		content := utils.PromptInput("Enter post content:")
		err := postService.CreatePost(UId, title, content, "shopping")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(constants.Green+"Post created:", title)
			err := userService.NotifyUsers(UId, title)
			if err != nil {
				fmt.Println(err)
			}
		}
	case 4:
		title := utils.PromptInput("Enter post title:")
		content := utils.PromptInput("Enter post content:")
		err := postService.CreatePost(UId, title, content, "other")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(constants.Green+"Post created:", title)
			err := userService.NotifyUsers(UId, title)
			if err != nil {
				fmt.Println(err)
			}
		}
	default:
		fmt.Println(constants.Red + "invalid choice" + constants.Reset)
	}
}
