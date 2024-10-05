package postgres

import (
	"context"
	"github.com/google/uuid"
	"time"

	"github.com/Sant1s/blogBack/internal/domain"
	"github.com/Sant1s/blogBack/internal/storage"
)

// ----------------------------
//		REGISTER USER
// ----------------------------

const queryRegisterUser = `
INSERT INTO users(id,
                  nickname,
                  password_hash,
                  created_at,
                  permission)
VALUES ($1, $2, $3, $4, $5);
`

func (p *Postgres) Register(ctx context.Context, request *domain.RegisterUserRequest) error {
	_, err := p.getUserId(ctx, request.Login)
	if err == nil {
		return storage.ErrAlreadyExists
	}

	userNewUUid := uuid.New().String()
	_, err = p.db.ExecContext(
		ctx,
		queryRegisterUser,
		userNewUUid,
		request.Login,
		request.PasswordHash,
		time.Now(),
	)
	if err != nil {
		return err
	}

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
