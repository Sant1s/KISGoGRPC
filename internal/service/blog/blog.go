package blog

import (
	"log/slog"
	"time"
)

type BlogService struct {
	log       *slog.Logger
	blogPosts BlogPosts
	blogLikes BlogLikes
	tokenTTL  time.Duration
}

type BlogPosts interface {
	// todo: дописать (data-layer)
}

type BlogLikes interface {
	// todo: дописать (data-layer)
}

func New(
	log *slog.Logger,
	blogPosts BlogPosts,
	blogLikes BlogLikes,
	tokenTTL time.Duration,
) *BlogService {
	return &BlogService{
		log:       log,
		blogPosts: blogPosts,
		blogLikes: blogLikes,
		tokenTTL:  tokenTTL,
	}
}
