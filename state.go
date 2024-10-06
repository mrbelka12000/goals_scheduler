package goals_scheduler

type State string

func (s State) MarshalBinary() ([]byte, error) {
	return []byte(s), nil
}

const (
	MessageStateStart       State = "start"
	MessageStateText        State = "waiting_for_text"
	MessageStateDeadline    State = "waiting_for_deadline"
	MessageStateChoseMethod State = "waiting_for_method"
	MessageStateTimer       State = "waiting_for_timer"
	MessageStateNotify      State = "waiting_for_notify"
	MessageStateDone        State = "done"
)
