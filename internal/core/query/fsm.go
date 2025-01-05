package query

import (
	"regexp"
	"strings"

	"goVault/internal/pkg/btypes"
	"goVault/internal/pkg/pointer"
)

type state int

const (
	// FMS internal states
	BEGIN state = iota
	FINISH

	// user states
	CmdSetState
	CmdGetState
	CmdDelState
)

var validQuerySyntaxRegex = regexp.MustCompile(`[^a-zA-Z0-9*_/. ]`) // that is thread-safe.
var manySpacesInRawQuery = regexp.MustCompile(`\s+`)

type FSM struct {
	state state
	err   error

	in  []string
	out *Query
}

func newFSM(rawQuery string) FSM {
	if isValidQuerySyntax(rawQuery) {
		return FSM{state: BEGIN, in: strings.Fields(*formatQueryRaw(&rawQuery))}
	}

	return FSM{state: FINISH, err: ErrQueryUnknowCommand}
}

func (fsm *FSM) next() {
	switch fsm.state {
	case BEGIN:
		fsm.parseCmd()

	case CmdSetState:
		if isValidArgumentSetCmd(fsm.in) {
			fsm.state = FINISH
			cmdStArgs := []string{fsm.in[1], fsm.in[2]}
			fsm.out.Arguments = cmdStArgs
		} else {
			fsm.state = FINISH
			fsm.err = ErrQueryInvalidArgumentSyntaxSet
		}

	case CmdGetState:
		if isValidArgumentGetCmd(fsm.in) {
			fsm.state = FINISH
			cmdStArgs := btypes.ToArray(fsm.in[1])
			fsm.out.Arguments = cmdStArgs
		} else {
			fsm.state = FINISH
			fsm.err = ErrQueryInvalidArgumentSyntaxGet
		}

	case CmdDelState:
		if isValidArgumentDelCmd(fsm.in) {
			fsm.state = FINISH
			cmdStArgs := btypes.ToArray(fsm.in[1])
			fsm.out.Arguments = cmdStArgs
		} else {
			fsm.state = FINISH
			fsm.err = ErrQueryInvalidArgumentSyntaxDel
		}

	default:
		fsm.state = FINISH
		fsm.err = ErrQueryUnknowState
	}
}

func (fsm *FSM) parseCmd() bool {
	var query Query

	switch strings.ToLower(fsm.in[0]) {
	case "set":
		fsm.state = CmdSetState
		query = NewQuery(SET)

	case "get":
		fsm.state = CmdGetState
		query = NewQuery(GET)
	case "del":
		fsm.state = CmdDelState
		query = NewQuery(DEL)
	default:
		fsm.state = FINISH
		fsm.err = ErrQueryUnknowCommand
		return false
	}

	fsm.out = &query
	return true
}

func isValidQuerySyntax(queryRaw string) bool {
	return !validQuerySyntaxRegex.MatchString(queryRaw) // Are the query has any forbidden symbols
}

func isValidArgumentSetCmd(args []string) bool {
	return len(args) >= 3
}

func isValidArgumentGetCmd(args []string) bool {
	return len(args) >= 2
}

func isValidArgumentDelCmd(args []string) bool {
	return len(args) >= 2
}

func formatQueryRaw(rawQuery *string) *string {
	return pointer.To(strings.TrimSpace(manySpacesInRawQuery.ReplaceAllString(*rawQuery, " ")))
}
