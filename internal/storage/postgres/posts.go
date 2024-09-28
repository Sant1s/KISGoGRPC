package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/Sant1s/blogBack/internal/storage"
	"log/slog"

	"github.com/Sant1s/blogBack/internal/domain"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// ----------------------------
//
//	CREATE POST
//
// ----------------------------
const queryCreatePost = `
INSERT INTO posts (id, author_id, data)
VALUES ($1, $2, $3)
RETURNING id;
`

// CreatePost implements blog.BlogPosts.
func (p *Postgres) CreatePost(ctx context.Context, post *domain.Post) error {
	const op = "postgres.CreatePost"

	uuid, err := p.getUserId(ctx, post.Author)
	if err != nil {
		p.logger.Error(
			"failed to execute query",
			slog.String("query", queryCreatePost),
			slog.String("op", op),
			slog.Any("err", err),
		)

		return err
	}

	// execute create query
	res := p.db.QueryRowContext(ctx, queryCreatePost, post.Id, uuid, post.Body)
	var id int64

	err = res.Scan(&id)

	if err != nil {
		p.logger.Error(
			"failed to execute query",
			slog.String("query", queryCreatePost),
			slog.String("op", op),
			slog.Any("err", err),
		)

		return fmt.Errorf("%w: %v", storage.ErrInternal, err)
	}

	return nil
}

// ----------------------------
// 		DELETE POST
// ----------------------------

/*

		Мнение автора таково, что при удалении поста комментарии не должны удаляться.
	Потому что вдруг нам будет нужна выгрузка пользователей, которые пишут плохие
	комментарии, на наша модель удаления запускается раз в ночь, когда маленькая нагрузка.
	Таких пользователей нужно забанить. А вдруг этот гад написал много плохого и резко всё удалил.

		Это причина, по которой автор не использовал здесь удаление комментариев к посту и транзакции,
	а не потому что он ленивый =)

*/

const queryDeletePost = `DELETE FROM posts WHERE id=$1;`

// DeletePost implements blog.BlogPosts.
func (p *Postgres) DeletePost(ctx context.Context, postId int64) error {
	const op = "postgres.DeletePost"

	p.logger.Info(fmt.Sprintf("executing query: %s", queryDeletePost), slog.String("op", op))

	res, err := p.db.ExecContext(ctx, queryDeletePost, postId)

	if err != nil {
		p.logger.Error(fmt.Sprintf("database error %s", err.Error()), slog.String("op", op))
		return fmt.Errorf("%w: %v", storage.ErrDoesNotExists, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		p.logger.Error(fmt.Sprintf("database error %v", err), slog.String("op", op))
		if rowsAffected == 0 {
			p.logger.Error(
				"failed to execute query",
				slog.String("op", op),
				slog.String("query", queryDeletePost),
				slog.Any("err", err),
			)

			return fmt.Errorf(
				"%w: %v",
				storage.ErrDoesNotExists,
				fmt.Sprintf("post with id=%d not found", postId),
			)
		}
		return fmt.Errorf("%w: %v", storage.ErrInternal, err)
	}

	return nil
}

// ----------------------------
//		GET LIST POSTS
// ----------------------------

const queryGetListPosts = `
SELECT p.id,
    u.nickname,
    p.data,
    p.created_at,
    p.comments_count,
    p.likes_count
FROM posts p
    LEFT JOIN users u ON p.author_id = u.id
ORDER BY p.created_at
LIMIT $1 OFFSET $2;`

// GetListPosts implements blog.BlogPosts.
func (p *Postgres) GetListPosts(ctx context.Context, limit int32, offset int32) ([]domain.Post, error) {
	const op = "postgres.GetListPost"

	p.logger.Info(fmt.Sprintf("executing query: %s", queryGetListPosts), slog.String("op", op))

	rows, err := p.db.QueryxContext(ctx, queryGetListPosts, limit, offset)
	if err != nil {
		p.logger.Error(
			"failed to execute query",
			slog.String("op", op),
			slog.String("query", queryGetListPosts),
			slog.Any("err", err),
		)

		return nil, fmt.Errorf("%w: %v", storage.ErrInternal, err)
	}

	result := make([]domain.Post, 0)
	for rows.Next() {
		var post domain.Post

		err = rows.StructScan(&post)

		if err != nil {
			p.logger.Error(
				"failed to scan row from query",
				slog.String("op", op),
				slog.String("query", queryGetListPosts),
				slog.Any("err", err),
			)
			return nil, fmt.Errorf("%w: %v", storage.ErrInternal, err)
		}

		result = append(result, post)
	}

	return result, nil
}

// ----------------------------
//		UPDATE POSTS
// ----------------------------

const queryUpdatePosts = `
UPDATE posts
SET data=$1
WHERE id=$2
RETURNING id;
`

// UpdatePost implements blog.BlogPosts.
func (p *Postgres) UpdatePost(ctx context.Context, request *domain.PostUpdateRequest) error {
	const op = "postgres.UpdatePost"

	p.logger.Info(fmt.Sprintf("executing query: %s", queryUpdatePosts), slog.String("op", op))

	res := p.db.QueryRowContext(ctx, queryUpdatePosts, request.Body, request.Id)

	var resId int64
	err := res.Scan(&resId)
	if err != nil {
		p.logger.Error(
			"failed executing query",
			slog.String("query", queryUpdatePosts),
			slog.String("op", op),
			slog.Any("err", err),
		)

		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf(
				"%w: %v",
				storage.ErrDoesNotExists,
				err,
			)
		}

		return fmt.Errorf(
			"%w: %v",
			storage.ErrInternal,
			err,
		)
	}

	return nil
}

// ----------------------------
//		UPDATE LIKES COUNT
// ----------------------------

const queryUpdateLikesCountOnPost = `
UPDATE posts
SET likes_count=likes_count + 1
WHERE id=$1
RETURNING id;
`

func (p *Postgres) UpdateLikesCountOnPost(ctx context.Context, postId int64, userName string) error {
	const op = "postgres.UpdateLikesCountOnPost"

	if _, err := p.getUserId(ctx, userName); err != nil {
		return err
	}

	p.logger.Info(fmt.Sprintf("executing query: %s", queryUpdateLikesCountOnPost), slog.String("op", op))

	res := p.db.QueryRowContext(ctx, queryUpdateLikesCountOnPost, postId)

	var resId int64
	err := res.Scan(&resId)
	if err != nil {
		p.logger.Error(
			"failed executing query",
			slog.String("query", queryUpdateLikesCountOnPost),
			slog.String("op", op),
			slog.Any("err", err),
		)

		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf(
				"%w: %v",
				storage.ErrDoesNotExists,
				err,
			)
		}

		return fmt.Errorf(
			"%w: %v",
			storage.ErrInternal,
			err,
		)
	}

	return nil
}
