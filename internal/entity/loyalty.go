package entity

import (
	"time"
)

const (
	NEW        = "NEW"
	PROCESSING = "PROCESSING"
	INVALID    = "INVALID"
	PROCESSED  = "PROCESSED"
)

// Loyalty .-
type Loyalty struct {
	User       *User       `json:"user"`
	UserOrders []UserOrder `json:"user_orders"` //nolint:tagliatelle
}

// UserOrder .-
type UserOrder struct {
	ID         int64     `json:"id"`
	Number     string    `json:"number" valid:"required"`
	Status     string    `json:"status"`
	UserID     int64     `json:"user_id"`           //nolint:tagliatelle
	AccrualSum float64   `json:"accrual,omitempty"` //nolint:tagliatelle
	UploadedAt time.Time `json:"uploaded_at"`       //nolint:tagliatelle
}

// User .-
type User struct {
	ID       int64   `json:"id"`
	Login    string  `json:"login" valid:"required"`
	Password string  `json:"password" valid:"required"`
	Balance  float64 `json:"balance"`
	Token    string  `json:"token"`
}
