package helpers

import (
	"fmt"
	"strings"
	"time"
)

type Time struct {
	time.Time
}

var nilTime = time.Time{}
var dateFormat = "2006-01-02 15:04:05"

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

func (t *Time) MarshalJSON() ([]byte, error) {
	if t.Time == nilTime {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", t.Time.Format(dateFormat))), nil
}

func (t *Time) IsSet() bool {
	return t.Time != nilTime
}
