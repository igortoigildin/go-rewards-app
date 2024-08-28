package withdrawal

import "time"

type Withdrawal struct {
	Order  string    `json:"order"`
	Sum    float64   `json:"sum"`
	Date   time.Time `json:"date"`
	UserID int64     `json:"-"`
}
