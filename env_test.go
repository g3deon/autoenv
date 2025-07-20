package autoenv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnvFile(t *testing.T) {
	testCases := []struct {
		name        string
		setupFile   func() (string, error)
		expectError bool
	}{
		{
			name: "valid env file",
			setupFile: func() (string, error) {
				return createTempEnvFile("KEY1=value1\nKEY2=value2")
			},
			expectError: false,
		},
		{
			name: "non-existent file",
			setupFile: func() (string, error) {
				return "non_existent_file.env", nil
			},
			expectError: false,
		},
		{
			name: "empty file",
			setupFile: func() (string, error) {
				return createTempEnvFile("")
			},
			expectError: false,
		},
		{
			name: "file with comments and empty lines",
			setupFile: func() (string, error) {
				return createTempEnvFile("# This is a comment\n\nKEY3=value3\n# Another comment\nKEY4=value4")
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filePath, err := tc.setupFile()
			if err != nil {
				t.Fatalf("Failed to setup test file: %v", err)
			}

			if filepath.Base(filePath) != "non_existent_file.env" {
				defer os.Remove(filePath)
			}

			err = loadEnvFile(filePath)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
				return
			}

			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
				return
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	testCases := []struct {
		name         string
		setupFile    func() (string, error)
		expectExists bool
		expectError  bool
	}{
		{
			name: "existing file",
			setupFile: func() (string, error) {
				return createTempEnvFile("test content")
			},
			expectExists: true,
			expectError:  false,
		},
		{
			name: "non-existent file",
			setupFile: func() (string, error) {
				return "definitely_does_not_exist.txt", nil
			},
			expectExists: false,
			expectError:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filePath, err := tc.setupFile()
			if err != nil {
				t.Fatalf("Failed to setup test: %v", err)
			}

			if tc.expectExists {
				defer os.Remove(filePath)
			}

			exists, err := fileExists(filePath)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
				return
			}

			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
				return
			}

			if exists != tc.expectExists {
				t.Errorf("Expected exists=%v but got %v", tc.expectExists, exists)
				return
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	testCases := []struct {
		name        string
		content     string
		expectError bool
	}{
		{
			name:        "valid env file",
			content:     "KEY1=value1\nKEY2=value2",
			expectError: false,
		},
		{
			name:        "empty file",
			content:     "",
			expectError: false,
		},
		{
			name:        "file with comments",
			content:     "# Comment\nKEY=value\n# Another comment",
			expectError: false,
		},
		{
			name:        "file with export statements",
			content:     "export KEY1=value1\nexport KEY2=value2",
			expectError: false,
		},
		{
			name:        "file with quoted values",
			content:     `KEY1="quoted value"\nKEY2='single quoted'`,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filePath, err := createTempEnvFile(tc.content)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(filePath)

			err = parseFile(filePath)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
				return
			}

			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
				return
			}
		})
	}
}

func TestProcessLine(t *testing.T) {
	testCases := []struct {
		name          string
		line          string
		expectedKey   string
		expectedValue string
		shouldSetEnv  bool
	}{
		{
			name:          "simple key-value pair",
			line:          "KEY=value",
			expectedKey:   "KEY",
			expectedValue: "value",
			shouldSetEnv:  true,
		},
		{
			name:          "key-value with spaces",
			line:          "  KEY  =  value  ",
			expectedKey:   "KEY",
			expectedValue: "value",
			shouldSetEnv:  true,
		},
		{
			name:          "export statement",
			line:          "export KEY=value",
			expectedKey:   "KEY",
			expectedValue: "value",
			shouldSetEnv:  true,
		},
		{
			name:          "quoted value with double quotes",
			line:          `KEY="quoted value"`,
			expectedKey:   "KEY",
			expectedValue: "quoted value",
			shouldSetEnv:  true,
		},
		{
			name:          "quoted value with single quotes",
			line:          `KEY='single quoted'`,
			expectedKey:   "KEY",
			expectedValue: "single quoted",
			shouldSetEnv:  true,
		},
		{
			name:         "comment line",
			line:         "# This is a comment",
			shouldSetEnv: false,
		},
		{
			name:         "empty line",
			line:         "",
			shouldSetEnv: false,
		},
		{
			name:         "whitespace only line",
			line:         "   ",
			shouldSetEnv: false,
		},
		{
			name:         "invalid format - no equals",
			line:         "INVALID_LINE",
			shouldSetEnv: false,
		},
		{
			name:          "value with equals sign",
			line:          "KEY=value=with=equals",
			expectedKey:   "KEY",
			expectedValue: "value=with=equals",
			shouldSetEnv:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.shouldSetEnv {
				return
			}

			originalValue := os.Getenv(tc.expectedKey)
			defer func() {
				if originalValue == "" {
					os.Unsetenv(tc.expectedKey)
					return
				}

				os.Setenv(tc.expectedKey, originalValue)
			}()

			os.Unsetenv(tc.expectedKey)
			processLine(tc.line)

			actualValue := os.Getenv(tc.expectedKey)
			if actualValue != tc.expectedValue {
				t.Errorf("Expected env var %s=%s but got %s", tc.expectedKey, tc.expectedValue, actualValue)
			}
		})
	}
}

func createTempEnvFile(content string) (string, error) {
	tmpFile, err := os.CreateTemp("", "test_env_*.env")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if content != "" {
		if _, err := tmpFile.WriteString(content); err != nil {
			os.Remove(tmpFile.Name())
			return "", err
		}
	}

	return tmpFile.Name(), nil
}
