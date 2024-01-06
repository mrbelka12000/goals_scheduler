package models

import (
	"time"

	"goals_scheduler/internal/cns"
)

type Notifier struct {
	ID          int64              `json:"id,omitempty"`
	UsrID       int64              `json:"usr_id"`
	Text        string             `json:"text,omitempty"`
	Ticker      time.Duration      `json:"ticker,omitempty"`
	Status      cns.StatusNotifier `json:"status,omitempty"`
	LastUpdated time.Time          `json:"last_updated"`
	Expires     time.Time          `json:"expires"`
}

type NotifierCU struct {
	UsrID   *int64              `json:"usr_id"`
	Text    *string             `json:"text,omitempty"`
	Ticker  *time.Duration      `json:"ticker,omitempty"`
	Status  *cns.StatusNotifier `json:"status,omitempty"`
	Expires *time.Time          `json:"expires"`
}

type NotifierPars struct {
	ID     *int64
	UsrID  *int64
	Status *cns.StatusNotifier
}
