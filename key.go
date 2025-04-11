package goals_scheduler

import (
	"fmt"
)

type Key string

const (
	KeyMessage      Key = "message"
	KeyDeadline     Key = "deadline"
	KeyScheduleType Key = "schedule_type"
	KeyState        Key = "state"
	KeyInterval     Key = "interval_seconds"
	KeyTimeOfDay    Key = "time_of_day"
	KeyDays         Key = "days"
	KeyHour         Key = "hour"
	KeyMinute       Key = "minute"
	KeyNotify       Key = "notify"
	KeyDate         Key = "date"
)

func (k Key) MarshalBinary() ([]byte, error) {
	return []byte(k), nil
}

var KeysToGoal = []Key{
	KeyMessage,
	KeyDeadline,
	KeyScheduleType,
	KeyState,
	KeyInterval,
	KeyTimeOfDay,
	KeyDays,
	KeyHour,
	KeyMinute,
	KeyNotify,
	KeyDate,
}

func getKey(key Key, userID int) string {
	return fmt.Sprintf("%v:%v", key, userID)
}

func GetKeyDate(userID int) string {
	return getKey(KeyDate, userID)
}

func GetKeyScheduleTime(userID int) string {
	return getKey(KeyScheduleType, userID)
}

func GetKeyInterval(userID int) string {
	return getKey(KeyInterval, userID)
}

func GetKeyTimeOfDay(userID int) string {
	return getKey(KeyTimeOfDay, userID)
}

func GetKeyDays(userID int) string {
	return getKey(KeyDays, userID)
}

func GetKeyMessage(userID int) string {
	return getKey(KeyMessage, userID)
}

func GetKeyDeadline(userID int) string {
	return getKey(KeyDeadline, userID)
}

func GetKeyHour(userID int) string {
	return getKey(KeyHour, userID)
}

func GetKeyMinute(userID int) string {
	return getKey(KeyMinute, userID)
}

func GetKeyState(userID int) string {
	return getKey(KeyState, userID)
}
