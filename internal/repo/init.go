package repo

import "database/sql"

type Repo struct {
	Goal *goal
}

func New(db *sql.DB) *Repo {
	return &Repo{
		Goal: newGoal(db),
	}
}
