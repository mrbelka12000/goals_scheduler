package goals_scheduler

import "fmt"

type Key string

const (
	KeyText         Key = "text"
	KeyDeadline     Key = "deadline"
	KeyState        Key = "state"
	KeyTimer        Key = "timer"
	KeyHour         Key = "hour"
	KeyMinute       Key = "minute"
	KeyNotify       Key = "notify"
	KeyDay          Key = "day"
	KeyTimerEnabled Key = "timer_enabled"
	KeyChoose       Key = "choose"
)

func (k Key) MarshalBinary() ([]byte, error) {
	return []byte(k), nil
}

var KeysToGoal = []Key{
	KeyText,
	KeyDeadline,
	KeyState,
	KeyTimer,
	KeyHour,
	KeyMinute,
	KeyNotify,
	KeyDay,
	KeyTimerEnabled,
	KeyChoose,
}

func getKey(key Key, userID int) string {
	return fmt.Sprintf("%v:%v", key, userID)
}

func GetKeyText(userID int) string {
	return getKey(KeyText, userID)
}

func GetKeyDeadline(userID int) string {
	return getKey(KeyDeadline, userID)
}

func GetKeyTimer(userID int) string {
	return getKey(KeyTimer, userID)
}

func GetKeyHour(userID int) string {
	return getKey(KeyHour, userID)
}

func GetKeyMinute(userID int) string {
	return getKey(KeyMinute, userID)
}

func GetKeyDay(userID int) string {
	return getKey(KeyDay, userID)
}

func GetKeyState(userID int) string {
	return getKey(KeyState, userID)
}

func GetKeyChoose(userID int) string {
	return getKey(KeyChoose, userID)
}
