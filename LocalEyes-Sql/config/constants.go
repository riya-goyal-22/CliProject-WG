package config

const (
	UserTable            = "users"
	PostTable            = "posts"
	QuestionTable        = "questions"
	Select               = "SELECT %s FROM %s"
	SelectWithCondition  = "SELECT %s FROM %s WHERE %s = ?"
	SelectWith2Condition = "SELECT %s FROM %s WHERE %s = ? AND %s = ?"
	Insert               = "INSERT INTO %s (%s) VALUES (%s)"
	Update               = "UPDATE %s SET %s WHERE %s = ?"
	UpdateWith2Condition = "UPDATE %s SET %s WHERE %s = ? AND %s = ?"
	Delete               = "DELETE FROM %s WHERE %s = ?"
	DeleteWith2Condition = "DELETE FROM %s WHERE %s= ? AND %s= ?"
)
