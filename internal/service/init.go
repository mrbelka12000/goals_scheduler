package service

import (
	"goals_scheduler/internal/repo"
	goalservice "goals_scheduler/internal/service/goal"
	notifyservice "goals_scheduler/internal/service/notifier"
)

type Service struct {
	Goal     *goalservice.Goal
	Notifier *notifyservice.Notifier
}

func New(repo *repo.Repo) *Service {
	return &Service{
		Goal:     goalservice.NewGoal(repo.Goal),
		Notifier: notifyservice.NewNotifier(repo.Notifier),
	}
}
