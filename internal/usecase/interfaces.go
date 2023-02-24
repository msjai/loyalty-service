package usecase

import (
	"github.com/msjai/loyalty-service/internal/entity"
)

type (
	// Loyalty -.
	Loyalty interface {
		PostRegUser(*entity.Loyalty) (*entity.Loyalty, error)
		PostLoginUser(*entity.Loyalty) (*entity.Loyalty, error)
		PostUserOrder(*entity.UserOrder) (*entity.UserOrder, error)
		GetUserOrders(*entity.User) ([]*entity.UserOrder, error)
		GetUserBalance(*entity.User) (*entity.UserBalance, error)
		PostUserWithDrawBalance(*entity.WithDraw) (*entity.WithDraw, error)
		GetUserWithdrawals(*entity.User) ([]*entity.WithDraw, error)
		RefreshOrderInfo(*entity.UserOrder) (*entity.UserOrder, error)
		CatchOrdersRefresh() ([]*entity.UserOrder, error)
	}

	// LoyaltyRepo -.
	LoyaltyRepo interface {
		AddNewUser(*entity.Loyalty) (*entity.Loyalty, error)
		FindUser(*entity.Loyalty) (*entity.Loyalty, error)
		AddOrder(*entity.UserOrder) (*entity.UserOrder, error)
		FindOrder(*entity.UserOrder) (*entity.UserOrder, error)
		CatchOrdersRefresh() ([]*entity.UserOrder, error)
		UpdateOrder(*entity.UserOrder) (*entity.UserOrder, error)
		FindOrders(*entity.User) ([]*entity.UserOrder, error)
		GetUserBalance(user *entity.User) (*entity.User, error)
		WithDraw(*entity.WithDraw) (*entity.WithDraw, error)
		GetUserWithdrawals(*entity.User) ([]*entity.WithDraw, error)
	}

	// LoyaltyWebAPI -.
	LoyaltyWebAPI interface {
		RefreshOrderInfo(*entity.UserOrder) (*entity.UserOrder, error)
	}
)
