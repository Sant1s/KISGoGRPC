package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidUser     = status.Errorf(codes.Unauthenticated, "invalid token")
)

type Auth interface {
	ValidateUser(ctx context.Context, credentials string) error
}

type Interceptors struct {
	l    *slog.Logger
	auth Auth
}

func New(logger *slog.Logger, auth Auth) *Interceptors {
	return &Interceptors{
		l:    logger,
		auth: auth,
	}
}

func (l *Interceptors) LoggingUnaryInterceptor(
	ctx context.Context,
	request interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	fmt.Printf("loggig interceptor on method: %s\n", info.FullMethod)
	h, err := handler(ctx, request)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (l *Interceptors) AuthUnaryInterceptor(
	ctx context.Context,
	request interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	fmt.Printf("auth interceptor on method: %s\n", info.FullMethod)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}

	authCreds, ok := md["authorization"]
	if !ok {
		return nil, errMissingMetadata
	}

	dbCtx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
	defer cancel()
	if err := l.auth.ValidateUser(dbCtx, authCreds[0]); err != nil {
		return nil, errInvalidUser
	}

	m, err := handler(ctx, request)
	if err != nil {
		l.l.Info("RPC failed with error ", slog.String("err", err.Error()))
	}
	return m, err
}
