package search

import (
	"fmt"
	"time"
)

func convertToRFC3339(date string) (string, error) {
	layout := "2006-01-02"

	parsedDate, err := time.Parse(layout, date)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err)
	}
	parsedDateString := parsedDate.Format(time.RFC3339)
	return parsedDateString, nil
}
