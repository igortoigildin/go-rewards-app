package user

type User struct {
	ID       int64    `json:"id"`
	Login    string   `json:"login"`
	Password Password `json:"-"`
	Token    Token
}
