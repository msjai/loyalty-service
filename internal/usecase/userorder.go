package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/msjai/loyalty-service/internal/entity"
	"github.com/msjai/loyalty-service/internal/usecase/repo"
)

// PostUserOrder -.
func (luc *LoyaltyUseCase) PostUserOrder(ctx context.Context, userOrder *entity.UserOrder) (*entity.UserOrder, error) {
	if !ValidOrderNumber(userOrder.Number) {
		return nil, fmt.Errorf("usecase - PostUserOrder - validNumber: %w", ErrInvalidOrderNumber)
	}

	userOrder, err := luc.repo.AddOrder(ctx, userOrder)
	if err != nil {
		if errors.Is(err, repo.ErrOrderNumExists) {
			// Если заказ с таким номером уже существует в базе,
			// то делаем повторный запрос, чтобы узнать ID пользователя, под которым был внесен заказ

			return nil, fmt.Errorf("usecase - PostRegUser - AddNewUser: %w", ErrLoginAlreadyTaken)
		}
		return nil, fmt.Errorf("usecase - PostRegUser - AddNewUser: %w", err)
	}

	return nil, nil
}

// GetUserOrders -.
func (luc *LoyaltyUseCase) GetUserOrders(context.Context, *entity.Loyalty) (*entity.Loyalty, error) {
	return nil, nil
}
