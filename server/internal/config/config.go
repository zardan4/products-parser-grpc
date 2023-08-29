package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type DBConfig struct {
	ConnLine string `envconfig:"MONGO_CONNECTIONLINE"`
	DBName   string `envconfig:"MONGO_DATABASE"`
}

type GrpcConfig struct {
	Port int64 `envconfig:"SERVER_INTERNAL_PORT"`
}

type Config struct {
	GrpcConfig
	DBConfig
}

func New() (*Config, error) {
	godotenv.Load()

	cfg := new(Config)

	err := envconfig.Process("mongo", &cfg.DBConfig)
	if err != nil {
		return nil, err
	}

	err = envconfig.Process("server", &cfg.GrpcConfig)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
