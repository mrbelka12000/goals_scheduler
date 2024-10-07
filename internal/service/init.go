package service

import (
	"github.com/mrbelka12000/goals_scheduler/internal/repo"
	goalservice "github.com/mrbelka12000/goals_scheduler/internal/service/goal"
	notifyservice "github.com/mrbelka12000/goals_scheduler/internal/service/notify"
)

type Service struct {
	Goal   *goalservice.Goal
	Notify *notifyservice.Notify
}

func New(repo *repo.Repo) *Service {
	return &Service{
		Goal:   goalservice.New(repo.Goal),
		Notify: notifyservice.New(repo.Notify),
	}
}
