package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/AlekSi/pointer"

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
		uc.log.Error().Err(err).Msg("start goal")
	}

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

	currSchema := getNextSchema(state)

	if currSchema.needInput && msg.Text == "" {
		return gs.SomethingWentWrong, "-", fmt.Errorf("no input provided")
	}

	if currSchema.waitingForText {
		if err = uc.cache.Set(gs.GetKeyText(msg.UserID), msg.Text, blockTime); err != nil {
			return gs.SomethingWentWrong, "-", fmt.Errorf("set text message in cache: %w", err)
		}
	}

	if currSchema.waitingForDeadline {
		if err = uc.cache.Set(gs.GetKeyDeadline(msg.UserID), msg.Text, blockTime); err != nil {
			return gs.SomethingWentWrong, "-", fmt.Errorf("set deadline message in cache: %w", err)
		}
	}

	if currSchema.waitingForTimer {
		if err = uc.cache.Set(gs.GetKeyTimer(msg.UserID), msg.Text, blockTime); err != nil {
			return gs.SomethingWentWrong, "-", fmt.Errorf("set timer message in cache: %w", err)
		}
	}

	if currSchema.waitingForTime {
		hours, minutes, err := validateTime(msg.Text)
		if err != nil {
			return gs.SomethingWentWrong, "-", fmt.Errorf("validate time message: %w", err)
		}

		if err = uc.cache.Set(gs.GetKeyHour(msg.UserID), fmt.Sprint(hours), blockTime); err != nil {
			return gs.SomethingWentWrong, "-", fmt.Errorf("set hour in cache: %w", err)
		}
		if err = uc.cache.Set(gs.GetKeyMinute(msg.UserID), fmt.Sprint(minutes), blockTime); err != nil {
			return gs.SomethingWentWrong, "-", fmt.Errorf("set minutes in cache: %w", err)
		}
	}

	if currSchema.waitingForDay {
		if err = uc.cache.Set(gs.GetKeyDay(msg.UserID), msg.Text, blockTime); err != nil {
			return gs.SomethingWentWrong, "-", fmt.Errorf("set day message in cache: %w", err)
		}
	}

	if currSchema.needToChangeState {
		if err = uc.cache.Set(gs.GetKeyState(msg.UserID), currSchema.nextState, blockTime); err != nil {
			return gs.SomethingWentWrong, "-", fmt.Errorf("set state %s in cache: %w", currSchema.nextState, err)
		}
	}

	if currSchema.isFinal {
		err = uc.BuildGoal(msg.UserID, msg.ChatID)
		if err != nil {
			return gs.SomethingWentWrong, "-", fmt.Errorf("build goal: %w", err)
		}
	}

	return currSchema.msg, currSchema.nextState, nil
}

func (uc *UseCase) ChooseMethod(userID int, val gs.State, choose gs.Key) error {
	if err := uc.cache.Set(gs.GetKeyState(userID), val, blockTime); err != nil {
		return fmt.Errorf("set state %s in cache: %w", val, err)
	}

	if err := uc.cache.Set(gs.GetKeyChoose(userID), choose, blockTime); err != nil {
		return fmt.Errorf("set choose %s in cache: %w", val, err)
	}

	return nil
}

func (uc *UseCase) BuildGoal(userID int, chatID string) error {
	var (
		goal   models.GoalCU
		notify models.NotifyCU
		err    error
	)

	val, ok := uc.cache.Get(gs.GetKeyDeadline(userID))
	if !ok {
		return errors.New("no deadline in cache")
	}

	parsedTime, err := time.Parse(gs.DateFormat, val)
	if err != nil {
		return fmt.Errorf("parse time: %w", err)
	}
	goal.Deadline = &parsedTime

	val, ok = uc.cache.Get(gs.GetKeyChoose(userID))
	if ok {
		switch val {
		case string(gs.KeyNotify):
			val, ok = uc.cache.Get(gs.GetKeyHour(userID))
			if !ok {
				return errors.New("no hours in cache")
			}
			hours, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("parse hour: %w", err)
			}
			notify.Hour = &hours

			val, ok = uc.cache.Get(gs.GetKeyMinute(userID))
			if !ok {
				return errors.New("no minutes in cache")
			}
			minutes, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("parse hour: %w", err)
			}
			notify.Minute = &minutes

			val, ok = uc.cache.Get(gs.GetKeyDay(userID))
			if !ok {
				return errors.New("no day in cache")
			}
			var obj models.DayInfo
			if err = json.Unmarshal([]byte(val), &obj); err != nil {
				return fmt.Errorf("parse day info: %w", err)
			}
			notify.DayInfo = obj
			goal.NotifyEnabled = true
		case string(gs.KeyTimer):
			val, ok = uc.cache.Get(gs.GetKeyTimer(userID))
			if !ok {
				return errors.New("no timer in cache")
			}

			durStr := val
			if durStr != "" && durStr != "-" {
				dur, err := time.ParseDuration(durStr)
				if err != nil {
					return fmt.Errorf("can not parse duration %s: %w", durStr, err)
				}
				goal.Timer = &dur
				goal.TimerEnabled = true
			}
		default:
			return fmt.Errorf("unknown choose method: %v", val)
		}
	}

	val, ok = uc.cache.Get(gs.GetKeyText(userID))
	if !ok {
		return errors.New("no deadline in cache")
	}

	goal.Text = &val
	goal.UsrID = pointer.ToInt(userID)
	goal.ChatID = pointer.ToString(chatID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	goalID, err := uc.CreateGoal(ctx, goal)
	if err != nil {
		return fmt.Errorf("goal create: %w", err)
	}

	if goal.NotifyEnabled {
		notify.GoalID = &goalID
		var errorTrigered bool
		for i, v := range notify.DayInfo.Mark {
			if v {
				wd := gs.Day(i)
				notify.Weekday = &wd
				_, err = uc.NotifyCreate(ctx, notify)
				if err != nil {
					uc.log.Err(err).Msg("notify create")
					errorTrigered = true
					break
				}
			}
		}
		if errorTrigered {
			return uc.GoalDelete(ctx, goalID)
		}
	}

	return nil
}

func validateTime(timeStr string) (hours int, minutes int, err error) {
	if timeStr == "" {
		return 0, 0, errors.New("time is empty")
	}
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0, 0, errors.New("time is invalid")
	}

	hours, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("hour is not number")
	}
	if hours < 0 || hours > 23 {
		return 0, 0, fmt.Errorf("hour must be between 0 and 23")
	}

	minutes, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("minute is not number")
	}
	if minutes < 0 || minutes > 59 {
		return 0, 0, fmt.Errorf("minute must be between 0 and 59")
	}

	return hours, minutes, nil
}
