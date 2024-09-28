package domain

import "time"

type Post struct {
	Id           int64     `db:"id" json:"id"`
	Author       string    `db:"nickname" json:"author"`
	Body         string    `db:"data" json:"body"`
	CreateTime   time.Time `db:"created_at" json:"create_time"`
	CommentCount int64     `db:"comments_count" json:"comments_count"`
	LikesCount   int64     `db:"likes_count" json:"likes_count"`
}

type PostUpdateRequest struct {
	Id   int64  `json:"id"`
	Body string `json:"body"`
}

type Comment struct {
	Id           int64     `db:"id" json:"id"`
	Author       string    `db:"nickname" json:"author"`
	Body         string    `db:"data" json:"body"`
	CreateTime   time.Time `db:"created_at" json:"create_time"`
	CommentCount int64     `db:"comments_count" json:"comments_count"`
	LikesCount   int64     `db:"likes_count" json:"likes_count"`
}

type CommentUpdateRequest struct {
	Id   int64  `json:"id"`
	Body string `json:"body"`
}
