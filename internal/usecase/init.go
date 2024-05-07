package usecase

import (
	"github.com/rs/zerolog"

	"goals_scheduler/internal/service"
)

type UseCase struct {
	log   zerolog.Logger
	srv   *service.Service
	cache cacher
}

func New(log zerolog.Logger, srv *service.Service, cache cacher) *UseCase {
	return &UseCase{
		log:   log,
		srv:   srv,
		cache: cache,
	}
}
