package withdrawal

import "time"

type Withdrawal struct {
	Order  string    `json:"order"`
	Sum    int       `json:"sum"`
	Date   time.Time `json:"date"`
	UserID int64     `json:"user_id"`
}
