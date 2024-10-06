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
	MessageStateTime        State = "waiting_for_time"
	MessageStateDay         State = "waiting_for_day"
	MessageStateDone        State = "done"
)
