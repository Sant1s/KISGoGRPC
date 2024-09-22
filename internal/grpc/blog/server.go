package bloggrpc

import (
	"context"
	"fmt"
	"time"

	"github.com/Sant1s/blogBack/internal/domain"
	blogService "github.com/Sant1s/blogBack/internal/gen"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
)

var (
	errInternal = fmt.Errorf("internal server error")
)

type Blog interface {
	// todo: add service interface
	GetPosts(ctx context.Context, limit, offset int32) (domain.Posts, error)
	CreatePost(ctx context.Context, post domain.Post) error
	UpdatePost(ctx context.Context, post domain.Post) (int64, error)
	DeletePost(ctx context.Context, post_id int64) error
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
		return &blogService.ListPostsResponse{}, errInternal
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
			Liked:        item.LikedByUser,
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
		LikedByUser:  false,
	}

	err := s.blog.CreatePost(reqCtx, post)
	if err != nil {
		return &blogService.Response{
			PostId:  0,
			Message: "error",
		}, errInternal
	}

	return &blogService.Response{
		PostId:  randomInt,
		Message: "ok",
	}, nil
}

func (s *serverAPI) UpdatePost(ctx context.Context, req *blogService.UpdatePostRequest) (*blogService.Response, error) {
	reqCtx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	post := domain.Post{
		Author:       req.GetAuthor(),
		Body:         req.GetBody(),
		CommentCount: req.GetCommentCount(),
		LikesCount:   req.GetLikesCount(),
		LikedByUser:  req.GetLiked(),
	}

	postId, err := s.blog.UpdatePost(reqCtx, post)
	if err != nil {
		return &blogService.Response{
			PostId:  0,
			Message: "error",
		}, errInternal
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
		}, errInternal
	}

	return &blogService.Response{
		PostId:  req.PostId,
		Message: "ok",
	}, nil
}
