package repo

import (
	"database/sql"
	"fmt"

	"github.com/msjai/loyalty-service/internal/entity"
)

// AddNewUser -.
func (r *LoyaltyRepoS) AddNewUser(loyalty *entity.Loyalty) (*entity.Loyalty, error) {
	if r.repo == nil {
		return nil, fmt.Errorf("repo - AddNewUser - repo: %w", ErrConnectionNotOpen)
	}

	tx, err := r.repo.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo - AddNewUser - repo.Begin: %w", err)
	}

	stmt, err := tx.Prepare(`INSERT INTO users (login, password) values ($1, $2) RETURNING id`)
	if err != nil {
		return nil, fmt.Errorf("repo - AddNewUser - tx.PrepareContext: %w", err)
	}
	defer stmt.Close()

	var (
		row *sql.Row
		id  int64
	)

	row = stmt.QueryRow(loyalty.User.Login, loyalty.User.Password)
	err = row.Scan(&id)
	if err != nil {
		return nil, handleInsertUserError(tx, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("repo - AddNewUser - tx.Commit: %w", err)
	}

	loyalty.User.ID = id
	return loyalty, nil
}

// FindUser -.
func (r *LoyaltyRepoS) FindUser(loyalty *entity.Loyalty) (*entity.Loyalty, error) {
	if r.repo == nil {
		return nil, fmt.Errorf("repo - FindUser - repo: %w", ErrConnectionNotOpen)
	}

	tx, err := r.repo.Begin()
	if err != nil {
		return nil, fmt.Errorf("repo - FindUser - repo.Begin: %w", err)
	}

	stmt, err := tx.Prepare(`SELECT id FROM users WHERE login=$1 and password= $2`)
	if err != nil {
		return nil, fmt.Errorf("repo - FindUser - tx.PrepareContext: %w", err)
	}
	defer stmt.Close()

	var (
		row *sql.Row
		id  int64
	)

	row = stmt.QueryRow(loyalty.User.Login, loyalty.User.Password)
	err = row.Scan(&id)
	if err != nil {
		return nil, handleFindUserError(tx, err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("repo - FindUser - tx.Commit: %w", err)
	}

	loyalty.User.ID = id

	return loyalty, nil
}
