package goals

import (
	"context"
	"fmt"
	"time"

	"github.com/AlekSi/pointer"

	gs "github.com/mrbelka12000/goals_scheduler"
)

type (
	Service interface {
		Create(ctx context.Context, obj GoalCU) (int64, error)
		Get(ctx context.Context, id int64) (Goal, error)
		List(ctx context.Context, pars GoalPars) ([]Goal, int64, error)
		DeleteAllOfUsers(ctx context.Context, usrID int) error
		Update(ctx context.Context, obj GoalCU, id int64) error
		Delete(ctx context.Context, id int64) error
	}

	service struct {
		repository Repository
	}
)

func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}

func (s *service) Create(ctx context.Context, obj GoalCU) (int64, error) {
	obj.Status = pointer.To(gs.StatusGoalStarted)
	if obj.Timer == nil {
		obj.Timer = pointer.ToDuration(365 * 24 * time.Hour)
	}

	obj.LastUpdated = pointer.To(time.Now().Add(*obj.Timer))

	err := s.validate(obj)
	if err != nil {
		return 0, fmt.Errorf("validate goal: %w", err)
	}

	id, err := s.repository.Create(ctx, obj)
	if err != nil {
		return 0, fmt.Errorf("create goal in db: %w", err)
	}

	return id, nil
}

func (s *service) Get(ctx context.Context, id int64) (Goal, error) {
	return s.repository.Get(ctx, id)
}

func (s *service) List(ctx context.Context, pars GoalPars) ([]Goal, int64, error) {
	return s.repository.List(ctx, pars)
}

func (s *service) DeleteAllOfUsers(ctx context.Context, usrID int) error {
	return s.repository.DeleteAllUsersGoals(ctx, usrID)
}

func (s *service) Update(ctx context.Context, obj GoalCU, id int64) error {
	return s.repository.Update(ctx, obj, id)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}

func (s *service) validate(goal GoalCU) error {
	if goal.Status == nil {
		return fmt.Errorf("invalid goal status")
	}
	if goal.Text == nil {
		return fmt.Errorf("invalid goal text")
	}
	if goal.ChatID == nil {
		return fmt.Errorf("invalid goal chatid")
	}
	if goal.UsrID == nil {
		return fmt.Errorf("invalid goal userid")
	}
	if goal.Timer == nil {
		goal.Timer = pointer.ToDuration(0)
	}

	return nil
}
