package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/mrbelka12000/goals_scheduler/goals"
	"github.com/mrbelka12000/goals_scheduler/messages"
	"github.com/mrbelka12000/goals_scheduler/pkg/cache/redis"
	"github.com/mrbelka12000/goals_scheduler/pkg/config"
	"github.com/mrbelka12000/goals_scheduler/pkg/database"
	"github.com/mrbelka12000/goals_scheduler/pkg/sender/tg"
	"github.com/mrbelka12000/goals_scheduler/scheme"
	"github.com/mrbelka12000/goals_scheduler/telegram"
)

func main() {
	log := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

	time.Local = time.UTC

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

	sender := tg.NewSender(cfg.TelegramMonitoringBot, cfg.TelegramMonitoringChatID, cfg.ServiceName, log)

	//scheme service
	schemeSvc := scheme.NewService()

	// goals service
	goalsRepo := goals.NewRepo(db)
	goalsSvc := goals.NewService(goalsRepo)
	goalsSvc = goals.NewErrorMW(goalsSvc, sender)

	messageSvc := messages.NewService(
		cache,
		goalsSvc,
		schemeSvc,
	)

	_, err = telegram.Connect(cfg, messageSvc, goalsSvc, log)
	if err != nil {
		log.Fatal().Err(err).Msg("connect to telegram")
	}
	log.Info().Msg("Bot started")

	gs := make(chan os.Signal, 1)
	signal.Notify(gs, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-gs:
		log.Info().Msg(fmt.Sprintf("Received signal: %d", sig))
		log.Info().Msg("Server stopped properly")
		close(gs)
	}
}
