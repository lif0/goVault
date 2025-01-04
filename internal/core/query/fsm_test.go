package query

import (
	"testing"
)

func TestNewFSM(t *testing.T) {
	tests := []struct {
		name          string
		rawQuery      string
		expectedState state
		expectedError error
	}{
		{"Valid query", "SET key value", BEGIN, nil},
		{"Invalid query syntax", "INVALID@QUERY", FINISH, ErrQueryUnknowCommand},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := newFSM(tt.rawQuery)

			if fsm.state != tt.expectedState {
				t.Fatalf("expected state %v, got %v", tt.expectedState, fsm.state)
			}
			if fsm.err != tt.expectedError {
				t.Fatalf("expected error %v, got %v", tt.expectedError, fsm.err)
			}
		})
	}
}

func TestFSM_Next(t *testing.T) {
	tests := []struct {
		name          string
		rawQuery      string
		expectedState state
		expectedError error
		expectedArgs  []string
	}{
		{
			name:          "Valid SET command",
			rawQuery:      "SET key value",
			expectedState: FINISH,
			expectedError: nil,
			expectedArgs:  []string{"key", "value"},
		},
		{
			name:          "Invalid SET command (missing arguments)",
			rawQuery:      "SET key",
			expectedState: FINISH,
			expectedError: ErrQueryInvalidArgumentSyntaxSet,
		},
		{
			name:          "Valid GET command",
			rawQuery:      "GET key",
			expectedState: FINISH,
			expectedError: nil,
			expectedArgs:  []string{"key"},
		},
		{
			name:          "Invalid GET command (missing arguments)",
			rawQuery:      "GET",
			expectedState: FINISH,
			expectedError: ErrQueryInvalidArgumentSyntaxGet,
		},
		{
			name:          "Valid DEL command",
			rawQuery:      "DEL key",
			expectedState: FINISH,
			expectedError: nil,
			expectedArgs:  []string{"key"},
		},
		{
			name:          "Invalid DEL command (missing arguments)",
			rawQuery:      "DEL",
			expectedState: FINISH,
			expectedError: ErrQueryInvalidArgumentSyntaxDel,
		},
		{
			name:          "Unknown command",
			rawQuery:      "UNKNOWN key",
			expectedState: FINISH,
			expectedError: ErrQueryUnknowState,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := newFSM(tt.rawQuery)
			fsm.next()
			fsm.next()

			if fsm.state != tt.expectedState {
				t.Fatalf("expected state %v, got %v", tt.expectedState, fsm.state)
			}
			if fsm.err != tt.expectedError {
				t.Fatalf("expected error %v, got %v", tt.expectedError, fsm.err)
			}
			if tt.expectedArgs != nil && !equal(fsm.out.Arguments, tt.expectedArgs) {
				t.Fatalf("expected arguments %v, got %v", tt.expectedArgs, fsm.out.Arguments)
			}
		})
	}
}

func TestFSM_ParseCmd(t *testing.T) {
	tests := []struct {
		name          string
		rawQuery      string
		expectedState state
		expectedError error
	}{
		{"Valid SET command", "SET key value", CmdSetState, nil},
		{"Valid GET command", "GET key", CmdGetState, nil},
		{"Valid DEL command", "DEL key", CmdDelState, nil},
		{"Unknown command", "UNKNOWN key", FINISH, ErrQueryUnknowCommand},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := newFSM(tt.rawQuery)
			fsm.parseCmd()

			if fsm.state != tt.expectedState {
				t.Fatalf("expected state %v, got %v", tt.expectedState, fsm.state)
			}
			if fsm.err != tt.expectedError {
				t.Fatalf("expected error %v, got %v", tt.expectedError, fsm.err)
			}
		})
	}
}

func TestValidationFunctions(t *testing.T) {
	t.Run("isValidQuerySyntax", func(t *testing.T) {
		if !isValidQuerySyntax("SET key value") {
			t.Fatal("expected valid query syntax")
		}
		if isValidQuerySyntax("INVALID@QUERY") {
			t.Fatal("expected invalid query syntax")
		}
	})

	t.Run("isValidArgumentSetCmd", func(t *testing.T) {
		if !isValidArgumentSetCmd([]string{"SET", "key", "value"}) {
			t.Fatal("expected valid SET arguments")
		}
		if isValidArgumentSetCmd([]string{"SET", "key"}) {
			t.Fatal("expected invalid SET arguments")
		}
	})

	t.Run("isValidArgumentGetCmd", func(t *testing.T) {
		if !isValidArgumentGetCmd([]string{"GET", "key"}) {
			t.Fatal("expected valid GET arguments")
		}
		if isValidArgumentGetCmd([]string{"GET"}) {
			t.Fatal("expected invalid GET arguments")
		}
	})

	t.Run("isValidArgumentDelCmd", func(t *testing.T) {
		if !isValidArgumentDelCmd([]string{"DEL", "key"}) {
			t.Fatal("expected valid DEL arguments")
		}
		if isValidArgumentDelCmd([]string{"DEL"}) {
			t.Fatal("expected invalid DEL arguments")
		}
	})

	t.Run("formatQueryRaw", func(t *testing.T) {
		raw := "   SET    key    value   "
		formatted := formatQueryRaw(&raw)
		if *formatted != "SET key value" {
			t.Fatalf("expected 'SET key value', got %q", *formatted)
		}
	})
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
