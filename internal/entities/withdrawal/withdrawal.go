package withdrawal

import "time"

type Withdrawal struct {
	Id     string    `json:"id"`
	Sum    int       `json:"sum"`
	Date   time.Time `json:"date"`
	UserID int64     `json:"user_id"`
}
