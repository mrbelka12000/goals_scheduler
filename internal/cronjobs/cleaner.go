package cronjobs

import (
	"context"
	"fmt"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/yanzay/tbot/v2"

	"goals_scheduler/internal/cns"
	"goals_scheduler/internal/delivery/bot"
	"goals_scheduler/internal/models"
)

const (
	goalEndedMessage = "Пришло время подводить итогу по цели:\n%s"
)

func cleaner(app *bot.Application) {
	goals, _, err := app.Uc.GoalList(context.Background(), models.GoalPars{
		StatusID: pointer.To(cns.StatusGoalStarted),
	})
	if err != nil {
		app.Log.Err(err).Msg("failed to get goal list in cleaner")
		return
	}

	for _, goal := range goals {
		if goal.Deadline.Before(time.Now()) {
			app.Client.SendMessage(goal.ChatID, fmt.Sprintf(goalEndedMessage, goal.Text), tbot.OptInlineKeyboardMarkup(bot.GetGoalActions(goal.ID)))

			err = app.Uc.GoalUpdate(context.Background(), models.GoalCU{
				Status: pointer.To(cns.StatusGoalFailed),
			}, goal.ID)
			if err != nil {
				app.Log.Err(err).Msg("failed to update goal status in cleaner")
				return
			}
		}
	}
}
