package usecase

import (
	"context"

	"github.com/mrbelka12000/goals_scheduler/internal/models"
)

func (uc *UseCase) NotifyCreate(ctx context.Context, obj models.NotifyCU) (int64, error) {
	return uc.srv.Notify.Create(ctx, obj)
}

func (uc *UseCase) NotifyGet(ctx context.Context, pars models.NotifyPars) (models.Notify, error) {
	return uc.srv.Notify.Get(ctx, pars)
}
