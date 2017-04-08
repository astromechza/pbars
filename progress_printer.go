package pbars

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

type countingWriter struct {
	Next  io.Writer
	Count int64
}

func (w *countingWriter) Write(p []byte) (int, error) {
	n, err := w.Next.Write(p)
	w.Count += int64(n)
	return n, err
}

var _ io.Writer = (*countingWriter)(nil)

type ProgressPrinter struct {
	Title            string
	Output           io.Writer
	Width            int
	ShowPercentage   bool
	ShowRate         bool
	ShowTimeEstimate bool
	UnitFunc         UnitFormatFunc
	Ratewatcher      RateWatcher
	Bardrawer        BarDrawerFunction

	lastprinttime  time.Time
	lastdrawnwidth int
}

const minRedrawDelay = 16 * time.Millisecond

func (pp *ProgressPrinter) Update(position, length int64) {
	pp.Ratewatcher.Update(position, length)
	if pp.Output != nil && (position >= length || time.Now().Sub(pp.lastprinttime) > minRedrawDelay) {
		pp.Reprint()
	}
}

func (pp *ProgressPrinter) Reprint() {
	pp.Output.Write([]byte("\r"))
	drawwidth := 0
	if len(pp.Title) > 0 {
		pp.Output.Write([]byte(pp.Title + " "))
		drawwidth += utf8.RuneCountInString(pp.Title) + 1
	}
	pp.Output.Write([]byte(pp.Bardrawer(pp.Ratewatcher.PercentageComplete(), pp.Width)))
	pp.Output.Write([]byte(" "))
	if pp.ShowPercentage {
		n, _ := fmt.Fprintf(pp.Output, "%.2f%% ", pp.Ratewatcher.PercentageComplete()*100.0)
		drawwidth += n
	}
	if pp.Ratewatcher.HasEstimate() {
		if pp.ShowRate {
			out := fmt.Sprintf("%s/s ", pp.UnitFunc(float64(pp.Ratewatcher.EstimatedUnitsPerSecond())))
			drawwidth += utf8.RuneCountInString(out)
			fmt.Fprint(pp.Output, out)
		}
		if pp.ShowTimeEstimate {
			n, _ := fmt.Fprintf(pp.Output, "%s ", FormatDuration(pp.Ratewatcher.EstimatedRemaining()))
			drawwidth += n
		}
	}
	d := pp.lastdrawnwidth - drawwidth
	if d > 0 {
		pp.Output.Write([]byte(strings.Repeat(" ", d)))
	}
	pp.lastdrawnwidth = drawwidth
	pp.lastprinttime = time.Now()
}

// Clear just clears the line using a \r character and spaces
func (pp *ProgressPrinter) Clear() {
	pp.Output.Write([]byte("\r"))
	pp.Output.Write([]byte(strings.Repeat(" ", pp.Width+3+pp.lastdrawnwidth)))
	pp.Output.Write([]byte("\r"))
}

// Done just writes a newline so we move on
func (pp *ProgressPrinter) Done() {
	pp.Output.Write([]byte("\n"))
}

// Interruptf interrupts the bar enough to print a message and then reprint the bar
func (pp *ProgressPrinter) Interruptf(format string, args ...interface{}) {
	pp.Output.Write([]byte("\r"))
	text := fmt.Sprintf(format, args...)
	if len(text) < pp.lastdrawnwidth {
		text += strings.Repeat(" ", pp.Width+3+pp.lastdrawnwidth-len(text))
	}
	fmt.Fprint(pp.Output, text+"\n")
	pp.Reprint()
}

var _ ProgressReceiver = (*ProgressPrinter)(nil)

// NewProgressPrinter constructs a progress bar printer for most use cases
func NewProgressPrinter(title string, width int, utf8 bool) *ProgressPrinter {
	barfunc := DrawBarASCII
	if utf8 {
		barfunc = DrawBarUTF8
	}

	return &ProgressPrinter{
		Title:            title,
		Output:           os.Stdout,
		Width:            width,
		ShowPercentage:   true,
		ShowRate:         true,
		ShowTimeEstimate: true,
		UnitFunc:         NoUnitFunc,
		Bardrawer:        barfunc,
		Ratewatcher:      NewUpdateAveragingRateWatcher(5),
	}
}
