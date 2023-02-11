package usecase

import (
	"context"

	"github.com/msjai/loyalty-service/internal/entity"
)

type (
	// Loyalty -.
	Loyalty interface {
		PostRegUser(context.Context, *entity.Loyalty) (*entity.Loyalty, error)
		PostLoginUser(context.Context, *entity.Loyalty) (*entity.Loyalty, error)
		PostUserOrder(context.Context, *entity.UserOrder) (*entity.UserOrder, error)
		GetUserOrders(context.Context, *entity.Loyalty) (*entity.Loyalty, error)
		GetUserBalance(context.Context, *entity.Loyalty) (*entity.Loyalty, error)
		PostUserWithDrawBalance(context.Context, *entity.Loyalty) (*entity.Loyalty, error)
		GetUserWithdrawals(context.Context, *entity.Loyalty) (*entity.Loyalty, error)
	}

	// LoyaltyRepo -.
	LoyaltyRepo interface {
		AddNewUser(context.Context, *entity.Loyalty) (*entity.Loyalty, error)
		FindUser(context.Context, *entity.Loyalty) (*entity.Loyalty, error)
		AddOrder(context.Context, *entity.UserOrder) (*entity.UserOrder, error)
	}

	// LoyaltyWebAPI -.
	LoyaltyWebAPI interface {
		GetOrderInfo(*entity.Loyalty) (*entity.Loyalty, error)
	}
)
