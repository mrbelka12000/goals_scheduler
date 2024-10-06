package usecase

import (
	"errors"
	"fmt"
	"time"

	gs "github.com/mrbelka12000/goals_scheduler"
	"github.com/mrbelka12000/goals_scheduler/internal/models"
)

const (
	blockTime = 5 * time.Minute
)

func (uc *UseCase) StartGoal(msg models.Message) string {
	// delete states
	for _, k := range gs.KeysToGoal {
		key := fmt.Sprintf("%v:%v", k, msg.UserID)
		uc.cache.Delete(key)
	}

	text, _, err := uc.handleGoalCreate(msg, true)
	if err != nil {
		uc.log.Error().Err(err).Msg("start goal create")
	}

	fmt.Println(text)
	return text
}

func (uc *UseCase) HandleMessage(msg models.Message) (string, gs.State) {
	text, state, err := uc.handleGoalCreate(msg, false)
	if err != nil {
		uc.log.Error().Err(err).Msg("handle goal create")
	}

	return text, state
}

func (uc *UseCase) handleGoalCreate(
	msg models.Message,
	isStart bool,
) (
	text string,
	nextState gs.State,
	err error,
) {
	var state gs.State
	stateStr, ok := uc.cache.Get(gs.GetKeyState(msg.UserID))
	if !ok {
		if !isStart {
			return "", "", errors.New("no need to handle")
		}
		state = gs.MessageStateStart
	} else {
		state = gs.State(stateStr)
	}

	fmt.Println(state)

	currSchema := getNextSchema(state)

	if currSchema.isStart {
		if err = uc.cache.Set(gs.GetKeyState(msg.UserID), currSchema.nextState, blockTime); err != nil {
			return gs.SomethingWentWrong, "-", fmt.Errorf("set state %s in cache: %w", currSchema.nextState, err)
		}
	}

	if currSchema.waitingForText {
		if err = uc.cache.Set(gs.GetKeyText(msg.UserID), msg.Text, blockTime); err != nil {
			return gs.SomethingWentWrong, "-", fmt.Errorf("set text message in cache: %w", err)
		}
		if err = uc.cache.Set(gs.GetKeyState(msg.UserID), currSchema.nextState, blockTime); err != nil {
			return gs.SomethingWentWrong, "-", fmt.Errorf("set state %s in cache: %w", currSchema.nextState, err)
		}
	}

	if currSchema.waitingForDeadline {
		if err = uc.cache.Set(gs.GetKeyDeadline(msg.UserID), msg.Text, blockTime); err != nil {
			return gs.SomethingWentWrong, "-", fmt.Errorf("set deadline message in cache: %w", err)
		}
		if err = uc.cache.Set(gs.GetKeyState(msg.UserID), currSchema.nextState, blockTime); err != nil {
			return gs.SomethingWentWrong, "-", fmt.Errorf("set state %s in cache: %w", currSchema.nextState, err)
		}
	}

	if currSchema.waitingForNotify {
		uc.cache.Set(gs.GetKeyNotify(msg.UserID), msg.Text, blockTime)

		// delete states
		for _, k := range gs.KeysToGoal {
			key := fmt.Sprintf("%v:%v", k, msg.UserID)
			uc.cache.Delete(key)
		}
	}

	return currSchema.msg, currSchema.nextState, nil
	//
	//switch state {
	//case cns.MessageStateText:
	//
	//	uc.cache.Set(gs.GetKeyText(msg.UserID), msg.Text, blockTime)
	//	uc.cache.Set(gs.GetKeyState(msg.UserID), cns.MessageStateDeadline, blockTime)
	//	return "Введите крайний срок для цели", cns.MessageStateDeadline, nil
	//
	//case cns.MessageStateDeadline:
	//
	//	return cns.NotifyFormat, cns.MessageStateTimer, nil
	//
	//case cns.MessageStateTimer:
	//	if msg.Text != "-" {
	//		uc.cache.Set(gs.GetKeyTimer(msg.UserID), msg.Text, blockTime)
	//	}
	//	mp := make(map[gs.Key]interface{})
	//
	//	// collect states
	//	for _, k := range gs.KeysToGoal {
	//		key := fmt.Sprintf("%v:%v", k, msg.UserID)
	//		val, _ := uc.cache.Get(key)
	//		mp[k] = val
	//	}
	//
	//	// parse time from request
	//	{
	//		parsedTime, err := time.Parse(cns.DateFormat, mp[gs.KeyDeadline].(string))
	//		if err != nil {
	//			uc.log.Err(err).Msg(fmt.Sprintf("parse time: %v", mp[gs.KeyDeadline]))
	//			return cns.SomethingWentWrong, "", nil
	//		}
	//		mp[gs.KeyDeadline] = parsedTime
	//	}
	//
	//	// parse ticker duration
	//	{
	//		durStr := mp[gs.KeyTimer].(string)
	//		if durStr != "" && durStr != "-" {
	//			dur, err := time.ParseDuration(durStr)
	//			if err != nil {
	//				uc.log.Err(err).Msg(fmt.Sprintf("parse duration: %v", mp[gs.KeyTimer]))
	//				return cns.SomethingWentWrong, "", nil
	//			}
	//			mp[gs.KeyTimer] = dur
	//			mp[gs.KeyTimerEnabled] = true
	//		} else {
	//			delete(mp, gs.KeyTimer)
	//		}
	//	}
	//
	//	{
	//		var goal models.GoalCU
	//
	//		jsonBody, _ := json.Marshal(mp)
	//		err := json.Unmarshal(jsonBody, &goal)
	//		if err != nil {
	//			uc.log.Err(err).Msg("goal from map")
	//			return cns.SomethingWentWrong, "", nil
	//		}
	//
	//		goal.UsrID = pointer.ToInt(msg.UserID)
	//		goal.ChatID = pointer.ToString(msg.ChatID)
	//
	//		_, err = uc.GoalCreate(context.Background(), goal)
	//		if err != nil {
	//			uc.log.Err(err).Msg("goal create")
	//			return cns.SomethingWentWrong, "", nil
	//		}
	//	}
	//
	//	return "Цель сохранилась", cns.MessageStateDone, nil
	//}
	//
	//return "Не туда", "", nil
}
