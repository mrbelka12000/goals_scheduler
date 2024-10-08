package goals_scheduler

type CallbackType string

const (
	CallbackTypeCalendar   CallbackType = "calendar"
	CallbackTypeGoal       CallbackType = "goal"
	CallbackTypeGoalCreate CallbackType = "goal_create"
	CallbackTypeDay        CallbackType = "day"
)
