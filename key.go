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

func GetKey(key Key, userID int64) string {
	return fmt.Sprintf("%v:%v", key, userID)
}

func GetKeyDate(userID int64) string {
	return GetKey(KeyDate, userID)
}

func GetKeyScheduleTime(userID int64) string {
	return GetKey(KeyScheduleType, userID)
}

func GetKeyInterval(userID int64) string {
	return GetKey(KeyInterval, userID)
}

func GetKeyTimeOfDay(userID int64) string {
	return GetKey(KeyTimeOfDay, userID)
}

func GetKeyDays(userID int64) string {
	return GetKey(KeyDays, userID)
}

func GetKeyMessage(userID int64) string {
	return GetKey(KeyMessage, userID)
}

func GetKeyDeadline(userID int64) string {
	return GetKey(KeyDeadline, userID)
}

func GetKeyHour(userID int64) string {
	return GetKey(KeyHour, userID)
}

func GetKeyMinute(userID int64) string {
	return GetKey(KeyMinute, userID)
}

func GetKeyState(userID int64) string {
	return GetKey(KeyState, userID)
}
