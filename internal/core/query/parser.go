package query

import (
	"errors"

	"goVault/internal"
)

type FMSMachine struct {
	logger internal.Logger
}

func NewParser(logger internal.Logger) (Parser, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &FMSMachine{
		logger: logger,
	}, nil
}

func (p *FMSMachine) Transition(query string) (out *Query, err error) {
	fsm := newFSM(query)
	for fsm.state != FINISH {
		fsm.next()
	}

	if fsm.err != nil {
		return nil, fsm.err
	}

	return fsm.out, nil
}
