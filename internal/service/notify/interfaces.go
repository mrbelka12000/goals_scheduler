package notify

import (
	"context"

	"github.com/mrbelka12000/goals_scheduler/internal/models"
)

type repo interface {
	Create(ctx context.Context, obj models.NotifyCU) (int64, error)
	Get(ctx context.Context, pars models.NotifyPars) (models.Notify, error)
}
