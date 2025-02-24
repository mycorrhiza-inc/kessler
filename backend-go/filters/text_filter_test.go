package filters

import (
	"regexp"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestTextFilter(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	filter := NewTextFilter(logger)

	tests := []struct {
		name     string
		input    interface{}
		expected string
		wantErr  bool
	}{
		{
			name:     "Basic text query",
			input:    "field:value",
			expected: "field:value",
			wantErr:  false,
		},
		{
			name:     "Query with special characters",
			input:    "field:value+test",
			expected: `field:value\+test`,
			wantErr:  false,
		},
		{
			name:     "Phrase query",
			input:    `field:"hello world"`,
			expected: `field:"hello world"`,
			wantErr:  false,
		},
		{
			name:     "SQL injection attempt",
			input:    "field:value; DROP TABLE users;",
			expected: `field:value\;\ DROP\ TABLE\ users\;`,
			wantErr:  false,
		},
		{
			name:     "Empty query",
			input:    "",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Non-string input",
			input:    123,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Query exceeding max length",
			input:    strings.Repeat("a", 2000),
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Complex query with operators",
			input:    `(field1:"test phrase" AND field2:value) OR field3:*`,
			expected: `\(field1:"test phrase"\ AND\ field2:value\)\ OR\ field3:\*`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := filter(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFieldNameValidation(t *testing.T) {
	tf := &TextFilter{
		fieldNameRegex: regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_\.]*$`),
	}

	tests := []struct {
		name      string
		fieldName string
		valid     bool
	}{
		{"Valid simple field", "field", true},
		{"Valid nested field", "user.name", true},
		{"Valid with underscore", "user_name", true},
		{"Invalid starts with number", "1field", false},
		{"Invalid special chars", "field@name", false},
		{"Invalid SQL injection", "field;DROP", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tf.validateFieldName(tt.fieldName); got != tt.valid {
				t.Errorf("validateFieldName(%q) = %v, want %v", tt.fieldName, got, tt.valid)
			}
		})
	}
}
