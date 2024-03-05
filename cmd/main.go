package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/yanzay/tbot/v2"

	"goals_scheduler/internal/client/webhooker"
	"goals_scheduler/internal/cronjobs"
	"goals_scheduler/internal/delivery/bot"
	v1 "goals_scheduler/internal/delivery/http/v1"
	"goals_scheduler/internal/repo"
	"goals_scheduler/internal/service"
	"goals_scheduler/internal/usecase"
	"goals_scheduler/pkg/cache/redis"
	"goals_scheduler/pkg/config"
	"goals_scheduler/pkg/database"
	"goals_scheduler/pkg/server"
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
	webHookerCl := webhooker.NewClient(cfg)

	rp := repo.New(db)
	srv := service.New(rp)
	uc := usecase.New(log, srv, cache, webHookerCl)

	telBot := tbot.New(cfg.TelegramToken)
	app := bot.NewApp(telBot.Client(), uc, log)

	server := server.NewServer(v1.RegisterHandlers(uc), cfg)
	defer server.Shutdown()

	go cronjobs.Start(app)

	log.Info().Msg("Bot started")
	if err := bot.Start(telBot, app); err != nil {
		log.Error().Err(err).Msg("start bot")
		return
	}
}
