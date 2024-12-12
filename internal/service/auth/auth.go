package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Sant1s/blogBack/internal/domain"
	"github.com/Sant1s/blogBack/internal/storage"
	"github.com/Sant1s/blogBack/internal/storage/postgres"
)

var (
	ErrInternal   = errors.New("iternal server error")
	ErrNotFound   = errors.New("object not found")
	ErrExists     = errors.New("object already exists")
	ErrBadRequest = errors.New("bad request")
)

type Service struct {
	log  *slog.Logger
	auth postgres.Auth
}

func New(
	log *slog.Logger,
	auth postgres.Auth,
) *Service {
	return &Service{
		log:  log,
		auth: auth,
	}
}

func (s *Service) RegisterUser(
	ctx context.Context,
	request *domain.RegisterUserRequest,
) (*domain.RegisterUserResponse, error) {
	if request.Login == "" || request.PasswordHash == "" || request.Permission == "" {
		return nil, ErrBadRequest
	}

	resp, err := s.auth.Register(ctx, request)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return nil, ErrExists
		}

		return nil, ErrInternal
	}

	return &domain.RegisterUserResponse{
		Id: resp.Id,
	}, nil
}

func (s *Service) LoginUser(
	ctx context.Context,
	request *domain.LoginUserRequest,
) (*domain.LoginUserResponse, error) {
	if request.Login == "" || request.PasswordHash == "" {
		return nil, ErrBadRequest
	}

	resp, err := s.auth.Login(ctx, request)
	if err != nil {
		if errors.Is(err, storage.ErrDoesNotExists) {
			return nil, ErrNotFound
		}

		return nil, ErrInternal
	}

	return &domain.LoginUserResponse{
		Id:         resp.Id,
		Permission: resp.Permission,
	}, nil
}
