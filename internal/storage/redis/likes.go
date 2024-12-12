package redis

import (
	"context"
	"fmt"
	"github.com/Sant1s/blogBack/internal/storage"
	"github.com/go-redis/redis"
	"log"
	"log/slog"
)

type Likes interface {
	LikePost(ctx context.Context, userName string, postId int64) error
	LikeComment(ctx context.Context, userName string, commentId int64) error

	RollbackLikePost(ctx context.Context, userId string, postId int64) error
	RollbackLikeComment(ctx context.Context, userName string, commentId int64) error
}

type Redis struct {
	logger *slog.Logger
	db     *redis.Client
}

func New(log *slog.Logger, storageCredentials map[string]any) (*Redis, error) {
	const op = "redis.New"

	rdb := redis.NewClient(&redis.Options{
		Addr:     storageCredentials["address"].(string),
		Password: storageCredentials["password"].(string),
		DB:       storageCredentials["db_number"].(int),
	})

	err := rdb.Ping().Err()
	if err != nil {
		log.Error(
			"can not connect to redis",
			slog.String("op", op),
			slog.Any("err", err),
		)

		panic(err)
	}

	log.Info(
		"redis client created successfully",
		slog.String("op", op),
	)

	return &Redis{
		logger: log,
		db:     rdb,
	}, nil
}

func (r Redis) LikePost(ctx context.Context, userName string, postId int64) error {
	const op = "redis.LikePost"

	r.logger.Info(
		"executing cache query",
		slog.String("op", op),
	)

	exists, err := r.db.SIsMember(fmt.Sprintf("likes:posts:%s", userName), postId).Result()
	if err != nil {
		log.Fatalf("Ошибка при проверке элемента: %v", err)
	}

	if exists {
		r.logger.Error(
			"error executing query",
			slog.String("op", op),
		)

		return storage.ErrAlreadyExists
	}

	result := r.db.SAdd(
		fmt.Sprintf("likes:posts:%s", userName),
		postId,
	)

	if err := result.Err(); err != nil {
		r.logger.Error(
			"error executing query",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)

		return storage.ErrInternal
	}

	return nil
}

func (r Redis) LikeComment(ctx context.Context, userName string, commentId int64) error {
	const op = "redis.LikeComment"

	r.logger.Info(
		"executing cache query",
		slog.String("op", op),
	)

	exists, err := r.db.SIsMember(fmt.Sprintf("likes:comments:%s", userName), commentId).Result()
	if err != nil {
		log.Fatalf("Ошибка при проверке элемента: %v", err)
	}

	if exists {
		r.logger.Error(
			"error executing query",
			slog.String("op", op),
		)

		return storage.ErrAlreadyExists
	}

	result := r.db.SAdd(
		fmt.Sprintf("likes:comments:%s", userName),
		commentId,
	)

	if err = result.Err(); err != nil {
		r.logger.Error(
			"error executing query",
			slog.String("op", op),
			slog.String("err", err.Error()),
		)

		return storage.ErrInternal
	}

	return nil
}

func (r Redis) RollbackLikePost(ctx context.Context, userName string, postId int64) error {
	const op = "redis.RollbackLikePost"

	r.logger.Info(
		"executing cache query",
		slog.String("op", op),
	)

	res := r.db.SRem(
		fmt.Sprintf("likes:posts:%s", userName),
		postId,
	)

	if err := res.Err(); err != nil {
		r.logger.Error(
			"error executing query",
			slog.String("op", op),
			slog.Any("err", err),
		)

		return storage.ErrInternal
	}

	return nil
}

func (r Redis) RollbackLikeComment(ctx context.Context, userName string, commentId int64) error {
	const op = "redis.RollbackLikeComment"

	r.logger.Info(
		"executing cache query",
		slog.String("op", op),
	)

	res := r.db.SRem(
		fmt.Sprintf("likes:comments:%s", userName),
		commentId,
	)

	if err := res.Err(); err != nil {
		r.logger.Error(
			"error executing query",
			slog.String("op", op),
			slog.Any("err", err),
		)

		return storage.ErrInternal
	}

	return nil
}

func (r Redis) GetLikedPosts(rdb *redis.Client, username string) ([]string, error) {
	return rdb.SMembers(
		fmt.Sprintf("likes:posts:%s", username),
	).Result()
}

func (r Redis) GetLikedComments(rdb *redis.Client, username string) ([]string, error) {
	return rdb.SMembers(
		fmt.Sprintf("likes:comments:%s", username),
	).Result()
}
