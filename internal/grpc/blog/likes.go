package bloggrpc

import (
	"context"
	"errors"
	blogService "github.com/Sant1s/blogBack/internal/gen"
	"github.com/Sant1s/blogBack/internal/service/blog"
	"log/slog"
	"time"
)

func (s *serverAPI) RemoveLikePost(ctx context.Context, request *blogService.LikePostRequest) (*blogService.Response, error) {
	const op = "bloggrpc.RemoveLikePost"

	reqCtx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	err := s.likes.RemoveLikePost(reqCtx, request.GetUserName(), request.GetPostId())
	if err != nil {
		s.logger.Error(
			"executing failed with error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		if errors.Is(err, blog.ErrNotFound) {
			return &blogService.Response{
				Id:      request.GetPostId(),
				Message: ErrNotFound.Error(),
			}, ErrNotFound
		}

		return &blogService.Response{
			Id:      request.GetPostId(),
			Message: ErrInternal.Error(),
		}, ErrInternal
	}

	return &blogService.Response{
		Id:      request.GetPostId(),
		Message: "ok",
	}, nil
}

func (s *serverAPI) RemoveLikeComment(ctx context.Context, request *blogService.LikeCommentRequest) (*blogService.Response, error) {
	const op = "bloggrpc.RemoveLikeComment"

	reqCtx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	err := s.likes.RemoveLikeComment(reqCtx, request.GetUserName(), request.GetCommentId())
	if err != nil {
		s.logger.Error(
			"executing failed with error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		if errors.Is(err, blog.ErrNotFound) {
			return &blogService.Response{
				Id:      request.GetCommentId(),
				Message: ErrNotFound.Error(),
			}, ErrNotFound
		}
		return nil, ErrInternal
	}

	return &blogService.Response{
		Id:      request.GetCommentId(),
		Message: "ok",
	}, nil
}
