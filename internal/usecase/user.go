package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/msjai/loyalty-service/internal/entity"
	"github.com/msjai/loyalty-service/internal/usecase/repo"
)

// PostRegUser -.
func (luc *LoyaltyUseCase) PostRegUser(ctx context.Context, loyalty *entity.Loyalty) (*entity.Loyalty, error) {
	// TODO Перенести вызов функции на уровень репозитория
	loyalty.User.Password = hashPassword(loyalty.User.Password)

	loyalty, err := luc.repo.AddNewUser(ctx, loyalty)
	if err != nil {
		if errors.Is(err, repo.ErrLoginAlreadyTaken) {
			return nil, fmt.Errorf("usecase - PostRegUser - AddNewUser: %w", ErrLoginAlreadyTaken)
		}
		return nil, fmt.Errorf("usecase - PostRegUser - AddNewUser: %w", err)
	}

	loyalty.User.Token, err = getToken(loyalty)
	if err != nil {
		return nil, fmt.Errorf("usecase - PostRegUser - getToken: %w", err)
	}

	return loyalty, nil
}

// PostLoginUser -.
func (luc *LoyaltyUseCase) PostLoginUser(ctx context.Context, loyalty *entity.Loyalty) (*entity.Loyalty, error) {
	loyalty.User.Password = hashPassword(loyalty.User.Password)

	loyalty, err := luc.repo.FindUser(ctx, loyalty)
	if err != nil {
		if errors.Is(err, repo.ErrInvalidLogPass) {
			return nil, fmt.Errorf("usecase - PostLoginUser - FindUser: %w", ErrInvalidLogPass)
		}
		return nil, fmt.Errorf("usecase - PostLoginUser - FindUser: %w", err)
	}

	loyalty.User.Token, err = getToken(loyalty)
	if err != nil {
		return nil, fmt.Errorf("usecase - PostLoginUser - getToken: %w", err)
	}

	return loyalty, nil
}
