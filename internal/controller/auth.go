package controller

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/asaskevich/govalidator"

	"github.com/msjai/loyalty-service/internal/entity"
	"github.com/msjai/loyalty-service/internal/usecase"
)

// clearUserFields -.
func clearUserFields(user entity.User) entity.User {
	user.ID = 0
	user.Balance = 0
	user.Token = ""
	return user
}

// PostRegUHandler -.
func (routes *loyaltyRoutes) PostRegUHandler(w http.ResponseWriter, r *http.Request) {
	var User entity.User

	reader := r.Body
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
	User = clearUserFields(User)

	// Проверяем формат структуры и обязательные для заполнения поля
	_, err = govalidator.ValidateStruct(User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loyalty, err := routes.loyalty.PostRegUser(&entity.Loyalty{User: &User})
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
}

// PostLogUHandler -.
func (routes *loyaltyRoutes) PostLogUHandler(w http.ResponseWriter, r *http.Request) {
	var User entity.User

	reader := r.Body
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
	User = clearUserFields(User)

	// Проверяем формат структуры и обязательные для заполнения поля
	_, err = govalidator.ValidateStruct(User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loyalty, err := routes.loyalty.PostLoginUser(&entity.Loyalty{User: &User})
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
}
