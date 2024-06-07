package main

import (
	"database/sql"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/yanzay/tbot/v2"

	"goals_scheduler/internal/cns"
	"goals_scheduler/internal/cronjobs"
	"goals_scheduler/internal/delivery/bot"
	"goals_scheduler/internal/models"
	"goals_scheduler/internal/repo"
	"goals_scheduler/internal/service"
	"goals_scheduler/internal/usecase"
	"goals_scheduler/pkg/cache/redis"
	"goals_scheduler/pkg/config"
	"goals_scheduler/pkg/database"
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
	if err := doMigrates(db); err != nil {
		log.Error().Err(err).Msg("do migrates")
		return
	}

	rp := repo.New(db)
	srv := service.New(rp)
	uc := usecase.New(log, srv, cache)

	telBot := tbot.New(cfg.TelegramToken)
	app := bot.NewApp(telBot.Client(), uc, log)

	go cronjobs.Start(app)

	log.Info().Msg("Bot started")
	if err := bot.Start(telBot, app); err != nil {
		log.Error().Err(err).Msg("start bot")
		return
	}
}

func doMigrates(db *sql.DB) error {
	rows, err := db.Query("SELECT id, status_id from goals")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		goal := models.Goal{}

		err := rows.Scan(&goal.ID, &goal.Status)
		if err != nil {
			return err
		}

		if goal.Status == "Started" {
			goal.Status = cns.StatusGoalStarted
		}
		if goal.Status == "Ended" {
			goal.Status = cns.StatusGoalEnded
		}

		_, err = db.Exec(`update goals set status_id = $1 where id = $2`, goal.Status, goal.ID)
		if err != nil {
			return err
		}
	}
	return nil
}
