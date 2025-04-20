package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestSetLevelName(t *testing.T) {
	t.Parallel()

	type testCase struct {
		level    Level
		name     string
		expected string
	}

	tests := map[string]testCase{
		"set-custom-name": {
			level:    LevelInfo,
			name:     "INFORMATION",
			expected: "INFORMATION",
		},
		"override-existing": {
			level:    LevelError,
			name:     "ERR",
			expected: "ERR",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			origName := test.level.String()
			SetLevelName(test.level, test.name)
			assert.Equal(t, test.expected, test.level.String())

			// Восстанавливаем оригинальное имя
			SetLevelName(test.level, origName)
		})
	}
}

func TestLevel_Level(t *testing.T) {
	t.Parallel()

	type testCase struct {
		level    Level
		expected Level
	}

	tests := map[string]testCase{
		"trace": {
			level:    LevelTrace,
			expected: LevelTrace,
		},
		"debug": {
			level:    LevelDebug,
			expected: LevelDebug,
		},
		"custom": {
			level:    Level(42),
			expected: Level(42),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, test.level.Level())
		})
	}
}

func TestLevel_String(t *testing.T) {
	t.Parallel()

	type testCase struct {
		level    Level
		expected string
	}

	tests := map[string]testCase{
		"known-level": {
			level:    LevelInfo,
			expected: "INFO",
		},
		"unknown-level": {
			level:    Level(42),
			expected: "Level<42>",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expected, test.level.String())
		})
	}
}

func TestLevel_MarshalText(t *testing.T) {
	t.Parallel()

	type testCase struct {
		level         Level
		expectedText  string
		expectedError error
	}

	tests := map[string]testCase{
		"standard-level": {
			level:         LevelWarn,
			expectedText:  "WARN",
			expectedError: nil,
		},
		"custom-level": {
			level:         Level(42),
			expectedText:  "Level<42>",
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			text, err := test.level.MarshalText()
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedText, string(text))
		})
	}
}

func TestLevel_UnmarshalText(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input         string
		expectedLevel Level
		expectError   bool
	}

	tests := map[string]testCase{
		"valid-level": {
			input:         "INFO",
			expectedLevel: LevelInfo,
			expectError:   false,
		},
		"unknown-level": {
			input:         "INVALID",
			expectedLevel: Level(0),
			expectError:   true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var level Level
			err := level.UnmarshalText([]byte(test.input))

			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedLevel, level)
			}
		})
	}
}

func TestLevel_MarshalJSON(t *testing.T) {
	t.Parallel()

	type testCase struct {
		level         Level
		expectedJSON  string
		expectedError error
	}

	escapeHtmlJson := func(s string) string {
		var buf bytes.Buffer
		json.HTMLEscape(&buf, []byte(s))
		return buf.String()
	}

	tests := map[string]testCase{
		"standard-level": {
			level:         LevelError,
			expectedJSON:  `"ERROR"`,
			expectedError: nil,
		},
		"custom-level": {
			level:         Level(42),
			expectedJSON:  fmt.Sprintf(`"%s"`, escapeHtmlJson("Level<42>")),
			expectedError: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			data, err := json.Marshal(test.level)
			assert.Equal(t, test.expectedError, err)
			assert.Equal(t, test.expectedJSON, string(data))
		})
	}
}

func TestLevel_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input         string
		directCall    bool
		expectedLevel Level
		expectError   bool
	}

	tests := map[string]testCase{
		"valid-level": {
			input:         `"DEBUG"`,
			expectedLevel: LevelDebug,
			expectError:   false,
		},
		"invalid-json": {
			input:         `invalid`,
			expectedLevel: Level(0),
			expectError:   true,
		},
		"unknown-level": {
			input:         `"UNKNOWN"`,
			expectedLevel: Level(0),
			expectError:   true,
		},
		"bad-quotes": {
			input:       `"ERROR`,
			directCall:  true,
			expectError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var (
				level Level
				err   error
			)
			if test.directCall {
				err = level.UnmarshalJSON([]byte(test.input))
			} else {
				err = json.Unmarshal([]byte(test.input), &level)
			}

			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedLevel, level)
			}
		})
	}
}
