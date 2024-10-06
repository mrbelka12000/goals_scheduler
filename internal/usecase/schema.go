package usecase

import (
	gs "github.com/mrbelka12000/goals_scheduler"
)

type schema struct {
	state     gs.State
	nextState gs.State
	msg       string

	isStart            bool
	waitingForText     bool
	waitingForDeadline bool
	waitingForNotify   bool
	waitingForTimer    bool
}

func initSchema() []schema {
	return []schema{
		{
			state:     gs.MessageStateStart,
			nextState: gs.MessageStateText,
			isStart:   true,
			msg:       "Введите текст цели",
		},
		{
			state:          gs.MessageStateText,
			nextState:      gs.MessageStateDeadline,
			msg:            "Введите крайний срок для цели",
			waitingForText: true,
		},
		{
			state:              gs.MessageStateDeadline,
			nextState:          gs.MessageStateChoseMethod,
			waitingForDeadline: true,
		},
		{
			state: gs.MessageStateChoseMethod,
		},
		{
			state:            gs.MessageStateNotify,
			nextState:        gs.MessageStateDone,
			msg:              "Цель сохранилась",
			waitingForNotify: true,
		},
		{
			state:            gs.MessageStateTimer,
			nextState:        gs.MessageStateDone,
			msg:              "Цель сохранилась",
			waitingForNotify: true,
		},
	}
}

func getNextSchema(currentState gs.State) (next schema) {
	schemas := initSchema()

	for _, s := range schemas {
		if currentState == s.state {
			return s
		}
	}

	return schemas[0]
}
