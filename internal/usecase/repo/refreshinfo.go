package repo

import (
	"fmt"

	"github.com/msjai/loyalty-service/internal/entity"
)

// CatchOrdersRefresh - Функция получает из базы заказы, статус по которым не является окончательным.
// Это все статусы, кроме (PROCESSING, INVALID)
func (r *LoyaltyRepoS) CatchOrdersRefresh() ([]*entity.UserOrder, error) {
	if r.repo == nil {
		return nil, fmt.Errorf("repo - CatchOrdersRefresh - repo: %w", ErrConnectionNotOpen)
	}

	tx, err := r.repo.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo - CatchOrdersRefresh - repo.Begin: %w", err)
	}

	stmt, err := tx.Prepare(`SELECT id, number, status, user_id, accrual_sum, uploaded_at FROM orders WHERE status<>$1 and status<>$2`)
	if err != nil {
		return nil, fmt.Errorf("repo - CatchOrdersRefresh - tx.PrepareContext: %w", err)
	}
	defer stmt.Close()

	var orders []*entity.UserOrder

	rows, err := stmt.Query(entity.PROCESSED, entity.INVALID)
	if err != nil {
		return orders, fmt.Errorf("repo - CatchOrdersRefresh - stmt.QueryContext: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var userOrder entity.UserOrder
		err = rows.Scan(&userOrder.ID, &userOrder.Number, &userOrder.Status, &userOrder.UserID, &userOrder.AccrualSum, &userOrder.UploadedAt)
		if err != nil {
			if errRollBack := tx.Rollback(); errRollBack != nil {
				return orders, fmt.Errorf("repo - CatchOrdersRefresh - row.Scan: %w - tx.RollBack(): %v", err, errRollBack)
			}

			return orders, fmt.Errorf("repo - CatchOrdersRefresh - row.Scan: %w", err)
		}

		orders = append(orders, &userOrder)
	}

	if err = tx.Commit(); err != nil {
		return orders, fmt.Errorf("repo - CatchOrdersRefresh - tx.Commit: %w", err)
	}

	err = rows.Err()
	if err != nil {
		return orders, fmt.Errorf("repo - CatchOrdersRefresh - rows.Err(): %w", err)
	}

	return orders, nil
}

func (r *LoyaltyRepoS) UpdateOrder(userOrder *entity.UserOrder) (*entity.UserOrder, error) {
	if r.repo == nil {
		return nil, fmt.Errorf("repo - AddOrder - repo: %w", ErrConnectionNotOpen)
	}

	tx, err := r.repo.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo - UpdateOrder - repo.Begin: %w", err)
	}

	stmt, err := tx.Prepare(`UPDATE orders SET status=$1, accrual_sum =$2 WHERE number=$3 RETURNING id, number, status, user_id, accrual_sum, uploaded_at`)
	if err != nil {
		return nil, fmt.Errorf("repo - UpdateOrder - tx.PrepareContext: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(userOrder.Status, userOrder.AccrualSum, userOrder.Number)
	err = row.Scan(&userOrder.ID, &userOrder.Number, &userOrder.Status, &userOrder.UserID, &userOrder.AccrualSum, &userOrder.UploadedAt)
	if err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return userOrder, fmt.Errorf("repo - UpdateOrder - row.Scan: %w - tx.RollBack(): %v", err, errRollBack)
		}

		return userOrder, fmt.Errorf("repo - UpdateOrder - row.Scan: %w", err)
	}

	// Здесь в одной транзакции увеличиваем баланс пользователя, в таблице users
	stmt1, err := tx.Prepare(`UPDATE users SET balance=balance+$1 WHERE id=$2`)
	if err != nil {
		return nil, fmt.Errorf("repo - UpdateOrder - tx.PrepareContext: %w", err)
	}
	defer stmt1.Close()

	res, err := stmt1.Exec(userOrder.AccrualSum, userOrder.UserID)
	if err != nil {
		if errRollBack := tx.Rollback(); errRollBack != nil {
			return userOrder, fmt.Errorf("repo - UpdateOrder - stmt1.Exec: %w - tx.RollBack(): %v", err, errRollBack)
		}

		return userOrder, fmt.Errorf("repo - UpdateOrder - stmt1.Exec: %w", err)
	}

	rowsAf, err := res.RowsAffected()
	err = handleRABalance(rowsAf, tx, userOrder, err)
	if err != nil {
		return userOrder, fmt.Errorf("repo - UpdateOrder - RowsAffected: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return userOrder, fmt.Errorf("repo - UpdateOrder - tx.Commit: %w", err)
	}

	return userOrder, nil
}
