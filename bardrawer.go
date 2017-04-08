package pbars

import (
	"bytes"
	"strings"
)

type BarDrawerFunction func(float32, int) string

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
