//go:build !test
// +build !test

package cli

import (
	"fmt"
	"localEyes/constants"
	"localEyes/internal/services"
	"localEyes/utils"
)

func openPost(questionService *services.QuestionService, postService *services.PostService, PId, UId int) {
	boolVal, err := postService.PostIdExist(PId)
	if err != nil {
		fmt.Println(constants.Red + err.Error() + constants.Reset)
	}
	if !boolVal {
		fmt.Println(constants.Red + "Post Id does not exist" + constants.Reset)
		return
	}

	for {
		fmt.Println(constants.Blue + "\n1.Add Question")
		fmt.Println("2.Answer a Question")
		fmt.Println("3.View Questions")
		fmt.Println("4.Delete Question")
		fmt.Println("5 Return" + constants.Reset)
		choice := utils.GetChoice()
		switch choice {
		case 1:
			text := utils.PromptInput("Enter your Question:")
			err := questionService.AskQuestion(UId, PId, text)
			if err != nil {
				fmt.Println(constants.Red + "Error Adding question:" + err.Error() + constants.Reset)
			} else {
				fmt.Println(constants.Green + "Question added" + constants.Reset)
			}
		case 2:
			QId, err := utils.PromptIntInput("Enter QId:")
			answer := utils.PromptInput("Enter your answer:")
			err = questionService.AddAnswer(QId, answer)
			if err != nil {
				fmt.Println(constants.Red + "Error Adding answer:" + err.Error() + constants.Reset)
			} else {
				fmt.Println(constants.Green + "Answer added" + constants.Reset)
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
				fmt.Println(constants.Red + "Error deleting question:" + err.Error() + constants.Reset)
			} else {
				fmt.Println(constants.Green + "Question deleted" + constants.Reset)
			}
		case 5:
			return
		default:
			fmt.Println(constants.Red + "Invalid Choice" + constants.Reset)
		}
	}
}
