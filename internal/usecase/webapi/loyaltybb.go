package webapi

import (
	"github.com/msjai/loyalty-service/internal/config"
	"github.com/msjai/loyalty-service/internal/entity"
)

// LoyaltyWebAPI -.
type LoyaltyWebAPI struct {
}

// New -.
func New(config *config.Config) *LoyaltyWebAPI {
	return &LoyaltyWebAPI{}
}

// GetOrderInfo -.
func (wa *LoyaltyWebAPI) GetOrderInfo(loyalty *entity.Loyalty) (*entity.Loyalty, error) {
	return loyalty, nil
}
