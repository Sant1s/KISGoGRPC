package grpcapplication

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	bloggrpc "github.com/Sant1s/blogBack/internal/grpc/blog"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

type Auth interface {
	ValidateUser(ctx context.Context, nickname, passwordHash string) error
}

func New(
	log *slog.Logger,
	blogService bloggrpc.Blog,
	likesService bloggrpc.Likes,
	authService bloggrpc.Auth,
	auth Auth,
	port int,
) *App {
	middlewares := middleware.New(log, auth)

	gRPCServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middlewares.LoggingUnaryInterceptor,
			middlewares.AuthUnaryInterceptor,
		),
	)

	bloggrpc.Register(log, gRPCServer, blogService, likesService, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) run() error {
	const op = "grpcapplication.run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("port", fmt.Sprintf("%v", a.port)),
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", slog.String("addr", lis.Addr().String()))

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) MustRun() {
	if err := a.run(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() {
	const op = "grpcapplication.Stop"

	a.log.With(
		slog.String("op", op),
		slog.String("port", fmt.Sprintf("%v", a.port)),
	).Info("stopping gRPC server")

	a.gRPCServer.GracefulStop()
}
