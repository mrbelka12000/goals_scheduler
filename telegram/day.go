package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	gs "github.com/mrbelka12000/goals_scheduler"
	"github.com/mrbelka12000/goals_scheduler/messages"
	"github.com/mrbelka12000/goals_scheduler/pkg/ptr"
)

const (
	chooseEmoji = "🎯"
)

type (
	DayInfo struct {
		Mark []bool
	}
)

func (b *Bot) generateDays(info *DayInfo, chatID int64, messageID int) {
	markup := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				{
					Text: "Пн",
					CallbackData: ptr.Ref(callbackDataBuilder(CallbackData{
						Type: gs.CallbackTypeDay,
						Day: &Day{
							Action:  ActionDayMark,
							Weekday: gs.Monday,
						},
					})),
				},
			},
			{
				{
					Text: "Вт",
					CallbackData: ptr.Ref(callbackDataBuilder(CallbackData{
						Type: gs.CallbackTypeDay,
						Day: &Day{
							Action:  ActionDayMark,
							Weekday: gs.Tuesday,
						},
					})),
				},
			},
			{
				{
					Text: "Ср",
					CallbackData: ptr.Ref(callbackDataBuilder(CallbackData{
						Type: gs.CallbackTypeDay,
						Day: &Day{
							Action:  ActionDayMark,
							Weekday: gs.Wednesday,
						},
					})),
				},
			},
			{
				{
					Text: "Чт",
					CallbackData: ptr.Ref(callbackDataBuilder(CallbackData{
						Type: gs.CallbackTypeDay,
						Day: &Day{
							Action:  ActionDayMark,
							Weekday: gs.Thursday,
						},
					})),
				},
			},
			{
				{
					Text: "Пт",
					CallbackData: ptr.Ref(callbackDataBuilder(CallbackData{
						Type: gs.CallbackTypeDay,
						Day: &Day{
							Action:  ActionDayMark,
							Weekday: gs.Friday,
						},
					})),
				},
			},
			{
				{
					Text: "Сб",
					CallbackData: ptr.Ref(callbackDataBuilder(CallbackData{
						Type: gs.CallbackTypeDay,
						Day: &Day{
							Action:  ActionDayMark,
							Weekday: gs.Saturday,
						},
					})),
				},
			},
			{
				{
					Text: "Вс",
					CallbackData: ptr.Ref(callbackDataBuilder(CallbackData{
						Type: gs.CallbackTypeDay,
						Day: &Day{
							Action:  ActionDayMark,
							Weekday: gs.Sunday,
						},
					})),
				},
			},
			{
				{
					Text: "Выбрать",
					CallbackData: ptr.Ref(callbackDataBuilder(CallbackData{
						Type: gs.CallbackTypeDay,
						Day: &Day{
							Action:  ActionDaySubmit,
							Weekday: gs.Wednesday,
						},
					})),
				},
			},
		},
	}

	if info == nil {
		msg := tgbotapi.NewMessage(chatID, gs.MessageChooseDay)
		msg.ReplyMarkup = markup
		b.handleSendMessageError(b.bot.Send(msg))
		return
	}

	for i, v := range info.Mark {
		text := markup.InlineKeyboard[i][0].Text
		enabled := strings.Contains(text, chooseEmoji)

		if !v {
			if enabled {
				text = strings.Replace(text, chooseEmoji, "", 1)
			}
		} else {
			if !enabled {
				text = fmt.Sprintf("%s %s", text, chooseEmoji)
			}
		}
		markup.InlineKeyboard[i][0].Text = text
	}

	editMsg := tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, markup)
	b.handleSendMessageError(b.bot.Send(editMsg))
	return
}

func (b *Bot) handleCallbackDay(ctx context.Context, cq *tgbotapi.CallbackQuery, data *Day) {
	key := fmt.Sprintf("%s:%d", cq.Message.Chat.ID, cq.Message.MessageID)
	b.dayMx.Lock()
	defer b.dayMx.Unlock()

	switch data.Action {
	case ActionDayMark:
		info, ok := b.dayStore[key]
		if !ok {
			info = DayInfo{
				Mark: make([]bool, 7),
			}
		}

		if !info.Mark[data.Weekday] {
			info.Mark[data.Weekday] = true
		} else {
			info.Mark[data.Weekday] = false
		}

		b.generateDays(&info, cq.Message.Chat.ID, cq.Message.MessageID)
		b.dayStore[key] = info

	case ActionDaySubmit:
		info, ok := b.dayStore[key]
		if !ok || !anyTrue(info.Mark) {
			callback := tgbotapi.NewCallbackWithAlert(cq.ID, gs.MessageChooseDay)
			b.handleSendMessageError(b.bot.Send(callback))
			return
		}

		jsonData, err := json.Marshal(info)
		if err != nil {
			b.log.Err(err).Msg("failed to marshal info")
			callback := tgbotapi.NewCallbackWithAlert(cq.ID, gs.SomethingWentWrong)
			b.handleSendMessageError(b.bot.Send(callback))
			return
		}

		delete(b.dayStore, key)
		deleteMsg := tgbotapi.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
		b.handleSendMessageError(b.bot.Send(deleteMsg))

		if err = b.messageSvc.SetValue(gs.KeyDays, jsonData, cq.Message.From.ID); err != nil {
			b.log.Err(err).Msg("failed to set key value")
			callback := tgbotapi.NewCallbackWithAlert(cq.ID, gs.SomethingWentWrong)
			b.handleSendMessageError(b.bot.Send(callback))
			return
		}

		if err := b.messageSvc.SetValue(gs.KeyState, gs.StateAlarmTime, cq.Message.From.ID); err != nil {
			b.log.Err(err).Msg("failed to set key value")
			callback := tgbotapi.NewCallbackWithAlert(cq.ID, gs.SomethingWentWrong)
			b.handleSendMessageError(b.bot.Send(callback))
			return
		}

		response, err := b.messageSvc.HandleMessage(ctx, messages.Message{
			UserID: cq.Message.From.ID,
			ChatID: cq.Message.Chat.ID,
		})
		if err != nil {
			b.log.Err(err).Msg("failed to handle message")
			callback := tgbotapi.NewCallbackWithAlert(cq.ID, gs.SomethingWentWrong)
			b.handleSendMessageError(b.bot.Send(callback))
			return
		}

		b.handleSendMessageError(b.bot.Send(tgbotapi.NewMessage(cq.Message.Chat.ID, response)))
	}
	return
}

func anyTrue(arr []bool) bool {
	for _, v := range arr {
		if v {
			return true
		}
	}
	return false
}
