package config

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var EmptyConfigPathErr = errors.New("empty config path")

type Config struct {
	Server struct {
		Env    string `yaml:"env" env:"ENV" env-default:"production"`
		Host   string `yaml:"host" env:"HOST" env-default:"localhost"`
		Port   int    `yaml:"port" env:"PORT" env-default:"8000"`
		LogDir string `yaml:"log_dir" env:"LOG_DIR" env-default:"/opt/blog/logs"`
	} `yaml:"server"`

	Gateway struct {
		Host            string `yaml:"host" env:"GATEWAY_HOST" env-default:"localhost"`
		Port            int    `yaml:"port" env:"GATEWAY_PORT" env-default:"8001"`
		GrpcGatewayHost string `yaml:"grpc_gateway_host" env:"GRPC_GATEWAY_HOST" env-default:"localhost"`
		GrpcGatewayPort int    `yaml:"grpc_gateway_port" env:"GRPC_GATEWAY_PORT" env-default:"7070"`
	} `yaml:"gateway"`

	Database struct {
		Postgres struct {
			Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
			Port     int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
			User     string `yaml:"user" env:"DB_USER" env-default:"postgres"`
			Password string `yaml:"password" env:"DB_PASSWORD" env-default:"postgres"`
			Db       string `yaml:"db" env:"DB_DB" env-default:"postgres"`
		} `yaml:"postgres"`

		Redis struct {
			Host     string `yaml:"host" env:"CACHE_HOST" env-default:"localhost"`
			Port     int    `yaml:"port" env:"CACHE_PORT" env-default:"6379"`
			Password string `yaml:"password" env:"CACHE_PASSWORD" env-default:"password"`
			DbNumer  int    `yaml:"db_number" env:"DB_NUMBER" env-default:"0"`
		} `yaml:"redis"`
	} `yaml:"database"`
}

func MustLoad() *Config {
	cfg, err := fetchConfig()
	if err != nil {
		panic("panic load config")
	}
	return cfg
}

func fetchConfig() (*Config, error) {
	path := getConfigPath()
	if path == "" {
		return nil, EmptyConfigPathErr
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func getConfigPath() string {
	path := os.Getenv("CONFIG_PATH")

	fmt.Println(path)

	if path == "" {
		flag.StringVar(&path, "config_path", "", "Config path for application")
		flag.Parse()
	}
	return path
}
