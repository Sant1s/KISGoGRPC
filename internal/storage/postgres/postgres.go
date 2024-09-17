package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	_ "github.com/jackc/pgx"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUserDoesNotExists = errors.New("user does not exists")
)

type Blog interface {
}

type Auth interface {
	ValidateUser(ctx context.Context, credentials string) error
}

type Postgres struct {
	logger *slog.Logger
	db     *sqlx.DB
}

func New(logger *slog.Logger, storagePath string) (*Postgres, error) {
	const op = "storage.postgres.New"

	db, err := sqlx.Connect("postgers", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Postgres{db: db}, nil
}

func (p *Postgres) ValidateUser(ctx context.Context, credentials string) error {
	const op = "storage.postgres.ValidateUser"
	stmt, err := p.db.Prepare("select * from users where credentials=$1")

	if err != nil {
		p.logger.Error("statement error", slog.String("op", op), slog.Any("err", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, credentials)

	if err != nil {
		return err
	}

	// todo: вот тут надо дописать логику
	return nil
}
