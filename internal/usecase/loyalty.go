package usecase

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/msjai/loyalty-service/internal/config"
)

var (
	ErrLoginAlreadyTaken    = errors.New("login is already taken")
	ErrInvalidLogPass       = errors.New("invalid username/password pair")
	ErrInvalidSigningMethod = errors.New("invalid signing method")

	ErrInvalidOrderNumber           = errors.New("invalid order number format")
	ErrOrderAlreadyRegByAnotherUser = errors.New("order already registered by another user")
	ErrOrderAlreadyRegByCurrUser    = errors.New("order already registered by current user")

	// ErrNoUserOdersUCL - это 204 ошибка, дляслучая когда нет данных ни по одному заказу пользователя
	ErrNoUserOdersUCL = errors.New("no data to response")
	// ErrNoUserWithdraw - это 204 ошибка, дляслучая когда нет данных ни по одному списанию пользователя
	ErrNoUserWithdrawUCL = errors.New("no data to response")

	ErrInsufficientFund = errors.New("insufficient funds to withdraw")
)

const (
	salt         = "dkY6#dgb&jdg"
	signTokenKey = "frefrefn@fW#csafssdfs" //nolint:gosec
	tokenTL      = 24 * time.Hour
)

// LoyaltyUseCase -.
type LoyaltyUseCase struct {
	repo   LoyaltyRepo
	webAPI LoyaltyWebAPI
	cfg    *config.Config
}

// New -.
func New(r LoyaltyRepo, w LoyaltyWebAPI, c *config.Config) *LoyaltyUseCase {
	return &LoyaltyUseCase{
		repo:   r,
		webAPI: w,
		cfg:    c,
	}
}

// TODO Перенести функцию на уровень репозитория
func hashPassword(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass))
	dst := h.Sum([]byte(salt))

	return fmt.Sprintf("%x", dst)
}
