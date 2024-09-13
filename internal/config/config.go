package config

import (
	"errors"
	"flag"
	"os"

	"github.com/joho/godotenv"
)

var EmptyConfigPathErr = errors.New("empty config path")

func MustLoad() {
	if err := fetchConfig(); err != nil {
		panic("panic load config")
	}
}

func fetchConfig() error {
	path := getConfigPath()
	if path == "" {
		return EmptyConfigPathErr
	}

	if err := godotenv.Load(path); err != nil {
		return err
	}

	return nil
}

func getConfigPath() string {
	path := os.Getenv("CONFIG_PATH")

	if path == "" {
		flag.StringVar(&path, "config_path", "", "Config path for application")
		flag.Parse()
	}
	return path
}
