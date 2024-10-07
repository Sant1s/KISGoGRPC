package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Sant1s/blogBack/internal/domain"
	"github.com/Sant1s/blogBack/internal/storage"
)

var (
	ErrInternal = errors.New("iternal server error")
	ErrNotFound = errors.New("object not found")
	ErrExists   = errors.New("object already exists")
)

// Auth data-layer postgres
type Auth interface {
	Register(ctx context.Context, request *domain.RegisterUserRequest) (*domain.RegisterUserResponse, error)
	Login(ctx context.Context, request *domain.LoginUserRequest) (*domain.LoginUserResponse, error)
}

type Service struct {
	log  *slog.Logger
	auth Auth
}

func New(
	log *slog.Logger,
	auth Auth,
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
