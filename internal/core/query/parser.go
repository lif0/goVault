package query

type FMSMachine struct {
}

func NewParser() Parser {
	return &FMSMachine{}
}

func (p *FMSMachine) Transition(query string) (out *Query, err error) {
	fsm := NewFSM(query)
	for fsm.state != FINISH {
		fsm.Next()
	}

	if fsm.err != nil {
		return nil, fsm.err
	}

	return fsm.out, nil
}
