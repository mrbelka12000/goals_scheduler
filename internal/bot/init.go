package bot

import (
	"fmt"
	"strings"

	"github.com/yanzay/tbot/v2"

	"goals_scheduler/internal/cns"
	"goals_scheduler/internal/models"
	"goals_scheduler/internal/usecase"
	"goals_scheduler/pkg/config"
)

type application struct {
	client   *tbot.Client
	calendar *calendar
	uc       *usecase.UseCase
}

func newApp(client *tbot.Client, uc *usecase.UseCase) *application {
	return &application{
		client:   client,
		calendar: newCalendar(client),
		uc:       uc,
	}
}

func Start(cfg config.Config, uc *usecase.UseCase) error {
	bot := tbot.New(cfg.TelegramToken)
	app := newApp(bot.Client(), uc)

	bot.HandleMessage("/start", app.handleStart)
	bot.HandleMessage("/goal", app.handleCreateGoal)
	bot.HandleMessage("/c", app.calendar.calendarHandler)
	bot.HandleMessage(".*", app.handleAllMessages)

	bot.HandleCallback(app.handleCallbacks)

	return bot.Start()
}

func (a *application) handleStart(m *tbot.Message) {
	a.client.SendMessage(m.Chat.ID, fmt.Sprintf("Привет %v", m.From.Username))
}

func (a *application) handleCreateGoal(m *tbot.Message) {
	msg := a.uc.StartGoal(models.Message{
		UserID: m.From.ID,
		Text:   m.Text,
	})
	a.client.SendMessage(m.Chat.ID, msg)
}

func (a *application) handleAllMessages(m *tbot.Message) {
	msg, state := a.uc.HandleMessage(models.Message{
		UserID: m.From.ID,
		Text:   m.Text,
	})
	if state == cns.MessageStateDeadline {
		a.calendar.calendarHandler(m)
		return
	}
	a.client.SendMessage(m.Chat.ID, msg)
}

func (a *application) handleCallbacks(cq *tbot.CallbackQuery) {
	data := cq.Data
	if data == "-" {
		return
	}
	if strings.Contains(data, CallbackCalendar) {
		msg := a.calendar.handleCallback(cq)
		if msg != "" {
			msg, _ = a.uc.HandleMessage(models.Message{
				UserID: cq.From.ID,
				Text:   msg,
			})
			a.client.SendMessage(cq.Message.Chat.ID, msg)
		}
		return
	}
}
