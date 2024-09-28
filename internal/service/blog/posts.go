package blog

import (
	"context"
	"errors"
	"github.com/Sant1s/blogBack/internal/domain"
	"github.com/Sant1s/blogBack/internal/storage"
)

func (b *Service) GetPosts(ctx context.Context, limit, offset int32) ([]domain.Post, error) {
	posts, err := b.blogPosts.GetListPosts(ctx, limit, offset)

	if err != nil {
		return nil, ErrInternal
	}

	return posts, nil
}

func (b *Service) CreatePost(ctx context.Context, post *domain.Post) error {
	err := b.blogPosts.CreatePost(ctx, post)

	if err != nil {
		return ErrInternal
	}
	return nil
}

func (b *Service) UpdatePost(ctx context.Context, post *domain.PostUpdateRequest) error {
	err := b.blogPosts.UpdatePost(ctx, post)

	if err != nil {
		if errors.Is(err, storage.ErrDoesNotExists) {
			return ErrNotFound
		}
		return ErrInternal
	}
	return nil
}

func (b *Service) DeletePost(ctx context.Context, postId int64) error {
	err := b.blogPosts.DeletePost(ctx, postId)

	if err != nil {
		if errors.Is(err, storage.ErrDoesNotExists) {
			return ErrNotFound
		}

		return ErrInternal
	}
	return nil
}
