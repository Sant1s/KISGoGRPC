package application

import (
	"fmt"
	"log/slog"

	grpcapplication "github.com/Sant1s/blogBack/internal/application/grpc"
	"github.com/Sant1s/blogBack/internal/service/blog"
	"github.com/Sant1s/blogBack/internal/storage/postgres"
	"github.com/Sant1s/blogBack/internal/storage/redis"
)

type App struct {
	GRPCSrv *grpcapplication.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	cacheStoragePath map[string]any,
) *App {
	const op = "application.New"

	pgStorage, err := postgres.New(log, storagePath) // TODO: code data-layer (postgres module)
	if err != nil {
		panic(fmt.Sprintf("can not connect to postgres: %s", op))
	}

	redisStorage, err := redis.New(log, cacheStoragePath) // todo: code data-layer (redis modeule)
	if err != nil {
		panic(fmt.Sprintf("can not connect to redis %s", op))
	}

	blogService := blog.New(log, pgStorage, redisStorage)

	return &App{
		GRPCSrv: grpcapplication.New(log, blogService, grpcPort, pgStorage),
	}
}
