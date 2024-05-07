package cronjobs

import (
	"time"

	"github.com/go-co-op/gocron"

	"goals_scheduler/internal/delivery/bot"
)

func Start(app *bot.Application) {
	s := gocron.NewScheduler(time.UTC)

	s.Every(30).Second().Do(func() {
		sender(app)
	})

	s.Every(15).Second().Do(func() {
		cleaner(app)
	})

	s.StartBlocking()
}
