package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Sant1s/blogBack/internal/domain"
	"github.com/jmoiron/sqlx"
)

type BlogPosts interface {
	GetListPosts(ctx context.Context, limit, offset int32) ([]domain.Post, error)
	CreatePost(ctx context.Context, post *domain.Post) error
	UpdatePost(ctx context.Context, updates *domain.PostUpdateRequest) (int64, error)
	DeletePost(ctx context.Context, postId int64) error

	GetListComments(ctx context.Context, limit, offset int32) ([]domain.Comment, error)
	CreateComment(ctx context.Context, comment *domain.Comment) error
	UpdateComment(ctx context.Context, updates *domain.CommentUpdateRequest) (int64, error)
	DeleteComment(ctx context.Context, commentId int64) error
}

type Auth interface {
	ValidateUser(ctx context.Context, nickname, passwordHash string) error
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
