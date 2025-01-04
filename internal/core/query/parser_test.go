package query

import (
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	"goVault/internal"
	internal_mock "goVault/mocks"
)

func TestNewParser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := internal_mock.NewMockLogger(ctrl)

	tests := []struct {
		name          string
		logger        internal.Logger
		expectedError string
	}{
		{"Valid logger", logger, ""},
		{"Nil logger", nil, "logger is invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.logger)

			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Fatalf("expected error %q, got %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if parser == nil {
				t.Fatal("expected parser instance, got nil")
			}
		})
	}
}

func TestFMSMachine_Transition(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := internal_mock.NewMockLogger(ctrl)

	tests := []struct {
		name          string
		query         string
		expectedError string
		expectedQuery *Query
	}{
		{
			name:          "Valid query",
			query:         "SET key value",
			expectedError: "",
			expectedQuery: &Query{CommandID: SET, Arguments: []string{"key", "value"}},
		},
		{
			name:          "Invalid query",
			query:         "INVALID QUERY",
			expectedError: "query: unknow db command",
			expectedQuery: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(logger)
			if err != nil {
				t.Fatalf("unexpected error while creating parser: %v", err)
			}

			out, err := parser.Transition(tt.query)

			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Fatalf("expected error %q, got %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if out == nil || !reflect.DeepEqual(*out, *tt.expectedQuery) {
				t.Fatalf("expected query %v, got %v", tt.expectedQuery, out)
			}
		})
	}
}
