package postgres

import (
	"context"

	"github.com/Sant1s/blogBack/internal/domain"
	"github.com/Sant1s/blogBack/internal/storage"
)

// ----------------------------
//		REGISTER USER
// ----------------------------

const queryRegisterUser = ``

func (p *Postgres) Register(ctx context.Context, request *domain.RegisterUserRequest) error {
	_, err := p.getUserId(ctx, request.Login)
	if err == nil {
		return storage.ErrDoesNotExists
	}

	//TODO implement me
	return nil
}

// ----------------------------
//		LOGIN USER
// ----------------------------

const queryLoginUser = ``

func (p *Postgres) Login(ctx context.Context, request *domain.LoginUserRequest) error {
	_, err := p.getUserId(ctx, request.Login)
	if err != nil {
		return storage.ErrDoesNotExists
	}

	//TODO implement me
	return nil
}
