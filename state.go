package goals_scheduler

type State string

func (s State) MarshalBinary() ([]byte, error) {
	return []byte(s), nil
}

const (
	StateStart         State = "start"
	StateEnterGoal     State = "enter_goal" // need input
	StateChoseReminder State = "choose_reminder"

	StateTimerInterval State = "waiting_for_timer_interval" // need input
	StateAlarmDays     State = "alarm_days"
	StateAlarmTime     State = "alarm_time" // need input
	StateDeadline      State = "waiting_for_deadline"

	StateReminderDate State = "reminder_date"
	StateReminderTime State = "reminder_time" // need input
	StateDone         State = "done"
)
