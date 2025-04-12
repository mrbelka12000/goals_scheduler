package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCommands(ctx context.Context, msg *tgbotapi.Message) (response string, err error) {

	switch msg.Command() {
	case "calendar":
		b.generateCalendar(msg.Chat.ID)

	case "goal":
		b.getGoalCreateActions(msg.Chat.ID)
	}
	return
}
