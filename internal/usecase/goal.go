package usecase

import (
	"context"

	"github.com/AlekSi/pointer"

	"goals_scheduler/internal/cns"
	"goals_scheduler/internal/models"
)

func (uc *UseCase) GoalCreate(ctx context.Context, obj models.GoalCU) (int64, error) {
	id, err := uc.srv.Goal.Create(ctx, &obj)
	if err != nil {
		return 0, err
	}

	if obj.NotifyEnabled {
		_, err = uc.NotifierCreate(ctx, &models.NotifierCU{
			UsrID:   obj.UsrID,
			ChatID:  obj.ChatID,
			GoalID:  &id,
			Notify:  obj.NotifyTime,
			Status:  pointer.To(cns.StatusNotifierStarted),
			EndTime: *obj.Deadline,
		})
		if err != nil {
			uc.log.Err(err).Msg("notifier create")
		}
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
