package bot

import (
	"encoding/json"
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
	bot.HandleMessage("/goals", app.handleGetGoals)
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

	cbData := models.CallbackData{}
	err := json.Unmarshal([]byte(data), &cbData)
	if err != nil {
		a.Log.Err(err).Msg("failed to unmarshal callback data")
		return
	}

	var msg string

	switch cbData.Type {
	case cns.TypeGoal:
		msg = a.handleCallbackGoal(cq, cbData.Goal)
	case cns.TypeCalendar:
		msg = a.calendar.handleCallback(cq, cbData.Calendar)
		if cbData.Calendar != nil && cbData.Calendar.Data != "" {
			msg, _ = a.Uc.HandleMessage(models.Message{
				UserID: cq.From.ID,
				Text:   cbData.Calendar.Data,
			})
		}
	}

	msgData := strings.Split(msg, "|")
	if len(msgData) == 2 {
		a.Client.SendMessage(
			cq.Message.Chat.ID,
			fmt.Sprintf("Цель: %v", msgData[1]),
			tbot.OptInlineKeyboardMarkup(GetGoalActions(cbData.Goal.ID)),
		)
		return
	}

	if msg != "" {
		a.Client.SendMessage(cq.Message.Chat.ID, msg)
	}
}

func callbackDataBuilder(cbData models.CallbackData) string {
	jsonData, err := json.Marshal(cbData)
	if err != nil {
		return "-"
	}
	return string(jsonData)
}
