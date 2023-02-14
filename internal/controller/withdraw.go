package controller

import (
	"net/http"
)

// PostUOrder -.
func (routes *loyaltyRoutes) PostUWithdraw(w http.ResponseWriter, r *http.Request) {
	// var UserOrder *entity.UserOrder
	// Через контекст получаем reader
	// В случае необходимости тело было распаковано в middleware
	// Далее передаем этот же контекст в UseCase
	// ctx := r.Context()
	// reader := ctx.Value(middleware.KeyReader).(io.Reader)
	// userID := ctx.Value(middleware.KeyUserID).(int64)
	//
	// b, err := io.ReadAll(reader)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

	// orderNumberS := string(b)
	// _, err = routes.loyalty.PostUserOrder(&entity.UserOrder{
	// 	UserID: userID,
	// 	Number: orderNumberS,
	// })
	//
	// if err != nil {
	// 	if errors.Is(err, usecase.ErrInvalidOrderNumber) {
	// 		http.Error(w, usecase.ErrInvalidOrderNumber.Error(), http.StatusUnprocessableEntity)
	// 		return
	// 	}
	//
	// 	if errors.Is(err, usecase.ErrOrderAlreadyRegByAnotherUser) {
	// 		http.Error(w, usecase.ErrOrderAlreadyRegByAnotherUser.Error(), http.StatusConflict)
	// 		return
	// 	}
	//
	// 	if errors.Is(err, usecase.ErrOrderAlreadyRegByCurrUser) {
	// 		w.Header().Set("Content-Type", ApplicationJSON)
	// 		w.WriteHeader(http.StatusOK)
	// 		w.Write([]byte(orderAlreadyRegByCurrentU)) //nolint:errcheck
	// 		return
	// 	}
	//
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	//
	// w.Header().Set("Content-Type", ApplicationJSON)
	// w.WriteHeader(http.StatusAccepted)
}
