package user

type User struct {
	UserID   int64    `json:"user_id"`
	Login    string   `json:"login"`
	Password Password `json:"-"`
	Balance  *float64 `json:"balance"`
	Token    Token
}
