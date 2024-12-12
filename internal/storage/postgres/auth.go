package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Sant1s/blogBack/internal/domain"
	"github.com/Sant1s/blogBack/internal/storage"
)

// ----------------------------
//		REGISTER USER
// ----------------------------

const queryRegisterUser = `
INSERT INTO users(id,
                  nickname,
                  password_hash,
                  created_at,
                  permission)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;
` // todo: Пока будет стоит только 'user'

func (p *Postgres) Register(
	ctx context.Context,
	request *domain.RegisterUserRequest,
) (*domain.RegisterUserResponse, error) {

	const op = "postgres.Register"

	_, err := p.getUserId(ctx, request.Login)

	if err == nil {
		p.logger.Error(
			"failed executing query: user already exists",
			slog.String("op", op),
			slog.String("query", queryRegisterUser),
			slog.Any("err", storage.ErrAlreadyExists),
		)
		return nil, storage.ErrAlreadyExists
	}

	p.logger.Info(
		"executing query: ",
		slog.String("op", op),
		slog.String("query", queryRegisterUser),
	)

	userNewUUid := uuid.New()
	res := p.db.QueryRowContext(
		ctx,
		queryRegisterUser,
		userNewUUid,
		request.Login,
		request.PasswordHash,
		time.Now(),
		request.Permission,
	)

	var uuid string
	if err = res.Scan(&uuid); err != nil {
		p.logger.Error(
			"failed executing query",
			slog.String("op", op),
			slog.String("query", queryRegisterUser),
			slog.Any("err", res.Err()),
		)

		return nil, fmt.Errorf("%w: %v", storage.ErrInternal, err)
	}

	return &domain.RegisterUserResponse{Id: uuid}, nil
}

// ----------------------------
//		LOGIN USER
// ----------------------------

const queryLoginUser = `
SELECT id, permission FROM users
WHERE nickname=$1 AND password_hash=$2;
`

func (p *Postgres) Login(ctx context.Context, request *domain.LoginUserRequest) (*domain.LoginUserResponse, error) {
	const op = "postgres.Login"

	_, err := p.getUserId(ctx, request.Login)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", storage.ErrDoesNotExists, err)
	}

	p.logger.Info(
		"executing query:",
		slog.String("op", op),
		slog.String("query", queryLoginUser),
	)

	res := p.db.QueryRowContext(
		ctx,
		queryLoginUser,
		request.Login,
		request.PasswordHash,
	)

	var (
		uuid       string
		permission string
	)
	if err = res.Scan(&uuid, &permission); err != nil {
		p.logger.Error(
			"failed executing query",
			slog.String("op", op),
			slog.String("query", queryLoginUser),
			slog.Any("err", err),
		)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", storage.ErrDoesNotExists, err)
		}

		return nil, fmt.Errorf("%w: %v", storage.ErrInternal, err)
	}

	return &domain.LoginUserResponse{
		Id:         uuid,
		Permission: permission,
	}, nil
}
