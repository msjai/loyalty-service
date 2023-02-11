package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/msjai/loyalty-service/internal/entity"
)

func (r *LoyaltyRepoS) AddOrder(ctx context.Context, userOrder *entity.UserOrder) (*entity.UserOrder, error) {
	if r.repo == nil {
		return nil, fmt.Errorf("repo - AddOrder - repo: %w", ErrConnectionNotOpen)
	}

	tx, err := r.repo.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo - AddOrder - repo.Begin: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO orders (number, status, user_id, uploaded_at)
									           values ($1, $2, $3, $4) RETURNING id`)
	if err != nil {
		return nil, fmt.Errorf("repo - AddOrder - tx.PrepareContext: %w", err)
	}
	defer stmt.Close()

	var (
		row *sql.Row
		id  int64
	)

	row = stmt.QueryRowContext(ctx, userOrder.Number, entity.NEW, userOrder.UserID, time.Now())
	err = row.Scan(&id)
	if err != nil {
		return nil, handleInsertOrderError(tx, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("repo - AddOrder - tx.Commit: %w", err)
	}

	userOrder.ID = id

	return userOrder, nil
}
