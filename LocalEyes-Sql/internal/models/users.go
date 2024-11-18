package models

type User struct {
	UId              string      `json:"id"`
	Username         string      `json:"username"`
	Password         string      `json:"password"`
	City             string      `json:"city"`
	DwellingAge      float64     `json:"dwelling_age"`
	IsActive         bool        `json:"is_active"`
	Notification     []string    `json:"notification"`
	Tag              string      `json:"tag"`
	SecurityQuestion string      `json:"security_question"`
	SecurityAnswer   string      `json:"security_answer"`
	NotifyChannel    chan string `json:"-"` //ignore
}

type ResetPasswordUser struct {
	Username       string `json:"username"`
	NewPassword    string `json:"new_password"`
	SecurityAnswer string `json:"security_answer"`
}
