package main

import (
	"context"
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Sant1s/blogBack/internal/application"
	"github.com/Sant1s/blogBack/internal/config"
	blogService "github.com/Sant1s/blogBack/internal/gen"
	metrics "github.com/Sant1s/blogBack/internal/metrics/prometeus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

//go:embed proto/api/blogService.swagger.json
var swaggerData []byte

//go:embed swagger-ui
var embededFs embed.FS

var (
	currentLogFile     *os.File
	currentLogFileName string
)

func init() {
	prometheus.MustRegister(metrics.RequestCounter)
}

func main() {

	cfg := config.MustLoad()

	logger, err := setupLogger(cfg)
	if err != nil {
		panic(err)
	}
	defer currentLogFile.Close()

	storagePath := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.Database.Postgres.User,
		cfg.Database.Postgres.Password,
		cfg.Database.Postgres.Host,
		cfg.Database.Postgres.Port,
		cfg.Database.Postgres.Db,
	)
	cacheStoragePath := map[string]any{
		"address":   fmt.Sprintf("%s:%d", cfg.Database.Redis.Host, cfg.Database.Redis.Port),
		"password":  cfg.Database.Redis.Password,
		"db_number": cfg.Database.Redis.DbNumer,
	}

	application := application.New(logger, cfg.Server.Port, storagePath, cacheStoragePath)

	go application.GRPCSrv.MustRun()

	doneCh, err := runRest(cfg, logger)
	if err != nil {
		panic(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	s := <-stop
	logger.Info("stopping service", slog.String("signal", s.String()))

	application.GRPCSrv.Stop()

	logger.Info("stopping gateway service", slog.String("signal", s.String()))

	doneCh <- struct{}{}

	logger.Info("service stopped")
}

func setupLogger(cfg *config.Config) (*slog.Logger, error) {
	if cfg.Server.Env == "local" {
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})), nil
	}

	if err := updateLogFile(cfg.Server.LogDir); err != nil {
		return nil, err
	}

	go func() {
		for range time.Tick(24 * time.Hour) {
			if err := updateLogFile(cfg.Server.LogDir); err != nil {
				fmt.Printf("Ошибка при обновлении файла логов: %v\n", err)
			}
		}
	}()

	switch cfg.Server.Env {
	case "production":
		return slog.New(slog.NewJSONHandler(currentLogFile, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})), nil

	case "development":
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})), nil
	}

	return nil, errors.New("can not parse env level")
}

func updateLogFile(logDir string) error {
	currentDate := time.Now().Format("2006-01-02")
	if currentLogFile == nil || currentLogFileName != currentDate {
		if currentLogFile != nil {
			err := currentLogFile.Close()
			if err != nil {
				return err
			}
		}

		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			return fmt.Errorf("can not make logs directory: %w", err)
		}

		logFilePath := filepath.Join(logDir, currentDate+".log")
		logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("can not open log-file: %w", err)
		}

		currentLogFile = logFile
		currentLogFileName = currentDate
	}
	return nil
}

func runRest(cfg *config.Config, logger *slog.Logger) (chan struct{}, error) {
	const op = "main.runRest"

	doneChan := make(chan struct{})

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	clientConn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", cfg.Server.Port))
	err = blogService.RegisterBlogServiceHandler(ctx, mux, clientConn)

	fsSwagger, err := fs.Sub(embededFs, "swagger-ui")
	if err != nil {
		panic(err)
	}

	swaggerMux := http.NewServeMux()

	swaggerMux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.FS(fsSwagger))))

	swaggerMux.HandleFunc("/swagger-ui/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(swaggerData)
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Gateway.Port),
		Handler: swaggerMux,
	}

	go func() {
		logger.Info(fmt.Sprintf("server listening at %d", cfg.Gateway.Port))
		if err := server.ListenAndServe(); err != nil {
			logger.Error(
				"error register gateway server",
				slog.String("op", op),
				slog.Any("err", err),
			)

		}
	}()

	go func() {
		select {
		case <-doneChan:
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
			defer cancel()

			logger.Info(
				"stoping gateway server",
				slog.String("op", op),
			)

			if err := server.Shutdown(ctx); err != nil {
				logger.Info(
					"error stoping gateway server",
					slog.String("op", op),
					slog.Any("err", err),
				)
			}
		}
	}()

	return doneChan, nil
}
