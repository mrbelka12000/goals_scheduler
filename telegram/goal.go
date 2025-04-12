package telegram

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	gs "github.com/mrbelka12000/goals_scheduler"
	"github.com/mrbelka12000/goals_scheduler/goals"
	"github.com/mrbelka12000/goals_scheduler/messages"
	"github.com/mrbelka12000/goals_scheduler/pkg/ptr"
)

func (b *Bot) handleCallbackGoal(ctx context.Context, cq *tgbotapi.CallbackQuery, data *GoalData) string {
	deleteMsg := tgbotapi.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	if data.Action == "-" {
		b.handleSendMessageError(b.bot.Send(deleteMsg))
		return ""
	}

	switch data.Action {
	case ActionGoalDelete:
		err := b.goalsSvc.Delete(
			ctx,
			data.ID,
		)
		if err != nil {
			b.log.Err(err).Msg("delete goal by id")
			return gs.SomethingWentWrong
		}
		b.handleSendMessageError(b.bot.Send(deleteMsg))

		return "Цель удалена"
	case ActionGoalUpdate:
		err := b.goalsSvc.Update(
			ctx,
			goals.Goal{
				Status: data.Status,
			},
			data.ID,
		)
		if err != nil {
			b.log.Err(err).Msg("update goal by id")
			return gs.SomethingWentWrong
		}

		b.handleSendMessageError(b.bot.Send(deleteMsg))
		return "Цель обновлена"
	case ActionGoalSelect:
		goal, err := b.goalsSvc.Get(ctx, data.ID)
		if err != nil {
			return gs.SomethingWentWrong
		}

		return fmt.Sprintf("%v|%v", "SelectGoal", goal.Message)
	}

	return ""
}

func GetGoalActions(id int64) *tgbotapi.InlineKeyboardMarkup {
	return &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				{
					Text: gs.StatusMapper(gs.StatusGoalFailed),
					CallbackData: ptr.Ref(callbackDataBuilder(CallbackData{
						Type: gs.CallbackTypeGoal,
						Goal: &GoalData{
							ID:     id,
							Action: ActionGoalUpdate,
							Status: gs.StatusGoalFailed,
						},
					})),
				},
				{
					Text: gs.StatusMapper(gs.StatusGoalEnded),
					CallbackData: ptr.Ref(callbackDataBuilder(CallbackData{
						Type: gs.CallbackTypeGoal,
						Goal: &GoalData{
							ID:     id,
							Action: ActionGoalUpdate,
							Status: gs.StatusGoalEnded,
						},
					})),
				},
			},
			{
				{
					Text: "Отмена",
					CallbackData: ptr.Ref(callbackDataBuilder(CallbackData{
						Type: gs.CallbackTypeGoal,
						Goal: &GoalData{
							Action: "-",
						},
					})),
				},
			},
		},
	}
}

func (b *Bot) handleCallbackGoalCreate(ctx context.Context, cq *tgbotapi.CallbackQuery, data *GoalCreateData) {
	deleteMsg := tgbotapi.NewDeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	b.handleSendMessageError(b.bot.Send(deleteMsg))
	switch data.Action {

	case ActionGoalCreateTimer:

		if err := b.messageSvc.SetValue(gs.KeyScheduleType, gs.ScheduleTypeTimer, cq.From.ID); err != nil {
			b.log.Err(err).Msg("choose method")
			b.handleSendMessageError(b.bot.Send(tgbotapi.NewMessage(cq.Message.Chat.ID, gs.SomethingWentWrong)))
			return
		}

		if err := b.messageSvc.SetValue(gs.KeyState, gs.StateTimerInterval, cq.From.ID); err != nil {
			b.log.Err(err).Msg("choose method")
			b.handleSendMessageError(b.bot.Send(tgbotapi.NewMessage(cq.Message.Chat.ID, gs.SomethingWentWrong)))
			return
		}

	case ActionGoalCreateDaily:

		if err := b.messageSvc.SetValue(gs.KeyScheduleType, gs.ScheduleTypeDaily, cq.From.ID); err != nil {
			b.log.Err(err).Msg("choose method")
			b.handleSendMessageError(b.bot.Send(tgbotapi.NewMessage(cq.Message.Chat.ID, gs.SomethingWentWrong)))
			return
		}

		if err := b.messageSvc.SetValue(gs.KeyState, gs.StateAlarmDays, cq.From.ID); err != nil {
			b.log.Err(err).Msg("choose method")
			b.handleSendMessageError(b.bot.Send(tgbotapi.NewMessage(cq.Message.Chat.ID, gs.SomethingWentWrong)))
			return
		}

		b.generateDays(nil, cq.Message.Chat.ID, 0)
		return

	case ActionGoalCreateOnce:

		if err := b.messageSvc.SetValue(gs.KeyScheduleType, gs.ScheduleTypeOnce, cq.From.ID); err != nil {
			b.log.Err(err).Msg("choose method")
			b.handleSendMessageError(b.bot.Send(tgbotapi.NewMessage(cq.Message.Chat.ID, gs.SomethingWentWrong)))
			return
		}

		if err := b.messageSvc.SetValue(gs.KeyState, gs.StateReminderDate, cq.From.ID); err != nil {
			b.log.Err(err).Msg("choose method")
			b.handleSendMessageError(b.bot.Send(tgbotapi.NewMessage(cq.Message.Chat.ID, gs.SomethingWentWrong)))
			return
		}

		b.generateCalendar(cq.Message.Chat.ID)
		return

	case "-":

		if err := b.messageSvc.SetValue(gs.KeyScheduleType, gs.ScheduleTypeNone, cq.From.ID); err != nil {
			b.log.Err(err).Msg("choose method")
			b.handleSendMessageError(b.bot.Send(tgbotapi.NewMessage(cq.Message.Chat.ID, gs.SomethingWentWrong)))
			return
		}

	}

	response, err := b.messageSvc.HandleMessage(ctx, messages.Message{
		UserID: cq.From.ID,
		ChatID: cq.Message.Chat.ID,
	})
	if err != nil {
		b.log.Err(err).Msg("choose method")
		b.handleSendMessageError(b.bot.Send(tgbotapi.NewMessage(cq.Message.Chat.ID, gs.SomethingWentWrong)))
		return
	}

	b.handleSendMessageError(b.bot.Send(tgbotapi.NewMessage(cq.Message.Chat.ID, response)))
}

func (b *Bot) getGoalCreateActions(chatID int64) {
	markup := &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				{
					Text: "Уведомления",
					CallbackData: ptr.Ref(callbackDataBuilder(CallbackData{
						Type: gs.CallbackTypeGoalCreate,
						GoalCreate: &GoalCreateData{
							Action: ActionGoalCreateDaily,
						},
					})),
				},
				{
					Text: "Таймер",
					CallbackData: ptr.Ref(callbackDataBuilder(CallbackData{
						Type: gs.CallbackTypeGoalCreate,
						GoalCreate: &GoalCreateData{
							Action: ActionGoalCreateTimer,
						},
					})),
				},
				{
					Text: "Дата",
					CallbackData: ptr.Ref(callbackDataBuilder(CallbackData{
						Type: gs.CallbackTypeGoalCreate,
						GoalCreate: &GoalCreateData{
							Action: ActionGoalCreateTimer,
						},
					})),
				},
			},
			{
				{
					Text: "Не напоминать",
					CallbackData: ptr.Ref(callbackDataBuilder(CallbackData{
						Type: gs.CallbackTypeGoalCreate,
						GoalCreate: &GoalCreateData{
							Action: "-",
						},
					})),
				},
			},
		},
	}

	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      chatID,
			ReplyMarkup: markup,
		},
		Text: gs.MessageChooseMethod,
	}

	b.handleSendMessageError(b.bot.Send(msg))
}

func generateGoalButtons(list []goals.Goal, useCallback bool, action string, status gs.StatusGoal) *tgbotapi.InlineKeyboardMarkup {
	val := &tgbotapi.InlineKeyboardMarkup{}
	callbackData := ptr.Ref(callbackDataBuilder(CallbackData{
		Type: gs.CallbackTypeGoal,
		Goal: &GoalData{
			Action: "-",
		},
	}))

	for _, l := range list {
		var row []tgbotapi.InlineKeyboardButton
		if useCallback {
			callbackData = ptr.Ref(callbackDataBuilder(CallbackData{
				Type: gs.CallbackTypeGoal,
				Goal: &GoalData{
					Action: action,
					ID:     l.ID,
					Status: status,
				},
			}))
		}

		row = append(row, tgbotapi.InlineKeyboardButton{
			Text:         fmt.Sprintf("%v   |   %v  |   %v", l.Message, ptr.Value(l.Deadline).Format(gs.DateFormat), gs.StatusMapper(l.Status)),
			CallbackData: callbackData,
		})

		val.InlineKeyboard = append(val.InlineKeyboard, row)
	}

	if useCallback {
		val.InlineKeyboard = append(val.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
			{
				Text: "Отмена",
				CallbackData: ptr.Ref(callbackDataBuilder(CallbackData{
					Type: gs.CallbackTypeGoal,
					Goal: &GoalData{
						Action: "-",
					},
				})),
			},
		})
	}

	return val
}
