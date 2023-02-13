package usecase

import (
	"context"
	"fmt"

	"github.com/msjai/loyalty-service/internal/entity"
)

// CatchOrdersRefresh - Функция на отвечает за получение списка заказов, статусы по котором не окончательные.
// Функция обращается на уровень репозитория, чтобы получить информацию из базы
func (luc *LoyaltyUseCase) CatchOrdersRefresh(ctx context.Context) ([]*entity.UserOrder, error) {
	orders, err := luc.repo.CatchOrdersRefresh(ctx)
	if err != nil {
		return orders, fmt.Errorf("usecase - CatchOrdersRefresh - repo.CatchOrdersRefresh: %w", err)
	}

	return orders, nil
}

// RefreshOrderInfo - Функци обновляет ифно в базе по 1 заказу.
// Сначала получает информацию по 1 заказу в черному ящике, далее обновляет инфо в базе данных
func (luc *LoyaltyUseCase) RefreshOrderInfo(ctx context.Context, userOrder *entity.UserOrder) (*entity.UserOrder, error) {
	oldStatus := userOrder.Status
	userOrder, err := luc.webAPI.RefreshOrderInfo(ctx, userOrder)
	if err != nil {
		return userOrder, fmt.Errorf("usecase - RefreshOrderInfo - webAPI.RefreshOrderInfo: %w", err)
	}

	// Если статус обновился, обновляем запись в базе
	if oldStatus != userOrder.Status {
		userOrder, err = luc.RefreshOrderInfo(ctx, userOrder)
		if err != nil {
			return userOrder, fmt.Errorf("usecase - RefreshOrderInfo - RefreshOrderInfo: %w", err)
		}
	}

	// После получения инфо по заказу из черного ящика, необходимо обновить информацию по заказу в базе

	return userOrder, nil
}
