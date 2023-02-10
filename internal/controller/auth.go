package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"

	"github.com/msjai/loyalty-service/internal/controller/middleware"
	"github.com/msjai/loyalty-service/internal/entity"
	"github.com/msjai/loyalty-service/internal/usecase"
)

type UserID struct {
	ID string `json:"id"`
}

// PostRegUHandler -.
func (routes *loyaltyRoutes) PostRegUHandler(w http.ResponseWriter, r *http.Request) {
	var User entity.User

	// Через контекст получаем reader
	// В случае необхоимости тело было распаковано в middleware
	// Далее передаем этот же контекст в UseCase
	ctx := r.Context()
	reader := ctx.Value(middleware.ReaderContextKey("reader")).(io.Reader)

	b, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(b, &User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обнуляем поля, которые не должны прийти в запросе
	User.ID = 0
	User.Balance = 0

	// Проверяем формат структуры и обязательные для заполнения поля
	_, err = govalidator.ValidateStruct(User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loyalty, err := routes.loyalty.PostRegUser(ctx, &entity.Loyalty{User: &User})
	if err != nil {
		if errors.Is(err, usecase.ErrLoginAlreadyTaken) {
			http.Error(w, "login is already taken", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", ApplicationJSON)
	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(UserID{ID: strconv.FormatInt(loyalty.User.ID, 10)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ctx := r.Context()
	// key := ctx.Value(middlewa.FavContextKey("key"))

	// log.Println(key)
	w.Write(response) //nolint:errcheck
}

// PostLogUHandler -.
func (routes *loyaltyRoutes) PostLogUHandler(w http.ResponseWriter, r *http.Request) {

}
