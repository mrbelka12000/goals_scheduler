package telegram

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	gs "github.com/mrbelka12000/goals_scheduler"
)

func (b *Bot) handleUpdates() {
	for update := range b.ch {
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Minute)

		if update.CallbackQuery != nil {
			err := b.handleCallbackQuery(ctx, update.CallbackQuery)
			if err != nil {
				b.handleSendMessageError(b.bot.Send(tgbotapi.MessageConfig{
					BaseChat: tgbotapi.BaseChat{
						ChatID: update.CallbackQuery.Message.Chat.ID,
					},
					Text: gs.SomethingWentWrong,
				}))
			}
			continue
		}

		if update.Message.Command() != "" {
			b.handleCommands(ctx, update.Message)
		}

		if update.Message != nil {
		}
	}
}
