package usecase

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/msjai/loyalty-service/internal/entity"
	"github.com/msjai/loyalty-service/internal/usecase/repo"
)

// PostUserOrder -.
func (luc *LoyaltyUseCase) PostUserOrder(userOrder *entity.UserOrder) (*entity.UserOrder, error) {
	uintNumber, _ := strconv.ParseUint(userOrder.Number, 10, 64)
	if !ValidOrderNumber(uintNumber) {
		return nil, fmt.Errorf("usecase - PostUserOrder - validNumber: %w", ErrInvalidOrderNumber)
	}

	userOrder, err := luc.repo.AddOrder(userOrder)
	if err != nil {
		// Если заказ с таким номером уже существует в базе,
		// то делаем повторный запрос, чтобы узнать ID пользователя, под которым был внесен заказ
		if errors.Is(err, repo.ErrOrderNumExists) {
			userOrder, err = luc.repo.FindOrder(userOrder)
			if err != nil {
				if errors.Is(err, repo.ErrOrderAlreadyRegByAnotherUser) {
					return userOrder, fmt.Errorf("usecase - PostUserOrder - FindOrder: %w", ErrOrderAlreadyRegByAnotherUser)
				}
				if errors.Is(err, repo.ErrOrderAlreadyRegByCurrUser) {
					return userOrder, fmt.Errorf("usecase - PostUserOrder - FindOrder: %w", ErrOrderAlreadyRegByCurrUser)
				}
				return userOrder, fmt.Errorf("usecase - PostUserOrder - FindOrder: %w", err)
			}
			return userOrder, nil
		}
		return userOrder, fmt.Errorf("usecase - PostUserOrder - AddOrder: %w", err)
	}

	return userOrder, nil
}

// GetUserOrders -.
func (luc *LoyaltyUseCase) GetUserOrders(user *entity.User) ([]*entity.UserOrder, error) {
	userOrders, err := luc.repo.FindOrders(user)
	if err != nil {
		if errors.Is(err, repo.ErrNoUserOdersRL) {
			return userOrders, fmt.Errorf("usecase - GetUserOrders - FindOrders: %w", ErrNoUserOdersUCL)
		}
		return userOrders, fmt.Errorf("usecase - GetUserOrders - FindOrders: %w", err)
	}

	return userOrders, nil
}

// GetUserBalance -.
func (luc *LoyaltyUseCase) GetUserBalance(user *entity.User) (*entity.UserBalance, error) {
	balance := &entity.UserBalance{}

	user, err := luc.repo.GetUserBalance(user)
	if err != nil {
		return balance, fmt.Errorf("usecase - GetUserBalance - repo.GetUserBalance: %w", err)
	}

	// Здесь делим суммы на 100, потому что в базе храним в копейках
	balance.Current = user.Balance / 100
	balance.Withdrawn = user.Withdrawn / 100

	return balance, nil
}
