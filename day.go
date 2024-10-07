package goals_scheduler

import "time"

type Day int

const (
	Monday Day = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

// CastWeekdayToDay because of week in golang start from sunday...
func CastWeekdayToDay(wd time.Weekday) Day {
	wd = wd - 1
	if wd < 0 {
		return Sunday
	}

	return Day(wd)
}
