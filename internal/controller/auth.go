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
	ID    string `json:"id"`
	Token string `json:"token"`
}

// clearUserFields -.
func clearUserFields(user *entity.User) *entity.User {
	user.ID = 0
	user.Balance = 0
	user.Token = ""
	return user
}

// PostRegUHandler -.
func (routes *loyaltyRoutes) PostRegUHandler(w http.ResponseWriter, r *http.Request) {
	var User entity.User
	// Через контекст получаем reader
	// В случае необхоимости тело было распаковано в middleware
	// Далее передаем этот же контекст в UseCase
	ctx := r.Context()
	reader := ctx.Value(middleware.KeyReader).(io.Reader)

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
	clearUserFields(&User)

	// Проверяем формат структуры и обязательные для заполнения поля
	_, err = govalidator.ValidateStruct(User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loyalty, err := routes.loyalty.PostRegUser(ctx, &entity.Loyalty{User: &User})
	if err != nil {
		if errors.Is(err, usecase.ErrLoginAlreadyTaken) {
			http.Error(w, usecase.ErrLoginAlreadyTaken.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", ApplicationJSON)
	w.Header().Set("Authorization", "Bearer "+loyalty.User.Token)
	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(UserID{ID: strconv.FormatInt(loyalty.User.ID, 10), Token: loyalty.User.Token})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(response) //nolint:errcheck
}

// PostLogUHandler -.
func (routes *loyaltyRoutes) PostLogUHandler(w http.ResponseWriter, r *http.Request) {
	var User entity.User
	// Через контекст получаем reader
	// В случае необхоимости тело было распаковано в middleware
	// Далее передаем этот же контекст в UseCase
	// var k middleware.KeyForReader
	// k = "reader"
	ctx := r.Context()
	reader := ctx.Value(middleware.KeyReader).(io.Reader)

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
	clearUserFields(&User)

	// Проверяем формат структуры и обязательные для заполнения поля
	_, err = govalidator.ValidateStruct(User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loyalty, err := routes.loyalty.PostLoginUser(ctx, &entity.Loyalty{User: &User})
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidLogPass) {
			http.Error(w, usecase.ErrInvalidLogPass.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", ApplicationJSON)
	w.Header().Set("Authorization", "Bearer "+loyalty.User.Token)
	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(UserID{ID: strconv.FormatInt(loyalty.User.ID, 10), Token: loyalty.User.Token})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(response) //nolint:errcheck
}
