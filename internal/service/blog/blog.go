package blog

import (
	"log/slog"
)

type BlogService struct {
	log       *slog.Logger
	blogPosts BlogPosts
	blogLikes BlogLikes
}

type BlogPosts interface {
	// todo: дописать (data-layer postgres)
}

type BlogLikes interface {
	// todo: дописать (data-layer redis)
}

func New(
	log *slog.Logger,
	blogPosts BlogPosts,
	blogLikes BlogLikes,
) *BlogService {
	return &BlogService{
		log:       log,
		blogPosts: blogPosts,
		blogLikes: blogLikes,
	}
}

// todo: дописать функции, сервисного уровня
