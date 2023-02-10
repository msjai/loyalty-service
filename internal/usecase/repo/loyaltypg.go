package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/msjai/loyalty-service/internal/entity"
)

var (
	ErrLoginAlreadyTaken = errors.New("login is already taken")
	ErrConnectionNotOpen = errors.New("data base pgsql connection not opened")
)

// LoyaltyRepoS -.
type LoyaltyRepoS struct {
	repo *sql.DB
}

// New -.
func New(db *sql.DB) *LoyaltyRepoS {
	return &LoyaltyRepoS{repo: db}
}

// AddNewUser -.
func (r *LoyaltyRepoS) AddNewUser(ctx context.Context, loyalty *entity.Loyalty) (*entity.Loyalty, error) {
	if r.repo == nil {
		return nil, fmt.Errorf("repo - AddNewUser - repo: %w", ErrConnectionNotOpen)
	}

	tx, err := r.repo.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo - AddNewUser - repo.Begin: %w", err)
	}

	stmt, err := tx.PrepareContext(ctx, `INSERT INTO users (login, password) values ($1, $2) RETURNING id`)
	if err != nil {
		return nil, fmt.Errorf("repo - AddNewUser - tx.PrepareContext: %w", err)
	}
	defer stmt.Close()

	var (
		row *sql.Row
		id  int64
	)

	row = stmt.QueryRowContext(ctx, loyalty.User.Login, loyalty.User.Password)
	err = row.Scan(&id)
	if err != nil {
		// Здесь err только для условия отката транзакции, не перезаписывает исходную ошибку
		if err := tx.Rollback(); err != nil {
			return nil, fmt.Errorf("repo - AddNewUser - tx.Rollback: %w", err)
		}
		// Обрабатываем ошибку, что такой логин уже есть в базе
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			// Отдаем ошибку исходную ошибку, о том что такой пользователь уже есть
			return nil, fmt.Errorf("repo - AddNewUser - Scan: %w", ErrLoginAlreadyTaken)
		}
		// Отдаем исходную ошибку, если она не про уникальность логина
		return nil, fmt.Errorf("repo - AddNewUser - stmt.Exec: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("repo - AddNewUser - tx.Commit: %w", err)
	}

	loyalty.User.ID = id
	return loyalty, nil
}

// FindUser -.
func (r *LoyaltyRepoS) FindUser(context.Context) (*entity.Loyalty, error) {
	return nil, nil
}
