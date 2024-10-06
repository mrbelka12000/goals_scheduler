package goal

import (
	"context"

	"github.com/mrbelka12000/goals_scheduler/internal/models"
)

type repo interface {
	Create(ctx context.Context, obj *models.GoalCU) (int64, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (models.Goal, error)
	List(ctx context.Context, pars models.GoalPars) ([]models.Goal, int64, error)
	DeleteAllUsersGoals(ctx context.Context, usrID int) error
	Update(ctx context.Context, obj models.GoalCU, id int64) error
}
