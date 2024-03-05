package usecase

import (
	"context"
	"fmt"

	"goals_scheduler/internal/client/webhooker"
	"goals_scheduler/internal/models"
)

func (uc *UseCase) NotifierCreate(ctx context.Context, obj *models.NotifierCU) (int64, error) {
	id, err := uc.srv.Notifier.Create(ctx, obj)
	if err != nil {
		return 0, err
	}

	err = uc.webHooker.CreateWebHook(ctx, webhooker.CreateWebHookRequest{
		Params: map[string]string{
			"id": fmt.Sprint(id),
		},
		EndTime: obj.EndTime,
	})
	if err != nil {
		uc.log.Err(err).Msg("can not create web hooker notification")
	}

	return id, nil
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
