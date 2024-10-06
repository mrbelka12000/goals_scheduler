package usecase

import (
	gs "github.com/mrbelka12000/goals_scheduler"
)

type schema struct {
	state     gs.State
	nextState gs.State
	msg       string

	isStart            bool
	isFinal            bool
	waitingForText     bool // text of goal
	waitingForDeadline bool // date 15-02-2002
	waitingForTime     bool // time 15:13, 23:55...
	waitingForTimer    bool // timer 5s, 10h...
	waitingForDay      bool // weekday Monday, Tuesday...
	needInput          bool // need text from user or not
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
			needInput:      true,
			waitingForText: true,
		},
		{
			state:              gs.MessageStateDeadline,
			nextState:          gs.MessageStateChoseMethod,
			needInput:          true,
			waitingForDeadline: true,
		},
		{
			state: gs.MessageStateChoseMethod,
		},
		{
			state:          gs.MessageStateTime,
			nextState:      gs.MessageStateDay,
			needInput:      true,
			waitingForTime: true,
		},
		{
			state:     gs.MessageStateDay,
			nextState: gs.MessageStateDone,
			needInput: true,
			isFinal:   true,
		},
		{
			state:           gs.MessageStateTimer,
			nextState:       gs.MessageStateDone,
			msg:             "Цель сохранилась",
			needInput:       true,
			waitingForTimer: true,
			isFinal:         true,
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
