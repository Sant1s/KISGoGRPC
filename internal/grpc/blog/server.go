package bloggrpc

import (
	"context"
	"github.com/Sant1s/blogBack/internal/domain"
	blogService "github.com/Sant1s/blogBack/internal/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

var (
	ErrInternal      = status.Errorf(codes.InvalidArgument, "internal server error")
	ErrNotFound      = status.Errorf(codes.NotFound, "object not found")
	ErrAlreadyExists = status.Errorf(codes.AlreadyExists, "object alredy exists")
)

// todo: add service interface

type Blog interface {
	GetPosts(ctx context.Context, limit, offset int32) ([]domain.Post, error)
	CreatePost(ctx context.Context, post *domain.Post) error
	UpdatePost(ctx context.Context, post *domain.PostUpdateRequest) error
	DeletePost(ctx context.Context, postId int64) error

	GetComments(ctx context.Context, limit, offset int32, postId int64) ([]domain.Comment, error)
	CreateComment(ctx context.Context, comment *domain.Comment) error
	UpdateComment(ctx context.Context, comment *domain.CommentUpdateRequest) error
	DeleteComment(ctx context.Context, commentId, postId int64) error
}

type Likes interface {
	LikePost(ctx context.Context, userName string, postId int64) error
	LikeComment(ctx context.Context, userName string, commentId int64) error

	RemoveLikePost(ctx context.Context, userName string, postId int64) error
	RemoveLikeComment(ctx context.Context, userName string, commentId int64) error
}

type serverAPI struct {
	logger *slog.Logger
	blogService.UnimplementedBlogServiceServer
	blog  Blog
	likes Likes
}

func Register(l *slog.Logger, gRPC *grpc.Server, blog Blog, likes Likes) {
	blogService.RegisterBlogServiceServer(gRPC, &serverAPI{
		logger: l,
		blog:   blog,
		likes:  likes,
	})
}
