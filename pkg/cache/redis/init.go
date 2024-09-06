package redis

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"

	"goals_scheduler/pkg/config"
)

type Cache struct {
	store *redis.Client
}

func New(cfg config.Config) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &Cache{
		store: client,
	}, nil
}

func (c *Cache) Set(key string, value interface{}, dur time.Duration) error {
	return c.store.Set(key, value, dur).Err()
}

func (c *Cache) Get(key string) (string, bool) {
	value, err := c.store.Get(key).Result()
	if err != nil {
		return "", false
	}

	return value, true
}

func (c *Cache) GetInt64(key string) (int64, bool) {
	strValue, ok := c.Get(key)
	if !ok {
		return 0, false
	}

	value, err := strconv.ParseInt(strValue, 10, 0)
	if err != nil {
		return 0, false
	}

	return value, true
}

func (c *Cache) Delete(key string) {
	c.store.Del(key)
}
