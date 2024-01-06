package goal

type Goal struct {
	repo repo
}

func NewGoal(repo repo) *Goal {
	return &Goal{
		repo: repo,
	}
}
