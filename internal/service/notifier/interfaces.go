package notifier

import (
	"context"

	"goals_scheduler/internal/models"
)

type repo interface {
	Create(ctx context.Context, obj *models.NotifierCU) (int64, error)
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (models.Notifier, error)
	List(ctx context.Context, pars models.NotifierPars) ([]models.Notifier, int64, error)
}
