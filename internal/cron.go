package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog"
	"github.com/yanzay/tbot/v2"

	gs "github.com/mrbelka12000/goals_scheduler"
	"github.com/mrbelka12000/goals_scheduler/internal/delivery/bot"
	"github.com/mrbelka12000/goals_scheduler/internal/models"
)

const (
	goalEndedMessage = "Пришло время подводить итогу по цели:\n%s"
	notifMessage     = "Как обстоят дела с целью:\n%s"
	interval         = 15
)

type (
	Cron struct {
		client *tbot.Client
		uc     uc
		log    zerolog.Logger
	}

	uc interface {
		GoalList(ctx context.Context, pars models.GoalPars) ([]models.Goal, int64, error)
		GoalUpdate(ctx context.Context, obj models.GoalCU, id int64) error
	}
)

func NewCron(client *tbot.Client, uc uc, log zerolog.Logger) *Cron {
	return &Cron{
		client: client,
		uc:     uc,
		log:    log,
	}
}

func (c *Cron) Start() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(interval).Second().Do(func() {
		c.sender()
	})

	s.Every(interval).Second().Do(func() {
		c.cleaner()
	})

	s.StartBlocking()
}

func (c *Cron) cleaner() {
	goals, _, err := c.uc.GoalList(context.Background(), models.GoalPars{
		StatusID: pointer.To(gs.StatusGoalStarted),
	})
	if err != nil {
		c.log.Err(err).Msg("failed to get goal list in cleaner")
		return
	}

	for _, goal := range goals {
		if goal.Deadline.Before(time.Now()) {
			c.client.SendMessage(goal.ChatID, fmt.Sprintf(goalEndedMessage, goal.Text), tbot.OptInlineKeyboardMarkup(bot.GetGoalActions(goal.ID)))

			err = c.uc.GoalUpdate(context.Background(), models.GoalCU{
				Status: pointer.To(gs.StatusGoalFailed),
			}, goal.ID)
			if err != nil {
				c.log.Err(err).Msg("failed to update goal status in cleaner")
				return
			}
		}
	}
}

func (c *Cron) sender() {

	goals, _, err := c.uc.GoalList(context.Background(), models.GoalPars{
		StatusID:     pointer.To(gs.StatusGoalStarted),
		TimerEnabled: pointer.To(true),
	})

	if err != nil {
		c.log.Err(err).Msg("failed to get goals in sender")
		return
	}

	for _, goal := range goals {
		if goal.LastUpdated.Before(time.Now()) && goal.Deadline.After(time.Now()) {
			c.client.SendMessage(goal.ChatID, fmt.Sprintf(notifMessage, goal.Text))

			err = c.uc.GoalUpdate(context.Background(), models.GoalCU{
				Timer: goal.Timer,
			}, goal.ID)
			if err != nil {
				c.log.Err(err).Msg("failed to update goal in sender")
				return
			}
		}
	}
}
