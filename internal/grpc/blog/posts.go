package bloggrpc

import (
	"context"
	"errors"
	"github.com/Sant1s/blogBack/internal/domain"
	blogService "github.com/Sant1s/blogBack/internal/gen"
	"github.com/Sant1s/blogBack/internal/service/blog"
	"github.com/Sant1s/blogBack/internal/storage"
	"golang.org/x/exp/rand"
	"log/slog"
	"time"
)

// ----------------------------
// 	POST SERVICE IMPLEMENTATION
// ----------------------------

func (s *serverAPI) ListPosts(ctx context.Context, req *blogService.ListPostsRequest) (*blogService.ListPostsResponse, error) {
	const op = "bloggrpc.ListPosts"

	reqCtx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	posts, err := s.blog.GetPosts(reqCtx, req.GetLimit(), req.GetOffset())
	if err != nil {
		s.logger.Error(
			"executing failed with error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		return &blogService.ListPostsResponse{
			Posts:   nil,
			Message: "internal server error",
		}, ErrInternal
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
	const op = "blogprpc.CreatePost"

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
		s.logger.Error(
			"executing failed with error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		return &blogService.Response{
			Id:      0,
			Message: "internal server error",
		}, ErrInternal
	}

	return &blogService.Response{
		Id:      randomInt,
		Message: "ok",
	}, nil
}

func (s *serverAPI) UpdatePost(ctx context.Context, req *blogService.UpdatePostRequest) (*blogService.Response, error) {
	const op = "bloggrpc.UpdatePost"

	reqCtx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	post := &domain.PostUpdateRequest{
		Id:   req.Id,
		Body: req.Data,
	}

	err := s.blog.UpdatePost(reqCtx, post)
	if err != nil {
		s.logger.Error(
			"executing failed with error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		if errors.Is(err, storage.ErrDoesNotExists) {
			return &blogService.Response{
				Id:      req.GetId(),
				Message: "error not found",
			}, ErrNotFound
		}

		return &blogService.Response{
			Id:      req.GetId(),
			Message: "internal server error",
		}, ErrInternal
	}

	return &blogService.Response{
		Id:      req.GetId(),
		Message: "ok",
	}, nil
}

func (s *serverAPI) DeletePost(ctx context.Context, req *blogService.DeletePostRequest) (*blogService.Response, error) {
	const op = "bloggrpc.DeletePost"

	reqCtx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	err := s.blog.DeletePost(reqCtx, req.PostId)
	if err != nil {
		s.logger.Error(
			"executing failed with error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		if errors.Is(err, storage.ErrDoesNotExists) {
			return &blogService.Response{
				Id:      0,
				Message: "error not found",
			}, ErrNotFound
		}

		return &blogService.Response{
			Id:      req.PostId,
			Message: "internal server error",
		}, ErrInternal
	}

	return &blogService.Response{
		Id:      req.PostId,
		Message: "ok",
	}, nil
}

func (s *serverAPI) LikePost(ctx context.Context, request *blogService.LikePostRequest) (*blogService.Response, error) {
	const op = "bloggrpc.LikePost"

	reqCtx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	err := s.likes.LikePost(reqCtx, request.GetUserName(), request.GetPostId())
	if err != nil {
		s.logger.Error(
			"executing failed with error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		if errors.Is(err, blog.ErrNotFound) {
			return &blogService.Response{
				Id:      request.GetPostId(),
				Message: "error: post not found",
			}, ErrNotFound
		}

		if errors.Is(err, blog.ErrExists) {
			return &blogService.Response{
				Id:      request.GetPostId(),
				Message: "error: like on post already exists",
			}, ErrAlreadyExists
		}

		return &blogService.Response{
			Id:      request.GetPostId(),
			Message: "internal server error",
		}, ErrInternal
	}

	return &blogService.Response{
		Id:      request.GetPostId(),
		Message: "ok",
	}, nil
}
