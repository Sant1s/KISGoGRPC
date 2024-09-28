package blog

import (
	"context"
	"errors"
	"github.com/Sant1s/blogBack/internal/storage"
)

func (b *Service) LikePost(ctx context.Context, userName string, postId int64) error {
	err := b.blogLikes.LikePost(ctx, userName, postId)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return ErrExists
		}
		return ErrInternal
	}

	err = b.blogPosts.UpdateLikesCountOnPost(ctx, postId, userName)
	if err != nil {
		if errRollback := b.blogLikes.RollbackLikePost(ctx, userName, postId); errRollback != nil {
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
	err := b.blogLikes.LikeComment(ctx, userName, commentId)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return ErrExists
		}
		return ErrInternal
	}

	err = b.blogPosts.UpdateLikesCountOnComment(ctx, commentId, userName)
	if err != nil {
		if errRollback := b.blogLikes.RollbackLikeComment(ctx, userName, commentId); errRollback != nil {
			return ErrInternal
		}

		if errors.Is(err, storage.ErrDoesNotExists) {
			return ErrNotFound
		}

		return ErrInternal
	}

	return nil
}
