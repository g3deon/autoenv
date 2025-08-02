package autoenv

import (
	"bytes"
	"os"
	"testing"
)

func TestLoadEnvFile(t *testing.T) {
	tests := []struct {
		name         string
		fileContents string
		expectedEnv  map[string]string
		expectError  bool
	}{
		{
			name:         "EmptyFile",
			fileContents: "",
			expectedEnv:  map[string]string{},
			expectError:  false,
		},
		{
			name: "ValidFileContents",
			fileContents: `
			VAR1=value1
			VAR2=value2
			`,
			expectedEnv: map[string]string{
				"VAR1": "value1",
				"VAR2": "value2",
			},
			expectError: false,
		},
		{
			name: "FileWithMalformedLines",
			fileContents: `
			VAR1=value1
			malformed_line
			VAR2=value2
			`,
			expectedEnv: map[string]string{
				"VAR1": "value1",
				"VAR2": "value2",
			},
			expectError: false,
		},
		{
			name: "FileWithExportPrefix",
			fileContents: `
			export VAR1="value1"
			export VAR2='value2'
			`,
			expectedEnv: map[string]string{
				"VAR1": "value1",
				"VAR2": "value2",
			},
			expectError: false,
		},
		{
			name: "FileWithComments",
			fileContents: `
			# This is a comment
			VAR1=value1   # inline comment
			VAR2=value2
			`,
			expectedEnv: map[string]string{
				"VAR1": "value1",
				"VAR2": "value2",
			},
			expectError: false,
		},
		{
			name: "FileWithSpacesInValues",
			fileContents: `
			VAR1=" value with spaces "
			VAR2='another value'
			`,
			expectedEnv: map[string]string{
				"VAR1": " value with spaces ",
				"VAR2": "another value",
			},
			expectError: false,
		},
		{
			name: "FileWithEmptyLines",
			fileContents: `

			VAR1=value1

			VAR2=value2

			`,
			expectedEnv: map[string]string{
				"VAR1": "value1",
				"VAR2": "value2",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", ".env")
			if err != nil {
				t.Fatalf("loadEnvFile() = failed to create temp file: %v", err)
			}
			defer func(name string) {
				_ = os.Remove(name)
			}(tmpFile.Name())

			if tt.fileContents != "" {
				if _, err := tmpFile.WriteString(tt.fileContents); err != nil {
					t.Fatalf("loadEnvFile() = failed to write to temp file: %v", err)
				}
			}
			if err := tmpFile.Close(); err != nil {
				t.Fatalf("loadEnvFile() = failed to close temp file: %v", err)
			}

			loader := &Loader{}
			err = loader.loadEnvFile(tmpFile.Name())

			if tt.expectError {
				if err == nil {
					t.Errorf("loadEnvFile() = got no error, want error")
				}
				return
			}
			if err != nil {
				t.Errorf("loadEnvFile() = got unexpected error: %v", err)
				return
			}

			for key, want := range tt.expectedEnv {
				if got := os.Getenv(key); got != want {
					t.Errorf("loadEnvFile() = key, %q got %q, want %q", key, got, want)
				}
			}
		})
	}
}

func TestTrimSpaces(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name:     "EmptyInput",
			input:    []byte(""),
			expected: nil,
		},
		{
			name:     "NoSpaces",
			input:    []byte("test"),
			expected: []byte("test"),
		},
		{
			name:     "LeadingSpaces",
			input:    []byte("   test"),
			expected: []byte("test"),
		},
		{
			name:     "TrailingSpaces",
			input:    []byte("test   "),
			expected: []byte("test"),
		},
		{
			name:     "LeadingAndTrailingSpaces",
			input:    []byte("   test   "),
			expected: []byte("test"),
		},
		{
			name:     "OnlySpaces",
			input:    []byte("       "),
			expected: nil,
		},
		{
			name:     "TabsOnly",
			input:    []byte("\t\t\t"),
			expected: nil,
		},
		{
			name:     "TabsAndText",
			input:    []byte("\ttest\t"),
			expected: []byte("test"),
		},
		{
			name:     "SpacesAndTabs",
			input:    []byte("  \t test \t  "),
			expected: []byte("test"),
		},
		{
			name:     "SpecialCharacters",
			input:    []byte("  !@#$%^&*()  "),
			expected: []byte("!@#$%^&*()"),
		},
		{
			name:     "SingleCharacter",
			input:    []byte("   a   "),
			expected: []byte("a"),
		},
		{
			name:     "MixedSpaces",
			input:    []byte(" \t\t a  \t"),
			expected: []byte("a"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := trimSpaces(tt.input)
			want := tt.expected
			if !bytes.Equal(got, want) {
				t.Errorf("trimSpaces() = got %q, want %q", got, want)
			}
		})
	}
}
