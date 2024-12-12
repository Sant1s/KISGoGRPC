package blog

import (
	"context"
	"errors"
	"github.com/Sant1s/blogBack/internal/storage"
	"time"
)

func (b *Service) LikePost(ctx context.Context, userName string, postId int64) error {
	ctxDB, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()

	err := b.blogLikes.LikePost(ctxDB, userName, postId)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return ErrExists
		}
		return ErrInternal
	}

	err = b.blogPosts.UpdateLikesCountOnPost(ctx, postId, userName, 1)

	ctxDB, cancel = context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()
	if err != nil {
		if errRollback := b.blogLikes.RollbackLikePost(ctxDB, userName, postId); errRollback != nil {
			return ErrInternal
		}

		if errors.Is(err, storage.ErrDoesNotExists) {
			return ErrNotFound
		}

		return ErrInternal
	}

	return nil
}

func (b *Service) LikeComment(ctx context.Context, userName string, commentId int64) error {
	ctxDB, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()

	err := b.blogLikes.LikeComment(ctxDB, userName, commentId)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return ErrExists
		}
		return ErrInternal
	}

	err = b.blogPosts.UpdateLikesCountOnComment(ctx, commentId, userName, 1)

	ctxDB, cancel = context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()
	if err != nil {
		if errRollback := b.blogLikes.RollbackLikeComment(ctxDB, userName, commentId); errRollback != nil {
			return ErrInternal
		}

		if errors.Is(err, storage.ErrDoesNotExists) {
			return ErrNotFound
		}

		return ErrInternal
	}

	return nil
}

func (b *Service) RemoveLikePost(ctx context.Context, userName string, postId int64) error {
	ctxDB, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()

	err := b.blogLikes.RollbackLikePost(ctxDB, userName, postId)
	if err != nil {
		return ErrInternal
	}

	ctxDB, cancel = context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()

	err = b.blogPosts.UpdateLikesCountOnPost(ctx, postId, userName, -1)
	if err != nil {
		if errRollback := b.blogLikes.LikePost(ctxDB, userName, postId); errRollback != nil {
			return ErrInternal
		}

		if errors.Is(err, storage.ErrDoesNotExists) {
			return ErrNotFound
		}

		return ErrInternal
	}

	return nil
}

func (b *Service) RemoveLikeComment(ctx context.Context, userName string, commentId int64) error {
	ctxDB, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()

	err := b.blogLikes.RollbackLikeComment(ctxDB, userName, commentId)
	if err != nil {
		return ErrInternal
	}

	ctxDB, cancel = context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()

	err = b.blogPosts.UpdateLikesCountOnComment(ctx, commentId, userName, -1)
	if err != nil {
		if errRollback := b.blogLikes.LikeComment(ctxDB, userName, commentId); errRollback != nil {
			return ErrInternal
		}

		if errors.Is(err, storage.ErrDoesNotExists) {
			return ErrNotFound
		}

		return ErrInternal
	}

	return nil
}
