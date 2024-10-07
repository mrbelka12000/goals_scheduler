package models

import gs "github.com/mrbelka12000/goals_scheduler"

type (
	NotifyCU struct {
		DayInfo DayInfo `json:"day_info"`
		Hour    *int    `json:"hour,omitempty"`
		Minute  *int    `json:"minute,omitempty"`
		Weekday *gs.Day `json:"weekday,omitempty"`
		GoalID  *int64  `json:"goal_id,omitempty"`
	}

	Notify struct {
		ID      int64  `json:"id,omitempty"`
		Hour    int    `json:"hour,omitempty"`
		Minute  int    `json:"minute,omitempty"`
		Weekday gs.Day `json:"weekday,omitempty"`
		GoalID  int    `json:"goal_id,omitempty"`
	}

	NotifyPars struct {
		ID      *int64
		GoalID  *int64
		WeekDay *gs.Day
	}
)
