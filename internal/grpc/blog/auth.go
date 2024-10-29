package bloggrpc

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/Sant1s/blogBack/internal/domain"
	blogService "github.com/Sant1s/blogBack/internal/gen"
	"github.com/Sant1s/blogBack/internal/service/auth"
)

func (s *serverAPI) Register(ctx context.Context, request *blogService.RegisterRequest) (*blogService.RegisterResponse, error) {
	const op = "bloggrpc.Register"

	//passwordHash, err := HashPassword(request.GetPassword())
	//if err != nil {
	//	s.logger.Error(
	//		"executing failed with error",
	//		slog.String("op", op),
	//		slog.Any("err", err),
	//	)
	//
	//	return nil, ErrInternal
	//}

	reqCtx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	response, err := s.auth.RegisterUser(reqCtx, &domain.RegisterUserRequest{
		Login:        request.GetLogin(),
		PasswordHash: request.Password, // todo: resolve password hash problem
		Permission:   request.Permission,
	})

	if err != nil {
		s.logger.Error(
			"executing failed with error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		if errors.Is(err, auth.ErrExists) {
			return nil, ErrAlreadyExists
		}

		return nil, ErrInternal
	}

	return &blogService.RegisterResponse{Id: response.Id}, nil
}

func (s *serverAPI) Login(ctx context.Context, request *blogService.LoginRequest) (*blogService.LoginResponse, error) {
	const op = "bloggrpc.Login"

	//passwordHash, err := HashPassword(request.GetPassword())
	//if err != nil {
	//	s.logger.Error(
	//		"executing failed with error",
	//		slog.String("op", op),
	//		slog.Any("err", err),
	//	)
	//
	//	return nil, ErrInternal
	//}

	reqCtx, cancel := context.WithTimeout(ctx, time.Millisecond*500)
	defer cancel()

	response, err := s.auth.LoginUser(reqCtx, &domain.LoginUserRequest{
		Login:        request.GetLogin(),
		PasswordHash: request.Password, // todo: resolve password hash problem
	})

	if err != nil {
		s.logger.Error(
			"executing failed with error",
			slog.String("op", op),
			slog.Any("err", err),
		)

		if errors.Is(err, auth.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, ErrInternal
	}

	return &blogService.LoginResponse{
		Id:        response.Id,
		Permisson: response.Permission,
	}, nil
}

//func HashPassword(password string) (string, error) {
//	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//	if err != nil {
//		return "", err
//	}
//	return string(hashedPassword), nil
//}
//
//func CheckPasswordHash(password, hash string) bool {
//	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
//	return err == nil
//}
