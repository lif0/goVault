package in_memory

import (
	"testing"
)

func TestHashTable_SetAndGet(t *testing.T) {
	tests := []struct {
		name       string
		operations []struct {
			op       string
			key      string
			value    string
			expected string
			found    bool
		}
	}{
		{
			name: "Set and Get single value",
			operations: []struct {
				op       string
				key      string
				value    string
				expected string
				found    bool
			}{
				{op: "set", key: "key1", value: "value1"},
				{op: "get", key: "key1", expected: "value1", found: true},
			},
		},
		{
			name: "Set multiple values and Get them",
			operations: []struct {
				op       string
				key      string
				value    string
				expected string
				found    bool
			}{
				{op: "set", key: "key1", value: "value1"},
				{op: "set", key: "key2", value: "value2"},
				{op: "get", key: "key1", expected: "value1", found: true},
				{op: "get", key: "key2", expected: "value2", found: true},
			},
		},
		{
			name: "Get non-existent key",
			operations: []struct {
				op       string
				key      string
				value    string
				expected string
				found    bool
			}{
				{op: "get", key: "nonexistent", expected: "", found: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewHashTable()

			for _, op := range tt.operations {
				switch op.op {
				case "set":
					table.Set(op.key, op.value)
				case "get":
					value, found := table.Get(op.key)
					if value != op.expected || found != op.found {
						t.Errorf("Expected value: %v, found: %v, got value: %v, found: %v",
							op.expected, op.found, value, found)
					}
				}
			}
		})
	}
}

func TestHashTable_Delete(t *testing.T) {
	tests := []struct {
		name       string
		operations []struct {
			op       string
			key      string
			value    string
			expected string
			found    bool
		}
	}{
		{
			name: "Delete existing key",
			operations: []struct {
				op       string
				key      string
				value    string
				expected string
				found    bool
			}{
				{op: "set", key: "key1", value: "value1"},
				{op: "del", key: "key1"},
				{op: "get", key: "key1", expected: "", found: false},
			},
		},
		{
			name: "Delete non-existent key",
			operations: []struct {
				op       string
				key      string
				value    string
				expected string
				found    bool
			}{
				{op: "del", key: "key1"},
				{op: "get", key: "key1", expected: "", found: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewHashTable()

			for _, op := range tt.operations {
				switch op.op {
				case "set":
					table.Set(op.key, op.value)
				case "del":
					table.Del(op.key)
				case "get":
					value, found := table.Get(op.key)
					if value != op.expected || found != op.found {
						t.Errorf("Expected value: %v, found: %v, got value: %v, found: %v",
							op.expected, op.found, value, found)
					}
				}
			}
		})
	}
}
