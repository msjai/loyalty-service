package usecase

import (
	"fmt"

	"github.com/msjai/loyalty-service/internal/entity"
)

// CatchOrdersRefresh - Функция на отвечает за получение списка заказов, статусы по котором не окончательные.
// Функция обращается на уровень репозитория, чтобы получить информацию из базы
func (luc *LoyaltyUseCase) CatchOrdersRefresh() ([]*entity.UserOrder, error) {
	orders, err := luc.repo.CatchOrdersRefresh()
	if err != nil {
		return orders, fmt.Errorf("usecase - CatchOrdersRefresh - repo.CatchOrdersRefresh: %w", err)
	}

	return orders, nil
}

// RefreshOrderInfo - Функци обновляет ифно в базе по 1 заказу.
// Сначала получает информацию по 1 заказу в черному ящике, далее обновляет инфо в базе данных
func (luc *LoyaltyUseCase) RefreshOrderInfo(userOrder *entity.UserOrder) (*entity.UserOrder, error) {
	oldStatus := userOrder.Status
	userOrder, err := luc.webAPI.RefreshOrderInfo(userOrder)
	if err != nil {
		return userOrder, fmt.Errorf("usecase - RefreshOrderInfo - webAPI.RefreshOrderInfo: %w", err)
	}

	// Если статус обновился, обновляем запись в базе
	if oldStatus != userOrder.Status {
		// userOrder, err = luc.RefreshOrderInfo(userOrder)
		userOrder, err = luc.repo.UpdateOrder(userOrder)
		if err != nil {
			return userOrder, fmt.Errorf("usecase - RefreshOrderInfo - repo.UpdateOrder: %w", err)
		}
	}

	return userOrder, nil
}
