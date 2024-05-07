package bot

import (
	"context"
	"strconv"

	"github.com/AlekSi/pointer"
	"github.com/yanzay/tbot/v2"

	"goals_scheduler/internal/cns"
	"goals_scheduler/internal/models"
)

const (
	CallbackGoal = "goal"
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

	a.Client.SendMessage(m.Chat.ID, "Выберите цель для удаления", tbot.OptInlineKeyboardMarkup(generateGoalBottons(list, true)))
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

func (a *Application) handleGetGoal(m *tbot.Message) {
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

	a.Client.SendMessage(m.Chat.ID, "Цели", tbot.OptInlineKeyboardMarkup(generateGoalBottons(list, false)))
}

func (a *Application) handleCallbackGoal(cq *tbot.CallbackQuery) string {
	idStr := cq.Data[len(CallbackGoal)+1:]

	if idStr == "-" {
		a.Client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
		return ""
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		a.Log.Err(err).Msg("parse goal id in callback")
		return cns.SomethingWentWrong
	}

	err = a.Uc.GoalDelete(context.Background(), id)
	if err != nil {
		a.Log.Err(err).Msg("delete goal by id")
		return cns.SomethingWentWrong
	}

	a.Client.DeleteMessage(cq.Message.Chat.ID, cq.Message.MessageID)
	return "Цель удалена"
}
