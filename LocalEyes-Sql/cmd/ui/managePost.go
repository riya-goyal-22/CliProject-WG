package ui

import (
	"fmt"
	"localEyes/config"
	"localEyes/internal/services"
	"localEyes/utils"
)

func managePost(postService *services.PostService, questionService *services.QuestionService, userService *services.UserService, uId int) {
	fmt.Println(config.Blue + "1.Create post")
	fmt.Println("2.Update Post")
	fmt.Println("3.View Posts")
	fmt.Println("4.Open Post")
	fmt.Println("5.Like Post")
	fmt.Println("6.Delete Post" + config.Reset)
	choice := utils.GetChoice()
	switch choice {
	case 1:
		postCreate(postService, userService, uId)
	case 2:
		myPosts, err := postService.GiveMyPosts(uId)
		if err != nil {
			utils.Logger.Println("ERROR: Error loading posts: " + err.Error())
			fmt.Println(config.Red + "Error loading posts:" + err.Error() + config.Reset)
		} else {
			displayPosts(myPosts)
		}
		PId, err := utils.PromptIntInput("Enter post id to update:")
		if err != nil {
			fmt.Println(config.Red + err.Error() + config.Reset)
		}
		title := utils.PromptInput("Enter new post title:")
		content := utils.PromptInput("Enter new post content:")
		err = postService.UpdateMyPost(PId, uId, title, content)
		if err != nil {
			fmt.Println(config.Red + "Error updating post:" + err.Error() + config.Reset)
		} else {
			fmt.Println(config.Green+"Post updated:", title)
		}

	case 3:
		var filterType string
		for {
			filterType = utils.PromptInput("Enter filter [food/travel/shopping/other/blank for no filter]:")
			if utils.ValidateFilter(filterType) {
				break
			} else {
				fmt.Println("Invalid filter type:", filterType)
			}
		}
		if filterType == "" {
			posts, err := postService.GiveAllPosts()
			if err != nil {
				utils.Logger.Println("ERROR: Error loading posts: " + err.Error())
				fmt.Println(config.Red + "Error loading posts:" + err.Error() + config.Reset)
			} else {
				displayPosts(posts)
			}
		} else {
			posts, err := postService.GiveFilteredPosts(filterType)
			if err != nil {
				utils.Logger.Println("ERROR: Error loading posts: " + err.Error())
				fmt.Println(config.Red + "Error loading posts:" + err.Error() + config.Reset)
			} else {
				displayPosts(posts)
			}
		}

	case 4:
		pId, err := utils.PromptIntInput("Enter post id to open:")
		if err != nil {
			fmt.Println(config.Red + err.Error() + config.Reset)
			break
		}
		openPost(questionService, postService, pId, uId)

	case 5:
		pId, err := utils.PromptIntInput("Enter post id to like:")
		err = postService.Like(pId)
		if err != nil {
			fmt.Println(config.Red + "Error liking post:" + err.Error() + config.Reset)
		} else {
			fmt.Println(config.Green + "Post Liked" + config.Reset)
		}

	case 6:
		myPosts, err := postService.GiveMyPosts(uId)
		if err != nil {
			utils.Logger.Println("ERROR: Error loading posts: " + err.Error())
			fmt.Println(config.Red + "Error loading posts:" + err.Error() + config.Reset)
		} else {
			displayPosts(myPosts)
		}
		pId, err := utils.PromptIntInput("Enter post id to delete:")
		if err != nil {
			fmt.Println(config.Red + "Error taking postId input:" + err.Error() + config.Reset)
		}
		err = postService.DeleteMyPost(uId, pId)
		if err != nil {
			fmt.Println(config.Red + "Error deleting post:" + err.Error() + config.Reset)
		} else {
			err = questionService.DeleteQuesByPId(pId)
			if err != nil {
				fmt.Println(config.Red + "Error deleting question:" + err.Error() + config.Reset)
			} else {
				fmt.Println(config.Green + "Post deleted successfully" + config.Reset)
			}
		}

	}
}

func postCreate(postService *services.PostService, userService *services.UserService, uId int) {
	fmt.Println(config.Blue + "1.Create Food post")
	fmt.Println("2.Create Travel post")
	fmt.Println("3.Create Shopping post")
	fmt.Println("4.Create Other post" + config.Reset)
	choice := utils.GetChoice()
	switch choice {
	case 1:
		title := utils.PromptInput("Enter post title:")
		content := utils.PromptInput("Enter post content:")
		err := postService.CreatePost(uId, title, content, "food")
		if err != nil {
			utils.Logger.Println("ERROR: Error creating post: " + err.Error())
			fmt.Println(err)
		} else {
			fmt.Println(config.Green+"Post created:", title)
			utils.Logger.Println("INFO: Post created:", title)
			err := userService.NotifyUsers(uId, title)
			if err != nil {
				utils.Logger.Println("ERROR: Error Notifying user: " + err.Error())
				fmt.Println(err)
			}
		}
	case 2:
		title := utils.PromptInput("Enter post title:")
		content := utils.PromptInput("Enter post content:")
		err := postService.CreatePost(uId, title, content, "travel")
		if err != nil {
			utils.Logger.Println("ERROR: Error creating post: " + err.Error())
			fmt.Println(err)
		} else {
			fmt.Println(config.Green+"Post created:", title)
			utils.Logger.Println("INFO: Post created:", title)
			err := userService.NotifyUsers(uId, title)
			if err != nil {
				utils.Logger.Println("ERROR: Error UnNotifying user: " + err.Error())
				fmt.Println(err)
			}
		}
	case 3:
		title := utils.PromptInput("Enter post title:")
		content := utils.PromptInput("Enter post content:")
		err := postService.CreatePost(uId, title, content, "shopping")
		if err != nil {
			utils.Logger.Println("ERROR: Error creating post: " + err.Error())
			fmt.Println(err)
		} else {
			fmt.Println(config.Green+"Post created:", title)
			utils.Logger.Println("INFO: Post created:", title)
			err := userService.NotifyUsers(uId, title)
			if err != nil {
				utils.Logger.Println("ERROR: Error UnNotifying user: " + err.Error())
				fmt.Println(err)
			}
		}
	case 4:
		title := utils.PromptInput("Enter post title:")
		content := utils.PromptInput("Enter post content:")
		err := postService.CreatePost(uId, title, content, "other")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(config.Green+"Post created:", title)
			err := userService.NotifyUsers(uId, title)
			if err != nil {
				utils.Logger.Println("ERROR: Error UnNotifying user: " + err.Error())
				fmt.Println(err)
			}
		}
	default:
		fmt.Println(config.Red + "invalid choice" + config.Reset)
	}
}
