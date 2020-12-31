package irapi

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// StringifiedUint64 is a uint64 represented as a human-readable string
type StringifiedUint64 uint64

func (u *StringifiedUint64) UnmarhsalJSON(b []byte) error {
	// Get the string
	str := string(b)

	// Remove the commas
	str = strings.ReplaceAll(str, ",", "")

	// Parse
	val, err := strconv.ParseUint(str, 10, 64)

	if err != nil {
		return err
	}

	*u = StringifiedUint64(val)

	return nil
}

func (u StringifiedUint64) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatUint(uint64(u), 10)), nil
}

// Laptime represents a laptime as a Duration, but with a textual representation of a wall-clock
type Laptime time.Duration

var laptimeRegexp = regexp.MustCompile(`^((?P<minutes>\d+):)?(?P<seconds>\d|[1-5]\d)\.(?P<millis>\d{3})$`)

func (l *Laptime) UnmarhsalJSON(b []byte) error {

	str := string(b)

	str, err := strconv.Unquote(str)

	if err != nil {
		return err
	}

	parts := laptimeRegexp.FindStringSubmatch(str)

	var minutes uint64

	if parts[2] == "" {
		minutes = 0
	} else {
		minutes, err = strconv.ParseUint(parts[2], 10, 64)

		if err != nil {
			return err
		}
	}

	seconds, err := strconv.ParseUint(parts[3], 10, 64)

	if err != nil {
		return err
	}

	millis, err := strconv.ParseUint(parts[4], 10, 64)

	if err != nil {
		return err
	}

	*l = Laptime((minutes * uint64(time.Minute)) + (seconds * uint64(time.Second)) + (millis * uint64(time.Millisecond)))

	return nil
}

func (l Laptime) String() string {

	d := time.Duration(l)

	minutes := int64(math.Trunc(d.Minutes()))
	seconds := int64(math.Trunc(d.Seconds())) - (60 * minutes)
	millis := (d.Milliseconds()) - (1000 * seconds) - (60 * 1000 * minutes)

	if minutes > 0 {
		return fmt.Sprintf("%d:%d.%03d", minutes, seconds, millis)
	}

	return fmt.Sprintf("%d.%03d", seconds, millis)
}

func (l Laptime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + l.String() + `"`), nil
}
