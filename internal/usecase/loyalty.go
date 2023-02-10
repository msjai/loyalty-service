package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/msjai/loyalty-service/internal/config"
	"github.com/msjai/loyalty-service/internal/entity"
	"github.com/msjai/loyalty-service/internal/usecase/repo"
)

var ErrLoginAlreadyTaken = errors.New("login is already taken")

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

// PostRegUser -.
func (luc *LoyaltyUseCase) PostRegUser(ctx context.Context, loyalty *entity.Loyalty) (*entity.Loyalty, error) {
	loyalty, err := luc.repo.AddNewUser(ctx, loyalty)
	if err != nil {
		if errors.Is(err, repo.ErrLoginAlreadyTaken) {
			return nil, fmt.Errorf("usecase - PostRegUser - AddNewUser: %w", ErrLoginAlreadyTaken)
		}
		return nil, fmt.Errorf("usecase - PostRegUser - AddNewUser: %w", err)
	}

	return loyalty, nil
}

// PostLoginUser -.
func (luc *LoyaltyUseCase) PostLoginUser(context.Context, *entity.Loyalty) (*entity.Loyalty, error) {
	return nil, nil
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
