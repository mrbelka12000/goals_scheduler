package cns

import (
	"fmt"
)

const (
	DateFormat         = "2006-01-02"
	SomethingWentWrong = "Что то пошло не так"
	NotifyFormat       = `
		Введите время для напоминания
Допустимые единицы времени: "ns", "us" (или "µs"), "ms", "s", "m", "h".
Отправьте - , в случае если не нужно напоминать`
)

type StatusGoal string

const (
	StatusGoalStarted StatusGoal = "Started"
	StatusGoalEnded   StatusGoal = "Ended"
)

const (
	MessageStateText     = "waiting_for_text"
	MessageStateDeadline = "waiting_for_deadline"
	MessageStateNotifier = "waiting_for_timer"
	MessageStateDone     = "done"
)

const (
	KeyText         = "text"
	KeyDeadline     = "deadline"
	KeyState        = "state"
	KeyTimer        = "timer"
	KeyTimerEnabled = "timer_enabled"
)

var KeysToGoal = []string{KeyText, KeyDeadline, KeyTimer, KeyState}

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
	return getKey(KeyTimer, userID)
}

func GetKeyState(userID int) string {
	return getKey(KeyState, userID)
}

func StatusMapper(status StatusGoal) string {
	switch status {
	case StatusGoalStarted:
		return "В процессе"
	case StatusGoalEnded:
		return "Завершена"
	default:
		return ""
	}
}
