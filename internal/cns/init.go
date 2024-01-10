package cns

import "fmt"

const (
	DateFormat = "2006-01-02"
)

type StatusGoal string

const (
	StatusGoalStarted StatusGoal = "Started"
	StatusGoalEnded   StatusGoal = "Ended"
)

type StatusNotifier string

const (
	StatusNotifierStarted StatusNotifier = "Started"
	StatusNotifierEnded   StatusNotifier = "Ended"
)

const (
	MessageStateText     = "waiting_for_text"
	MessageStateDeadline = "waiting_for_deadline"
	MessageStateNotifier = "waiting_for_notify_time"
)

const (
	KeyText          = "text"
	KeyDeadline      = "deadline"
	KeyState         = "state"
	KeyNotify        = "notify"
	KeyNotifyEnabled = "notify_enabled"
)

var KeysToGoal = []string{KeyText, KeyDeadline, KeyNotify, KeyState}

func getKey(key string, userID int) string {
	return fmt.Sprintf("%v:%v", key, userID)
}

func GetKeyText(userID int) string {
	return getKey(KeyText, userID)
}

func GetKeyDeadline(userID int) string {
	return getKey(KeyDeadline, userID)
}

func GetKeyNotify(userID int) string {
	return getKey(KeyNotify, userID)
}

func GetKeyState(userID int) string {
	return getKey(KeyState, userID)
}
