package repo

import (
	"fmt"
	"time"

	"github.com/msjai/loyalty-service/internal/entity"
)

func (r *LoyaltyRepoS) WithDraw(withDraw *entity.WithDraw) (*entity.WithDraw, error) {
	if r.repo == nil {
		return nil, fmt.Errorf("repo - WithDraw - repo: %w", ErrConnectionNotOpen)
	}

	tx, err := r.repo.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo - WithDraw - repo.Begin: %w", err)
	}

	procAt := time.Now()

	stmt, err := tx.Prepare(`INSERT INTO writes_off (order_woff_num, sum, user_id, date)
								values ($1, $2, $3, $4)
                    	RETURNING id, order_woff_num, sum, user_id, date`)
	if err != nil {
		return nil, fmt.Errorf("repo - WithDraw - tx.PrepareContext: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(withDraw.Number, withDraw.Sum, withDraw.UserID, procAt)
	err = row.Scan(&withDraw.ID, &withDraw.Number, &withDraw.Sum, &withDraw.UserID, &withDraw.ProcessedAt)
	if err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return withDraw, fmt.Errorf("repo - UpdateOrder - row.Scan: %w - tx.RollBack(): %v", err, errRollBack)
		}

		return withDraw, fmt.Errorf("repo - WithDraw - row.Scan: %w", err)
	}

	// Здесь в одной транзакции уменьшаем баланс пользователя, в таблице users
	// Если баланс пользователя после списания стал меньше 0, то откатываем транзакцию
	stmt1, err := tx.Prepare(`UPDATE users SET balance=balance-$1, withdrawn=withdrawn+$1 WHERE id=$2
								RETURNING id, balance, withdrawn`)
	if err != nil {
		return nil, fmt.Errorf("repo - WithDraw - tx.PrepareContext: %w", err)
	}
	defer stmt1.Close()

	var (
		userID    int64
		balance   float64
		withDrawn float64
	)

	res := stmt1.QueryRow(withDraw.Sum, withDraw.UserID)
	err = res.Scan(&userID, &balance, &withDrawn)
	if err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return withDraw, fmt.Errorf("repo - WithDraw - stmt1.QueryRow: %w - tx.RollBack(): %v", err, errRollBack)
		}

		return withDraw, fmt.Errorf("repo - WithDraw - stmt1.Exec: %w", err)
	}

	if balance < 0 {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return withDraw, fmt.Errorf("repo - WithDraw - stmt1.QueryRow: %w - tx.RollBack(): %v", ERRInsufficientFund, errRollBack)
		}

		return withDraw, fmt.Errorf("repo - WithDraw - stmt1.QueryRow: %w", ERRInsufficientFund)
	}

	// rowsAf, err := res.RowsAffected()
	// err = handleRABalance(rowsAf, tx, userOrder, err)
	// if err != nil {
	// 	return userOrder, fmt.Errorf("repo - WithDraw - RowsAffected: %w", err)
	// }

	if err = tx.Commit(); err != nil {
		return withDraw, fmt.Errorf("repo - WithDraw - tx.Commit: %w", err)
	}

	return withDraw, nil
}
