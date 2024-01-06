package models

import (
	"time"

	"goals_scheduler/internal/cns"
)

type Goal struct {
	ID       int64          `json:"id,omitempty"`
	Text     string         `json:"text,omitempty"`
	Deadline time.Time      `json:"deadline"`
	Status   cns.StatusGoal `json:"status,omitempty"`
	UsrID    int64          `json:"usr_id,omitempty"`

	NotifierID int64 `json:"notifier_id,omitempty"`
	Notifier   `json:"notifier"`
}

type GoalCU struct {
	Text     *string         `json:"text,omitempty"`
	Deadline *time.Time      `json:"deadline,omitempty"`
	UsrID    *int64          `json:"usr_id,omitempty"`
	Status   *cns.StatusGoal `json:"-"`

	NotifierID     *int64         `json:"-"`
	NotifierText   *string        `json:"notifier_text,omitempty"`
	NotifierTicker *time.Duration `json:"notifier_ticker,omitempty"`
}

type GoalPars struct {
	ID       *int64
	UsrID    *int64
	StatusID *cns.StatusGoal
}
