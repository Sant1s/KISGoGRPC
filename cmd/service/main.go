package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/Sant1s/blogBack/internal/application"
	"github.com/Sant1s/blogBack/internal/config"
)

var (
	currentLogFile     *os.File
	currentLogFileName string
)

func init() {
	config.MustLoad()
}

func main() {
	logger, err := setupLogger(os.Getenv("ENV"))
	if err != nil {
		panic(err)
	}
	defer currentLogFile.Close()

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	tokenTTL, _ := time.ParseDuration(os.Getenv("TOKEN_TTL"))
	application := application.New(logger, port, os.Getenv("STORAGE_PATH"), tokenTTL)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	s := <-stop
	logger.Info("stopping service", slog.String("signal", s.String()))

	application.GRPCSrv.Stop()

	logger.Info("service stopped")
}

func setupLogger(env string) (*slog.Logger, error) {
	if env == "local" {
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})), nil
	}

	if err := updateLogFile(os.Getenv("LOG_DIR")); err != nil {
		return nil, err
	}

	go func() {
		for range time.Tick(24 * time.Hour) {
			if err := updateLogFile(os.Getenv("LOG_DIR")); err != nil {
				fmt.Printf("Ошибка при обновлении файла логов: %v\n", err)
			}
		}
	}()

	switch env {
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
