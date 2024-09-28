package blog

import (
	"context"
	"errors"
	"github.com/Sant1s/blogBack/internal/domain"
	"log/slog"
)

var (
	ErrInternal = errors.New("iternal server error")
	ErrNotFound = errors.New("object not found")
	ErrExists   = errors.New("object already exists")
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
	UpdatePost(ctx context.Context, updates *domain.PostUpdateRequest) error
	DeletePost(ctx context.Context, postId int64) error
	UpdateLikesCountOnPost(ctx context.Context, postId int64, userName string) error

	GetListComments(ctx context.Context, limit, offset int32, postId int64) ([]domain.Comment, error)
	CreateComment(ctx context.Context, comment *domain.Comment) error
	UpdateComment(ctx context.Context, updates *domain.CommentUpdateRequest) error
	DeleteComment(ctx context.Context, commentId, postId int64) error
	UpdateLikesCountOnComment(ctx context.Context, commentId int64, userName string) error
}

// (data-layer redis)

type Likes interface {
	LikePost(ctx context.Context, userId string, postId int64) error
	LikeComment(ctx context.Context, userId string, commentId int64) error

	RollbackLikePost(ctx context.Context, userId string, postId int64) error
	RollbackLikeComment(ctx context.Context, userName string, commentId int64) error
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
