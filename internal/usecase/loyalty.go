package usecase

import (
	"context"

	"github.com/msjai/loyalty-service/internal/config"
	"github.com/msjai/loyalty-service/internal/entity"
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

// PostRegUser -.
func (luc *LoyaltyUseCase) PostRegUser(context.Context, *entity.Loyalty) (*entity.Loyalty, error) {
	return nil, nil
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
