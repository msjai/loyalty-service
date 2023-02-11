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
	// Обрабатываем ошибку, что такой логин уже есть в базе
	// var pgErr *pgconn.PgError
	// if errors.As(err, &pgErr) && pgErr.ConstraintName == "no rows in result set" {
	if errors.Is(err, sql.ErrNoRows) {
		// Отдаем ошибку исходную ошибку, о том что такой пользователь уже есть
		return fmt.Errorf("repo - FindUser - Scan: %w", ErrInvalidLogPass)
	}
	// Отдаем исходную ошибку, если она не про уникальность логина
	return fmt.Errorf("repo - FindUser - stmt.QueryRowContext: %w", err)
}

func handleInsertOrderError(tx *sql.Tx, err error) error {
	// Здесь err только для условия отката транзакции, не перезаписывает исходную ошибку
	if err := tx.Rollback(); err != nil {
		return fmt.Errorf("repo - AddOrder - tx.Rollback: %w", err)
	}
	// Обрабатываем ошибку, что заказ с таким номером уже есть в базе
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		// Отдаем ошибку исходную ошибку, о том что такой заказ уже есть
		return fmt.Errorf("repo - AddOrder - Scan: %w", ErrOrderNumExists)
	}
	// Отдаем исходную ошибку, если она не про уникальность логина
	return fmt.Errorf("repo - AddOrder - stmt.QueryRowContext: %w", err)
}
