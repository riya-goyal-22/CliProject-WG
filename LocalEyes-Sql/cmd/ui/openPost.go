//go:build !test
// +build !test

package ui

import (
	"fmt"
	"localEyes/config"
	"localEyes/internal/services"
	"localEyes/utils"
)

func openPost(questionService *services.QuestionService, postService *services.PostService, PId, UId int) {
	boolVal, err := postService.PostIdExist(PId)
	if err != nil {
		fmt.Println(config.Red + err.Error() + config.Reset)
	}
	if !boolVal {
		fmt.Println(config.Red + "Post Id does not exist" + config.Reset)
		return
	}

	for {
		fmt.Println(config.Blue + "\n1.Add Question")
		fmt.Println("2.Answer a Question")
		fmt.Println("3.View Questions")
		fmt.Println("4.Delete Question")
		fmt.Println("5 Return" + config.Reset)
		choice := utils.GetChoice()
		switch choice {
		case 1:
			text := utils.PromptInput("Enter your Question:")
			err := questionService.AskQuestion(UId, PId, text)
			if err != nil {
				fmt.Println(config.Red + "Error Adding question:" + err.Error() + config.Reset)
			} else {
				fmt.Println(config.Green + "Question added" + config.Reset)
			}
		case 2:
			QId, err := utils.PromptIntInput("Enter QId:")
			answer := utils.PromptInput("Enter your answer:")
			err = questionService.AddAnswer(QId, answer)
			if err != nil {
				fmt.Println(config.Red + "Error Adding answer:" + err.Error() + config.Reset)
			} else {
				fmt.Println(config.Green + "Answer added" + config.Reset)
			}
		case 3:
			questions, err := questionService.GetPostQuestions(PId)
			if err != nil {
				fmt.Println(err)
			} else {
				displayQuestions(questions)
			}
		case 4:
			questions, err := questionService.GetPostQuestions(PId)
			if err != nil {
				fmt.Println(err)
			} else {
				displayQuestions(questions)
			}
			QId, err := utils.PromptIntInput("Enter Question Id to delete:")
			err = questionService.DeleteUserQues(UId, QId)
			if err != nil {
				fmt.Println(config.Red + "Error deleting question:" + err.Error() + config.Reset)
			} else {
				fmt.Println(config.Green + "Question deleted" + config.Reset)
			}
		case 5:
			return
		default:
			fmt.Println(config.Red + "Invalid Choice" + config.Reset)
		}
	}
}
