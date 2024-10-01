package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Sant1s/blogBack/internal/storage"
	"github.com/jackc/pgx/v5"

	"github.com/Sant1s/blogBack/internal/domain"
	"github.com/jmoiron/sqlx"
)

type BlogPosts interface {
	GetListPosts(ctx context.Context, limit, offset int32) ([]domain.Post, error)
	CreatePost(ctx context.Context, post *domain.Post) error
	UpdatePost(ctx context.Context, updates *domain.PostUpdateRequest) error
	DeletePost(ctx context.Context, postId int64) error

	GetListComments(ctx context.Context, limit, offset int32, postId int64) ([]domain.Comment, error)
	CreateComment(ctx context.Context, comment *domain.Comment) error
	UpdateComment(ctx context.Context, updates *domain.CommentUpdateRequest) error
	DeleteComment(ctx context.Context, commentId, postId int64) error
}

type Auth interface {
	ValidateUser(ctx context.Context, nickname, passwordHash string) error
	Register(ctx context.Context, request *domain.RegisterUserRequest) error
	Login(ctx context.Context, request *domain.LoginUserRequest) error
}

type Postgres struct {
	logger *slog.Logger
	db     *sqlx.DB
}

var _ BlogPosts = (*Postgres)(nil)

func New(logger *slog.Logger, storagePath string) (*Postgres, error) {
	const op = "storage.postgres.New"

	db, err := sqlx.Connect("pgx", storagePath)
	if err != nil {
		logger.Error(
			"error init postgres client",
			slog.String("op", op),
		)

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info(
		"postgres client init successfully",
		slog.String("op", op),
	)

	return &Postgres{
		db:     db,
		logger: logger,
	}, nil
}

// ----------------------------
//		VALIDATE USER
// ----------------------------

const queryValidateUser = `SELECT id FROM users WHERE nickname=$1 AND password_hash=$2`

func (p *Postgres) ValidateUser(ctx context.Context, nickname, passwordHash string) error {
	const op = "storage.postgres.ValidateUser"

	p.logger.Info(fmt.Sprintf("executing query: %s", queryValidateUser), slog.String("op", op))

	var id string
	res := p.db.QueryRowContext(ctx, queryValidateUser, nickname, passwordHash)

	err := res.Scan(&id)

	if err != nil {
		p.logger.Error(
			"failed to execute query",
			slog.String("query", queryValidateUser),
			slog.String("op", op),
			slog.Any("err", err),
		)

		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w: %v", storage.ErrDoesNotExists, err)
		}
		return fmt.Errorf("%w: %w", storage.ErrInternal, err)
	}

	return nil
}

// ----------------------------
// 		GET AUTHOR UUID
// ----------------------------

const queryGetAuthorId = `SELECT id FROM users WHERE nickname=$1;`

func (p *Postgres) getUserId(ctx context.Context, name string) (string, error) {
	const op = "postgres.getUserId"

	// Get author uuid
	var authorUuid string

	p.logger.Info(fmt.Sprintf("executing query: %s", queryGetAuthorId), slog.String("op", op))

	resAuthorId := p.db.QueryRowContext(ctx, queryGetAuthorId, name)

	err := resAuthorId.Scan(&authorUuid)
	if err != nil {
		p.logger.Error(
			"failed to execute query: %s",
			slog.String("query", queryGetAuthorId),
			slog.String("op", op),
			slog.Any("err", err),
		)

		if errors.Is(err, pgx.ErrNoRows) || authorUuid == "" {
			return "", fmt.Errorf("%w: %v", storage.ErrDoesNotExists, err)
		}

		return "", fmt.Errorf("%w: %v", storage.ErrInternal, err)
	}

	return authorUuid, nil
}
