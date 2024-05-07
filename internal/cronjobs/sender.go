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

func sender(app *bot.Application) {

	goals, _, err := app.Uc.GoalList(context.Background(), models.GoalPars{
		StatusID:     pointer.To(cns.StatusGoalStarted),
		TimerEnabled: pointer.To(true),
	})

	if err != nil {
		app.Log.Err(err).Msg("failed to get goals in sender")
		return
	}

	for _, goal := range goals {
		fmt.Println(goal.LastUpdated)
		//now := time.Now().In(loc)
		if goal.LastUpdated.Before(time.Now()) && goal.Deadline.After(time.Now()) {
			app.Client.SendMessage(goal.ChatID, fmt.Sprintf(notifMessage, goal.Text))

			err = app.Uc.GoalUpdate(context.Background(), models.GoalCU{
				Timer: goal.Timer,
			}, goal.ID)
			if err != nil {
				app.Log.Err(err).Msg("failed to update goal in sender")
				return
			}
		}
	}
}
