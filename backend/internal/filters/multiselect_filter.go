package filters

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
)

type MultiSelectInput struct {
	Values    []string `json:"values"`     // Array of selected values
	Field     string   `json:"field"`      // Field name to filter on
	Inclusive bool     `json:"inclusive"`  // Whether to include (IN) or exclude (NOT IN) values
	MaxValues int      `json:"max_values"` // Maximum number of values allowed (0 for unlimited)
}

type MultiSelectOutput struct {
	Values      []string `json:"values"`       // Sanitized, deduplicated values
	QueryString string   `json:"query_string"` // Quickwit compatible query string
}

// NewMultiSelectFilter creates a new multi-select filter implementation
func NewMultiSelectFilter(logger *otelzap.Logger) FilterFunc {
	logger = logger.Named("multi_select_filter")

	return func(input interface{}) (interface{}, error) {
		// Type assertion for input
		multiSelect, ok := input.(MultiSelectInput)
		if !ok {
			logger.Error("invalid input type",
				zap.String("expected_type", "MultiSelectInput"),
				zap.Any("received_input", input))
			return nil, fmt.Errorf("invalid input type for multi-select filter")
		}

		// Validate field name

		if !ValidateMultiselectField(multiSelect.Field) {
			logger.Error("invalid field name",
				zap.String("field", multiSelect.Field))
			return nil, fmt.Errorf("invalid field name")
		}

		// Handle empty values array
		if len(multiSelect.Values) == 0 {
			logger.Debug("empty values array received")
			return MultiSelectOutput{
				Values:      []string{},
				QueryString: "*:*", // Match all in Quickwit
			}, nil
		}

		// Enforce max values limit if specified
		if multiSelect.MaxValues > 0 && len(multiSelect.Values) > multiSelect.MaxValues {
			logger.Error("too many values",
				zap.Int("max_allowed", multiSelect.MaxValues),
				zap.Int("received", len(multiSelect.Values)))
			return nil, fmt.Errorf("too many values: maximum allowed is %d", multiSelect.MaxValues)
		}

		// Sanitize and deduplicate values
		sanitizedValues := make([]string, 0, len(multiSelect.Values))
		seenValues := make(map[string]bool)

		for _, value := range multiSelect.Values {
			// Basic sanitization
			sanitized := sanitizeValue(value)
			if sanitized == "" {
				continue
			}

			// Deduplicate
			if !seenValues[sanitized] {
				seenValues[sanitized] = true
				sanitizedValues = append(sanitizedValues, sanitized)
			}
		}

		// Build Quickwit query string
		var queryString string
		if multiSelect.Inclusive {
			queryString = fmt.Sprintf("%s:(%s)",
				multiSelect.Field,
				strings.Join(sanitizedValues, " OR "))
		} else {
			queryString = fmt.Sprintf("NOT %s:(%s)",
				multiSelect.Field,
				strings.Join(sanitizedValues, " OR "))
		}

		output := MultiSelectOutput{
			Values:      sanitizedValues,
			QueryString: queryString,
		}

		return output, nil
	}
}

// isValidFieldName checks if the field name contains only allowed characters
func ValidateMultiselectField(field string) bool {
	// Allow only alphanumeric characters, underscores, and dots
	for _, char := range field {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' ||
			char == '.') {
			return false
		}
	}
	return len(field) > 0
}

// sanitizeValue performs basic sanitization on input values
func sanitizeValue(value string) string {
	// Remove any single quotes (SQL injection prevention)
	value = strings.ReplaceAll(value, "'", "")
	// Trim whitespace
	value = strings.TrimSpace(value)
	// Basic length validation
	if len(value) > 255 {
		return value[:255]
	}
	return value
}
