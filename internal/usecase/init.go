package usecase

import (
	"github.com/rs/zerolog"

	"goals_scheduler/internal/service"
)

type UseCase struct {
	log       zerolog.Logger
	srv       *service.Service
	cache     cacher
	webHooker webHooker
}

func New(log zerolog.Logger, srv *service.Service, cache cacher, wh webHooker) *UseCase {
	return &UseCase{
		log:       log,
		srv:       srv,
		cache:     cache,
		webHooker: wh,
	}
}
