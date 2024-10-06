package bot

import (
	"context"
	"fmt"

	"github.com/AlekSi/pointer"
	"github.com/yanzay/tbot/v2"

	gs "github.com/mrbelka12000/goals_scheduler"
	"github.com/mrbelka12000/goals_scheduler/internal/models"
)

const (
	ActionGoalDelete = "delete"
	ActionGoalUpdate = "update"
	ActionGoalSelect = "select"

	ActionGoalCreateTimer  = "timer"
	ActionGoalCreateNotify = "notify"
)

func (a *Application) deleteGoal(m *tbot.Message) {
	list, _, err := a.uc.GoalList(context.Background(), models.GoalPars{UsrID: pointer.ToInt(m.From.ID)})
	if err != nil {
		a.client.SendMessage(m.Chat.ID, gs.SomethingWentWrong)
		a.log.Err(err).Msg("get goals list")
		return
	}

	if len(list) == 0 {
		a.client.SendMessage(m.Chat.ID, "Пока что у вас нет целей")
		return
	}

	a.client.SendMessage(
		m.Chat.ID,
		"Выберите цель для удаления",
		tbot.OptInlineKeyboardMarkup(generateGoalBottons(list, true, ActionGoalDelete, "-")))
}

func (a *Application) deleteUsersGoals(m *tbot.Message) {
	err := a.uc.GoalDeleteAllOfUsers(context.Background(), m.From.ID)
	if err != nil {
		a.log.Err(err).Msg("delete user`s goals")
		a.client.SendMessage(m.Chat.ID, gs.SomethingWentWrong)
		return
	}

	a.client.SendMessage(m.Chat.ID, "Все удалено")
}

func (a *Application) handleCreateGoal(m *tbot.Message) {
	msg := a.uc.StartGoal(models.Message{
		UserID: m.From.ID,
		Text:   m.Text,
	})

	a.client.SendMessage(m.Chat.ID, msg)
}

func (a *Application) handleGetGoals(m *tbot.Message) {
	list, _, err := a.uc.GoalList(context.Background(), models.GoalPars{UsrID: pointer.ToInt(m.From.ID)})
	if err != nil {
		a.client.SendMessage(m.Chat.ID, gs.SomethingWentWrong)
		a.log.Err(err).Msg("get goals list")
		return
	}

	if len(list) == 0 {
		a.client.SendMessage(m.Chat.ID, "Пока что у вас нет целей")
		return
	}

	a.client.SendMessage(
		m.Chat.ID,
		"Цели",
		tbot.OptInlineKeyboardMarkup(generateGoalBottons(list, true, ActionGoalSelect, "-")))
}

func (a *Application) handleCallbackGoal(cq *tbot.CallbackQuery, data *models.GoalData) string {
	if data.Action == "-" {
		a.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
		return ""
	}

	switch data.Action {
	case ActionGoalDelete:
		err := a.uc.GoalDelete(
			context.Background(),
			data.ID,
		)
		if err != nil {
			a.log.Err(err).Msg("delete goal by id")
			return gs.SomethingWentWrong
		}

		a.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)

		return "Цель удалена"
	case ActionGoalUpdate:
		err := a.uc.GoalUpdate(
			context.Background(),
			models.GoalCU{
				Status: &data.Status,
			},
			data.ID,
		)
		if err != nil {
			a.log.Err(err).Msg("update goal by id")
			return gs.SomethingWentWrong
		}

		a.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
		return "Цель обновлена"
	case ActionGoalSelect:
		goal, err := a.uc.GoalGet(context.Background(), data.ID)
		if err != nil {
			return gs.SomethingWentWrong
		}

		return fmt.Sprintf("%v|%v", "SelectGoal", goal.Text)
	}

	return ""
}

func GetGoalActions(id int64) *tbot.InlineKeyboardMarkup {
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			{
				{
					Text: gs.StatusMapper(gs.StatusGoalFailed),
					CallbackData: callbackDataBuilder(models.CallbackData{
						Type: gs.CallbackTypeGoal,
						Goal: &models.GoalData{
							ID:     id,
							Action: ActionGoalUpdate,
							Status: gs.StatusGoalFailed,
						},
					}),
				},
				{
					Text: gs.StatusMapper(gs.StatusGoalEnded),
					CallbackData: callbackDataBuilder(models.CallbackData{
						Type: gs.CallbackTypeGoal,
						Goal: &models.GoalData{
							ID:     id,
							Action: ActionGoalUpdate,
							Status: gs.StatusGoalEnded,
						},
					}),
				},
			},
			{
				{
					Text: "Отмена",
					CallbackData: callbackDataBuilder(models.CallbackData{
						Type: gs.CallbackTypeGoal,
						Goal: &models.GoalData{
							Action: "-",
						},
					}),
				},
			},
		},
	}
}

func (a *Application) handleGoalCreate(cq *tbot.CallbackQuery, data *models.GoalCreateData) string {
	a.client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	switch data.Action {
	case ActionGoalCreateTimer:
		err := a.uc.ChooseMethod(cq.From.ID, gs.MessageStateTimer, gs.KeyTimer)
		if err != nil {
			a.log.Err(err).Msg("choose method")
			return gs.SomethingWentWrong
		}

		return gs.MessageTimerFormat

	case ActionGoalCreateNotify:
		err := a.uc.ChooseMethod(cq.From.ID, gs.MessageStateTime, gs.KeyNotify)
		if err != nil {
			a.log.Err(err).Msg("choose method")
			return gs.SomethingWentWrong
		}

		return gs.MessageTimeFromat

	case "-":
		err := a.uc.BuildGoal(cq.From.ID, cq.Message.Chat.ID)
		if err != nil {
			a.log.Err(err).Msg("create goal")
			return gs.SomethingWentWrong
		}

		return gs.MessageDone
	}

	return ""
}

func (a *Application) GetGoalCreateActions() *tbot.InlineKeyboardMarkup {
	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			{
				{
					Text: "Уведомления",
					CallbackData: callbackDataBuilder(models.CallbackData{
						Type: gs.CallbackTypeGoalCreate,
						GoalCreate: &models.GoalCreateData{
							Action: ActionGoalCreateNotify,
						},
					}),
				},
				{
					Text: "Таймер",
					CallbackData: callbackDataBuilder(models.CallbackData{
						Type: gs.CallbackTypeGoalCreate,
						GoalCreate: &models.GoalCreateData{
							Action: ActionGoalCreateTimer,
						},
					}),
				},
			},
			{
				{
					Text: "Не напоминать",
					CallbackData: callbackDataBuilder(models.CallbackData{
						Type: gs.CallbackTypeGoalCreate,
						GoalCreate: &models.GoalCreateData{
							Action: "-",
						},
					}),
				},
			},
		},
	}

}
