package autoenv

import "testing"

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "HELLO_WORLD"},
		{"snakeCase", "SNAKE_CASE"},
		{"test123", "TEST123"},
		{"", ""},
		{"singleword", "SINGLEWORD"},
	}

	for _, test := range tests {
		result := toSnakeCase(test.input)
		if result != test.expected {
			t.Errorf("toSnakeCase(%q) = %q; want %q", test.input, result, test.expected)
		}
	}
}
