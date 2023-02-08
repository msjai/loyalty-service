package controller

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/asaskevich/govalidator"

	"github.com/msjai/loyalty-service/internal/controller/middleware"
	"github.com/msjai/loyalty-service/internal/entity"
)

type UserID struct {
	Id string `json:"id"`
}

// PostRegUHandler -.
func (routes *loyaltyRoutes) PostRegUHandler(w http.ResponseWriter, r *http.Request) {
	var User entity.User

	// Через контекст получаем reader
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

	_, err = govalidator.ValidateStruct(User)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// id, err := h.services.Authorization.CreateUser(input)
	// if err != nil {
	// 	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	w.Header().Set("Content-Type", ApplicationJSON)
	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(UserID{Id: "1"})
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
