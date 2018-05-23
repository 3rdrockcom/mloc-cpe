package helpers

import (
	"fmt"
	"strings"
	"time"
)

// Time is a custom time implementation
type Time struct {
	time.Time
}

// nilTime is an empty time object
var nilTime = time.Time{}

// dateFormat is the target time format
var dateFormat = "2006-01-02 15:04:05"

// UnmarshalJSON returns the JSON decoding of custom time
func (t *Time) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")

	if string(b) == "null" {
		return nil
	}
	tp, err := time.Parse(dateFormat, s)
	if err != nil {
		return err
	}
	*t = Time{tp}
	return nil
}

// MarshalJSON returns the JSON encoding of custom time
func (t *Time) MarshalJSON() ([]byte, error) {
	if t.Time == nilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", t.Time.Format(dateFormat))), nil
}

// IsSet checks if the time has been set
func (t *Time) IsSet() bool {
	return t.Time != nilTime
}
