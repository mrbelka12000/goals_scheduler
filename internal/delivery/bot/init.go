package bot

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/yanzay/tbot/v2"

	"goals_scheduler/internal/cns"
	"goals_scheduler/internal/models"
	"goals_scheduler/internal/usecase"
)

type Application struct {
	Client   *tbot.Client
	Uc       *usecase.UseCase
	Log      zerolog.Logger
	calendar *calendar
}

func NewApp(client *tbot.Client, uc *usecase.UseCase, log zerolog.Logger) *Application {
	return &Application{
		Client:   client,
		calendar: newCalendar(client, log),
		Uc:       uc,
		Log:      log,
	}
}

func Start(bot *tbot.Server, app *Application) error {

	bot.HandleMessage("/start", app.handleStart)
	bot.HandleMessage("/goals", app.handleGetGoal)
	bot.HandleMessage("/goal", app.handleCreateGoal)
	bot.HandleMessage("/c", app.calendar.calendarHandler)
	bot.HandleMessage("/delete_goal", app.deleteGoal)
	bot.HandleMessage("/delete_all_goals", app.deleteUsersGoals)
	bot.HandleMessage(".*", app.handleAllMessages)

	bot.HandleCallback(app.handleCallbacks)

	return bot.Start()
}

func (a *Application) handleStart(m *tbot.Message) {
	a.Client.SendMessage(m.Chat.ID, fmt.Sprintf("Привет %v", m.From.Username))
}

func (a *Application) handleAllMessages(m *tbot.Message) {
	msg, state := a.Uc.HandleMessage(models.Message{
		UserID: m.From.ID,
		ChatID: m.Chat.ID,
		Text:   m.Text,
	})

	if state == "" {
		return
	}

	if state == cns.MessageStateDeadline {
		a.calendar.calendarHandler(m)
		return
	}

	a.Client.SendMessage(m.Chat.ID, msg)
}

func (a *Application) handleCallbacks(cq *tbot.CallbackQuery) {
	data := cq.Data
	if data == "-" {
		return
	}

	if strings.Contains(data, CallbackCalendar) {
		msg := a.calendar.handleCallback(cq)
		if msg != "" {
			msg, _ = a.Uc.HandleMessage(models.Message{
				UserID: cq.From.ID,
				Text:   msg,
			})
			a.Client.SendMessage(cq.Message.Chat.ID, msg)
		}
		return
	} else if strings.Contains(data, CallbackGoal) {
		msg := a.handleCallbackGoal(cq)
		if msg != "" {
			a.Client.SendMessage(cq.Message.Chat.ID, msg)
		}
		return
	}
}
