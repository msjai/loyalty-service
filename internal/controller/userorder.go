package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/msjai/loyalty-service/internal/controller/middleware"
	"github.com/msjai/loyalty-service/internal/entity"
	"github.com/msjai/loyalty-service/internal/usecase"
)

const orderAlreadyRegByCurrentU = "order already registered by current user"

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
}

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
	_, err = routes.loyalty.PostUserOrder(&entity.UserOrder{
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
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", ApplicationJSON)
	w.WriteHeader(http.StatusAccepted)
}

// GerUOrders -.
func (routes *loyaltyRoutes) GerUOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(middleware.KeyUserID).(int64)

	userOrders, err := routes.loyalty.GetUserOrders(&entity.User{ID: userID})
	if err != nil {
		if errors.Is(err, usecase.ErrNoUserOdersUCL) {
			http.Error(w, usecase.ErrNoUserOdersUCL.Error(), http.StatusNoContent)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(userOrders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", ApplicationJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(response) //nolint:errcheck
}

func (routes *loyaltyRoutes) GetUBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ctx.Value(middleware.KeyUserID).(int64)

	User, err := routes.loyalty.GetUserBalance(&entity.User{ID: userID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", ApplicationJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(response) //nolint:errcheck
}
