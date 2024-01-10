package models

import (
	"time"

	"goals_scheduler/internal/cns"
)

type (
	Goal struct {
		ID       int64          `json:"id,omitempty"`
		Text     string         `json:"text,omitempty"`
		Deadline time.Time      `json:"deadline"`
		Status   cns.StatusGoal `json:"status,omitempty"`
		UsrID    int            `json:"usr_id,omitempty"`
		ChatID   string         `json:"chat_id"`
	}

	GoalCU struct {
		Text          *string         `json:"text,omitempty"`
		UsrID         *int            `json:"usr_id,omitempty"`
		ChatID        *string         `json:"chat_id"`
		Status        *cns.StatusGoal `json:"-"`
		Deadline      *time.Time      `json:"deadline,omitempty"`
		NotifyEnabled bool            `json:"notify_enabled"`
		NotifyTime    *time.Duration  `json:"notify"`
	}

	GoalPars struct {
		ID       *int64
		UsrID    *int
		StatusID *cns.StatusGoal
	}
)
