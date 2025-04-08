package messages

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	gs "github.com/mrbelka12000/goals_scheduler"
	"github.com/mrbelka12000/goals_scheduler/goals"
	"github.com/mrbelka12000/goals_scheduler/scheme"
)

const (
	defaultCacheDuration = 30 * time.Minute
)

type (
	Service interface {
		HandleMessage(msg Message) (response string, err error)
		HandleCallback(cb Callback) (response string, err error)
		HandleStart(msg Message) (response string, err error)
	}

	service struct {
		cache cacher
		goals goals.Service

		scheme scheme.Service

		timerInterval time.Duration
	}

	cacher interface {
		Delete(key string)
		Set(key string, value interface{}, dur time.Duration) error
		Get(key string) (string, bool)
		GetInt64(key string) (int64, bool)
	}
)

func NewService(
	cache cacher,
	goals goals.Service,
	scheme scheme.Service,
	timerInterval time.Duration,
) Service {
	s := &service{
		cache:         cache,
		goals:         goals,
		scheme:        scheme,
		timerInterval: timerInterval,
	}

	go s.sendNotification()

	return s
}

func (s *service) HandleMessage(msg Message) (response string, err error) {
	currentState, ok := s.cache.Get(gs.GetKeyState(msg.UserID))
	if !ok {
		return "", ErrStateNotFound
	}

	nextScheme, err := s.scheme.GetNextScheme(gs.State(currentState))
	if err != nil {
		return "", fmt.Errorf("get next scheme: %w", err)
	}

	if nextScheme.Action.NeedInput && msg.Text == "" {
		return gs.MessageNeedText, nil
	}

	if nextScheme.Action.NeedTimerInterval {
		if err = s.setState(msg.UserID, gs.StateTimerInterval); err != nil {
			return "", fmt.Errorf("set state: %w", err)
		}
	}

	if nextScheme.Action.NeedTime {
		hours, minutes, err := parseTime(msg.Text)
		if err != nil {
			return "", fmt.Errorf("validate time: %w", err)
		}

		if err := s.cache.Set(gs.GetKeyHour(msg.UserID), hours, defaultCacheDuration); err != nil {
			return "", fmt.Errorf("set hours: %w", err)
		}

		if err := s.cache.Set(gs.GetKeyMinute(msg.UserID), minutes, defaultCacheDuration); err != nil {
			return "", fmt.Errorf("set minutes: %w", err)
		}
	}

	err = s.setState(msg.UserID, nextScheme.Action.NextState)
	if err != nil {
		return "", fmt.Errorf("set next scheme: %w", err)
	}

	return nextScheme.Action.MessageToUser, nil
}

func (s *service) HandleCallback(cb Callback) (response string, err error) {
	return
}

func (s *service) HandleStart(msg Message) (response string, err error) {
	// delete previous states
	for _, k := range gs.KeysToGoal {
		key := fmt.Sprintf("%v:%v", k, msg.UserID)
		s.cache.Delete(key)
	}

	err = s.setState(msg.UserID, gs.StateStart)
	if err != nil {
		return "", fmt.Errorf("failed to set goal state: %w", err)
	}

	nextScheme, err := s.scheme.GetNextScheme(gs.StateStart)
	if err != nil {
		return "", fmt.Errorf("failed to get next scheme: %w", err)
	}

	return nextScheme.Action.MessageToUser, nil
}

func (s *service) setState(userID int, state gs.State) error {
	return s.cache.Set(gs.GetKeyState(userID), state, defaultCacheDuration)
}

func (s *service) sendNotification() {
	ticker := time.NewTicker(s.timerInterval)

	for range ticker.C {
		//TODO get goals
	}
}

func parseTime(timeStr string) (hours int, minutes int, err error) {
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
