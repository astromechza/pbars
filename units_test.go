package pbars

import (
	"testing"

	. "gopkg.in/go-playground/assert.v1"
)

func TestNoUnitFunc(t *testing.T) {
	Equal(t, NoUnitFunc(0), "0.00")
	Equal(t, NoUnitFunc(12312.788), "12312.79")
}

func TestByteFormatFunc(t *testing.T) {
	Equal(t, ByteFormatFunc(0), "0.00B")
	Equal(t, ByteFormatFunc(1000), "1000.00B")
	Equal(t, ByteFormatFunc(2512), "2.45KB")
	Equal(t, ByteFormatFunc(172638712), "164.64MB")
	Equal(t, ByteFormatFunc(17122638712), "15.95GB")
}
