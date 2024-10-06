package usecase

import (
	"github.com/rs/zerolog"

	"github.com/mrbelka12000/goals_scheduler/internal/service"
)

type UseCase struct {
	log     zerolog.Logger
	srv     *service.Service
	cache   cacher
	schemas []schema
}

func New(log zerolog.Logger, srv *service.Service, cache cacher) *UseCase {
	return &UseCase{
		log:     log,
		srv:     srv,
		cache:   cache,
		schemas: initSchema(),
	}
}
