//go:build !test
// +build !test

package cli

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
		UIdStr := strconv.Itoa(user.UId)
		Dwelling := strconv.Itoa(user.DwellingAge)
		activeStatus := "No"
		if user.IsActive {
			activeStatus = "Yes"
		}
		table.Append([]string{UIdStr, user.Username, user.City, Dwelling, activeStatus, user.Tag})
	}

	// Render the table
	table.Render()
}

func displayPosts(posts []*models.Post) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"PostId", "Title", "Type", "Content", "Likes", "Created At"})

	// Add rows to the table, only including Name and City
	for _, post := range posts {
		PIdStr := strconv.Itoa(post.PostId)
		Likes := strconv.Itoa(post.Likes)
		Time := post.CreatedAt.Format("2006-01-02 15:04:05")
		table.Append([]string{PIdStr, post.Title, post.Type, post.Content, Likes, Time})
	}

	// Render the table
	table.Render()
}

func displayQuestions(questions []*models.Question) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"QID", "Question", "Replies", "Created At"})

	// Add rows to the table, only including Name and City
	for _, question := range questions {
		QIdStr := strconv.Itoa(question.QId)
		Time := question.CreatedAt.Format("2006-01-02 15:04:05")
		Replies := strings.Join(question.Replies, ", ")
		table.Append([]string{QIdStr, question.Text, Replies, Time})
	}
	// Render the table
	table.Render()
}
