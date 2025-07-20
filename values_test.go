package autoenv

import (
	"reflect"
	"testing"
)

func TestAssignValuesToStruct(t *testing.T) {
	type SimpleStruct struct {
		Name string
		Age  int
	}

	type MixedStruct struct {
		Text   string
		Number int
		Flag   bool
		Score  float64
	}

	type PrivateStruct struct {
		Public  string
		private int
		Another bool
	}

	tests := []struct {
		name     string
		input    any
		values   []reflect.Value
		expected any
	}{
		{
			name:  "basic_assignment",
			input: &SimpleStruct{},
			values: []reflect.Value{
				reflect.ValueOf("John"),
				reflect.ValueOf(25),
			},
			expected: &SimpleStruct{Name: "John", Age: 25},
		},
		{
			name:  "partial_assignment",
			input: &SimpleStruct{},
			values: []reflect.Value{
				reflect.ValueOf("Alice"),
			},
			expected: &SimpleStruct{Name: "Alice", Age: 0},
		},
		{
			name:  "type_mismatch",
			input: &SimpleStruct{},
			values: []reflect.Value{
				reflect.ValueOf("Bob"),
				reflect.ValueOf("not_an_int"),
			},
			expected: &SimpleStruct{Name: "Bob", Age: 0},
		},
		{
			name:  "invalid_value",
			input: &SimpleStruct{},
			values: []reflect.Value{
				reflect.ValueOf("Valid"),
				{},
			},
			expected: &SimpleStruct{Name: "Valid", Age: 0},
		},
		{
			name:     "empty_values",
			input:    &SimpleStruct{Name: "Original", Age: 42},
			values:   []reflect.Value{},
			expected: &SimpleStruct{Name: "Original", Age: 42},
		},
		{
			name:  "more_values_than_fields",
			input: &SimpleStruct{},
			values: []reflect.Value{
				reflect.ValueOf("First"),
				reflect.ValueOf(100),
				reflect.ValueOf("Extra1"),
				reflect.ValueOf("Extra2"),
			},
			expected: &SimpleStruct{Name: "First", Age: 100},
		},
		{
			name:  "mixed_types_struct",
			input: &MixedStruct{},
			values: []reflect.Value{
				reflect.ValueOf("Hello"),
				reflect.ValueOf(42),
				reflect.ValueOf(true),
				reflect.ValueOf(99.9),
			},
			expected: &MixedStruct{Text: "Hello", Number: 42, Flag: true, Score: 99.9},
		},
		{
			name:  "struct_with_unexported_fields",
			input: &PrivateStruct{},
			values: []reflect.Value{
				reflect.ValueOf("PublicValue"),
				reflect.ValueOf(123),
				reflect.ValueOf(true),
			},
			expected: &PrivateStruct{Public: "PublicValue", private: 0, Another: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var structValue reflect.Value

			inputValue := reflect.ValueOf(tt.input)
			if inputValue.Kind() == reflect.Ptr {
				newStruct := reflect.New(inputValue.Elem().Type())
				newStruct.Elem().Set(inputValue.Elem())

				structValue = newStruct
				assignValuesToStruct(structValue, tt.values...)

				result := structValue.Interface()
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("assignValuesToStruct() = %+v; want %+v", result, tt.expected)
				}

				return
			}

			newStruct := reflect.New(inputValue.Type())
			newStruct.Elem().Set(inputValue)

			structValue = newStruct.Elem()
			assignValuesToStruct(structValue, tt.values...)

			result := structValue.Interface()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("assignValuesToStruct() = %+v; want %+v", result, tt.expected)
			}
		})
	}
}
