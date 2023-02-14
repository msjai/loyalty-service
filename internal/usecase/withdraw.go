package usecase

import (
	"fmt"
	"strconv"

	"github.com/msjai/loyalty-service/internal/entity"
)

// PostUserWithDrawBalance -.
func (luc *LoyaltyUseCase) PostUserWithDrawBalance(withDraw *entity.WithDraw) (*entity.WithDraw, error) {
	uintNumber, _ := strconv.ParseUint(withDraw.Number, 10, 64)
	if !ValidOrderNumber(uintNumber) {
		return withDraw, fmt.Errorf("usecase - PostUserWithDrawBalance - validNumber: %w", ErrInvalidOrderNumber)
	}

	withDraw, err := luc.repo.WithDraw(withDraw)
	if err != nil {

	}

	return withDraw, nil
}
