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

// User .-
type User struct {
	ID        int64   `json:"id"`
	Login     string  `json:"login" valid:"required"`
	Password  string  `json:"password" valid:"required"`
	Balance   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
	Token     string  `json:"token"`
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

type WithDraw struct {
	ID          int64     `json:"id"`
	Number      string    `json:"number" valid:"required"`
	Sum         float64   `json:"sum" valid:"required"` //nolint:tagliatelle
	UserID      int64     `json:"user_id"`              //nolint:tagliatelle
	ProcessedAt time.Time `json:"processed_at"`         //nolint:tagliatelle
}
