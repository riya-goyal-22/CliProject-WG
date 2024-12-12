package models

type User struct {
	UId          string   `json:"id"`
	Email        string   `json:"email"`
	Username     string   `json:"username"`
	Password     string   `json:"password"`
	City         string   `json:"city"`
	DwellingAge  float64  `json:"dwelling_age"`
	IsActive     bool     `json:"is_active"`
	Notification []string `json:"notification"`
	Tag          string   `json:"tag"`
}

type UserEmail struct {
	Email string `json:"email"`
}

type ResetPasswordUser struct {
	Email       string `json:"email"`
	OTP         string `json:"otp"`
	NewPassword string `json:"new_password"`
}
