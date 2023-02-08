package repo

import (
	"context"
	"database/sql"

	"github.com/msjai/loyalty-service/internal/entity"
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
func (r *LoyaltyRepoS) AddNewUser(context.Context, *entity.Loyalty) error {
	return nil
}

// FindUser -.
func (r *LoyaltyRepoS) FindUser(context.Context) (*entity.Loyalty, error) {
	return nil, nil
}
