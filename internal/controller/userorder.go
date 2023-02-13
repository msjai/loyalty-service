package controller

import (
	"errors"
	"io"
	"net/http"

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
			// go routes.refreshOrdersInfo()
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", ApplicationJSON)
	w.WriteHeader(http.StatusAccepted)
	// Здесь идем в черный ящик, получаем инфо по заказу в системе начисления баллов
	// go routes.refreshOrdersInfo()
}

// refreshOrdersInfo - Функция инициирует обновление информации по заказам, статусы по которым не окончательные.
// Функция обращается к уровню usecase.
// Далее по каждому заказу из списка инициируется обновление статуса
func (routes *loyaltyRoutes) refreshOrdersInfo() {
	// Воркер работает в цикле,пока не получит сигнал заснуть
	// Или пока не получит сигнал остановиться
	for {
		select {
		case <-routes.cfg.Done:
			// Если требуется завершить работу воркера, послыаем данные в канал
			return
		default:
			// Читаем из закрытого канала.
			// Если требуется остановить, то открываем пустой канал (операция становится блокирующей)
			<-routes.cfg.Sig

			l := routes.cfg.L

			orders, err := routes.loyalty.CatchOrdersRefresh()
			if err != nil {
				l.Errorf("controller - refreshOrdersInfo - CatchOrdersRefresh: %w", err)
			}

			for _, order := range orders {
				_, err = routes.loyalty.RefreshOrderInfo(order)
				if err != nil {
					l.Errorf("controller - refreshOrdersInfo - loyalty.RefreshOrderInfo: %w", err)
				}
			}
		}
	}

	// l := routes.cfg.L
	//
	// orders, err := routes.loyalty.CatchOrdersRefresh()
	// if err != nil {
	// 	l.Errorf("controller - refreshOrdersInfo - CatchOrdersRefresh: %w", err)
	// }
	//
	// for _, order := range orders {
	// 	_, err = routes.loyalty.RefreshOrderInfo(order)
	// 	if err != nil {
	// 		l.Errorf("controller - refreshOrdersInfo - loyalty.RefreshOrderInfo: %w", err)
	// 	}
	// }
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
