package repo

import (
	"database/sql"
	"errors"
)

var (
	ErrLoginAlreadyTaken = errors.New("login is already taken")
	ErrInvalidLogPass    = errors.New("invalid username/password pair")
	ErrConnectionNotOpen = errors.New("data base pgsql connection not opened")

	// ErrOrderNumExists - Если получили эту ошибку, то делаем повторный запрос,
	// чтобы выяснить ID пользователя, под которым заказ был уже зарегистрирован
	ErrOrderNumExists               = errors.New("order already exists")
	ErrOrderAlreadyRegByAnotherUser = errors.New("order already registered by another user")
	ErrOrderAlreadyRegByCurrUser    = errors.New("order already registered by current user")

	// ErrNoUserOders - это 204 ошибка, дляслучая когда нет данных ни по одному заказу пользователя
	ErrNoUserOdersRL = errors.New("no data to response")

	ErrUBalanceNotUpdAfterRegOrder = errors.New("after reg new order can't update user balance")
)

// LoyaltyRepoS -.
type LoyaltyRepoS struct {
	repo *sql.DB
}

// New -.
func New(db *sql.DB) *LoyaltyRepoS {
	return &LoyaltyRepoS{repo: db}
}
