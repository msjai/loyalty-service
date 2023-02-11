package usecase

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/msjai/loyalty-service/internal/config"
	"github.com/msjai/loyalty-service/internal/entity"
	"github.com/msjai/loyalty-service/internal/usecase/repo"
)

var (
	ErrLoginAlreadyTaken    = errors.New("login is already taken")
	ErrInvalidLogPass       = errors.New("invalid username/password pair")
	ErrInvalidSigningMethod = errors.New("invalid signing method")
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

// PostRegUser -.
func (luc *LoyaltyUseCase) PostRegUser(ctx context.Context, loyalty *entity.Loyalty) (*entity.Loyalty, error) {
	// TODO Перенести вызов функции на уровень репозитория
	loyalty.User.Password = hashPassword(loyalty.User.Password)

	loyalty, err := luc.repo.AddNewUser(ctx, loyalty)
	if err != nil {
		if errors.Is(err, repo.ErrLoginAlreadyTaken) {
			return nil, fmt.Errorf("usecase - PostRegUser - AddNewUser: %w", ErrLoginAlreadyTaken)
		}
		return nil, fmt.Errorf("usecase - PostRegUser - AddNewUser: %w", err)
	}

	loyalty.User.Token, err = getToken(loyalty)
	if err != nil {
		return nil, fmt.Errorf("usecase - PostRegUser - getToken: %w", err)
	}

	return loyalty, nil
}

// PostLoginUser -.
func (luc *LoyaltyUseCase) PostLoginUser(ctx context.Context, loyalty *entity.Loyalty) (*entity.Loyalty, error) {
	loyalty.User.Password = hashPassword(loyalty.User.Password)

	loyalty, err := luc.repo.FindUser(ctx, loyalty)
	if err != nil {
		if errors.Is(err, repo.ErrInvalidLogPass) {
			return nil, fmt.Errorf("usecase - PostLoginUser - FindUser: %w", ErrInvalidLogPass)
		}
		return nil, fmt.Errorf("usecase - PostLoginUser - FindUser: %w", err)
	}

	loyalty.User.Token, err = getToken(loyalty)
	if err != nil {
		return nil, fmt.Errorf("usecase - PostLoginUser - getToken: %w", err)
	}

	return loyalty, nil
}

// PostUserOrder -.
func (luc *LoyaltyUseCase) PostUserOrder(context.Context, *entity.Loyalty) (*entity.Loyalty, error) {
	return nil, nil
}

// GetUserOrders -.
func (luc *LoyaltyUseCase) GetUserOrders(context.Context, *entity.Loyalty) (*entity.Loyalty, error) {
	return nil, nil
}

// GetUserBalance -.
func (luc *LoyaltyUseCase) GetUserBalance(context.Context, *entity.Loyalty) (*entity.Loyalty, error) {
	return nil, nil
}

// PostUserWithDrawBalance -.
func (luc *LoyaltyUseCase) PostUserWithDrawBalance(context.Context, *entity.Loyalty) (*entity.Loyalty, error) {
	return nil, nil
}

// GetUserWithdrawals -.
func (luc *LoyaltyUseCase) GetUserWithdrawals(context.Context, *entity.Loyalty) (*entity.Loyalty, error) {
	return nil, nil
}
