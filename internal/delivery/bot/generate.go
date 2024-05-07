package bot

import (
	"fmt"

	"github.com/yanzay/tbot/v2"

	"goals_scheduler/internal/cns"
	"goals_scheduler/internal/models"
)

func generateGoalBottons(list []models.Goal, useCallback bool) *tbot.InlineKeyboardMarkup {
	val := &tbot.InlineKeyboardMarkup{}
	callbackData := "-"

	for _, l := range list {
		var row []tbot.InlineKeyboardButton
		if useCallback {
			callbackData = fmt.Sprintf("%v %v", CallbackGoal, l.ID)
		}

		row = append(row, tbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("Цель: %v   |   %v  |   Статус: %v", l.Text, l.Deadline.Format(cns.DateFormat), cns.StatusMapper(l.Status)),
			CallbackData: callbackData,
		})

		val.InlineKeyboard = append(val.InlineKeyboard, row)
	}

	if useCallback {
		val.InlineKeyboard = append(val.InlineKeyboard, []tbot.InlineKeyboardButton{
			{
				Text:         "Отмена",
				CallbackData: fmt.Sprintf("%v -", CallbackGoal),
			},
		})
	}

	return val
}
