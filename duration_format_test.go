package pbars

import (
	"testing"
	"time"

	. "gopkg.in/go-playground/assert.v1"
)

func TestFormatDuration(t *testing.T) {
	Equal(t, FormatDuration(time.Duration(123)), "123ns")
	Equal(t, FormatDuration(time.Duration(12345)), "12.34µs")
	Equal(t, FormatDuration(time.Duration(12345678)), "12.34ms")
	Equal(t, FormatDuration(time.Duration(5477*time.Second/100)), "54.77s")
	Equal(t, FormatDuration(time.Duration(360*time.Second)), "6m0.00s")
	Equal(t, FormatDuration(time.Duration(234*time.Minute)), "3h54m0.00s")
}
