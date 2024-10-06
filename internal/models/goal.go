package models

import (
	"time"

	gs "github.com/mrbelka12000/goals_scheduler"
)

type (
	Goal struct {
		ID           int64          `json:"id,omitempty"`
		Text         string         `json:"text,omitempty"`
		Deadline     time.Time      `json:"deadline"`
		Status       gs.StatusGoal  `json:"status,omitempty"`
		UsrID        int            `json:"usr_id,omitempty"`
		ChatID       string         `json:"chat_id"`
		Timer        *time.Duration `json:"timer,omitempty"`
		TimerEnabled bool           `json:"timer_enabled"`
		LastUpdated  time.Time      `json:"last_updated"`
	}

	GoalCU struct {
		Text         *string        `json:"text,omitempty"`
		UsrID        *int           `json:"usr_id,omitempty"`
		ChatID       *string        `json:"chat_id"`
		Status       *gs.StatusGoal `json:"-"`
		Deadline     *time.Time     `json:"deadline,omitempty"`
		TimerEnabled bool           `json:"timer_enabled"`
		Timer        *time.Duration `json:"timer"`
		LastUpdated  *time.Time     `json:"last_updated,omitempty"`
	}

	GoalPars struct {
		ID           *int64
		UsrID        *int
		StatusID     *gs.StatusGoal
		TimerEnabled *bool
	}
)
