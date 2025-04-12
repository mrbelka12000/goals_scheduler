package telegram

import (
	"encoding/json"
	"fmt"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"

	gs "github.com/mrbelka12000/goals_scheduler"
	"github.com/mrbelka12000/goals_scheduler/goals"
	"github.com/mrbelka12000/goals_scheduler/messages"
	"github.com/mrbelka12000/goals_scheduler/pkg/config"
)

const (
	ActionGoalDelete = "delete"
	ActionGoalUpdate = "update"
	ActionGoalSelect = "select"

	ActionGoalCreateTimer = "timer"
	ActionGoalCreateDaily = "daily"
	ActionGoalCreateOnce  = "once"

	ActionDayMark   = "mark"
	ActionDaySubmit = "submit"
)

type (
	Bot struct {
		bot *tgbotapi.BotAPI
		ch  tgbotapi.UpdatesChannel

		log        zerolog.Logger
		messageSvc messages.Service
		goalsSvc   goals.Service

		dayStore map[string]DayInfo
		dayMx    sync.Mutex

		calendarStore map[string]calendarInfo
		calendarMx    sync.Mutex
	}

	CallbackData struct {
		Type       gs.CallbackType `json:"type,omitempty"` // goal or calendar
		Calendar   *CalendarData   `json:"calendar,omitempty"`
		Goal       *GoalData       `json:"goal,omitempty"`
		GoalCreate *GoalCreateData `json:"goal_create,omitempty"`
		Day        *Day            `json:"day,omitempty"`
	}

	CalendarData struct {
		Action string `json:"action,omitempty"`
		Data   string `json:"data,omitempty"`
	}

	GoalData struct {
		Action string        `json:"action,omitempty"` // delete, update goal
		ID     int64         `json:"id,omitempty"`     // id of goal
		Status gs.StatusGoal `json:"status,omitempty"` // status to update
	}

	GoalCreateData struct {
		Action string `json:"action,omitempty"`
	}

	Day struct {
		Action  string `json:"action,omitempty"`
		Weekday gs.Day `json:"weekday"`
	}
)

func Connect(
	cfg config.Config,
	messageSvc messages.Service,
	goalsSvc goals.Service,
	log zerolog.Logger,
) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, fmt.Errorf("new bot: %w", err)
	}

	uCfg := tgbotapi.NewUpdate(0)
	uCfg.Timeout = 60

	b := &Bot{
		bot:        bot,
		ch:         bot.GetUpdatesChan(uCfg),
		log:        log,
		messageSvc: messageSvc,
		goalsSvc:   goalsSvc,

		dayStore:      make(map[string]DayInfo),
		calendarStore: make(map[string]calendarInfo),
	}

	go b.handleUpdates()

	return b, nil
}

func (b *Bot) handleSendMessageError(_ tgbotapi.Message, err error) {
	if err != nil {
		b.log.Err(err).Msg("handleSendMessageError")
	}
}

func callbackDataBuilder(cbData CallbackData) string {
	jsonData, err := json.Marshal(cbData)
	if err != nil {
		return "-"
	}
	return string(jsonData)
}
