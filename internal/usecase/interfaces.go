package usecase

import (
	"context"
	"time"

	"goals_scheduler/internal/client/webhooker"
)

type (
	cacher interface {
		Delete(key string)
		Set(key string, value interface{}, dur time.Duration) error
		Get(key string) (string, bool)
		GetInt64(key string) (int64, bool)
	}
	webHooker interface {
		CreateWebHook(ctx context.Context, req webhooker.CreateWebHookRequest) error
	}
)
