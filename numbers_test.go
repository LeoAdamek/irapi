package irapi

import (
	"testing"
	"time"
)

func TestLapTimeUnmarshalJSON(t *testing.T) {

	validTimes := []string{
		`"1:43.580"`,
		`"99:59.999"`,
		`"15.000"`,
		`"1.234"`,
		`"5:55.555"`,
	}

	for _, v := range validTimes {
		var l Laptime

		if err := l.UnmarhsalJSON([]byte(v)); err != nil {
			t.Log("Unexpected Parse Error:", err)
			t.Fail()
		}
	}
}

func TestLaptimeString(t *testing.T) {
	mapping := map[string]Laptime{
		"1:42.548": Laptime(1*time.Minute + 42*time.Second + 548*time.Millisecond),
		"1.234":    Laptime(1*time.Second + 234*time.Millisecond),
	}

	for expected, input := range mapping {

		actual := input.String()

		if string(actual) != expected {
			t.Logf("Expected '%s' but got '%s'", expected, string(actual))
			t.Fail()
		}

	}
}

func BenchmarkLaptimeString(b *testing.B) {
	input := Laptime(3*time.Minute + 28*time.Second + 544*time.Millisecond)

	for i := 0; i < b.N; i++ {
		input.String()
	}
}

func BenchmarkLaptimeUnmarshalJSON(b *testing.B) {
	validTimes := []string{
		`"1:43.580"`,
		`"99:59.999"`,
		`"15.000"`,
		`"1.234"`,
		`"5:55.555"`,
	}

	var l Laptime

	for i := 0; i < b.N; i++ {
		for _, v := range validTimes {
			l.UnmarhsalJSON([]byte(v))
		}
	}
}
