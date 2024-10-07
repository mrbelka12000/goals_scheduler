package repo

import "database/sql"

type Repo struct {
	Goal   *goal
	Notify *notify
}

func New(db *sql.DB) *Repo {
	return &Repo{
		Goal:   newGoal(db),
		Notify: newNotify(db),
	}
}
