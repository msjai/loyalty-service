package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/msjai/loyalty-service/internal/config"
	"github.com/msjai/loyalty-service/internal/entity"
)

// LoyaltyWebAPI -.
type LoyaltyWebAPI struct {
	cfg *config.Config
}

// New -.
func New(config *config.Config) *LoyaltyWebAPI {
	return &LoyaltyWebAPI{cfg: config}
}

// RefreshOrderInfo - Функция получает информацию из черного ящика по 1 определенному заказу
func (wa *LoyaltyWebAPI) RefreshOrderInfo(ctx context.Context, userOrder *entity.UserOrder) (*entity.UserOrder, error) {
	l := wa.cfg.L

	ctxGet, cancel := context.WithCancel(context.Background())
	defer cancel()

	request, err := http.NewRequestWithContext(ctxGet, http.MethodGet, "http://"+wa.cfg.AccrualSystemAddress+"/api/orders/"+fmt.Sprint(userOrder.Number), nil)
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

	return userOrder, nil
}
