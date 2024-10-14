//go:build !test
// +build !test

package ui

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"localEyes/config"
	"localEyes/internal/services"
	"localEyes/utils"
	"log"
	"strconv"
	"strings"
)

func signUp(userService *services.UserService) {
	fmt.Println(config.Blue + "==============================")
	fmt.Println("SIGN UP")
	fmt.Println("==============================" + config.Reset)
	var tag, username, password string
	for {
		username = utils.PromptInput("Enter your username:")
		if utils.ValidateUsername(username, userService.Repo) {
			break
		} else {
			fmt.Println(config.Red + "Username already taken" + config.Reset)
		}
	}
	for {
		prompt := &promptui.Prompt{
			Label:     config.Cyan + "Enter a strong password [6 characters long ,having special character and number]" + config.Reset,
			Mask:      '*',
			IsConfirm: false,
		}
		password = utils.PromptPassword(prompt)
		if utils.ValidatePassword(password) {
			break
		} else {
			fmt.Println(config.Red + "Password is weak" + config.Reset)
		}
	}
	if city := utils.PromptInput("Enter your city:"); strings.ToLower(city) != "delhi" {
		err := errors.New(config.Red + "You are not a vaid user for this application" + config.Reset)
		log.Fatal(err)
	}
	DwellingAge, _ := strconv.Atoi(utils.PromptInput("For how many years you are living here/lived here:"))
	if DwellingAge > 2 {
		tag = "resident"
	} else {
		tag = "newbie"
	}
	err := userService.Signup(username, password, DwellingAge, tag)
	if err != nil {
		fmt.Println(config.Red + "Error Signing Up\n" + err.Error() + config.Reset)
		return
	} else {
		fmt.Println("\n" + config.Green + "Successfully Signed Up!\n" + config.Reset)
	}
}
