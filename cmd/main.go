package main

import (
	"os"

	"github.com/rs/zerolog"

	"goals_scheduler/internal/bot"
	"goals_scheduler/internal/repo"
	"goals_scheduler/internal/service"
	"goals_scheduler/internal/usecase"
	"goals_scheduler/pkg/cache/redis"
	"goals_scheduler/pkg/config"
	"goals_scheduler/pkg/database"
)

func main() {
	log := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

	cfg, err := config.Get()
	if err != nil {
		log.Fatal().Err(err).Msg("get config")
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("database connect")
	}
	defer db.Close()

	cache, err := redis.New(cfg)
	if err != nil {
		log.Error().Err(err).Msg("connect to redis")
		return
	}

	rp := repo.New(db)
	srv := service.New(rp)
	uc := usecase.New(log, srv, cache)

	log.Info().Msg("Bot started")
	if err := bot.Start(cfg, uc); err != nil {
		log.Error().Err(err).Msg("start bot")
		return
	}
}
