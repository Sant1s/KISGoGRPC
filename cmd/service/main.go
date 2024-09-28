package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Sant1s/blogBack/internal/application"
	"github.com/Sant1s/blogBack/internal/config"
	metrics "github.com/Sant1s/blogBack/internal/metrics/prometeus"
	"github.com/prometheus/client_golang/prometheus"
)

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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	s := <-stop
	logger.Info("stopping service", slog.String("signal", s.String()))

	application.GRPCSrv.Stop()

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
