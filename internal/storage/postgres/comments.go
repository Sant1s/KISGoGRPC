package postgres

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Sant1s/blogBack/internal/domain"
)

// ----------------------------
//		GET LIST COMMENTS
// ----------------------------

func (p *Postgres) GetListComments(ctx context.Context, limit, offset int32) ([]domain.Comment, error) {
	//TODO implement me
	panic("implement me")
}

// ----------------------------
//		CREATE COMMENT
// ----------------------------

func (p *Postgres) CreateComment(ctx context.Context, comment *domain.Comment) error {
	//TODO implement me
	panic("implement me")
}

// ----------------------------
//		UPDATE COMMENT
// ----------------------------

func (p *Postgres) UpdateComment(ctx context.Context, updates *domain.CommentUpdateRequest) (int64, error) {
	//TODO implement me
	panic("implement me")
}

// ----------------------------
//		DELETE COMMENT
// ----------------------------

func (p *Postgres) DeleteComment(ctx context.Context, commentId int64) error {
	//TODO implement me
	panic("implement me")
}
