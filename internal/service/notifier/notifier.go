package notifier

type Notifier struct {
	repo repo
}

func NewNotifier(repo repo) *Notifier {
	return &Notifier{
		repo: repo,
	}
}
