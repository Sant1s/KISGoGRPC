package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"strings"
	"time"

	bloggrpc "github.com/Sant1s/blogBack/internal/grpc/blog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrMissingMetadata    = status.Errorf(codes.InvalidArgument, "missing metadata")
	ErrInvalidUser        = status.Errorf(codes.Unauthenticated, "invalid token")
	ErrInvalidCredentials = status.Errorf(codes.InvalidArgument, "invalid user Credentials")
	ErrInternal           = status.Error(codes.Internal, "internal")
)

type Auth interface {
	ValidateUser(ctx context.Context, nickname, passwordHash string) error
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
	const op = "middleware.LoggingUnaryInterceptor"
	l.l.Info(
		"start executing handler",
		slog.String("op", op),
		slog.String("handler", info.FullMethod),
	)
	start := time.Now()
	h, err := handler(ctx, request)
	duration := time.Since(start)

	if err != nil {
		l.l.Error(
			"executing handler ended with error",
			slog.String("op", op),
			slog.String("handler", info.FullMethod),
			slog.String("err", err.Error()),
			slog.String("time", duration.String()),
		)
		return nil, err
	}

	l.l.Info(
		"executing handler ended successfully",
		slog.String("op", op),
		slog.String("handler", info.FullMethod),
		slog.String("time", duration.String()),
	)

	return h, nil
}

func (l *Interceptors) AuthUnaryInterceptor(
	ctx context.Context,
	request interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	const op = "middleware.AuthUnaryInterceptor"

	if info.FullMethod != "/kis.blog.backend.BlogService/Login" &&
		info.FullMethod != "/kis.blog.backend.BlogService/Register" {

		l.l.Info(fmt.Sprintf("auth interceptor on method: %s\n", info.FullMethod), slog.String("op", op))

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, ErrMissingMetadata
		}

		authCredentials, ok := md["authorization"]
		if !ok {
			l.l.Error(
				"authorization header required",
				slog.String("op", op),
			)
			return nil, ErrMissingMetadata
		}

		dbCtx, cancel := context.WithTimeout(ctx, time.Millisecond*100)
		defer cancel()

		nickname, pass, err := l.parseCredentials(authCredentials[0])
		if err != nil {
			l.l.Error(
				"invalid credentials",
				slog.String("op", op),
			)
			return nil, ErrInvalidCredentials
		}

		passHash, err := bloggrpc.HashPassword(pass)
		if err != nil {
			l.l.Error(
				"password hashing error",
				slog.String("op", op),
			)
			return nil, ErrInternal
		}

		if err := l.auth.ValidateUser(dbCtx, nickname, passHash); err != nil {
			l.l.Error(
				fmt.Sprintf("can not validate user: %s", err.Error()),
				slog.String("op", op),
			)
			return nil, ErrInvalidUser
		}
	}

	m, err := handler(ctx, request)
	if err != nil {
		l.l.Info(
			fmt.Sprintf("rpc failed with error: %s", err.Error()),
			slog.String("op", op),
		)
	}
	return m, err
}

func (l *Interceptors) parseCredentials(Credentials string) (string, string, error) {
	const op = "middleware.parseCredentials"
	decodeCredentials := strings.Split(Credentials, " ")

	decodedBytes, err := base64.StdEncoding.DecodeString(decodeCredentials[1])
	if err != nil {
		l.l.Error(
			fmt.Sprintf("Error decoding base64 Credentials: %v", err),
			slog.String("op", op),
		)
		return "", "", err
	}

	decodedStr := string(decodedBytes)
	splitCredentials := strings.Split(decodedStr, ":")

	return splitCredentials[0], splitCredentials[1], nil
}
