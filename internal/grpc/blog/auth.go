package bloggrpc

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Sant1s/blogBack/internal/domain"
	blogService "github.com/Sant1s/blogBack/internal/gen"
	"github.com/Sant1s/blogBack/internal/service/auth"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
)

func (s *serverAPI) Register(ctx context.Context, request *blogService.RegisterRequest) (*blogService.RegisterResponse, error) {
	const op = "bloggrpc.Register"

	passwordHash, err := HashPassword(request.GetPassword())
	if err != nil {
		s.logger.Error(
			"executing failed with error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		return &blogService.RegisterResponse{
			Code:   int64(codes.Internal),
			Output: err.Error(),
		}, ErrInternal
	}

	reqCtx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	response, err := s.auth.RegisterUser(reqCtx, &domain.RegisterUserRequest{
		Login:        request.GetLogin(),
		PasswordHash: passwordHash,
	})

	if err != nil {
		s.logger.Error(
			"executing failed with error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		if errors.Is(err, auth.ErrExists) {
			return &blogService.RegisterResponse{
				Code:   int64(codes.NotFound),
				Output: err.Error(), // todo: мб надо отправлять что-то менее открытое
			}, ErrAlreadyExists
		}

		return &blogService.RegisterResponse{
			Code:   int64(codes.Internal),
			Output: err.Error(), // todo: мб надо отправлять что-то менее открытое
		}, ErrInternal
	}

	return &blogService.RegisterResponse{
		Code:   int64(codes.OK),
		Output: response.Output,
	}, nil
}

func (s *serverAPI) Login(ctx context.Context, request *blogService.LoginRequest) (*blogService.LoginResponse, error) {
	const op = "bloggrpc.Login"

	passwordHash, err := HashPassword(request.GetPassword())
	if err != nil {
		s.logger.Error(
			"executing failed with error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		return &blogService.LoginResponse{
			Code:   int64(codes.Internal),
			Output: err.Error(),
		}, ErrInternal
	}

	reqCtx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	response, err := s.auth.LoginUser(reqCtx, &domain.LoginUserRequest{
		Login:        request.GetLogin(),
		PasswordHash: passwordHash,
	})

	if err != nil {
		s.logger.Error(
			"executing failed with error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		if errors.Is(err, auth.ErrNotFound) {
			return &blogService.LoginResponse{
				Code:   int64(codes.AlreadyExists),
				Output: err.Error(), // todo: мб надо отправлять что-то менее открытое
			}, ErrNotFound
		}

		return &blogService.LoginResponse{
			Code:   int64(codes.Internal),
			Output: err.Error(), // todo: мб надо отправлять что-то менее открытое
		}, ErrInternal
	}

	return &blogService.LoginResponse{
		Code:   int64(codes.OK),
		Output: response.Output,
	}, nil
}

func HashPassword(password string) (string, error) {
	cost := 10
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
