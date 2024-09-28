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
//	LIKES SERVICE IMPLEMENTATION
// ----------------------------

func (s *serverAPI) ListComments(ctx context.Context, request *blogService.ListCommentsRequest) (*blogService.ListCommentsResponse, error) {
	const op = "bloggrpc.ListComments"

	reqCtx, cancel := context.WithTimeout(ctx, time.Second*500)
	defer cancel()

	comments, err := s.blog.GetComments(reqCtx, request.GetLimit(), request.GetOffset(), request.GetPostId())
	if err != nil {
		s.logger.Error(
			"executing failed with error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		return &blogService.ListCommentsResponse{
			Comments: nil,
			Message:  "internal server error",
		}, ErrInternal
	}

	var resp blogService.ListCommentsResponse

	for _, item := range comments {
		post := &blogService.Comment{
			Id:           item.Id,
			Author:       item.Author,
			Body:         item.Body,
			CreateTime:   item.CreateTime.String(),
			CommentCount: item.CommentCount,
			LikesCount:   item.LikesCount,
			ParentId:     item.ParentId,
		}

		resp.Comments = append(resp.Comments, post)
	}

	resp.Message = "ok"

	return &resp, nil
}

func (s *serverAPI) CreateComments(ctx context.Context, request *blogService.CreateCommentRequest) (*blogService.Response, error) {
	const op = "blogprpc.CreateComments"

	reqCtx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	rand.Seed(uint64(time.Now().UnixNano()))
	randomInt := rand.Int63()

	comment := domain.Comment{
		Id:           randomInt,
		Author:       request.GetAuthor(),
		PostId:       request.GetPostId(),
		Body:         request.GetBody(),
		CreateTime:   time.Now(),
		CommentCount: 0,
		LikesCount:   0,
		ParentId:     request.GetParentId(),
	}

	err := s.blog.CreateComment(reqCtx, &comment)
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

func (s *serverAPI) UpdateComments(ctx context.Context, request *blogService.UpdateCommentRequest) (*blogService.Response, error) {
	const op = "bloggrpc.UpdateComments"

	reqCtx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	commentUpdate := &domain.CommentUpdateRequest{
		Id:   request.GetCommentId(),
		Body: request.GetBody(),
	}

	err := s.blog.UpdateComment(reqCtx, commentUpdate)
	if err != nil {
		s.logger.Error(
			"executing failed with error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		if errors.Is(err, storage.ErrDoesNotExists) {
			return &blogService.Response{
				Id:      request.GetCommentId(),
				Message: "error not found",
			}, ErrNotFound
		}

		return &blogService.Response{
			Id:      request.GetCommentId(),
			Message: "internal server error",
		}, ErrInternal
	}

	return &blogService.Response{
		Id:      request.GetCommentId(),
		Message: "ok",
	}, nil
}

func (s *serverAPI) DeleteComment(ctx context.Context, request *blogService.DeleteCommentRequest) (*blogService.Response, error) {
	const op = "bloggrpc.DeleteComment"

	reqCtx, cancel := context.WithTimeout(ctx, time.Hour*500)
	defer cancel()

	err := s.blog.DeleteComment(reqCtx, request.GetCommentId(), request.GetPostId())
	if err != nil {
		s.logger.Error(
			"executing failed with error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		if errors.Is(err, storage.ErrDoesNotExists) {
			return &blogService.Response{
				Id:      request.GetCommentId(),
				Message: "error not found",
			}, ErrNotFound
		}

		return &blogService.Response{
			Id:      request.GetCommentId(),
			Message: "internal server error",
		}, ErrInternal
	}

	return &blogService.Response{
		Id:      request.GetCommentId(),
		Message: "ok",
	}, nil
}

func (s *serverAPI) LikeComment(ctx context.Context, request *blogService.LikeCommentRequest) (*blogService.Response, error) {
	const op = "bloggrpc.LikeComment"

	reqCtx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	err := s.likes.LikeComment(reqCtx, request.GetUserName(), request.GetCommentId())
	if err != nil {
		s.logger.Error(
			"executing failed with error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		if errors.Is(err, blog.ErrExists) {
			return &blogService.Response{
				Id:      request.GetCommentId(),
				Message: "error: post not found",
			}, ErrAlreadyExists
		}

		if errors.Is(err, blog.ErrNotFound) {
			return &blogService.Response{
				Id:      request.GetCommentId(),
				Message: "error: post not found",
			}, ErrNotFound
		}

		return &blogService.Response{
			Id:      request.GetCommentId(),
			Message: "internal server error",
		}, ErrInternal
	}

	return &blogService.Response{
		Id:      request.GetCommentId(),
		Message: "ok",
	}, nil
}
