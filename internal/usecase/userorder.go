package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/msjai/loyalty-service/internal/entity"
	"github.com/msjai/loyalty-service/internal/usecase/repo"
)

// PostUserOrder -.
func (luc *LoyaltyUseCase) PostUserOrder(ctx context.Context, userOrder *entity.UserOrder) (*entity.UserOrder, error) {
	uintNumber, _ := strconv.ParseUint(userOrder.Number, 10, 64)
	if !ValidOrderNumber(uintNumber) {
		return nil, fmt.Errorf("usecase - PostUserOrder - validNumber: %w", ErrInvalidOrderNumber)
	}

	userOrder, err := luc.repo.AddOrder(ctx, userOrder)
	// var userIDInDB int64
	if err != nil {
		// Если заказ с таким номером уже существует в базе,
		// то делаем повторный запрос, чтобы узнать ID пользователя, под которым был внесен заказ
		if errors.Is(err, repo.ErrOrderNumExists) {
			userOrder, err = luc.repo.FindOrder(ctx, userOrder)
			if err != nil {
				if errors.Is(err, repo.ErrOrderAlreadyRegByAnotherUser) {
					return nil, fmt.Errorf("usecase - PostUserOrder - FindOrder: %w", ErrOrderAlreadyRegByAnotherUser)
				}
				if errors.Is(err, repo.ErrOrderAlreadyRegByCurrUser) {
					return nil, fmt.Errorf("usecase - PostUserOrder - FindOrder: %w", ErrOrderAlreadyRegByCurrUser)
				}
				return nil, fmt.Errorf("usecase - PostUserOrder - FindOrder: %w", err)
			}
			return userOrder, nil
		}
		return nil, fmt.Errorf("usecase - PostUserOrder - AddOrder: %w", err)
	}

	return userOrder, nil
}

// GetUserOrders -.
func (luc *LoyaltyUseCase) GetUserOrders(context.Context, *entity.Loyalty) (*entity.Loyalty, error) {
	return nil, nil
}
