package messages

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	gs "github.com/mrbelka12000/goals_scheduler"
	"github.com/mrbelka12000/goals_scheduler/goals"
	"github.com/mrbelka12000/goals_scheduler/pkg/ptr"
	"github.com/mrbelka12000/goals_scheduler/scheme"
)

const (
	defaultCacheDuration = 30 * time.Minute
)

type (
	Service interface {
		HandleMessage(ctx context.Context, msg Message) (response string, err error)
		HandleCallback(ctx context.Context, cb Callback) (response string, err error)
		HandleStart(ctx context.Context, msg Message) (response string, err error)
	}

	service struct {
		cache cacher
		goals goals.Service

		scheme scheme.Service
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
) Service {
	s := &service{
		cache:  cache,
		goals:  goals,
		scheme: scheme,
	}

	return s
}

func (s *service) HandleMessage(ctx context.Context, msg Message) (response string, err error) {
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

	if nextScheme.Action.IsFinal {
		goal, err := s.collectInfoForGoal(msg.UserID, msg.ChatID)
		if err != nil {
			return "", fmt.Errorf("collect info for goals: %w", err)
		}

		s.clearCache(msg.UserID)

		_, err = s.goals.Create(ctx, goal)
		if err != nil {
			return "", fmt.Errorf("create goals: %w", err)
		}

		return nextScheme.Action.MessageToUser, nil
	}

	err = s.setState(msg.UserID, nextScheme.Action.NextState)
	if err != nil {
		return "", fmt.Errorf("set next scheme: %w", err)
	}

	return nextScheme.Action.MessageToUser, nil
}

func (s *service) HandleCallback(ctx context.Context, cb Callback) (response string, err error) {
	return
}

func (s *service) HandleStart(ctx context.Context, msg Message) (response string, err error) {
	s.clearCache(msg.UserID)

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

func (s *service) collectInfoForGoal(userID int, chatID string) (goals.Goal, error) {
	goal := goals.Goal{
		TelegramChatID: chatID,
	}

	values := s.collectAllKeyValuesFromCache(userID)

	// set default data
	{
		goal.ScheduleType = values[gs.KeyScheduleType].(gs.ScheduleType)
		goal.Message = values[gs.KeyMessage].(string)
		goal.Status = gs.StatusGoalStarted
		goal.CreatedAt = time.Now()
		goal.UpdatedAt = time.Now()
	}

	switch goal.ScheduleType {

	case gs.ScheduleTypeTimer:
		intervalAny, ok := values[gs.KeyInterval]
		if !ok {
			return goal, fmt.Errorf("interval not set")
		}

		interval, ok := intervalAny.(string)
		if !ok {
			return goal, fmt.Errorf("invalid interval")
		}

		duration, err := time.ParseDuration(interval)
		if err != nil {
			return goal, fmt.Errorf("parse duration: %w", err)
		}

		goal.IntervalSeconds = ptr.Ref(duration.Seconds())
		goal.NextExecution = ptr.Ref(time.Now().Add(duration))

	case gs.ScheduleTypeDaily:

		daysAny, ok := values[gs.KeyDays]
		if !ok {
			return goal, fmt.Errorf("days not set")
		}

		daysRaw, ok := daysAny.(string)
		if !ok {
			return goal, fmt.Errorf("invalid days")
		}

		_ = daysRaw

	case gs.ScheduleTypeOnce:

		scheduledTime, err := getScheduledTime(values)
		if err != nil {
			return goal, err
		}

		goal.ScheduledTime = ptr.Ref(scheduledTime)

	}

	return goal, nil
}

func (s *service) collectAllKeyValuesFromCache(userID int) map[gs.Key]any {
	values := make(map[gs.Key]any)

	for _, key := range gs.KeysToGoal {
		val, _ := s.cache.Get(fmt.Sprintf("%v:%v", key, userID))
		values[key] = val
	}

	return values
}

func (s *service) clearCache(userID int) {

	// delete previous states
	for _, k := range gs.KeysToGoal {
		key := fmt.Sprintf("%v:%v", k, userID)
		s.cache.Delete(key)
	}

}

func (s *service) setState(userID int, state gs.State) error {
	return s.cache.Set(gs.GetKeyState(userID), state, defaultCacheDuration)
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

func getScheduledTime(values map[gs.Key]any) (time.Time, error) {
	dateAny, ok := values[gs.KeyDate]
	if !ok {
		return time.Time{}, fmt.Errorf("date not set")
	}

	date, ok := dateAny.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("invalid date")
	}

	dateValues := strings.Split(date, "-")
	if len(dateValues) != 3 {
		return time.Time{}, fmt.Errorf("invalid date")
	}

	year, err := strconv.Atoi(dateValues[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid year")
	}

	month, err := strconv.Atoi(dateValues[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid month")
	}

	day, err := strconv.Atoi(dateValues[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day")
	}

	hourAny, ok := values[gs.KeyHour]
	if !ok {
		return time.Time{}, fmt.Errorf("hour not set")
	}

	hourStr, ok := hourAny.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("invalid hour")
	}

	hour, err := strconv.Atoi(hourStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid hour")
	}

	minuteAny, ok := values[gs.KeyMinute]
	if !ok {
		return time.Time{}, fmt.Errorf("minute not set")
	}

	minuteStr, ok := minuteAny.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("invalid minute")
	}

	minute, err := strconv.Atoi(minuteStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid minute")
	}

	scheduledTime := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.Local)

	return scheduledTime, nil
}
