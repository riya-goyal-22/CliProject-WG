package models

type Response struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
}

type ResponseUser struct {
	UId         string `json:"id"`
	Username    string `json:"username"`
	City        string `json:"city"`
	LivingSince int    `json:"living_since"`
	Tag         string `json:"tag"`
}

type ResponseQuestion struct {
	QId       string   `json:"question_id"`
	PostId    string   `json:"post_id"`
	UserId    string   `json:"user_id"`
	Text      string   `json:"text"`
	Replies   []string `json:"replies"`
	CreatedAt string   `json:"created_at"`
}

type ResponsePost struct {
	PostId    string `json:"post_id"`
	UId       string `json:"user_id"`
	Title     string `json:"title"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	Likes     int    `json:"likes"`
	CreatedAt string `json:"created_at"`
}
