package repo

import (
	"database/sql"
	"errors"
)

var (
	ErrLoginAlreadyTaken = errors.New("login is already taken")
	ErrInvalidLogPass    = errors.New("invalid username/password pair")
	ErrConnectionNotOpen = errors.New("data base pgsql connection not opened")
	ErrOrderNumExists    = errors.New("order already exists")
)

// LoyaltyRepoS -.
type LoyaltyRepoS struct {
	repo *sql.DB
}

// New -.
func New(db *sql.DB) *LoyaltyRepoS {
	return &LoyaltyRepoS{repo: db}
}
