package webapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

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
func (wa *LoyaltyWebAPI) RefreshOrderInfo(userOrder *entity.UserOrder) (*entity.UserOrder, error) {
	l := wa.cfg.L

	// ctxGet, cancel := context.WithCancel(context.Background())
	// defer cancel()

	//	request, err := http.NewRequest(http.MethodGet, "http://"+wa.cfg.AccrualSystemAddress+"/api/orders/"+fmt.Sprint(userOrder.Number), nil)
	request, err := http.NewRequest(http.MethodGet, wa.cfg.AccrualSystemAddress+"/api/orders/"+fmt.Sprint(userOrder.Number), nil)
	if err != nil {
		l.Errorf("webapi - RefreshOrderInfo - NewRequestWithContext: %v", err.Error())
	}
	request.Header.Set("Accept", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		l.Errorf("webapi - RefreshOrderInfo - DefaultClient.Do: %v", err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		b, err := io.ReadAll(response.Body)
		if err != nil {
			l.Errorf("webapi - RefreshOrderInfo - io.ReadAll: %v", err.Error())
		}

		err = json.Unmarshal(b, &userOrder)
		if err != nil {
			l.Errorf("webapi - RefreshOrderInfo - json.Unmarshal: %v", err.Error())
		}
	}

	if response.StatusCode == http.StatusTooManyRequests {
		b, err := io.ReadAll(response.Body)
		if err != nil {
			l.Errorf("webapi - RefreshOrderInfo - io.ReadAll: %v", err.Error())
		}

		l.Errorf("webapi - RefreshOrderInfo - json.Unmarshal: %v", err.Error())

		timeToWait, err := strconv.Atoi(string(b))
		if err != nil {
			l.Errorf("webapi - RefreshOrderInfo - strconv.Atoi: %v", err.Error())
			timeToWait = 10
		}

		l.Infof("sleeping for %v seconds", timeToWait)
		time.Sleep(time.Second * time.Duration(timeToWait))
	}

	if response.StatusCode == http.StatusInternalServerError {
		l.Infof("webapi - RefreshOrderInfo - response.StatusCode: %v", http.StatusInternalServerError)
	}

	if response.StatusCode == http.StatusNoContent {
		l.Infof("webapi - RefreshOrderInfo - response.StatusCode: %v - order: %v", http.StatusNoContent, *userOrder)
	}

	return userOrder, nil
}
