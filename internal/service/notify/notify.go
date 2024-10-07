package notify

import (
	"context"
	"errors"

	"github.com/mrbelka12000/goals_scheduler/internal/models"
)

type Notify struct {
	repo repo
}

func New(repo repo) *Notify {
	return &Notify{
		repo: repo,
	}
}

func (n *Notify) Create(ctx context.Context, obj models.NotifyCU) (int64, error) {
	return n.repo.Create(ctx, obj)
}

func (n *Notify) Get(ctx context.Context, pars models.NotifyPars) (models.Notify, error) {
	if pars.ID == nil && pars.GoalID == nil && pars.WeekDay == nil {
		return models.Notify{}, errors.New("no pars defined")
	}

	return n.repo.Get(ctx, pars)
}
