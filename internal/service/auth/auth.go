package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Sant1s/blogBack/internal/domain"
	"github.com/Sant1s/blogBack/internal/storage"
	"google.golang.org/grpc/codes"
)

var (
	ErrInternal = errors.New("iternal server error")
	ErrNotFound = errors.New("object not found")
	ErrExists   = errors.New("object already exists")
)

// Auth data-layer postgres
type Auth interface {
	Register(ctx context.Context, request *domain.RegisterUserRequest) error
	Login(ctx context.Context, request *domain.LoginUserRequest) error
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

	err := s.auth.Register(ctx, request)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return nil, ErrExists
		}

		return nil, ErrInternal
	}

	return &domain.RegisterUserResponse{
		Code:   int64(codes.OK),
		Output: "ok",
	}, nil
}

func (s *Service) LoginUser(
	ctx context.Context,
	request *domain.LoginUserRequest,
) (*domain.LoginUserResponse, error) {

	err := s.auth.Login(ctx, request)
	if err != nil {
		if errors.Is(err, storage.ErrDoesNotExists) {
			return nil, ErrNotFound
		}

		return nil, ErrInternal
	}

	return &domain.LoginUserResponse{
		Code:   int64(codes.OK),
		Output: "ok",
	}, nil
}
