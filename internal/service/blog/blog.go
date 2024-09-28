package blog

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Sant1s/blogBack/internal/domain"
)

type Service struct {
	log       *slog.Logger
	blogPosts Posts
	blogLikes Likes
}

// (data-layer postgres)

type Posts interface {
	GetListPosts(ctx context.Context, limit, offset int32) ([]domain.Post, error)
	CreatePost(ctx context.Context, post *domain.Post) error
	UpdatePost(ctx context.Context, updates *domain.PostUpdateRequest) (int64, error)
	DeletePost(ctx context.Context, postId int64) error

	GetListComments(ctx context.Context, limit, offset int32) ([]domain.Comment, error)
	CreateComment(ctx context.Context, comment *domain.Comment) error
	UpdateComment(ctx context.Context, updates *domain.CommentUpdateRequest) (int64, error)
	DeleteComment(ctx context.Context, commentId int64) error
}

type Likes interface {
	// todo: дописать (data-layer redis)
}

func New(
	log *slog.Logger,
	blogPosts Posts,
	blogLikes Likes,
) *Service {
	return &Service{
		log:       log,
		blogPosts: blogPosts,
		blogLikes: blogLikes,
	}
}

func (b *Service) GetPosts(ctx context.Context, limit, offset int32) ([]domain.Post, error) {
	posts, err := b.blogPosts.GetListPosts(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return posts, nil
}

func (b *Service) CreatePost(ctx context.Context, post *domain.Post) error {
	err := b.blogPosts.CreatePost(ctx, post)
	if err != nil {
		return fmt.Errorf("database error %w", err)
	}
	return nil
}

func (b *Service) UpdatePost(ctx context.Context, post *domain.PostUpdateRequest) (int64, error) {
	id, err := b.blogPosts.UpdatePost(ctx, post)
	if err != nil {
		return 0, fmt.Errorf("database error: %w", err)
	}
	return id, nil
}

func (b *Service) DeletePost(ctx context.Context, postId int64) error {
	err := b.blogPosts.DeletePost(ctx, postId)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

func (b *Service) GetComments(ctx context.Context, limit, offset int32) ([]domain.Comment, error) {
	//TODO implement me
	fmt.Println("GetComments")
	return nil, nil
}

func (b *Service) CreateComment(ctx context.Context, comment *domain.Comment) error {
	//TODO implement me
	fmt.Println("CreateComment")
	return nil
}

func (b *Service) UpdateComment(ctx context.Context, comment *domain.CommentUpdateRequest) (int64, error) {
	//TODO implement me
	fmt.Println("UpdateComment")
	return 0, nil
}

func (b *Service) DeleteComment(ctx context.Context, commentId int64) error {
	//TODO implement me
	fmt.Println("DeleteComment")
	return nil
}
