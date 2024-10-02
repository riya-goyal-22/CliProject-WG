package models

type LivingSince struct {
	Days   int `json:"days"`
	Months int `json:"months"`
	Years  int `json:"years"`
}

type Client struct {
	Username    string      `json:"username"`
	Password    string      `json:"password"`
	City        string      `json:"city"`
	LivingSince LivingSince `json:"living_since"`
}

type RequestPost struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

type RequestQuestion struct {
	Question string `json:"question"`
}

type RequestAnswer struct {
	Answer string `json:"answer"`
}
