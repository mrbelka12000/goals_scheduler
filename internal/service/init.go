package service

import (
	"goals_scheduler/internal/repo"
	goalservice "goals_scheduler/internal/service/goal"
	notifservice "goals_scheduler/internal/service/notifier"
)

type Service struct {
	Goal     *goalservice.Goal
	Notifier *notifservice.Notifier
}

func New(repo *repo.Repo) *Service {
	return &Service{
		Goal:     goalservice.NewGoal(repo.Goal),
		Notifier: notifservice.NewNotifier(repo.Notifier),
	}
}
