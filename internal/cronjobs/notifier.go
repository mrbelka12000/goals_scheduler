package cronjobs

import (
	"context"
	"time"

	"github.com/AlekSi/pointer"

	"goals_scheduler/internal/bot"
	"goals_scheduler/internal/cns"
	"goals_scheduler/internal/models"
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
			app.Client.SendMessage(l.ChatID, "Privet Privet")

			err = app.Uc.NotifierUpdate(context.Background(), models.NotifierCU{
				Notify: pointer.ToDuration(l.Notify),
			}, l.ID)
			if err != nil {
				app.Log.Err(err).Msg("notifier update")
				continue
			}
		}
	}
}
