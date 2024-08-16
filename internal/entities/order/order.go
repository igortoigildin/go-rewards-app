package order

import "time"

type Order struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *int      `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"date"`
	UserID     int64     `json:"user_id,omitempty"`
}
