package models

import "goals_scheduler/internal/cns"

type (
	Message struct {
		UserID int
		ChatID string
		Text   string
	}

	CallbackData struct {
		Type     string        `json:"type,omitempty"` // goal or calendar
		Calendar *CalendarData `json:"calendar,omitempty"`
		Goal     *GoalData     `json:"goal,omitempty"`
	}

	CalendarData struct {
		Action string `json:"action,omitempty"`
		Data   string `json:"data,omitempty"`
	}

	GoalData struct {
		Action string         `json:"action,omitempty"` // delete, update goal
		ID     int64          `json:"id,omitempty"`     // id of goal
		Status cns.StatusGoal `json:"status,omitempty"` // status to update
	}
)
