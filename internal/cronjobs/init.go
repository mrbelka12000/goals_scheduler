package cronjobs

import (
	"time"

	"github.com/go-co-op/gocron"

	"goals_scheduler/internal/delivery/bot"
)

func Start(app *bot.Application) {
	s := gocron.NewScheduler(time.UTC)

	s.Every(15).Second().Do(func() {
		notifier(app)
	})

	s.StartBlocking()
}
