package bloggrpc

import (
	blogService "github.com/Sant1s/blogBack/internal/gen"
	"google.golang.org/grpc"
)

type Blog interface {
	// todo: add service interface
}

type serverAPI struct {
	blogService.UnimplementedBlogPostServiceServer
	blog Blog
}

func Register(gRPC *grpc.Server, blog Blog) {
	blogService.RegisterBlogPostServiceServer(gRPC, &serverAPI{blog: blog})
}
