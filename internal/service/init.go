package service

import (
	"github.com/mrbelka12000/goals_scheduler/internal/repo"
	goalservice "github.com/mrbelka12000/goals_scheduler/internal/service/goal"
)

type Service struct {
	Goal *goalservice.Goal
}

func New(repo *repo.Repo) *Service {
	return &Service{
		Goal: goalservice.NewGoal(repo.Goal),
	}
}
