package pbars

import (
	"bytes"
	"strings"
)

// BarDrawerFunction is a function that takes in the progress value (0.0 -> 1.0) and a width
// and returns a string containing the progress bar. It does not handle any other information.
type BarDrawerFunction func(progress float32, width int) string

func clamp(v, bottom, top float32) float32 {
	if v < bottom {
		v = bottom
	} else if v > top {
		v = top
	}
	return v
}

var ratioBlocks = []rune{' ', '▏', '▎', '▍', '▌', '▋', '▊', '▉', '█'}

func ratioBlock(v float32) rune {
	index := int(v * 8.0)
	return ratioBlocks[index]
}

// DrawBarASCII prints a bar that looks like '[===============.....]' it should be compatible with
// any terminal but can only print whole characters and so has slightly less resolution than others.
func DrawBarASCII(perc float32, width int) string {
	perc = clamp(perc, 0, 1.0)
	innerwidth := width - 2
	filledwidth := int(float32(innerwidth) * perc)
	emptywidth := innerwidth - filledwidth
	buff := bytes.NewBufferString("[")
	buff.WriteString(strings.Repeat("=", filledwidth))
	buff.WriteString(strings.Repeat(".", emptywidth))
	buff.WriteString("]")
	return buff.String()
}

var _ BarDrawerFunction = DrawBarASCII

// DrawBarUTF8 prints a bar using the utf8 block characters. Because UTF8 has fractional block
// characters, this bar has a high resolution and is effective with less width than the ascii
// counterpart.
func DrawBarUTF8(perc float32, width int) string {
	perc = clamp(perc, 0, 1.0)
	innerwidth := width - 2
	filledwidth := int(float32(innerwidth) * perc)
	emptywidth := innerwidth - filledwidth
	buff := bytes.NewBufferString("|")
	buff.WriteString(strings.Repeat("█", filledwidth))
	subratio := (float32(innerwidth) * perc) - float32(filledwidth)
	if subratio > 0.0 && emptywidth > 0 {
		buff.WriteRune(ratioBlock(subratio))
		emptywidth--
	}
	buff.WriteString(strings.Repeat(" ", emptywidth))
	buff.WriteString("|")
	return buff.String()
}

var _ BarDrawerFunction = DrawBarUTF8
