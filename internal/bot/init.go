package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/AlekSi/pointer"
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
		calendar: newCalendar(client),
		Uc:       uc,
		Log:      log,
	}
}

func Start(bot *tbot.Server, app *Application) error {

	bot.HandleMessage("/start", app.handleStart)
	bot.HandleMessage("/goals", app.handleGetGoal)
	bot.HandleMessage("/goal", app.handleCreateGoal)
	bot.HandleMessage("/c", app.calendar.calendarHandler)
	bot.HandleMessage(".*", app.handleAllMessages)
	bot.HandleMessage("/delete_goals", app.deleteUsersGoals)

	bot.HandleCallback(app.handleCallbacks)

	return bot.Start()
}

func (a *Application) handleStart(m *tbot.Message) {
	a.Client.SendMessage(m.Chat.ID, fmt.Sprintf("Привет %v", m.From.Username))
}

func (a *Application) handleCreateGoal(m *tbot.Message) {
	msg := a.Uc.StartGoal(models.Message{
		UserID: m.From.ID,
		Text:   m.Text,
	})

	a.Client.SendMessage(m.Chat.ID, msg)
}

func (a *Application) handleGetGoal(m *tbot.Message) {
	list, _, err := a.Uc.GoalList(context.Background(), models.GoalPars{UsrID: pointer.ToInt(m.From.ID)})
	if err != nil {
		a.Client.SendMessage(m.Chat.ID, "Что то пошло не так")
		return
	}

	a.Client.SendMessage(m.Chat.ID, "Goals", tbot.OptInlineKeyboardMarkup(generateGoalBottons(list)))
}

func (a *Application) handleAllMessages(m *tbot.Message) {
	msg, state := a.Uc.HandleMessage(models.Message{
		UserID: m.From.ID,
		ChatID: m.Chat.ID,
		Text:   m.Text,
	})

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
	}
}

func (a *Application) deleteUsersGoals(m *tbot.Message) {
	err := a.Uc.GoalDeleteAllOfUsers(context.Background(), m.From.ID)
	if err != nil {
		a.Log.Err(err).Msg("delete user`s goals")
		a.Client.SendMessage(m.Chat.ID, "Что то пошло не так")
		return
	}

	a.Client.SendMessage(m.Chat.ID, "Все удалено")
}
