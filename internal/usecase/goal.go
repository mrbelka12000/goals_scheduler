package usecase

import (
	"context"

	"github.com/mrbelka12000/goals_scheduler/internal/models"
)

func (uc *UseCase) GoalCreate(ctx context.Context, obj models.GoalCU) (int64, error) {
	id, err := uc.srv.Goal.Create(ctx, &obj)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (uc *UseCase) GoalGet(ctx context.Context, id int64) (models.Goal, error) {
	return uc.srv.Goal.Get(ctx, id)
}

func (uc *UseCase) GoalList(ctx context.Context, pars models.GoalPars) ([]models.Goal, int64, error) {
	return uc.srv.Goal.List(ctx, pars)
}

func (uc *UseCase) GoalDeleteAllOfUsers(ctx context.Context, usrID int) error {
	return uc.srv.Goal.DeleteAllOfUsers(ctx, usrID)
}

func (uc *UseCase) GoalUpdate(ctx context.Context, obj models.GoalCU, id int64) error {
	return uc.srv.Goal.Update(ctx, obj, id)
}

func (uc *UseCase) GoalDelete(ctx context.Context, id int64) error {
	return uc.srv.Goal.Delete(ctx, id)
}
