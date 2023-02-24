package usecase

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/msjai/loyalty-service/internal/entity"
	"github.com/msjai/loyalty-service/internal/usecase/repo"
)

// PostUserWithDrawBalance -.
func (luc *LoyaltyUseCase) PostUserWithDrawBalance(withDraw *entity.WithDraw) (*entity.WithDraw, error) {
	uintNumber, _ := strconv.ParseUint(withDraw.Number, 10, 64)
	if !ValidOrderNumber(uintNumber) {
		return withDraw, fmt.Errorf("usecase - PostUserWithDrawBalance - validNumber: %w", ErrInvalidOrderNumber)
	}

	withDraw.Sum *= 100

	withDraw, err := luc.repo.WithDraw(withDraw)
	if err != nil {
		if errors.Is(err, repo.ErrInsufficientFund) {
			return withDraw, fmt.Errorf("usecase - PostUserWithDrawBalance - repo.WithDraw: %w", ErrInsufficientFund)
		}

		return withDraw, fmt.Errorf("usecase - PostUserWithDrawBalance - repo.WithDraw: %w", err)
	}

	return withDraw, nil
}

// GetUserWithdrawals -.
func (luc *LoyaltyUseCase) GetUserWithdrawals(user *entity.User) ([]*entity.WithDraw, error) {
	userWithdrawals, err := luc.repo.GetUserWithdrawals(user)
	if err != nil {
		if errors.Is(err, repo.ErrNoUserWithdrawRL) {
			return userWithdrawals, fmt.Errorf("usecase - GetUserWithdrawals - repo.GetUserWithdrawals: %w", ErrNoUserWithdrawUCL)
		}
		return userWithdrawals, fmt.Errorf("usecase - GetUserWithdrawals - repo.GetUserWithdrawals: %w", err)
	}

	return userWithdrawals, nil
}
