package usecase

import (
	"fmt"

	"goals_scheduler/internal/cns"
	"goals_scheduler/internal/models"
)

func (uc *UseCase) StartGoal(msg models.Message) string {
	err := uc.cache.Set(cns.GetKeyState(msg.UserID), cns.MessageStateText, 0)
	if err != nil {
		uc.log.Error().Err(err).Msg("set cache")
		return "Возникла ошибка, повторите позже"
	}
	return "Введите текст цели"
}

func (uc *UseCase) HandleMessage(msg models.Message) (string, string) {
	state, stateOk := uc.cache.Get(cns.GetKeyState(msg.UserID))
	if stateOk && state != "" {
		return uc.handleStates(msg, state)
	}

	return "Все ок", ""
}

func (uc *UseCase) handleStates(msg models.Message, state string) (string, string) {
	switch state {
	case cns.MessageStateText:
		uc.cache.Set(cns.GetKeyText(msg.UserID), msg.Text, 0)
		uc.cache.Set(cns.GetKeyState(msg.UserID), cns.MessageStateDeadline, 0)
		return "Введите крайний срок для цели", cns.MessageStateDeadline
	case cns.MessageStateDeadline:
		uc.cache.Set(cns.GetKeyDeadline(msg.UserID), msg.Text, 0)
		uc.cache.Set(cns.GetKeyState(msg.UserID), cns.MessageStateNotifier, 0)
		return "Введите время для напоминания", cns.MessageStateNotifier
	case cns.MessageStateNotifier:
		uc.cache.Set(cns.GetKeyNotify(msg.UserID), msg.Text, 0)
		mp := make(map[string]interface{})
		for _, k := range cns.KeysToGoal {
			key := fmt.Sprintf("%v:%v", k, msg.UserID)
			val, _ := uc.cache.Get(key)
			mp[k] = val
			uc.cache.Delete(key)
		}
		fmt.Println(mp)

		return "Цель сохранилась", ""
	}

	return "Не туда", ""
}
