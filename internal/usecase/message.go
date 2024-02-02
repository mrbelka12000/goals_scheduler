package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AlekSi/pointer"

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

		return fmt.Sprintf(`
		Введите время для напоминания
Допустимые единицы времени: "ns", "us" (или "µs"), "ms", "s", "m", "h".
Отправьте - , в случае если не нужно напоминать

		`), cns.MessageStateNotifier

	case cns.MessageStateNotifier:
		if msg.Text != "-" {
			uc.cache.Set(cns.GetKeyNotify(msg.UserID), msg.Text, 0)
		}
		mp := make(map[string]interface{})

		// collect states
		for _, k := range cns.KeysToGoal {
			key := fmt.Sprintf("%v:%v", k, msg.UserID)
			val, _ := uc.cache.Get(key)
			mp[k] = val
		}

		// parse time from request
		{
			parsedTime, err := time.Parse(cns.DateFormat, mp[cns.KeyDeadline].(string))
			if err != nil {
				uc.log.Err(err).Msg(fmt.Sprintf("parse time: %v", mp[cns.KeyDeadline]))
				return "Что то пошло не так", ""
			}
			mp[cns.KeyDeadline] = parsedTime
		}

		// parse ticker duration
		{
			durStr := mp[cns.KeyNotify].(string)
			if durStr != "" && durStr != "-" {
				dur, err := time.ParseDuration(durStr)
				if err != nil {
					uc.log.Err(err).Msg(fmt.Sprintf("parse duration: %v", mp[cns.KeyNotify]))
					return "Что то пошло не так", ""
				}
				mp[cns.KeyNotify] = dur
				mp[cns.KeyNotifyEnabled] = true
			} else {
				delete(mp, cns.KeyNotify)
			}
		}

		{
			var goal models.GoalCU

			jsonBody, _ := json.Marshal(mp)
			err := json.Unmarshal(jsonBody, &goal)
			if err != nil {
				uc.log.Err(err).Msg("goal from map")
				return "Что то пошло не так", ""
			}

			goal.UsrID = pointer.ToInt(msg.UserID)
			goal.ChatID = pointer.ToString(msg.ChatID)

			_, err = uc.GoalCreate(context.Background(), goal)
			if err != nil {
				uc.log.Err(err).Msg("goal create")
				return "Что то пошло не так", ""
			}
		}

		// delete states
		for _, k := range cns.KeysToGoal {
			key := fmt.Sprintf("%v:%v", k, msg.UserID)
			val, _ := uc.cache.Get(key)
			mp[k] = val
			uc.cache.Delete(key)
		}

		return "Цель сохранилась", ""
	}

	return "Не туда", ""
}
