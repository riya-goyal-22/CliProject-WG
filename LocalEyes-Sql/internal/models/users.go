package models

type User struct {
	UId           string      `json:"id"`
	Username      string      `json:"username"`
	Password      string      `json:"password"`
	City          string      `json:"city"`
	DwellingAge   int         `json:"dwelling_age"`
	IsActive      bool        `json:"is_active"`
	Notification  []string    `json:"notification"`
	Tag           string      `json:"tag"`
	NotifyChannel chan string `json:"-"` //ignore
}
