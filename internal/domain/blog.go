package domain

import "time"

type Post struct {
	Id           string    `json:"id"`
	Author       string    `json:"author"`
	Body         string    `json:"body"`
	CreateTime   time.Time `json:"create_time"`
	CommentCount int64     `json:"comment_count"`
	LikesCount   int64     `json:"likes_count"`
	LikedByUser  bool      `json:"liked"`
	ParentId     string    `json:"parent_id"`
}
type Posts []Post