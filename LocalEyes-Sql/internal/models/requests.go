package models

type LivingSince struct {
	Days   float64 `json:"days"`
	Months float64 `json:"months"`
	Years  float64 `json:"years"`
}

type Client struct {
	Username    string      `json:"username"`
	Password    string      `json:"password"`
	City        string      `json:"city"`
	LivingSince LivingSince `json:"living_since"`
	Answer      string      `json:"security_answer"`
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
