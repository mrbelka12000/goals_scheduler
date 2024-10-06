package goals_scheduler

type StatusGoal string

const (
	StatusGoalStarted StatusGoal = "st"
	StatusGoalFailed  StatusGoal = "fa"
	StatusGoalEnded   StatusGoal = "en"
)

func StatusMapper(status StatusGoal) string {
	switch status {
	case StatusGoalStarted:
		return "В процессе"
	case StatusGoalEnded:
		return "Завершена"
	case StatusGoalFailed:
		return "Провалена"
	default:
		return ""
	}
}
