package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Sant1s/blogBack/internal/storage"
	"github.com/jackc/pgx/v5"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Sant1s/blogBack/internal/domain"
)

// ----------------------------
//		GET LIST COMMENTS
// ----------------------------

const queryGetListComments = `
SELECT
    comments.id AS id,
    u.nickname AS nickname,
    p.id AS post_id,
    comments.data AS data,
    comments.created_at AS created_at,
    comments.comments_count AS comments_count,
    comments.likes_count AS likes_count,
    comments.parent_id AS parent_id
FROM comments
JOIN users u ON u.id=comments.author_id
JOIN posts p ON p.id=comments.post_id
WHERE post_id=$1
LIMIT $2 OFFSET $3;
`

func (p *Postgres) GetListComments(ctx context.Context, limit, offset int32, postId int64) ([]domain.Comment, error) {
	const op = "postgres.GetListComments"

	p.logger.Info(fmt.Sprintf("executing query: %s", queryGetListComments), slog.String("op", op))

	rows, err := p.db.QueryxContext(ctx, queryGetListComments, postId, limit, offset)
	if err != nil {
		p.logger.Error(
			"failed to execute query",
			slog.String("op", op),
			slog.String("query", queryGetListPosts),
			slog.Any("err", err),
		)

		return nil, fmt.Errorf("%w: %v", storage.ErrInternal, err)
	}

	result := make([]domain.Comment, 0)
	for rows.Next() {
		var comment domain.Comment

		err = rows.StructScan(&comment)

		if err != nil {
			p.logger.Error(
				"failed to scan row from query",
				slog.String("op", op),
				slog.String("query", queryGetListComments),
				slog.Any("err", err),
			)
			return nil, fmt.Errorf("%w: %v", storage.ErrInternal, err)
		}

		result = append(result, comment)
	}

	return result, nil
}

// ----------------------------
//		CREATE COMMENT
// ----------------------------

const queryCreateComment = `
INSERT INTO comments (id, author_id, parent_id, post_id, data)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;
`

const queryUpdatePostCommentCount = `
UPDATE posts
SET comments_count = comments_count + 1
WHERE id = $1;
`

func (p *Postgres) CreateComment(ctx context.Context, comment *domain.Comment) error {
	const op = "postgres.CreateComment"

	tx, err := p.db.Beginx()
	if err != nil {
		p.logger.Error(
			"failed creating transaction",
			slog.String("op", op),
			slog.Any("err", err),
		)

		return fmt.Errorf("%w: %v", storage.ErrInternal, err)
	}

	row := tx.QueryRowContext(ctx, queryGetAuthorId, comment.Author)

	var uuid string
	if err = row.Scan(&uuid); err != nil {
		p.logger.Error(
			"user with that nickname not found",
			slog.String("query", queryGetAuthorId),
			slog.String("op", op),
			slog.Any("err", err),
		)

		errRollback := tx.Rollback()
		if errRollback != nil {
			p.logger.Error(
				"error rollback transaction",
				slog.String("op", op),
				slog.Any("err", errRollback),
			)
			panic(err) // Паника просто по-приколу. Автор трезв, просто ему это показалось смешнявкой
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w: %v", storage.ErrDoesNotExists, err)
		}

		return fmt.Errorf("%w: %v", storage.ErrInternal, err)
	}

	// execute create query
	res := tx.QueryRowContext(
		ctx,
		queryCreateComment,
		comment.Id,
		uuid,
		comment.ParentId,
		comment.PostId,
		comment.Body,
	)
	var id int64

	err = res.Scan(&id)

	if err != nil {
		p.logger.Error(
			"failed executing query",
			slog.String("query", queryGetAuthorId),
			slog.String("op", op),
			slog.Any("err", err),
		)

		errRollback := tx.Rollback()
		if errRollback != nil {
			p.logger.Error(
				"error rollback transaction",
				slog.String("op", op),
				slog.Any("err", errRollback),
			)
			panic(err) // Паника просто по-приколу. Автор трезв, просто ему это показалось смешнявкой
		}

		return fmt.Errorf("%w: %v", storage.ErrInternal, err)
	}

	_, err = tx.ExecContext(ctx, queryUpdatePostCommentCount, comment.PostId)
	if err != nil {
		p.logger.Error(
			"failed executing query",
			slog.String("query", queryGetAuthorId),
			slog.String("op", op),
			slog.Any("err", err),
		)

		errRollback := tx.Rollback()
		if errRollback != nil {
			p.logger.Error(
				"error rollback transaction",
				slog.String("op", op),
				slog.Any("err", errRollback),
			)
			panic(err) // Паника просто по-приколу. Автор трезв, просто ему это показалось смешнявкой
		}

		return fmt.Errorf("%w: %v", storage.ErrInternal, err)
	}

	err = tx.Commit()
	if err != nil {
		p.logger.Error(
			"commit transaction error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		return fmt.Errorf("%w: %v", storage.ErrInternal, err)
	}

	return nil
}

// ----------------------------
//		UPDATE COMMENT
// ----------------------------

const queryUpdateComment = `
UPDATE comments
SET data=$1
WHERE id=$2
RETURNING id;
`

func (p *Postgres) UpdateComment(ctx context.Context, updates *domain.CommentUpdateRequest) error {
	const op = "postgres.UpdateComment"

	p.logger.Info(fmt.Sprintf("executing query: %s", queryUpdateComment), slog.String("op", op))

	res := p.db.QueryRowContext(ctx, queryUpdateComment, updates.Body, updates.Id)

	var resId int64
	err := res.Scan(&resId)
	if err != nil {
		p.logger.Error(
			"failed executing query",
			slog.String("query", queryUpdateComment),
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
//		DELETE COMMENT
// ----------------------------

const queryPostCommentCountUpdate = `
UPDATE posts
SET comments_count = comments_count - 1
WHERE id = $1 AND comments_count > 0;
`
const queryDeleteComment = `DELETE FROM comments WHERE id=$1;`

func (p *Postgres) DeleteComment(ctx context.Context, commentId, postId int64) error {
	const op = "postgres.DeleteComment"

	tx, err := p.db.Beginx()
	if err != nil {
		p.logger.Error(
			"failed creating transaction",
			slog.String("op", op),
			slog.Any("err", err),
		)

		return fmt.Errorf("%w: %v", storage.ErrInternal, err)
	}

	p.logger.Info(
		fmt.Sprintf(
			"executing transaction with queries: %s, %s",
			queryPostCommentCountUpdate,
			queryDeleteComment),
		slog.String("op", op),
	)

	_, err = tx.ExecContext(ctx, queryPostCommentCountUpdate, postId)
	if err != nil {
		p.logger.Error(
			"failed executing query",
			slog.String("query", queryPostCommentCountUpdate),
			slog.String("op", op),
			slog.Any("err", err),
		)

		errRollback := tx.Rollback()
		if errRollback != nil {
			p.logger.Error(
				"GG WP!",
				slog.String("op", op),
				slog.Any("err", err),
			)

			panic(err) // Паника просто по-приколу. Автор трезв, просто ему это показалось смешнявкой
		}

		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w: %v", storage.ErrDoesNotExists, err)
		}
		return fmt.Errorf("%w: %v", storage.ErrInternal, err)
	}

	row, err := tx.ExecContext(ctx, queryDeleteComment, commentId)
	if rowsAffected, errRows := row.RowsAffected(); err != nil || errRows != nil || rowsAffected == 0 {
		p.logger.Error(
			"failed executing query",
			slog.String("query", queryDeleteComment),
			slog.String("op", op),
			slog.Any("err", err),
		)

		errRollback := tx.Rollback()
		if errRollback != nil {
			p.logger.Error(
				"GG WP!",
				slog.String("op", op),
				slog.Any("err", err),
			)

			panic(err) // Паника просто по-приколу. Автор трезв, просто ему это показалось смешнявкой
		}

		if rowsAffected == 0 {
			return fmt.Errorf("%w: %v", storage.ErrDoesNotExists, err)
		}
		return fmt.Errorf("%w: %v", storage.ErrInternal, err)
	}

	err = tx.Commit()
	if err != nil {
		p.logger.Error(
			"commit transaction error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		return fmt.Errorf("%w: %v", storage.ErrInternal, err)
	}

	return nil
}

// ----------------------------
//		UPDATE LIKES COUNT
// ----------------------------

const queryUpdateLikesCountOnComment = `
UPDATE comments
SET likes_count=likes_count + $1
WHERE id=$2 AND likes_count >= 0
RETURNING id;
`

func (p *Postgres) UpdateLikesCountOnComment(ctx context.Context, commentId int64, userName string, delta int64) error {
	const op = "postgres.UpdateLikesCountOnComment"

	if _, err := p.getUserId(ctx, userName); err != nil {
		return err
	}

	p.logger.Info(fmt.Sprintf("executing query: %s", queryUpdateLikesCountOnComment), slog.String("op", op))

	res := p.db.QueryRowContext(ctx, queryUpdateLikesCountOnComment, delta, commentId)

	var resId int64
	err := res.Scan(&resId)
	if err != nil {
		p.logger.Error(
			"failed executing query",
			slog.String("query", queryUpdateLikesCountOnComment),
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
