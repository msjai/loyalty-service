package usecase

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/msjai/loyalty-service/internal/entity"
)

// tokenClaims -.
type tokenClaims struct {
	jwt.StandardClaims
	UserID int64 `json:"user_id"` //nolint:tagliatelle
}

func getToken(loyalty *entity.Loyalty) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: loyalty.User.ID,
	})

	signedToken, err := token.SignedString([]byte(signTokenKey))
	if err != nil {
		return "", fmt.Errorf("usecase - getToken - SignedString: %w", err)
	}

	return signedToken, nil
}

func ParseToken(token string) (int64, error) {
	t, err := jwt.ParseWithClaims(token, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("usecase - ParseToken - SigningMethodHMAC: %w", ErrInvalidSigningMethod)
		}
		return []byte(signTokenKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("usecase - ParseToken - ParseWithClaims: %w", err)
	}

	claims, ok := t.Claims.(*tokenClaims)
	if !ok {
		return 0, fmt.Errorf("usecase - ParseToken - Claims: %w", err)
	}

	if !t.Valid {
		return 0, fmt.Errorf("usecase - ParseToken - Valid: %w", err)
	}

	return claims.UserID, nil
}
