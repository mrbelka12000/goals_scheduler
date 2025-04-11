package goals

import (
	"time"

	gs "github.com/mrbelka12000/goals_scheduler"
)

type (
	Goal struct {
		ID             int64           `json:"id,omitempty"`
		TelegramChatID string          `json:"telegram_chat_id,omitempty"`
		ScheduleType   gs.ScheduleType `json:"schedule_type,omitempty"`

		// daily notifications
		Day    gs.Day `json:"day,omitempty"`
		Hour   int    `json:"hour,omitempty"`
		Minute int    `json:"minute,omitempty"`

		// interval notifications
		IntervalSeconds *float64   `json:"interval_seconds,omitempty"`
		NextExecution   *time.Time `json:"next_execution,omitempty"`

		// scheduled notification
		ScheduledTime *time.Time `json:"scheduled_time,omitempty"`

		Message   string        `json:"message,omitempty"`
		Status    gs.StatusGoal `json:"status,omitempty"`
		CreatedAt time.Time     `json:"created_at"`
		UpdatedAt time.Time     `json:"updated_at"`
	}

	GoalPars struct {
		ID           *int64
		ChatID       *int64
		ScheduleType *string
		Status       *gs.StatusGoal
	}
)
