package messages

import (
	"errors"
)

var (
	ErrStateNotFound = errors.New("state not found")
	ErrInvalidState  = errors.New("invalid state")
)

type (
	Message struct {
		UserID int64
		ChatID int64
		Text   string
	}

	DayInfo struct {
		Mark []bool
	}
)
