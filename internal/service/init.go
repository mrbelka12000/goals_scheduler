package service

import (
	"goals_scheduler/internal/repo"
	goalservice "goals_scheduler/internal/service/goal"
)

type Service struct {
	Goal *goalservice.Goal
}

func New(repo *repo.Repo) *Service {
	return &Service{
		Goal: goalservice.NewGoal(repo.Goal),
	}
}
