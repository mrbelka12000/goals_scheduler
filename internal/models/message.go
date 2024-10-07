package models

import (
	gs "github.com/mrbelka12000/goals_scheduler"
)

type (
	Message struct {
		UserID int
		ChatID string
		Text   string
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

	DayInfo struct {
		Mark []bool
	}
)
