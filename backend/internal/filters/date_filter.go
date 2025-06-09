package filters

import (
	"fmt"
	"kessler/pkg/logger"
	"time"

	"go.uber.org/zap"
)

type DateRangeInput struct {
	Start     string `json:"start"`     // RFC3339 formatted date string
	End       string `json:"end"`       // RFC3339 formatted date string
	Inclusive bool   `json:"inclusive"` // Whether the range is inclusive of boundaries
	TimeZone  string `json:"timezone"`  // Optional timezone (defaults to UTC)
}

type DateRangeOutput struct {
	StartTimestamp string `json:"start_timestamp"`
	EndTimestamp   string `json:"end_timestamp"`
	QueryString    string `json:"query_string"` // Quickwit compatible query string
}

// NewDateRangeFilter creates a new date range filter implementation
func NewDateRangeFilter() FilterFunc {
	log = logger.Named("date_range_filter")

	return func(input interface{}) (interface{}, error) {
		// Type assertion for input
		dateRange, ok := input.(DateRangeInput)
		if !ok {
			log.Error("invalid input type",
				zap.String("expected_type", "DateRangeInput"),
				zap.Any("received_input", input))
			return nil, fmt.Errorf("invalid input type for date range filter")
		}

		// Set default timezone if not specified
		if dateRange.TimeZone == "" {
			dateRange.TimeZone = "UTC"
			log.Debug("using default timezone", zap.String("timezone", "UTC"))
		}

		// Parse timezone
		location, err := time.LoadLocation(dateRange.TimeZone)
		if err != nil {
			log.Error("invalid timezone",
				zap.String("timezone", dateRange.TimeZone),
				zap.Error(err))
			return nil, fmt.Errorf("invalid timezone: %w", err)
		}

		// Parse and validate dates
		start, err := time.Parse(time.RFC3339, dateRange.Start)
		if err != nil {
			log.Error("invalid start date format",
				zap.String("date", dateRange.Start),
				zap.Error(err))
			return nil, fmt.Errorf("invalid start date format (RFC3339 required): %w", err)
		}

		end, err := time.Parse(time.RFC3339, dateRange.End)
		if err != nil {
			log.Error("invalid end date format",
				zap.String("date", dateRange.End),
				zap.Error(err))
			return nil, fmt.Errorf("invalid end date format (RFC3339 required): %w", err)
		}

		// Normalize dates to specified timezone
		start = start.In(location)
		end = end.In(location)

		// Ensure start is before end
		if start.After(end) {
			start, end = end, start
		}

		// Create Quickwit query string based on inclusivity
		var queryString string
		if dateRange.Inclusive {
			queryString = fmt.Sprintf("timestamp:[%s TO %s]",
				start.Format(time.RFC3339),
				end.Format(time.RFC3339))
		} else {
			queryString = fmt.Sprintf("timestamp:{%s TO %s}",
				start.Format(time.RFC3339),
				end.Format(time.RFC3339))
		}

		output := DateRangeOutput{
			StartTimestamp: start.Format(time.RFC3339),
			EndTimestamp:   end.Format(time.RFC3339),
			QueryString:    queryString,
		}

		return output, nil
	}
}
