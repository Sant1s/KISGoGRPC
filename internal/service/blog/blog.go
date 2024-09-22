package blog

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Sant1s/blogBack/internal/domain"
)

type BlogService struct {
	log       *slog.Logger
	blogPosts BlogPosts
	blogLikes BlogLikes
}

type BlogPosts interface {
	//! (data-layer postgres)
	GetListPosts(ctx context.Context, limit, offset int32) (domain.Posts, error)
	CreatePost(ctx context.Context, post *domain.Post) error
	UpdatePost(ctx context.Context, updates *domain.Post) (int64, error)
	DeletePost(ctx context.Context, post_id int64) error
}

type BlogLikes interface {
	// todo: дописать (data-layer redis)
}

func New(
	log *slog.Logger,
	blogPosts BlogPosts,
	blogLikes BlogLikes,
) *BlogService {
	return &BlogService{
		log:       log,
		blogPosts: blogPosts,
		blogLikes: blogLikes,
	}
}

func (b *BlogService) GetPosts(ctx context.Context, limit, offset int32) (domain.Posts, error) {
	posts, err := b.blogPosts.GetListPosts(ctx, limit, offset)
	if err != nil {
		return domain.Posts{}, fmt.Errorf("database error: %w", err)
	}

	return posts, nil
}

func (b *BlogService) CreatePost(ctx context.Context, post domain.Post) error {
	err := b.blogPosts.CreatePost(ctx, &post)
	if err != nil {
		return fmt.Errorf("database error %w", err)
	}
	return nil
}

func (b *BlogService) UpdatePost(ctx context.Context, post domain.Post) (int64, error) {
	id, err := b.blogPosts.UpdatePost(ctx, &post)
	if err != nil {
		return 0, fmt.Errorf("database error: %w", err)
	}
	return id, nil
}

func (b *BlogService) DeletePost(ctx context.Context, post_id int64) error {
	err := b.blogPosts.DeletePost(ctx, post_id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}
