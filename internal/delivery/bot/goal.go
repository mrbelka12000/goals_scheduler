package bot

import (
	"context"
	"fmt"

	"github.com/AlekSi/pointer"
	"github.com/yanzay/tbot/v2"

	"goals_scheduler/internal/cns"
	"goals_scheduler/internal/models"
)

const (
	ActionDelete = "delete"
	ActionUpdate = "update"
	ActionSelect = "select"
)

func (a *Application) deleteGoal(m *tbot.Message) {
	list, _, err := a.Uc.GoalList(context.Background(), models.GoalPars{UsrID: pointer.ToInt(m.From.ID)})
	if err != nil {
		a.Client.SendMessage(m.Chat.ID, cns.SomethingWentWrong)
		a.Log.Err(err).Msg("get goals list")
		return
	}

	if len(list) == 0 {
		a.Client.SendMessage(m.Chat.ID, "Пока что у вас нет целей")
		return
	}

	a.Client.SendMessage(
		m.Chat.ID,
		"Выберите цель для удаления",
		tbot.OptInlineKeyboardMarkup(generateGoalBottons(list, true, ActionDelete, "-")))
}

func (a *Application) deleteUsersGoals(m *tbot.Message) {
	err := a.Uc.GoalDeleteAllOfUsers(context.Background(), m.From.ID)
	if err != nil {
		a.Log.Err(err).Msg("delete user`s goals")
		a.Client.SendMessage(m.Chat.ID, cns.SomethingWentWrong)
		return
	}

	a.Client.SendMessage(m.Chat.ID, "Все удалено")
}

func (a *Application) handleCreateGoal(m *tbot.Message) {
	msg := a.Uc.StartGoal(models.Message{
		UserID: m.From.ID,
		Text:   m.Text,
	})

	a.Client.SendMessage(m.Chat.ID, msg)
}

func (a *Application) handleGetGoals(m *tbot.Message) {
	list, _, err := a.Uc.GoalList(context.Background(), models.GoalPars{UsrID: pointer.ToInt(m.From.ID)})
	if err != nil {
		a.Client.SendMessage(m.Chat.ID, cns.SomethingWentWrong)
		a.Log.Err(err).Msg("get goals list")
		return
	}

	if len(list) == 0 {
		a.Client.SendMessage(m.Chat.ID, "Пока что у вас нет целей")
		return
	}

	a.Client.SendMessage(
		m.Chat.ID,
		"Цели",
		tbot.OptInlineKeyboardMarkup(generateGoalBottons(list, true, ActionSelect, "-")))
}

func (a *Application) handleCallbackGoal(cq *tbot.CallbackQuery, data *models.GoalData) string {
	if data.Action == "-" {
		a.Client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
		return ""
	}

	switch data.Action {
	case ActionDelete:
		err := a.Uc.GoalDelete(
			context.Background(),
			data.ID,
		)
		if err != nil {
			a.Log.Err(err).Msg("delete goal by id")
			return cns.SomethingWentWrong
		}

		a.Client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)

		return "Цель удалена"
	case ActionUpdate:
		err := a.Uc.GoalUpdate(
			context.Background(),
			models.GoalCU{
				Status: &data.Status,
			},
			data.ID,
		)
		if err != nil {
			a.Log.Err(err).Msg("update goal by id")
			return cns.SomethingWentWrong
		}

		a.Client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
		return "Цель обновлена"
	case ActionSelect:
		goal, err := a.Uc.GoalGet(context.Background(), data.ID)
		if err != nil {
			return cns.SomethingWentWrong
		}

		return fmt.Sprintf("%v|%v", "SelectGoal", goal.Text)
	}

	return ""
}

func (a *Application) GetGoalEndChoose(id int64) *tbot.InlineKeyboardMarkup {
	result := &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			{
				{
					Text: cns.StatusMapper(cns.StatusGoalFailed),
					CallbackData: callbackDataBuilder(models.CallbackData{
						Type: cns.TypeGoal,
						Goal: &models.GoalData{
							ID:     id,
							Action: ActionUpdate,
							Status: cns.StatusGoalFailed,
						},
					}),
				},
				{
					Text: cns.StatusMapper(cns.StatusGoalEnded),
					CallbackData: callbackDataBuilder(models.CallbackData{
						Type: cns.TypeGoal,
						Goal: &models.GoalData{
							ID:     id,
							Action: ActionUpdate,
							Status: cns.StatusGoalEnded,
						},
					}),
				},
			},
		},
	}
	return result
}
