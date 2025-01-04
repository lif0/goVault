package pointer

import (
	"testing"
	"time"
)

func TestTo_Int(t *testing.T) {
	input := 42
	result := To(input)
	if *result != input {
		t.Errorf("Expected %v, got %v", input, *result)
	}
}

func TestTo_String(t *testing.T) {
	input := "hello"
	result := To(input)
	if *result != input {
		t.Errorf("Expected %v, got %v", input, *result)
	}
}

func TestTo_Struct(t *testing.T) {
	input := struct {
		Field string
	}{Field: "value"}
	result := To(input)
	if *result != input {
		t.Errorf("Expected %v, got %v", input, *result)
	}
}

// Тесты для функции ToOrNil
func TestToOrNil_NonZeroInt(t *testing.T) {
	input := 42
	result := ToOrNil(input)
	if result == nil || *result != input {
		t.Errorf("Expected pointer to %v, got %v", input, result)
	}
}

func TestToOrNil_ZeroInt(t *testing.T) {
	input := 0
	result := ToOrNil(input)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestToOrNil_NonZeroString(t *testing.T) {
	input := "hello"
	result := ToOrNil(input)
	if result == nil || *result != input {
		t.Errorf("Expected pointer to %v, got %v", input, result)
	}
}

func TestToOrNil_ZeroString(t *testing.T) {
	input := ""
	result := ToOrNil(input)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestToOrNil_StructWithIsZero(t *testing.T) {
	input := time.Time{}
	result := ToOrNil(input)
	if result != nil {
		t.Errorf("Expected nil for zero time.Time, got %v", result)
	}
}

func TestToOrNil_NonZeroStruct(t *testing.T) {
	input := struct {
		Field string
	}{Field: "value"}
	result := ToOrNil(input)
	if result == nil || *result != input {
		t.Errorf("Expected pointer to %v, got %v", input, result)
	}
}

func TestToOrNil_NonZeroCustomStruct(t *testing.T) {
	type CustomStruct struct {
		Field string
	}

	input := CustomStruct{Field: "value"}
	result := ToOrNil(input)

	if result == nil || *result != input {
		t.Errorf("Expected pointer to %v, got %v", input, result)
	}
}

type CustomStructWithIsZero struct {
	Field string
}

func (cs CustomStructWithIsZero) IsZero() bool {
	return cs.Field == "" // Считаем значение "нулевым", если Field пустой
}
func TestToOrNil_StructWithNonZeroIsZero(t *testing.T) {

	input := CustomStructWithIsZero{Field: "value"}
	result := ToOrNil(input)

	if result == nil || *result != input {
		t.Errorf("Expected pointer to %v, got %v", input, result)
	}
}

func TestValueOf_NonNilPointer(t *testing.T) {
	input := To(42)
	result := ValueOf(input)
	if result != 42 {
		t.Errorf("Expected %v, got %v", 42, result)
	}
}

func TestValueOf_NilPointer(t *testing.T) {
	var input *int
	result := ValueOf(input)
	if result != 0 {
		t.Errorf("Expected zero value (0), got %v", result)
	}
}

func TestValueOf_NonNilStringPointer(t *testing.T) {
	input := To("hello")
	result := ValueOf(input)
	if result != "hello" {
		t.Errorf("Expected %v, got %v", "hello", result)
	}
}

func TestValueOf_NilStringPointer(t *testing.T) {
	var input *string
	result := ValueOf(input)
	if result != "" {
		t.Errorf("Expected zero value (empty string), got %v", result)
	}
}
