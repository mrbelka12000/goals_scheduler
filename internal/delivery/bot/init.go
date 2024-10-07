package bot

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog"
	"github.com/yanzay/tbot/v2"

	gs "github.com/mrbelka12000/goals_scheduler"
	"github.com/mrbelka12000/goals_scheduler/internal/models"
	"github.com/mrbelka12000/goals_scheduler/internal/usecase"
)

const (
	ActionGoalDelete = "delete"
	ActionGoalUpdate = "update"
	ActionGoalSelect = "select"

	ActionGoalCreateTimer  = "timer"
	ActionGoalCreateNotify = "notify"

	ActionDayMark   = "mark"
	ActionDaySubmit = "submit"
)

type Application struct {
	client   *tbot.Client
	uc       *usecase.UseCase
	log      zerolog.Logger
	calendar *calendar
	day      *day
}

func NewApp(client *tbot.Client, uc *usecase.UseCase, log zerolog.Logger) *Application {
	return &Application{
		client:   client,
		calendar: newCalendar(client, log),
		day:      newDay(client, log),
		uc:       uc,
		log:      log,
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
	a.client.SendMessage(m.Chat.ID, fmt.Sprintf("Привет %v", m.From.Username))
}

func (a *Application) handleAllMessages(m *tbot.Message) {
	msg, state := a.uc.HandleMessage(models.Message{
		UserID: m.From.ID,
		ChatID: m.Chat.ID,
		Text:   m.Text,
	})

	if state == "" {
		return
	}

	switch state {
	case gs.MessageStateDeadline:
		a.calendar.calendarHandler(m)
	case gs.MessageStateDay:
		a.client.SendMessage(m.Chat.ID, msg, tbot.OptInlineKeyboardMarkup(a.day.getBaseKeyboard(nil)))
	default:
		a.client.SendMessage(m.Chat.ID, msg)
	}
}

func (a *Application) handleCallbacks(cq *tbot.CallbackQuery) {
	data := cq.Data
	if data == "-" {
		return
	}

	cbData := models.CallbackData{}
	err := json.Unmarshal([]byte(data), &cbData)
	if err != nil {
		a.log.Err(err).Msg("failed to unmarshal callback data")
		return
	}

	var msg string

	switch cbData.Type {
	case gs.CallbackTypeGoal:
		msg = a.handleCallbackGoal(cq, cbData.Goal)
	case gs.CallbackTypeCalendar:
		msg = a.calendar.handleCallback(cq, cbData.Calendar)
		if cbData.Calendar != nil && cbData.Calendar.Data != "" {
			msg, _ = a.uc.HandleMessage(models.Message{
				UserID: cq.From.ID,
				Text:   cbData.Calendar.Data,
			})

			a.client.SendMessage(cq.Message.Chat.ID, "Выберите метод:", tbot.OptInlineKeyboardMarkup(a.GetGoalCreateActions()))

			return
		}
	case gs.CallbackTypeGoalCreate:
		msg = a.handleCallbackGoalCreate(cq, cbData.GoalCreate)

	case gs.CallbackTypeDay:
		text := a.day.handleCallbackDay(cq, cbData.Day)
		if text != "" {
			msg, _ = a.uc.HandleMessage(models.Message{
				UserID: cq.From.ID,
				Text:   text,
			})
		}
	}

	// only for CallbackTypeGoal
	{
		msgData := strings.Split(msg, "|")
		if len(msgData) == 2 {
			a.client.SendMessage(
				cq.Message.Chat.ID,
				fmt.Sprintf("Цель: %v", msgData[1]),
				tbot.OptInlineKeyboardMarkup(GetGoalActions(cbData.Goal.ID)),
			)
			return
		}
	}

	if msg != "" {
		a.client.SendMessage(cq.Message.Chat.ID, msg)
	}
}

func callbackDataBuilder(cbData models.CallbackData) string {
	jsonData, err := json.Marshal(cbData)
	if err != nil {
		return "-"
	}
	return string(jsonData)
}
