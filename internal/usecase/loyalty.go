package usecase

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/msjai/loyalty-service/internal/config"
	"github.com/msjai/loyalty-service/internal/entity"
)

var (
	ErrLoginAlreadyTaken    = errors.New("login is already taken")
	ErrInvalidLogPass       = errors.New("invalid username/password pair")
	ErrInvalidSigningMethod = errors.New("invalid signing method")

	ErrInvalidOrderNumber           = errors.New("invalid order number format")
	ErrOrderAlreadyRegByAnotherUser = errors.New("order already registered by another user")
	ErrOrderAlreadyRegByCurrUser    = errors.New("order already registered by current user")

	// ErrNoUserOders - это 204 ошибка, дляслучая когда нет данных ни по одному заказу пользователя
	ErrNoUserOdersUCL = errors.New("no data to response")
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

// PostUserWithDrawBalance -.
func (luc *LoyaltyUseCase) PostUserWithDrawBalance(*entity.Loyalty) (*entity.Loyalty, error) {
	return nil, nil
}

// GetUserWithdrawals -.
func (luc *LoyaltyUseCase) GetUserWithdrawals(*entity.Loyalty) (*entity.Loyalty, error) {
	return nil, nil
}
