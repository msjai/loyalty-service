package controller

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/msjai/loyalty-service/internal/controller/middleware"
	"github.com/msjai/loyalty-service/internal/entity"
	"github.com/msjai/loyalty-service/internal/usecase"
)

const orderAlreadyRegByCurrentU = "order already registered by current user"

// PostUOrder -.
func (routes *loyaltyRoutes) PostUOrder(w http.ResponseWriter, r *http.Request) {
	// var UserOrder *entity.UserOrder
	// Через контекст получаем reader
	// В случае необходимости тело было распаковано в middleware
	// Далее передаем этот же контекст в UseCase
	ctx := r.Context()
	reader := ctx.Value(middleware.KeyReader).(io.Reader)
	userID := ctx.Value(middleware.KeyUserID).(int64)

	b, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderNumberS := string(b)
	_, err = routes.loyalty.PostUserOrder(ctx, &entity.UserOrder{
		UserID: userID,
		Number: orderNumberS,
	})

	if err != nil {
		if errors.Is(err, usecase.ErrInvalidOrderNumber) {
			http.Error(w, usecase.ErrInvalidOrderNumber.Error(), http.StatusUnprocessableEntity)
			return
		}

		if errors.Is(err, usecase.ErrOrderAlreadyRegByAnotherUser) {
			http.Error(w, usecase.ErrOrderAlreadyRegByAnotherUser.Error(), http.StatusConflict)
			return
		}

		if errors.Is(err, usecase.ErrOrderAlreadyRegByCurrUser) {
			w.Header().Set("Content-Type", ApplicationJSON)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(orderAlreadyRegByCurrentU)) //nolint:errcheck
			// Здесь идем в черный ящик, получаем инфо по заказу в системе начисления баллов
			go routes.refreshOrdersInfo(ctx)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", ApplicationJSON)
	w.WriteHeader(http.StatusAccepted)
	// Здесь идем в черный ящик, получаем инфо по заказу в системе начисления баллов
	go routes.refreshOrdersInfo(ctx)
}

// refreshOrdersInfo - Функция инициирует обновление информации по заказам, статусы по которым не окончательные.
// Функция обращается к уровню usecase.
// Далее по каждому заказу из списка инициируется обновление статуса
func (routes *loyaltyRoutes) refreshOrdersInfo(ctx context.Context) {
	l := routes.cfg.L
	ctxRefresh, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	orders, err := routes.loyalty.CatchOrdersRefresh(ctxRefresh)
	if err != nil {
		l.Errorf("repo - CatchOrdersRefresh - repo.Begin: %w", err)
	}

	for _, order := range orders {
		routes.loyalty.RefreshOrderInfo(ctxRefresh, order)
	}

}

// GerUOrders -.
func (routes *loyaltyRoutes) GerUOrders(w http.ResponseWriter, r *http.Request) {
	// UserOrders := entity.Loyalty{
	// 	UserID: "1",
	// 	UserOrders: []entity.UserOrder{
	// 		{Number: "1", Status: entity.NEW, Accrual: 100, UploadedAt: time.Now()},
	// 		{Number: "2", Status: entity.NEW, Accrual: 100, UploadedAt: time.Now()},
	// 		{Number: "3", Status: entity.NEW, Accrual: 100, UploadedAt: time.Now()},
	// 	},
	// }
	//
	// userorders, _ := json.Marshal(UserOrders)
	//
	// w.Header().Set("Content-Type", ApplicationJSON)
	// w.WriteHeader(http.StatusOK)
	// w.Write(userorders)
}
