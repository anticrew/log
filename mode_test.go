package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestSetPrintModeName(t *testing.T) {
	t.Parallel()

	type testCase struct {
		mode     PrintMode
		name     string
		expected string
	}

	tests := map[string]testCase{
		"set-custom-name": {
			mode:     PrintModePretty,
			name:     "PRETTIED",
			expected: "PRETTIED",
		},
		"override-existing": {
			mode:     PrintModeJson,
			name:     "JSONL",
			expected: "JSONL",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			origName := tc.mode.String()
			SetPrintModeName(tc.mode, tc.name)
			assert.Equal(t, tc.expected, tc.mode.String())

			// Восстанавливаем оригинальное имя
			SetPrintModeName(tc.mode, origName)
		})
	}
}

func TestPrintMode_String(t *testing.T) {
	t.Parallel()

	type testCase struct {
		mode     PrintMode
		expected string
	}

	tests := map[string]testCase{
		"pretty-mode": {
			mode:     PrintModePretty,
			expected: "PRETTY",
		},
		"json-mode": {
			mode:     PrintModeJson,
			expected: "JSON",
		},
		"text-mode": {
			mode:     PrintModeLogFmt,
			expected: "TEXT",
		},
		"unknown-mode": {
			mode:     PrintMode(42),
			expected: "PrintMode<42>",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, tc.mode.String())
		})
	}
}

func TestPrintMode_MarshalText(t *testing.T) {
	t.Parallel()

	type testCase struct {
		mode          PrintMode
		expectedText  string
		expectedError error
	}

	tests := map[string]testCase{
		"pretty-mode": {
			mode:          PrintModePretty,
			expectedText:  "PRETTY",
			expectedError: nil,
		},
		"custom-mode": {
			mode:          PrintMode(42),
			expectedText:  "PrintMode<42>",
			expectedError: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			text, err := tc.mode.MarshalText()
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedText, string(text))
		})
	}
}

func TestPrintMode_UnmarshalText(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input        string
		expectedMode PrintMode
		expectError  bool
	}

	tests := map[string]testCase{
		"valid-pretty": {
			input:        "PRETTY",
			expectedMode: PrintModePretty,
			expectError:  false,
		},
		"valid-json": {
			input:        "JSON",
			expectedMode: PrintModeJson,
			expectError:  false,
		},
		"valid-text": {
			input:        "TEXT",
			expectedMode: PrintModeLogFmt,
			expectError:  false,
		},
		"unknown-mode": {
			input:        "INVALID",
			expectedMode: PrintMode(0),
			expectError:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var mode PrintMode
			err := mode.UnmarshalText([]byte(tc.input))

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedMode, mode)
			}
		})
	}
}

func TestPrintMode_MarshalJSON(t *testing.T) {
	t.Parallel()

	type testCase struct {
		mode          PrintMode
		expectedJSON  string
		expectedError error
	}

	escapeHtmlJson := func(s string) string {
		var buf bytes.Buffer
		json.HTMLEscape(&buf, []byte(s))
		return buf.String()
	}

	tests := map[string]testCase{
		"pretty-mode": {
			mode:          PrintModePretty,
			expectedJSON:  `"PRETTY"`,
			expectedError: nil,
		},
		"custom-mode": {
			mode:          PrintMode(42),
			expectedJSON:  fmt.Sprintf(`"%s"`, escapeHtmlJson("PrintMode<42>")),
			expectedError: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			data, err := json.Marshal(tc.mode)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedJSON, string(data))
		})
	}
}

func TestPrintMode_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input        string
		directCall   bool
		expectedMode PrintMode
		expectError  bool
	}

	tests := map[string]testCase{
		"valid-pretty": {
			input:        `"PRETTY"`,
			expectedMode: PrintModePretty,
			expectError:  false,
		},
		"valid-json": {
			input:        `"JSON"`,
			expectedMode: PrintModeJson,
			expectError:  false,
		},
		"invalid-json": {
			input:        `invalid`,
			expectedMode: PrintMode(0),
			expectError:  true,
		},
		"unknown-mode": {
			input:        `"UNKNOWN"`,
			expectedMode: PrintMode(0),
			expectError:  true,
		},
		"bad-quotes": {
			input:       `"PRETTY`,
			directCall:  true,
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var (
				mode PrintMode
				err  error
			)
			if tc.directCall {
				err = mode.UnmarshalJSON([]byte(tc.input))
			} else {
				err = json.Unmarshal([]byte(tc.input), &mode)
			}

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedMode, mode)
			}
		})
	}
}
