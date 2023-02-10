package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/msjai/loyalty-service/internal/entity"
)

var ErrLoginAlreadyTaken = errors.New("login is already taken")

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
	var id int64
	query := fmt.Sprintf("INSERT INTO users (login, password) values ($1, $2) RETURNING id")

	row := r.repo.QueryRow(query, loyalty.User.Login, loyalty.User.Password)
	err := row.Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("repo - AddNewUser - Scan: %w", ErrLoginAlreadyTaken)
		}
		return nil, fmt.Errorf("repo - AddNewUser - Scan: %w", err)
	}

	return loyalty, nil
}

// FindUser -.
func (r *LoyaltyRepoS) FindUser(context.Context) (*entity.Loyalty, error) {
	return nil, nil
}
