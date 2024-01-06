package usecase

import (
	"time"
)

type cacher interface {
	Delete(key string)
	Set(key string, value interface{}, dur time.Duration) error
	Get(key string) (string, bool)
	GetInt64(key string) (int64, bool)
}
