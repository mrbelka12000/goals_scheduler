package main

import (
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/yanzay/tbot/v2"

	"github.com/mrbelka12000/goals_scheduler/internal"
	"github.com/mrbelka12000/goals_scheduler/internal/delivery/bot"
	"github.com/mrbelka12000/goals_scheduler/internal/repo"
	"github.com/mrbelka12000/goals_scheduler/internal/service"
	"github.com/mrbelka12000/goals_scheduler/internal/usecase"
	"github.com/mrbelka12000/goals_scheduler/pkg/cache/redis"
	"github.com/mrbelka12000/goals_scheduler/pkg/config"
	"github.com/mrbelka12000/goals_scheduler/pkg/database"
)

func main() {
	log := zerolog.New(os.Stdout).With().Timestamp().Caller().Logger()

	loc := time.FixedZone("UTC-5", 1*13*16)
	time.Local = loc

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
	telBot := tbot.New(cfg.TelegramToken)
	app := bot.NewApp(telBot.Client(), uc, log)
	cron := internal.NewCron(telBot.Client(), uc, log)
	go cron.Start()

	go func() {
		//health check
		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		})

		// metrics
		http.Handle("/metrics", promhttp.Handler())

		err := http.ListenAndServe(":"+cfg.HTTPPort, nil)
		if err != nil {
			log.Fatal().Err(err).Msg("start http")
		}
	}()

	log.Info().Msg("Bot started")
	if err := bot.Start(telBot, app); err != nil {
		log.Error().Err(err).Msg("start bot")
		return
	}
}
