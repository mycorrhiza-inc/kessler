package timestamp

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// RFC3339Time represents an RFC3339 DateTime
// @Description An RFC3339 DateTime
// @Schema {"type": "string", "example": "2024-02-27T12:34:56Z", "format": "date-time"}
type RFC3339Time time.Time

func (t RFC3339Time) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RFC3339))), nil
}

func (t *RFC3339Time) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = strings.Trim(str, "\"")
	if str == "" {
		*t = RFC3339Time{}
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}
	*t = RFC3339Time(parsed)
	return nil
}

func (t RFC3339Time) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t RFC3339Time) String() string {
	return time.Time(t).Format(time.RFC3339)
}

func KessTimeFromString(str string) (RFC3339Time, error) {
	kt := &RFC3339Time{}
	err := kt.UnmarshalJSON([]byte(fmt.Sprintf("\"%s\"", str)))
	if err != nil {
		return RFC3339Time{}, err
	}
	return *kt, nil
}

func KesslerTimeFromMMDDYYYY(dateStr string) (RFC3339Time, error) {
	if dateStr == "" {
		return RFC3339Time{}, errors.New("empty date string")
	}
	dateParts := strings.Split(dateStr, "/")
	if len(dateParts) != 3 {
		return RFC3339Time{}, errors.New("date string must be in the format MM/DD/YYYY")
	}
	month := dateParts[0]
	day := dateParts[1]
	year := dateParts[2]

	parsedDate, err := time.Parse("01/02/2006", fmt.Sprintf("%s/%s/%s", month, day, year))
	if err != nil {
		return RFC3339Time{}, err
	}
	return RFC3339Time(parsedDate), nil
}

func CreateRFC3339FromString(dateStr string) (string, error) {
	kt, err := KesslerTimeFromMMDDYYYY(dateStr)
	if err != nil {
		return "", err
	}
	jsonBytes, err := kt.MarshalJSON()
	if err != nil {
		return "", err
	}
	return strings.Trim(string(jsonBytes), "\""), nil
}
