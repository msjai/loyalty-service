package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

func handleInsertError(tx *sql.Tx, err error) error {
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
	return fmt.Errorf("repo - AddNewUser - stmt.Exec: %w", err)
}
