package models

import (
	"time"

	"goals_scheduler/internal/cns"
)

type (
	Goal struct {
		ID           int64          `json:"id,omitempty"`
		Text         string         `json:"text,omitempty"`
		Deadline     time.Time      `json:"deadline"`
		Status       cns.StatusGoal `json:"status,omitempty"`
		UsrID        int            `json:"usr_id,omitempty"`
		ChatID       string         `json:"chat_id"`
		Timer        *time.Duration `json:"timer,omitempty"`
		TimerEnabled bool           `json:"timer_enabled"`
		LastUpdated  time.Time      `json:"last_updated"`
	}

	GoalCU struct {
		Text         *string         `json:"text,omitempty"`
		UsrID        *int            `json:"usr_id,omitempty"`
		ChatID       *string         `json:"chat_id"`
		Status       *cns.StatusGoal `json:"-"`
		Deadline     *time.Time      `json:"deadline,omitempty"`
		TimerEnabled bool            `json:"timer_enabled"`
		Timer        *time.Duration  `json:"timer"`
	}

	GoalPars struct {
		ID           *int64
		UsrID        *int
		StatusID     *cns.StatusGoal
		TimerEnabled *bool
	}
)
