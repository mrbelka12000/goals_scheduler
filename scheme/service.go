package scheme

import gs "github.com/mrbelka12000/goals_scheduler"

type (
	Service interface {
		GetNextScheme(state gs.State) (Scheme, error)
	}

	service struct {
		scheme []Scheme
	}
)

func NewService() Service {
	return &service{
		scheme: initScheme(),
	}
}

func (s *service) GetNextScheme(state gs.State) (Scheme, error) {
	for _, scheme := range s.scheme {
		if scheme.CurrentState == state {
			return scheme, nil
		}
	}

	return Scheme{}, ErrSchemeNotFound
}

func initScheme() []Scheme {
	return []Scheme{
		{
			CurrentState: gs.StateStart,
			Action: Action{
				NextState:         gs.StateStart,
				NeedInput:         true,
				StrictChangeState: true,
				MessageToUser:     gs.MessageInputText,
			},
		},
		{
			CurrentState: gs.StateEnterGoal,
			Action: Action{
				NextState:         gs.StateChoseReminder,
				ChosenChangeState: true,
				MessageToUser:     gs.MessageChooseMethod,
			},
		},
		{
			CurrentState: gs.StateChoseReminder,
			Action: Action{
				ChosenChangeState: false,
			},
		},
		{
			CurrentState: gs.StateTimerInterval,
			Action: Action{
				NextState:         gs.StateDeadline,
				StrictChangeState: true,
				MessageToUser:     gs.MessageTimerFormat,
				NeedTimerInterval: true,
			},
		},
		{
			CurrentState: gs.StateAlarmDays,
			Action: Action{
				NextState:         gs.StateAlarmTime,
				StrictChangeState: true,
				MessageToUser:     gs.MessageChooseDay,
			},
		},
		{
			CurrentState: gs.StateAlarmTime,
			Action: Action{
				NextState:         gs.StateDeadline,
				NeedInput:         true,
				StrictChangeState: true,
				MessageToUser:     gs.MessageTimeFormat,
				NeedTime:          true,
			},
		},
		{
			CurrentState: gs.StateReminderTime,
			Action: Action{
				NextState:         gs.StateDone,
				NeedInput:         true,
				StrictChangeState: true,
				MessageToUser:     gs.MessageTimeFormat,
				NeedTime:          true,
			},
		},
	}
}
