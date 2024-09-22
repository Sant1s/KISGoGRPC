package middleware

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidUser     = status.Errorf(codes.Unauthenticated, "invalid token")
	errInvalidCreds    = status.Errorf(codes.InvalidArgument, "invalid user creds")
)

type Auth interface {
	ValidateUser(ctx context.Context, nickname, password_hash string) error
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

	dbCtx, cancel := context.WithTimeout(ctx, time.Second*100)
	defer cancel()

	nickname, pass_hash, err := l.parseCreds(authCreds[0])
	if err != nil {
		return nil, errInvalidCreds
	}

	if err := l.auth.ValidateUser(dbCtx, nickname, pass_hash); err != nil {
		return nil, errInvalidUser
	}

	m, err := handler(ctx, request)
	if err != nil {
		l.l.Info("RPC failed with error ", slog.String("err", err.Error()))
	}
	return m, err
}

func (l *Interceptors) parseCreds(creds string) (string, string, error) {
	decodeCreds := strings.Split(creds, " ")

	decodedBytes, err := base64.StdEncoding.DecodeString(decodeCreds[1])
	if err != nil {
		l.l.Error(fmt.Sprintf("Error decoding base64 creds: %v", err))
		return "", "", err
	}

	decodedStr := string(decodedBytes)
	splitCreds := strings.Split(decodedStr, ":")
	return splitCreds[0], splitCreds[1], nil
}
