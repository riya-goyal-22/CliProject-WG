//go:build !test
// +build !test

package cli

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"localEyes/constants"
	"localEyes/internal/services"
	"localEyes/utils"
	"log"
	"strconv"
	"strings"
)

func signUp(userService *services.UserService) {
	fmt.Println(constants.Blue + "==============================")
	fmt.Println("SIGN UP")
	fmt.Println("==============================" + constants.Reset)
	var tag, username, password string
	for {
		username = utils.PromptInput("Enter your username:")
		if utils.ValidateUsername(username, userService.Repo) {
			break
		} else {
			fmt.Println(constants.Red + "Username already taken" + constants.Reset)
		}
	}
	for {
		prompt := &promptui.Prompt{
			Label:     constants.Cyan + "Enter a strong password [6 characters long ,having special character and number]" + constants.Reset,
			Mask:      '*',
			IsConfirm: false,
		}
		password = utils.PromptPassword(prompt)
		if utils.ValidatePassword(password) {
			break
		} else {
			fmt.Println(constants.Red + "Password is weak" + constants.Reset)
		}
	}
	if city := utils.PromptInput("Enter your city:"); strings.ToLower(city) != "delhi" {
		err := errors.New(constants.Red + "You are not a vaid user for this application" + constants.Reset)
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
		fmt.Println(constants.Red + "Error Signing Up\n" + err.Error() + constants.Reset)
		return
	} else {
		fmt.Println("\n" + constants.Green + "Successfully Signed Up!\n" + constants.Reset)
	}
}
