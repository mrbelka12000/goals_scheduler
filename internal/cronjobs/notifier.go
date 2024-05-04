package cronjobs

import (
	"context"
	"fmt"
	"time"

	"github.com/AlekSi/pointer"

	"goals_scheduler/internal/cns"
	"goals_scheduler/internal/delivery/bot"
	"goals_scheduler/internal/models"
)

const (
	notifMessage = "Как обстоят дела с целью:\n%s"
)

func notifier(app *bot.Application) {
	pars := models.NotifierPars{
		Status: pointer.To(cns.StatusNotifierStarted),
	}

	list, _, err := app.Uc.NotifierList(context.Background(), pars)
	if err != nil {
		app.Log.Err(err).Msg("get notifier list")
		return
	}

	for _, l := range list {
		if l.LastUpdated.Before(time.Now()) && l.Notify != 0 {
			goal, err := app.Uc.GoalGet(context.Background(), *l.GoalID)
			if err != nil {
				app.Log.Err(err).Msg("get goal in notification")
				continue
			}
			app.Client.SendMessage(l.ChatID, fmt.Sprintf(notifMessage, goal.Text))

			err = app.Uc.NotifierUpdate(context.Background(), models.NotifierCU{
				Notify: pointer.ToDuration(l.Notify),
			}, l.ID)
			if err != nil {
				app.Log.Err(err).Msg("notifier update")
			}
		}
	}
}
