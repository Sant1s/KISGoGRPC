package bloggrpc

import (
	"context"
	"fmt"

	blogService "github.com/Sant1s/blogBack/internal/gen"
	"google.golang.org/grpc"
)

type Blog interface {
	// todo: add service interface
}

type serverAPI struct {
	blogService.UnimplementedBlogServiceServer
	blog Blog
}

func Register(gRPC *grpc.Server, blog Blog) {
	blogService.RegisterBlogServiceServer(gRPC, &serverAPI{blog: blog})
}

func (s *serverAPI) ListPosts(context.Context, *blogService.ListPostsRequest) (*blogService.ListPostsResponse, error) {
	fmt.Println("impl me")
	return nil, nil
}

func (s *serverAPI) CreatePost(context.Context, *blogService.CreatePostRequest) (*blogService.Response, error) {
	return nil, nil

}

func (s *serverAPI) UpdatePost(context.Context, *blogService.UpdatePostRequest) (*blogService.Response, error) {
	return nil, nil

}

func (s *serverAPI) DeletePost(context.Context, *blogService.DeletePostRequest) (*blogService.Response, error) {
	return nil, nil

}
