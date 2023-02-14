package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/msjai/loyalty-service/internal/entity"
)

// AddOrder -.
func (r *LoyaltyRepoS) AddOrder(userOrder *entity.UserOrder) (*entity.UserOrder, error) {
	if r.repo == nil {
		return nil, fmt.Errorf("repo - AddOrder - repo: %w", ErrConnectionNotOpen)
	}

	tx, err := r.repo.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo - AddOrder - repo.Begin: %w", err)
	}

	stmt, err := tx.Prepare(`INSERT INTO orders (number, status, user_id, uploaded_at)
									           values ($1, $2, $3, $4) RETURNING id, number, status, user_id, accrual_sum, uploaded_at`)
	if err != nil {
		return nil, fmt.Errorf("repo - AddOrder - tx.PrepareContext: %w", err)
	}
	defer stmt.Close()

	var (
		row        *sql.Row
		id         int64
		number     string
		status     string
		userID     int64
		accrualSUM float64
		uploadedAt time.Time
	)

	row = stmt.QueryRow(userOrder.Number, entity.NEW, userOrder.UserID, time.Now())
	err = row.Scan(&id, &number, &status, &userID, &accrualSUM, &uploadedAt)
	if err != nil {
		return userOrder, handleInsertOrderError(tx, err)
	}

	if err = tx.Commit(); err != nil {
		return userOrder, fmt.Errorf("repo - AddOrder - tx.Commit: %w", err)
	}

	userOrder.ID = id
	userOrder.Number = number
	userOrder.Status = status
	userOrder.UserID = userID
	userOrder.AccrualSum = accrualSUM
	userOrder.UploadedAt = uploadedAt

	return userOrder, nil
}

// FindOrder -.
func (r *LoyaltyRepoS) FindOrder(userOrder *entity.UserOrder) (*entity.UserOrder, error) {
	if r.repo == nil {
		return nil, fmt.Errorf("repo - FindOrder - repo: %w", ErrConnectionNotOpen)
	}

	tx, err := r.repo.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo - FindOrder - repo.Begin: %w", err)
	}

	stmt, err := tx.Prepare(`SELECT user_id FROM orders WHERE number=$1`)
	if err != nil {
		return nil, fmt.Errorf("repo - FindOrder - tx.PrepareContext: %w", err)
	}
	defer stmt.Close()

	var (
		row    *sql.Row
		userID int64
	)

	row = stmt.QueryRow(userOrder.Number)
	err = row.Scan(&userID)
	if err != nil {
		return nil, handleFindOrderError(tx, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("repo - FindOrder - tx.Commit: %w", err)
	}

	// Значит заказ был зарегистрирован другим пользователем
	if userID != userOrder.UserID {
		return userOrder, fmt.Errorf("repo - FindOrder - hand made err: %w", ErrOrderAlreadyRegByAnotherUser)
	}

	return userOrder, fmt.Errorf("repo - FindOrder - hand made err: %w", ErrOrderAlreadyRegByCurrUser)
}

func (r *LoyaltyRepoS) FindOrders(*entity.User) ([]*entity.UserOrder, error) {

	return nil, nil
}