package pbars

import (
	"testing"

	"unicode/utf8"

	. "gopkg.in/go-playground/assert.v1"
)

func TestDrawBarASCII(t *testing.T) {
	z := DrawBarASCII(0.5, 30)
	Equal(t, len(z), 30)
	Equal(t, z, "[==============..............]")
	z = DrawBarASCII(0.75, 30)
	Equal(t, len(z), 30)
	Equal(t, z, "[=====================.......]")
	z = DrawBarASCII(1, 30)
	Equal(t, len(z), 30)
	Equal(t, z, "[============================]")
}

func TestDrawBarUTF8(t *testing.T) {
	z := DrawBarUTF8(0.5, 30)
	Equal(t, z, "|██████████████              |")
	Equal(t, utf8.RuneCountInString(z), 30)
	z = DrawBarUTF8(0.75, 30)
	Equal(t, z, "|█████████████████████       |")
	Equal(t, utf8.RuneCountInString(z), 30)
	z = DrawBarUTF8(1, 30)
	Equal(t, z, "|████████████████████████████|")
	Equal(t, utf8.RuneCountInString(z), 30)
	z = DrawBarUTF8(1, 30)
	Equal(t, z, "|████████████████████████████|")
	Equal(t, utf8.RuneCountInString(z), 30)
	z = DrawBarUTF8(0.025, 12)
	Equal(t, z, "|▎         |")
	Equal(t, utf8.RuneCountInString(z), 12)
	z = DrawBarUTF8(0.05, 12)
	Equal(t, z, "|▌         |")
	Equal(t, utf8.RuneCountInString(z), 12)
	z = DrawBarUTF8(0.09, 12)
	Equal(t, z, "|▉         |")
	Equal(t, utf8.RuneCountInString(z), 12)
}

func BenchmarkUTF8BarDraw100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DrawBarUTF8(0.512321, 100)
	}
}

func BenchmarkASCIIBarDraw100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DrawBarASCII(0.512321, 100)
	}
}
