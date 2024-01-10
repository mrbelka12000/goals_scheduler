package models

import (
	"time"

	"goals_scheduler/internal/cns"
)

type (
	Notifier struct {
		ID          int64              `json:"id,omitempty"`
		UsrID       int                `json:"usr_id,omitempty"`
		ChatID      string             `json:"chat_id,omitempty"`
		Notify      time.Duration      `json:"notify,omitempty"`
		LastUpdated time.Time          `json:"last_updated"`
		Status      cns.StatusNotifier `json:"status"`
		GoalID      *int64             `json:"goal_id"`
	}

	NotifierCU struct {
		UsrID  *int
		ChatID *string
		GoalID *int64
		Notify *time.Duration
		Status *cns.StatusNotifier
	}

	NotifierPars struct {
		ID     *int64
		UsrID  *int64
		Status *cns.StatusNotifier
	}
)
