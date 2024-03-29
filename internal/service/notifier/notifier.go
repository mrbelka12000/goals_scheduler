package notifier

import (
	"context"
	"fmt"
	"time"

	"github.com/AlekSi/pointer"

	"goals_scheduler/internal/cns"
	"goals_scheduler/internal/models"
)

const (
	infDur = 24 * 365 * 10 * time.Hour
)

type Notifier struct {
	repo repo
}

func NewNotifier(repo repo) *Notifier {
	return &Notifier{
		repo: repo,
	}
}

func (n *Notifier) Create(ctx context.Context, obj *models.NotifierCU) (int64, error) {
	obj.Status = pointer.To(cns.StatusNotifierStarted)
	if obj.Notify == nil {
		obj.Notify = pointer.ToDuration(infDur)
	}
	obj.LastUpdated = pointer.ToTime(time.Now().Add(*obj.Notify))

	id, err := n.repo.Create(ctx, obj)
	if err != nil {
		return 0, fmt.Errorf("create goal in db: %w", err)
	}

	return id, nil
}

func (n *Notifier) Get(ctx context.Context, id int64) (models.Notifier, error) {
	return n.repo.Get(ctx, id)
}

func (n *Notifier) List(ctx context.Context, pars models.NotifierPars) ([]models.Notifier, int64, error) {
	return n.repo.List(ctx, pars)
}

func (n *Notifier) Update(ctx context.Context, obj models.NotifierCU, id int64) error {
	return n.repo.Update(ctx, obj, id)
}
