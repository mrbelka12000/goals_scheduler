package bot

import (
	"fmt"

	"github.com/yanzay/tbot/v2"

	gs "github.com/mrbelka12000/goals_scheduler"
	"github.com/mrbelka12000/goals_scheduler/internal/models"
)

func generateGoalBottons(list []models.Goal, useCallback bool, action string, status gs.StatusGoal) *tbot.InlineKeyboardMarkup {
	val := &tbot.InlineKeyboardMarkup{}
	callbackData := callbackDataBuilder(models.CallbackData{
		Type: gs.CallbackTypeGoal,
		Goal: &models.GoalData{
			Action: "-",
		},
	})

	for _, l := range list {
		var row []tbot.InlineKeyboardButton
		if useCallback {
			callbackData = callbackDataBuilder(models.CallbackData{
				Type: gs.CallbackTypeGoal,
				Goal: &models.GoalData{
					Action: action,
					ID:     l.ID,
					Status: status,
				},
			})
		}

		row = append(row, tbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%v   |   %v  |   %v", l.Text, l.Deadline.Format(gs.DateFormat), gs.StatusMapper(l.Status)),
			CallbackData: callbackData,
		})

		val.InlineKeyboard = append(val.InlineKeyboard, row)
	}

	if useCallback {
		val.InlineKeyboard = append(val.InlineKeyboard, []tbot.InlineKeyboardButton{
			{
				Text: "Отмена",
				CallbackData: callbackDataBuilder(models.CallbackData{
					Type: gs.CallbackTypeGoal,
					Goal: &models.GoalData{
						Action: "-",
					},
				}),
			},
		})
	}

	return val
}
