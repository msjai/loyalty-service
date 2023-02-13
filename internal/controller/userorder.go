package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

	// TODO Убрать лишние преобразования, номер заказа теперь string
	orderNumberS := string(b)
	// orderNumberI, err := strconv.Atoi(orderNumberS)
	// if err != nil {
	// 	http.Error(w, usecase.ErrInvalidOrderNumber.Error(), http.StatusUnprocessableEntity)
	// }
	// orderNumber := uint64(orderNumberI)

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
			//	routes.getOrderInfo(ctx, orderNumber)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", ApplicationJSON)
	w.WriteHeader(http.StatusAccepted)
	// Здесь идем в черный ящик, получаем инфо по заказу в системе начисления баллов
	//	routes.getOrderInfo(ctx, orderNumber)
}

func (routes *loyaltyRoutes) getOrderInfo(ctx context.Context, orderNumber uint64) {

	l := routes.cfg.L
	var userOrder entity.UserOrder

	request, err := http.NewRequestWithContext(ctx, http.MethodGet,
		"http://"+routes.cfg.AccrualSystemAddress+"/api/orders/"+fmt.Sprint(orderNumber), nil)
	if err != nil {
		l.Infof("controller - getOrderInfo - NewRequestWithContext: %v", err)
	}
	request.Header.Set("Accept", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		l.Infof("controller - getOrderInfo - DefaultClient.Do: %v", err)
	}
	defer response.Body.Close()

	b, err := io.ReadAll(response.Body)
	if err != nil {
		l.Infof("controller - getOrderInfo - io.ReadAll: %v", err)
	}

	err = json.Unmarshal(b, &userOrder)
	if err != nil {
		l.Infof("controller - getOrderInfo - json.Unmarshal: %v", err)
	}

	l.Info(userOrder)

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
