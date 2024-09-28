package config

import (
	"errors"
	"flag"
	"os"

	"gopkg.in/yaml.v2"
)

var EmptyConfigPathErr = errors.New("empty config path")

type Config struct {
	Server struct {
		Env    string `yaml:"env"`
		Host   string `yaml:"host"`
		Port   int    `yaml:"port"`
		LogDir string `yaml:"log_dir"`
	} `yaml:"server"`

	Database struct {
		Postgres struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			User     string `yaml:"user"`
			Password string `yaml:"password"`
			Db       string `yaml:"db"`
		} `yaml:"postgres"`

		Redis struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			Password string `yaml:"password"`
			DbNumer  int    `yaml:"db_number"`
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

	if path == "" {
		flag.StringVar(&path, "config_path", "", "Config path for application")
		flag.Parse()
	}
	return path
}
