//go:build !test
// +build !test

package ui

import (
	"github.com/olekukonko/tablewriter"
	"localEyes/internal/models"
	"os"
	"strconv"
	"strings"
)

func displayUsers(users []*models.User) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"UserId", "UserName", "City", "Resident Till", "ActiveStatus", "Tag"})

	// Add rows to the table, only including Name and City
	for _, user := range users {
		uIdStr := strconv.Itoa(user.UId)
		dwelling := strconv.Itoa(user.DwellingAge)
		activeStatus := "No"
		if user.IsActive {
			activeStatus = "Yes"
		}
		table.Append([]string{uIdStr, user.Username, user.City, dwelling, activeStatus, user.Tag})
	}

	// Render the table
	table.Render()
}

func displayPosts(posts []*models.Post) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"PostId", "Title", "Type", "Content", "Likes", "Created At"})

	// Add rows to the table, only including Name and City
	for _, post := range posts {
		pIdStr := strconv.Itoa(post.PostId)
		likes := strconv.Itoa(post.Likes)
		time := post.CreatedAt.Format("2006-01-02 15:04:05")
		table.Append([]string{pIdStr, post.Title, post.Type, post.Content, likes, time})
	}

	// Render the table
	table.Render()
}

func displayQuestions(questions []*models.Question) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"QID", "Question", "Replies", "Created At"})

	// Add rows to the table, only including Name and City
	for _, question := range questions {
		qIdStr := strconv.Itoa(question.QId)
		time := question.CreatedAt.Format("2006-01-02 15:04:05")
		replies := strings.Join(question.Replies, ", ")
		table.Append([]string{qIdStr, question.Text, replies, time})
	}
	// Render the table
	table.Render()
}
