package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

func handleInsertUserError(tx *sql.Tx, err error) error {
	// Здесь err только для условия отката транзакции, не перезаписывает исходную ошибку
	if err := tx.Rollback(); err != nil {
		return fmt.Errorf("repo - AddNewUser - tx.Rollback: %w", err)
	}
	// Обрабатываем ошибку, что такой логин уже есть в базе
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		// Отдаем ошибку исходную ошибку, о том что такой пользователь уже есть
		return fmt.Errorf("repo - AddNewUser - Scan: %w", ErrLoginAlreadyTaken)
	}
	// Отдаем исходную ошибку, если она не про уникальность логина
	return fmt.Errorf("repo - AddNewUser - stmt.QueryRowContext: %w", err)
}

func handleFindUserError(tx *sql.Tx, err error) error {
	// Здесь err только для условия отката транзакции, не перезаписывает исходную ошибку
	if err := tx.Rollback(); err != nil {
		return fmt.Errorf("repo - FindUser - tx.Rollback: %w", err)
	}

	// Если не нашли таких записей
	if errors.Is(err, sql.ErrNoRows) {
		// Отдаем ошибку о том что не правильное имя пользователя/пароль
		return fmt.Errorf("repo - FindUser - Scan: %w", ErrInvalidLogPass)
	}
	// Отдаем исходную ошибку, если она не про уникальность логина
	return fmt.Errorf("repo - FindUser - stmt.QueryRowContext: %w", err)
}

func handleInsertOrderError(tx *sql.Tx, err error) error {
	// Здесь err только для условия отката транзакции, не перезаписывает исходную ошибку
	rollbackERR := tx.Rollback()

	// Обрабатываем ошибку, что заказ с таким номером уже есть в базе
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		// Отдаем ошибку исходную ошибку, о том что такой заказ уже есть
		if rollbackERR != nil {
			return fmt.Errorf("repo - AddOrder - Scan: %w - tx.Rollback(): %v", ErrOrderNumExists, rollbackERR)
		}
		return fmt.Errorf("repo - AddOrder - Scan: %w", ErrOrderNumExists)
	}
	// Отдаем исходную ошибку, если она не про уникальность логина
	if rollbackERR != nil {
		return fmt.Errorf("repo - AddOrder - Scan: %w - tx.Rollback(): %v", ErrOrderNumExists, rollbackERR)
	}
	return fmt.Errorf("repo - AddOrder - stmt.QueryRowContext: %w", err)
}

func handleFindOrderError(tx *sql.Tx, err error) error {
	// Здесь err только для условия отката транзакции, не перезаписывает исходную ошибку
	if rollbackERR := tx.Rollback(); rollbackERR != nil {
		return fmt.Errorf("repo - AddOrder - stmt.QueryRowContext: %w - tx.Rollback(): %v", err, rollbackERR)
	}

	return fmt.Errorf("repo - AddOrder - stmt.QueryRowContext: %w", err)
}

func handleFindOrdersError(tx *sql.Tx, err error) error {
	if errRollBack := tx.Rollback(); errRollBack != nil {
		// Если не нашли таких записей
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("repo - FindOrders - row.Scan: %w - tx.RollBack(): %v", ErrNoUserOdersRL, errRollBack)
		}
		return fmt.Errorf("repo - FindOrders - row.Scan: %w - tx.RollBack(): %v", err, errRollBack)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("repo - FindOrders - row.Scan: %w", ErrNoUserOdersRL)
	}
	return fmt.Errorf("repo - FindOrders - row.Scan: %w", err)
}
