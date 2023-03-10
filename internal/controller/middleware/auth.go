package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/msjai/loyalty-service/internal/usecase"
)

const (
	AuthorizationHeader = "Authorization"
)

type KeyForUserID string

var KeyUserID KeyForUserID = "reader"

var (
	ErrNotAuthentificated          = errors.New("user not authenticated")
	ErrNotAuthEmptyHeader          = errors.New("user not authenticated; no authorization header")
	ErrNotAuthBadWrongHeaderFormat = errors.New("user not authenticated; wrong authorization header format")
)

func IdentifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		aHeader := r.Header.Get(AuthorizationHeader)
		if aHeader == "" {
			http.Error(w, ErrNotAuthEmptyHeader.Error(), http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(aHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, ErrNotAuthBadWrongHeaderFormat.Error(), http.StatusUnauthorized)
			return
		}

		userID, err := usecase.ParseToken(headerParts[1])
		if err != nil {
			http.Error(w, ErrNotAuthentificated.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), KeyUserID, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
