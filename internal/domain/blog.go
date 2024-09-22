package domain

import "time"

type Post struct {
	Id           int64     `db:"id"`
	Author       string    `db:"nickname"`
	Body         string    `db:"data"`
	CreateTime   time.Time `db:"created_at"`
	CommentCount int64     `db:"comments_count"`
	LikesCount   int64     `db:"likes_count"`
	LikedByUser  bool      `db:"liked"`
}

type Posts []Post
