package goals

import (
	"context"
	"fmt"
)

type (
	Service interface {
		Create(ctx context.Context, obj Goal) (int64, error)
		Get(ctx context.Context, id int64) (Goal, error)
		List(ctx context.Context, pars GoalPars) ([]Goal, int64, error)
		DeleteAllOfUsers(ctx context.Context, chatID string) error
		Update(ctx context.Context, obj Goal, id int64) error
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

func (s *service) Create(ctx context.Context, obj Goal) (int64, error) {

	id, err := s.repository.CreateNotification(ctx, obj)
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

func (s *service) DeleteAllOfUsers(ctx context.Context, chatID string) error {
	return s.repository.DeleteAllUsersGoals(ctx, chatID)
}

func (s *service) Update(ctx context.Context, obj Goal, id int64) error {
	return s.repository.Update(ctx, obj, id)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}
