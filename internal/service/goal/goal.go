package goal

import (
	"context"
	"fmt"

	"github.com/AlekSi/pointer"

	"goals_scheduler/internal/cns"
	"goals_scheduler/internal/models"
)

type Goal struct {
	repo repo
}

func NewGoal(repo repo) *Goal {
	return &Goal{
		repo: repo,
	}
}

func (g *Goal) Create(ctx context.Context, obj *models.GoalCU) (int64, error) {
	obj.Status = pointer.To(cns.StatusGoalStarted)

	id, err := g.repo.Create(ctx, obj)
	if err != nil {
		return 0, fmt.Errorf("create goal in db: %w", err)
	}

	return id, nil
}

func (g *Goal) Get(ctx context.Context, id int64) (models.Goal, error) {
	return g.repo.Get(ctx, id)
}

func (g *Goal) List(ctx context.Context, pars models.GoalPars) ([]models.Goal, int64, error) {
	return g.repo.List(ctx, pars)
}

func (g *Goal) DeleteAllOfUsers(ctx context.Context, usrID int) error {
	return g.repo.DeleteAllUsersGoals(ctx, usrID)
}
