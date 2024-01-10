package usecase

import (
	"context"

	"goals_scheduler/internal/models"
)

func (uc *UseCase) NotifierCreate(ctx context.Context, obj *models.NotifierCU) (int64, error) {
	return uc.srv.Notifier.Create(ctx, obj)
}

func (uc *UseCase) NotifierGet(ctx context.Context, id int64) (models.Notifier, error) {
	return uc.srv.Notifier.Get(ctx, id)
}

func (uc *UseCase) NotifierList(ctx context.Context, pars models.NotifierPars) ([]models.Notifier, int64, error) {
	return uc.srv.Notifier.List(ctx, pars)
}

func (uc *UseCase) NotifierUpdate(ctx context.Context, obj models.NotifierCU, id int64) error {
	return uc.srv.Notifier.Update(ctx, obj, id)
}
