package timestamp

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// KesslerTime represents an RFC3339 DateTime
// @Description A RFC3339 DateTime
// @Schema {"type": "string", "example": "2024-02-27T12:34:56Z", "format": "date-time"}
type KesslerTime time.Time

func (t KesslerTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RFC3339))), nil
}

func (t *KesslerTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	str = strings.Trim(str, "\"")
	if str == "" {
		*t = KesslerTime{}
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}
	*t = KesslerTime(parsed)
	return nil
}

func (t KesslerTime) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t KesslerTime) String() string {
	return time.Time(t).Format(time.RFC3339)
}

func KessTimeFromString(str string) (KesslerTime, error) {
	kt := &KesslerTime{}
	err := kt.UnmarshalJSON([]byte(fmt.Sprintf("\"%s\"", str)))
	if err != nil {
		return KesslerTime{}, err
	}
	return *kt, nil
}

func KesslerTimeFromMMDDYYYY(dateStr string) (KesslerTime, error) {
	if dateStr == "" {
		return KesslerTime{}, errors.New("empty date string")
	}
	dateParts := strings.Split(dateStr, "/")
	if len(dateParts) != 3 {
		return KesslerTime{}, errors.New("date string must be in the format MM/DD/YYYY")
	}
	month := dateParts[0]
	day := dateParts[1]
	year := dateParts[2]

	parsedDate, err := time.Parse("01/02/2006", fmt.Sprintf("%s/%s/%s", month, day, year))
	if err != nil {
		return KesslerTime{}, err
	}
	return KesslerTime(parsedDate), nil
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
