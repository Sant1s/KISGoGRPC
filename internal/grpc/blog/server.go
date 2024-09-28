package bloggrpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	"github.com/Sant1s/blogBack/internal/domain"
	blogService "github.com/Sant1s/blogBack/internal/gen"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
)

var (
	ErrInternal = status.Errorf(codes.InvalidArgument, "internal server error")
)

// todo: add service interface

type Blog interface {
	GetPosts(ctx context.Context, limit, offset int32) ([]domain.Post, error)
	CreatePost(ctx context.Context, post *domain.Post) error
	UpdatePost(ctx context.Context, post *domain.PostUpdateRequest) (int64, error)
	DeletePost(ctx context.Context, postId int64) error

	GetComments(ctx context.Context, limit, offset int32) ([]domain.Comment, error)
	CreateComment(ctx context.Context, comment *domain.Comment) error
	UpdateComment(ctx context.Context, comment *domain.CommentUpdateRequest) (int64, error)
	DeleteComment(ctx context.Context, commentId int64) error
}

type serverAPI struct {
	blogService.UnimplementedBlogServiceServer
	blog Blog
}

func Register(gRPC *grpc.Server, blog Blog) {
	blogService.RegisterBlogServiceServer(gRPC, &serverAPI{blog: blog})
}

func (s *serverAPI) ListPosts(ctx context.Context, req *blogService.ListPostsRequest) (*blogService.ListPostsResponse, error) {
	reqCtx, cancel := context.WithTimeout(ctx, time.Second*500)
	defer cancel()
	posts, err := s.blog.GetPosts(reqCtx, req.GetLimit(), req.GetOffset())
	if err != nil {
		//todo: правильно обработать ошибки
		return nil, ErrInternal
	}

	var resp blogService.ListPostsResponse

	for _, item := range posts {
		post := &blogService.Post{
			Id:           item.Id,
			Author:       item.Author,
			Body:         item.Body,
			CreateTime:   item.CreateTime.String(),
			CommentCount: item.CommentCount,
			LikesCount:   item.LikesCount,
		}
		resp.Posts = append(resp.Posts, post)
	}

	resp.Message = "ok"

	return &resp, nil
}

func (s *serverAPI) CreatePost(ctx context.Context, req *blogService.CreatePostRequest) (*blogService.Response, error) {
	reqCtx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	rand.Seed(uint64(time.Now().UnixNano()))
	randomInt := rand.Int63()

	post := domain.Post{
		Id:           randomInt,
		Author:       req.Author,
		Body:         req.Body,
		CommentCount: 0,
		LikesCount:   0,
	}

	err := s.blog.CreatePost(reqCtx, &post)
	if err != nil {
		return &blogService.Response{
			PostId:  0,
			Message: "error",
		}, ErrInternal
	}

	return &blogService.Response{
		PostId:  randomInt,
		Message: "ok",
	}, nil
}

func (s *serverAPI) UpdatePost(ctx context.Context, req *blogService.UpdatePostRequest) (*blogService.Response, error) {
	reqCtx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	post := &domain.PostUpdateRequest{
		Id:   req.Id,
		Body: req.Data,
	}

	postId, err := s.blog.UpdatePost(reqCtx, post)
	if err != nil {
		return &blogService.Response{
			PostId:  0,
			Message: "error",
		}, ErrInternal
	}

	return &blogService.Response{
		PostId:  postId,
		Message: "ok",
	}, nil
}

func (s *serverAPI) DeletePost(ctx context.Context, req *blogService.DeletePostRequest) (*blogService.Response, error) {
	reqCtx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	err := s.blog.DeletePost(reqCtx, req.PostId)
	if err != nil {
		return &blogService.Response{
			PostId:  req.PostId,
			Message: "error",
		}, ErrInternal
	}

	return &blogService.Response{
		PostId:  req.PostId,
		Message: "ok",
	}, nil
}

func (s *serverAPI) LikeComment(ctx context.Context, request *blogService.LikeCommentRequest) (*blogService.Response, error) {
	//TODO implement me
	fmt.Println("LikeComment")
	return nil, nil
}

func (s *serverAPI) ListComments(ctx context.Context, request *blogService.ListCommentsRequest) (*blogService.ListCommentsResponse, error) {
	//TODO implement me
	fmt.Println("LikeComment")
	return nil, nil
}

func (s *serverAPI) CreateComments(ctx context.Context, request *blogService.CreateCommentRequest) (*blogService.Response, error) {
	//TODO implement me
	fmt.Println("LikeComment")
	return nil, nil
}

func (s *serverAPI) UpdateComments(ctx context.Context, request *blogService.UpdateCommentRequest) (*blogService.Response, error) {
	//TODO implement me
	fmt.Println("LikeComment")
	return nil, nil
}

func (s *serverAPI) LikePost(ctx context.Context, request *blogService.LikePostRequest) (*blogService.Response, error) {
	//TODO implement me
	fmt.Println("LikeComment")
	return nil, nil
}

func (s *serverAPI) DeleteComment(ctx context.Context, request *blogService.DeleteCommentRequest) (*blogService.Response, error) {
	//TODO implement me
	fmt.Println("LikeComment")
	return nil, nil
}
