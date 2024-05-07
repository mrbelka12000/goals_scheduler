package config

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	//PathToDB       string `env:"PATH_TO_DB,required"`
	PGURL          string `env:"PG_URL,required"`
	TelegramToken  string `env:"TELEGRAM_TOKEN,required"`
	RedisAddr      string `env:"REDIS_ADDR"`
	MigrationsPath string `env:"MIGRATIONS_PATH, default=migrations/"`
	HttpPort       string `env:"HTTP_PORT, default=8080"`
	ServiceName    string `env:"SERVICE_NAME,required"`
	UseMigrates    bool   `env:"USE_MIGRATES,default=false"`
}

func Get() (Config, error) {
	return parseConfig()
}

func parseConfig() (cfg Config, err error) {
	godotenv.Load()

	err = envconfig.Process(context.Background(), &cfg)
	if err != nil {
		return cfg, fmt.Errorf("fill config: %w", err)
	}

	return cfg, nil
}
