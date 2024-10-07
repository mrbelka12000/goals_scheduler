package goal

import (
	"context"
	"fmt"
	"time"

	"github.com/AlekSi/pointer"

	gs "github.com/mrbelka12000/goals_scheduler"
	"github.com/mrbelka12000/goals_scheduler/internal/models"
)

type Goal struct {
	repo repo
}

func New(repo repo) *Goal {
	return &Goal{
		repo: repo,
	}
}

func (g *Goal) Create(ctx context.Context, obj *models.GoalCU) (int64, error) {
	obj.Status = pointer.To(gs.StatusGoalStarted)
	if obj.Timer == nil {
		obj.Timer = pointer.ToDuration(365 * 24 * time.Hour)
	}

	obj.LastUpdated = pointer.To(time.Now().Add(*obj.Timer))

	fmt.Printf("%+v\n", *obj)
	err := g.validate(obj)
	if err != nil {
		return 0, fmt.Errorf("validate goal: %w", err)
	}

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

func (g *Goal) Update(ctx context.Context, obj models.GoalCU, id int64) error {
	return g.repo.Update(ctx, obj, id)
}

func (g *Goal) Delete(ctx context.Context, id int64) error {
	return g.repo.Delete(ctx, id)
}

func (g *Goal) validate(goal *models.GoalCU) error {
	if goal == nil {
		return fmt.Errorf("invalid goal")
	}

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
