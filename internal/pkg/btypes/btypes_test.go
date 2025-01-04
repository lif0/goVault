package btypes

import (
	"testing"
)

func TestToArray(t *testing.T) {
	// Определяем таблицу тестов
	tests := []struct {
		name     string
		input    any
		expected any
	}{
		{
			name:     "int",
			input:    42,
			expected: []int{42},
		},
		{
			name:     "string",
			input:    "hello",
			expected: []string{"hello"},
		},
		{
			name: "struct",
			input: struct {
				Field string
			}{Field: "value"},
			expected: []struct {
				Field string
			}{
				{Field: "value"},
			},
		},
	}

	// Выполняем тесты
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToArray(tt.input)

			if len(result) != 1 {
				t.Errorf("Expected array length of 1, got %d", len(result))
			}

			if result[0] != tt.input {
				t.Errorf("Expected array element to be %+v, got %+v", tt.input, result[0])
			}
		})
	}
}
