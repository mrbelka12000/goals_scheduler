package repo

import "database/sql"

type Repo struct {
	Goal     *goal
	Notifier *notifier
}

func New(db *sql.DB) *Repo {
	return &Repo{
		Goal:     newGoal(db),
		Notifier: newNotifier(db),
	}
}
