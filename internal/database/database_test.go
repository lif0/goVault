package database

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"goVault/internal"
	"goVault/internal/core/query"
	"goVault/internal/core/vault"
	internal_mock "goVault/mocks"
	query_mock "goVault/mocks/core/query"
	vault_mock "goVault/mocks/core/vault"

	"go.uber.org/mock/gomock"
)

func TestNewDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := internal_mock.NewMockLogger(ctrl)
	parserMock := query_mock.NewMockParser(ctrl)
	vaultMock := vault_mock.NewMockVault(ctrl)

	tests := []struct {
		name          string
		parserLayer   parserLayer
		vaultLayer    vaultLayer
		logger        internal.Logger
		expectedError string
	}{
		{"Valid inputs", parserMock, vaultMock, logger, ""},
		{"Nil parserLayer", nil, vaultMock, logger, "parser is invalid"},
		{"Nil vaultLayer", parserMock, nil, logger, "vault is invalid"},
		{"Nil logger", parserMock, vaultMock, nil, "logger is invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := NewDatabase(tt.parserLayer, tt.vaultLayer, tt.logger)
			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Fatalf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if db == nil {
				t.Fatal("expected database instance, got nil")
			}
		})
	}
}

func TestDatabase_HandleQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := internal_mock.NewMockLogger(ctrl)
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	parserMock := query_mock.NewMockParser(ctrl)
	vaultMock := vault_mock.NewMockVault(ctrl)
	db, _ := NewDatabase(parserMock, vaultMock, logger)

	tests := []struct {
		name           string
		queryStr       string
		commandID      query.DBCommand
		arguments      []string
		parserError    error
		vaultError     error
		expectedResult string
	}{
		{
			name:           "Handle SET query successfully",
			queryStr:       "SET key value",
			commandID:      query.SET,
			arguments:      []string{"key", "value"},
			expectedResult: MsgDBOk,
		},
		{
			name:           "Handle GET query successfully",
			queryStr:       "GET key",
			commandID:      query.GET,
			arguments:      []string{"key"},
			expectedResult: MsgDBOk + " value",
		},
		{
			name:           "Handle DEL query successfully",
			queryStr:       "DEL key",
			commandID:      query.DEL,
			arguments:      []string{"key"},
			expectedResult: MsgDBOk,
		},
		{
			name:           "Handle parser error",
			queryStr:       "INVALID QUERY",
			parserError:    errors.New("parser error"),
			expectedResult: "[error] parser error",
		},
		{
			name:           "Handle GET with not found error",
			queryStr:       "GET keyHandle_GET_with_not_found_error",
			commandID:      query.GET,
			arguments:      []string{"keyHandle_GET_with_not_found_error"},
			vaultError:     vault.ErrVaultNotFound,
			expectedResult: "[not found]",
		},
		{
			name:           "Handle DEL with vault error",
			queryStr:       "DEL keyHandle_DEL_with_vault_error",
			commandID:      query.DEL,
			arguments:      []string{"keyHandle_DEL_with_vault_error"},
			vaultError:     errors.New("vault error"),
			expectedResult: "[error] vault error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			parserMock.EXPECT().Transition(tt.queryStr).Return(&query.Query{
				CommandID: tt.commandID,
				Arguments: tt.arguments,
			}, tt.parserError).AnyTimes()

			switch tt.commandID {
			case query.SET:
				vaultMock.EXPECT().Set(ctx, tt.arguments[0], tt.arguments[1]).Return(tt.vaultError).AnyTimes()
			case query.GET:
				vaultMock.EXPECT().Get(ctx, tt.arguments[0]).Return("value", tt.vaultError).AnyTimes()
			case query.DEL:
				vaultMock.EXPECT().Del(ctx, tt.arguments[0]).Return(tt.vaultError).AnyTimes()
			}

			result := db.HandleQuery(ctx, tt.queryStr)
			if result != tt.expectedResult {
				t.Fatalf("expected %q, got %q", tt.expectedResult, result)
			}
		})
	}
}

func TestDatabase_HandleQuery_ErrorCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := internal_mock.NewMockLogger(ctrl)
	logger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()
	parser := query_mock.NewMockParser(ctrl)
	vaultLayer := vault_mock.NewMockVault(ctrl)
	db, _ := NewDatabase(parser, vaultLayer, logger)

	tests := []struct {
		name          string
		queryStr      string
		parserError   error
		commandID     query.DBCommand
		vaultError    error
		expectedError string
	}{
		{
			name:          "Unknown Command ID",
			queryStr:      "UNKNOWN key",
			commandID:     query.DBCommand(128),
			expectedError: "compute layer is incorrect: command_id 128",
		},
		{
			name:          "Parser Error",
			queryStr:      "INVALID QUERY",
			parserError:   errors.New("parser error"),
			expectedError: "[error] parser error",
		},
		{
			name:          "Vault SET Error",
			queryStr:      "SET key value",
			commandID:     query.SET,
			vaultError:    errors.New("vault set error"),
			expectedError: "[error] vault set error",
		},
		{
			name:          "Vault GET Error",
			queryStr:      "GET key",
			commandID:     query.GET,
			vaultError:    errors.New("vault get error"),
			expectedError: "[error] vault get error",
		},
		{
			name:          "Vault DEL Error",
			queryStr:      "DEL key",
			commandID:     query.DEL,
			vaultError:    errors.New("vault del error"),
			expectedError: "[error] vault del error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Mock parser behavior
			parser.EXPECT().Transition(tt.queryStr).Return(&query.Query{
				CommandID: tt.commandID,
				Arguments: []string{"key", "value"},
			}, tt.parserError).AnyTimes()

			// Mock vault layer behavior
			switch tt.commandID {
			case query.SET:
				vaultLayer.EXPECT().Set(ctx, "key", "value").Return(tt.vaultError).AnyTimes()
			case query.GET:
				vaultLayer.EXPECT().Get(ctx, "key").Return("", tt.vaultError).AnyTimes()
			case query.DEL:
				vaultLayer.EXPECT().Del(ctx, "key").Return(tt.vaultError).AnyTimes()
			}

			// Mock logger behavior
			if tt.expectedError == fmt.Sprintf("compute layer is incorrect: command_id %v", tt.commandID) {
				logger.EXPECT().Error(tt.expectedError).Times(1)
			} else if tt.expectedError != "" {
				logger.EXPECT().Error(gomock.Any()).AnyTimes()
			}

			// Execute the query
			result := db.HandleQuery(ctx, tt.queryStr)

			// Validate the result
			if result != tt.expectedError {
				t.Fatalf("expected result %q, got %q", tt.expectedError, result)
			}
		})
	}
}
