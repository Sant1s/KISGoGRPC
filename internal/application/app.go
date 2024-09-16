package application

import (
	"log/slog"
	"time"

	grpcapplication "github.com/Sant1s/blogBack/internal/application/grpc"
	"github.com/Sant1s/blogBack/internal/service/blog"
	"github.com/Sant1s/blogBack/internal/storage/postgres"
)

type App struct {
	GRPCSrv *grpcapplication.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := postgres.New(storagePath) // TODO: code data-layer (postgres module)
	if err != nil {
		panic(err)
	}

	permissionService := blog.New(log, storage, storage, tokenTTL)

	return &App{
		GRPCSrv: grpcapplication.New(log, permissionService, grpcPort),
	}
}
