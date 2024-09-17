package redis

import (
	"log/slog"

	"github.com/go-redis/redis"
)

type Likes interface {
}

type Redis struct {
	logger *slog.Logger
	db     *redis.Client
}

func New(log *slog.Logger, storageCreds map[string]any) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     storageCreds["addres"].(string),
		Password: storageCreds["password"].(string),
		DB:       storageCreds["db_number"].(int),
	})

	return &Redis{
		logger: log,
		db:     rdb,
	}, nil
}
