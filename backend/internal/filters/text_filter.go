package filters

import (
	"context"
	"fmt"
	"kessler/internal/quickwit"
	"kessler/pkg/logger"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

const (
	maxQueryLength = 1024 // Maximum length of a query string
	// Special characters that need escaping in non-quoted terms
	specialChars = `+,^:{}\"[]()~!\\*\s`
)

type TextFilter struct {
	logger         *otelzap.Logger
	IndexName      string
	fieldNameRegex regexp.Regexp
}

type QueryFilter interface{}

func NewTextFilter() FilterFunc {
	tf := &TextFilter{
		logger: logger.GetLogger("text_filter"),
	}
	return tf.Apply
}

// func (tf *TextFilter) Apply(input QueryFilter) (interface{}, error) {
func (tf *TextFilter) Apply(input interface{}) (interface{}, error) {
	query, ok := input.(string)
	if !ok {
		tf.logger.Error("invalid input type", zap.Any("input", input))
		return nil, fmt.Errorf("input must be a string")
	}

	// Validate input
	if err := tf.validateInput(query); err != nil {
		tf.logger.Error("input validation failed",
			zap.String("query", query),
			zap.Error(err))
		return nil, err
	}

	// Process the query
	processed, err := tf.processQuery(query)
	if err != nil {
		tf.logger.Error("query processing failed",
			zap.String("query", query),
			zap.Error(err))
		return nil, err
	}

	tf.logger.Debug("query processed successfully",
		zap.String("original_query", query),
		zap.String("processed_query", processed))
	return processed, nil
}

func (tf *TextFilter) validateInput(query string) error {
	if query == "" {
		tf.logger.Warn("empty query received")
		return fmt.Errorf("query cannot be empty")
	}

	if !utf8.ValidString(query) {
		tf.logger.Error("invalid UTF-8 characters in query",
			zap.String("query", query))
		return fmt.Errorf("query contains invalid UTF-8 characters")
	}

	if utf8.RuneCountInString(query) > maxQueryLength {
		tf.logger.Error("query exceeds maximum length",
			zap.String("query", query),
			zap.Int("max_length", maxQueryLength),
			zap.Int("actual_length", utf8.RuneCountInString(query)))
		return fmt.Errorf("query exceeds maximum length of %d characters", maxQueryLength)
	}

	return nil
}

func (tf *TextFilter) processQuery(query string) (string, error) {
	var result strings.Builder
	inQuotes := false
	quoteChar := rune(0)

	for _, char := range query {
		switch {
		case char == '"' || char == '\'':
			if inQuotes && char == quoteChar {
				inQuotes = false
			} else if !inQuotes {
				inQuotes = true
				quoteChar = char
			}
			result.WriteRune(char)

		case inQuotes:
			// In quotes, only escape the quote character
			if char == quoteChar {
				result.WriteRune('\\')
			}
			result.WriteRune(char)

		case strings.ContainsRune(specialChars, char):
			// Outside quotes, escape special characters
			result.WriteRune('\\')
			result.WriteRune(char)

		default:
			result.WriteRune(char)
		}
	}

	return result.String(), nil
}

func (tf *TextFilter) validateFieldName(field string) error {
	client := quickwit.NewClient(quickwit.QuickwitURL, context.Background())
	client.ValidateFieldName(tf.IndexName, field)
	return nil
}
