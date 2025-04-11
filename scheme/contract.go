package scheme

import (
	"errors"

	gs "github.com/mrbelka12000/goals_scheduler"
)

var (
	ErrSchemeNotFound = errors.New("scheme not found")
)

type (
	Scheme struct {
		CurrentState gs.State
		Action       Action
	}

	Action struct {
		NextState gs.State

		NeedInput         bool
		StrictChangeState bool
		ChosenChangeState bool

		NeedTimerInterval bool
		NeedTime          bool
		IsFinal           bool

		MessageToUser string
	}
)
