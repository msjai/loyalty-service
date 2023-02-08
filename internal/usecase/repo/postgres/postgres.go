package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

var ErrAccrualSA = errors.New("data base URI not set")

const (
	maxOpenIdleConns        = 20
	maxIdleTimeConn         = 30 * time.Second
	maxLifeTimeConn         = 120 * time.Second
	pingAttempts            = 3
	migrationAttempts       = 3
	migrationAttemptTimeout = time.Second
)

// New -.
func New(dbURI string, l *zap.SugaredLogger) (*sql.DB, error) {
	l.Infow("Initializing URL repository postgresql...")

	db, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("repo - New - sql.Open: %w", err)
	}

	db.SetMaxOpenConns(maxOpenIdleConns)   // Ограничение количество соединений, используемых приложением
	db.SetMaxIdleConns(maxOpenIdleConns)   // Ограничение количества простаивающих соединений. Рекомендуется не меньше "SetMaxOpenConns"
	db.SetConnMaxIdleTime(maxIdleTimeConn) // Ограничение времени простаивания одного соединения
	db.SetConnMaxLifetime(maxLifeTimeConn) // Ограничение времени работы одного соединения

	err = pingRepo(db, l)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - pingRepo: %w", err)
	}

	err = migrateUP(dbURI, l)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - migrateUP: %w", err)
	}

	return db, nil
}

// pingRepo -.
func pingRepo(db *sql.DB, l *zap.SugaredLogger) error {
	var (
		err      error
		attempts = pingAttempts
	)

	start := time.Now()
	for attempts > 0 {
		err = db.Ping()

		if err == nil {
			l.Infow("repo: postgres ping succes...")
			break
		}

		attempts--
		l.Infof("repo: postgres is trying to connect, attempts left: %d", attempts)
	}
	d := time.Since(start)
	l.Debugf("repo: tried to ping %v", d.Truncate(time.Millisecond))

	if err != nil {
		return fmt.Errorf("repo - pingRepo - db.Ping: %w", err)
	}

	return nil
}

func migrateUP(dbURI string, l *zap.SugaredLogger) error {
	l.Infof("Migration started... db url is:%v", dbURI)

	if len(dbURI) == 0 {
		return fmt.Errorf("posgres - migrateUP - dbURI: %w", ErrAccrualSA)
	}

	var (
		attempts = migrationAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://internal/usecase/repo/postgres/migrations", dbURI)
		if err == nil {
			break
		}

		l.Infof("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(migrationAttemptTimeout)
		attempts--
	}

	if err != nil {
		return fmt.Errorf("posgres - migrateUP - migrate.New: %w", err)
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("posgres - migrateUP - m.UP(): %w", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		l.Infoln("Migrate: no change")
		return nil
	}

	l.Infoln("Migrate: up success")

	return nil
}
