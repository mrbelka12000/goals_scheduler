package bot

import (
	"fmt"

	"github.com/yanzay/tbot/v2"

	"goals_scheduler/internal/cns"
	"goals_scheduler/internal/models"
)

func generateGoalBottons(list []models.Goal) *tbot.InlineKeyboardMarkup {
	val := &tbot.InlineKeyboardMarkup{}

	for _, l := range list {
		var row []tbot.InlineKeyboardButton

		row = append(row, tbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%v", l.Text),
			CallbackData: "-",
		})
		row = append(row, tbot.InlineKeyboardButton{
			Text:         fmt.Sprintf("%v", l.Deadline.Format(cns.DateFormat)),
			CallbackData: "-",
		})

		val.InlineKeyboard = append(val.InlineKeyboard, row)
	}
	return val
}
