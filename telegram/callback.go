package telegram

import (
	"context"
	"encoding/json"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	gs "github.com/mrbelka12000/goals_scheduler"
)

func (b *Bot) handleCallbackQuery(ctx context.Context, cq *tgbotapi.CallbackQuery) (err error) {

	var cbData CallbackData
	if err := json.Unmarshal([]byte(cq.Data), &cbData); err != nil {
		return fmt.Errorf("unmarshal callback query: %w", err)
	}

	switch cbData.Type {

	case gs.CallbackTypeDay:

		b.handleCallbackDay(ctx, cq, cbData.Day)

	case gs.CallbackTypeGoal:

	case gs.CallbackTypeGoalCreate:

		fmt.Println("popal")
		b.handleCallbackGoalCreate(ctx, cq, cbData.GoalCreate)

	case gs.CallbackTypeCalendar:

		b.handleCallbackCalendar(cq, cbData.Calendar)

	}

	return nil
}
