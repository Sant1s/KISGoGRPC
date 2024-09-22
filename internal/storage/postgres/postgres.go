package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"

	"github.com/Sant1s/blogBack/internal/domain"
	_ "github.com/jackc/pgx/v5/stdlib" // Standard library bindings for pgx
	"github.com/jmoiron/sqlx"
)

var (
	ErrUserDoesNotExists = errors.New("user does not exists")
)

type BlogPosts interface {
	GetListPosts(ctx context.Context, limit, offset int32) (domain.Posts, error)
	CreatePost(ctx context.Context, post *domain.Post) error
	UpdatePost(ctx context.Context, updates *domain.Post) (int64, error)
	DeletePost(ctx context.Context, post_id int64) error
}

type Auth interface {
	ValidateUser(ctx context.Context, nickname, password_hash string) error
}

type Postgres struct {
	logger *slog.Logger
	db     *sqlx.DB
}

func New(logger *slog.Logger, storagePath string) (*Postgres, error) {
	const op = "storage.postgres.New"

	db, err := sqlx.Connect("pgx", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Postgres{db: db}, nil
}

func (p *Postgres) ValidateUser(ctx context.Context, nickname, password_hash string) error {
	const op = "storage.postgres.ValidateUser"
	stmt, err := p.db.Prepare("select id from users where nickname=$1 and password_hash=$2")

	if err != nil {
		p.logger.Error("statement error", slog.String("op", op), slog.Any("err", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	var id string
	err = stmt.QueryRowContext(ctx, nickname, password_hash).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserDoesNotExists
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// CreatePost implements blog.BlogPosts.
func (p *Postgres) CreatePost(ctx context.Context, post *domain.Post) error {

	panic("unimplemented")
}

// DeletePost implements blog.BlogPosts.
func (p *Postgres) DeletePost(ctx context.Context, post_id int64) error {
	panic("unimplemented")
}

// GetListPosts implements blog.BlogPosts.
func (p *Postgres) GetListPosts(ctx context.Context, limit int32, offset int32) (domain.Posts, error) {
	const op = "postgres.GetListPost"

	query := `SELECT p.id, u.nickname, p.data, p.created_at, p.comments_count, p.likes_count, false as liked
              FROM posts p
              LEFT JOIN users u ON p.author_id = u.id
              ORDER BY p.created_at
              LIMIT $1 OFFSET $2;`

	var posts domain.Posts
	err := p.db.Select(&posts, query, limit, offset)
	if err != nil {
		log.Fatalln(err)
	}

	return posts, nil
}

// UpdatePost implements blog.BlogPosts.
func (p *Postgres) UpdatePost(ctx context.Context, updates *domain.Post) (int64, error) {
	panic("unimplemented")
}
