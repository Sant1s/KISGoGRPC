package service

type BlogService interface {
}

type Blog struct {
}

func NewBlogService() *Blog {
	return &Blog{}
}
