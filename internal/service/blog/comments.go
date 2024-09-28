package blog

import (
	"context"
	"errors"
	"github.com/Sant1s/blogBack/internal/domain"
	"github.com/Sant1s/blogBack/internal/storage"
)

func (b *Service) GetComments(ctx context.Context, limit, offset int32, postId int64) ([]domain.Comment, error) {
	res, err := b.blogPosts.GetListComments(ctx, limit, offset, postId)

	if err != nil {
		return nil, ErrInternal
	}

	return res, nil
}

func (b *Service) CreateComment(ctx context.Context, comment *domain.Comment) error {
	err := b.blogPosts.CreateComment(ctx, comment)
	if err != nil {
		return ErrInternal
	}
	return nil
}

func (b *Service) UpdateComment(ctx context.Context, comment *domain.CommentUpdateRequest) error {
	err := b.blogPosts.UpdateComment(ctx, comment)
	if err != nil {
		if errors.Is(err, storage.ErrDoesNotExists) {
			return ErrNotFound
		}

		return ErrInternal
	}

	return nil
}

func (b *Service) DeleteComment(ctx context.Context, commentId, postId int64) error {
	err := b.blogPosts.DeleteComment(ctx, commentId, postId)

	if err != nil {
		if errors.Is(err, storage.ErrDoesNotExists) {
			return ErrNotFound
		}
		return ErrInternal
	}

	return nil
}
