package utils_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"localEyes/tests/mocks"
	"localEyes/utils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPromptInput(t *testing.T) {
	// Simulate user input by redirecting os.Stdin
	input := "test input"
	r, w, _ := os.Pipe()
	_, _ = w.Write([]byte(input + "\n"))
	w.Close()

	// Backup the real os.Stdin and restore it after the test
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = r
	fmt.Println("entered prompt input")
	// Call the function to test
	result := utils.PromptInput("Enter input: ")
	fmt.Println("exit prompt input")
	// Check the result
	assert.Equal(t, input, result)
}

func TestGetChoice(t *testing.T) {
	// Simulate user input for the choice
	input := "3"
	r, w, _ := os.Pipe()
	_, _ = w.Write([]byte(input + "\n"))
	w.Close()

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = r

	result := utils.GetChoice()

	assert.Equal(t, 3, result)
}

func TestPromptIntInput_ValidInput(t *testing.T) {
	// Simulate valid user input by redirecting os.Stdin
	input := "42"
	r, w, _ := os.Pipe()
	_, _ = w.Write([]byte(input + "\n"))
	w.Close()

	// Backup the real os.Stdin and restore it after the test
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	os.Stdin = r
	fmt.Println("entered prompt input")
	// Call the function to test
	result, err := utils.PromptIntInput("Enter a number: ")
	fmt.Println("exited prompt input")
	// Check the result
	expected := 42
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestPromptPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock instance
	mockPrompt := mocks.NewMockPromptInterface(ctrl)

	// Set up expectations
	mockPrompt.EXPECT().Run().Return("testpassword", nil)

	// Call the function under test
	result := utils.PromptPassword(mockPrompt)

	// Assert the results
	if result != "testpassword" {
		t.Errorf("expected testpassword, got %s", result)
	}
}
