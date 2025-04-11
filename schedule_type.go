package goals_scheduler

type ScheduleType string

const (
	ScheduleTypeTimer ScheduleType = "timer"
	ScheduleTypeDaily ScheduleType = "daily"
	ScheduleTypeOnce  ScheduleType = "once"
	ScheduleTypeNone  ScheduleType = "none"
)

func (s ScheduleType) MarshalBinary() ([]byte, error) {
	return []byte(s), nil
}
