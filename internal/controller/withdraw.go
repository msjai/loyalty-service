package controller

import (
	"net/http"
	"time"

	"github.com/msjai/loyalty-service/internal/entity"
)

const withDrawSucces = "withdraw success"

// clearWhithDrawFields -.
func clearWhithDrawFields(withDraw *entity.WithDraw) *entity.WithDraw {
	withDraw.ID = 0
	withDraw.UserID = 0
	withDraw.ProcessedAt = time.Time{}
	return withDraw
}

// PostUWithdraw -.
func (routes *loyaltyRoutes) PostUWithdraw(w http.ResponseWriter, r *http.Request) {
	// var withDraw = &entity.WithDraw{}
	// // Через контекст получаем reader
	// // В случае необходимости тело было распаковано в middleware
	// // Далее передаем этот же контекст в UseCase
	// ctx := r.Context()
	// reader := ctx.Value(middleware.KeyReader).(io.Reader)
	// userID := ctx.Value(middleware.KeyUserID).(int64)
	//
	// b, err := io.ReadAll(reader)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	// 	return
	// }
	//
	// err = json.Unmarshal(b, withDraw)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	// 	return
	// }
	//
	// withDraw = clearWhithDrawFields(withDraw)
	// withDraw.UserID = userID
	//
	// _, err = routes.loyalty.PostUserWithDrawBalance(withDraw)
	// if err != nil {
	// 	if errors.Is(err, usecase.ErrInvalidOrderNumber) {
	// 		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
	// 		return
	// 	}
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	//
	// w.Header().Set("Content-Type", ApplicationJSON)
	// w.WriteHeader(http.StatusOK)
	// w.Write([]byte(withDrawSucces)) //nolint:errcheck
}
