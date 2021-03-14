package irapi

import (
	"strconv"
	"time"
)

// times.go implements various wrapper types for the multiple ways time and
// duration values are represented in iRacing

// Timestamp represents a time as a milliseconds Unix timestamp in JSON
type Timestamp time.Time

// UnmarshalJSON decodes a unix timestamp as a time
func (t *Timestamp) UnmarshalJSON(b []byte) error {

	unix, err := strconv.ParseInt(string(b), 10, 64)

	if err != nil {
		return err
	}

	*t = Timestamp(time.Unix(unix/1000, 0))

	return nil
}

// MarshalJSON encodes a time as a millisecond Unix timestamp
func (t Timestamp) MarshalJSON() ([]byte, error) {
	ms := time.Time(t).Unix() * 1000

	return []byte(strconv.FormatInt(ms, 10)), nil
}

// SimTime represents a time as a string
type SimTime time.Time

const simTimeFormat = "2006-01-02 15:04:05"
const simTimeFormatLoose = "2006-01-02 15:04:05"

func (s *SimTime) UnmarshalJSON(b []byte) error {

	str, err := strconv.Unquote(string(b))

	if err != nil {
		return err
	}

	t, err := time.ParseInLocation(simTimeFormat, str, time.UTC)

	if err != nil {
		t, err = time.ParseInLocation(simTimeFormatLoose, str, time.UTC)

		if err != nil {
			return err
		}
	}

	*s = SimTime(t)
	return nil
}

// Milliseconds is a Duration represented in JSON as an integer number
// of milliseconds
type Milliseconds time.Duration

// UnmarshalJSON decodes a JSON value as Milliseconds
func (m *Milliseconds) UnmarhsalJSON(b []byte) error {

	ms, err := strconv.ParseInt(string(b), 10, 64)

	if err != nil {
		return err
	}

	*m = Milliseconds(time.Duration(ms) * time.Millisecond)

	return nil
}

// MarshalJSON encodes Milliseconds to JSON
func (m Milliseconds) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(m), 10)), nil
}
