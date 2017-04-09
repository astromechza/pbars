package pbars

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

// ProgressPrinter is the model that represents a single progress bar and is responsible for
// printing it to the terminal correctly. It ties together all of the formatting, rate
// calculations, and bar drawing.
type ProgressPrinter struct {

	// Title is the string printed to the left of the progress bar. It may be blank,
	// and it may contain utf8 characters. It is not included in the width.
	Title string

	// Output is the output writer to print the progress bar to. This is usually os.Stdout,
	// but for tests and certain other cases can be overriden or be nil.
	Output io.Writer

	// Width is the width of the actual bar in the progress bar.
	Width int

	// ShowPercentage sets whether the percentage should be visible on the right hand side.
	ShowPercentage bool

	// ShowRate sets whether the rate in units per second should be visible on the right hand side.
	ShowRate bool

	// ShowTimeEstimate sets whether the estimate of the remaining time should be visible on the right hand side.
	ShowTimeEstimate bool

	// UnitFunc is the formatter applied to the rate value in order to format it to subject-specific units.
	UnitFunc UnitFormatFunc

	// RateWatcher is a model that is used to track the rate, percentage, and time estimates. You probably don't
	// want to change this.
	Ratewatcher RateWatcher

	// Bardrawer is the function used to turn the percentage and width into a progress bar. This can be overriden to
	// pick between utf8 or ascii or to implement your own bar style.
	Bardrawer BarDrawerFunction

	// NonTTY can be used to set the progress bar into a nontty mode when writing output to files. The bar will be printed
	// only once when Done is called, and interrupt messages will function as normal.
	NonTTY bool

	done           bool
	lastprinttime  time.Time
	lastdrawnwidth int
}

const minRedrawDelay = 16 * time.Millisecond

// Update updates the progress bar, possibly redrawing it if the position has changed and there is
// and output stream attached.
func (pp *ProgressPrinter) Update(position, length int64) {
	pp.Ratewatcher.Update(position, length)
	if !pp.NonTTY && time.Now().Sub(pp.lastprinttime) > minRedrawDelay {
		pp.Output.Write([]byte("\r"))
		pp.Reprint()
	}
}

// Reprint will reprint the progress bar in place. You probably don't want to call this but its public just
// in case.
func (pp *ProgressPrinter) Reprint() {
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
	if pp.done {
		if pp.ShowRate {
			out := fmt.Sprintf("%s/s ", pp.UnitFunc(float64(pp.Ratewatcher.OverallUnitsPerSecond())))
			drawwidth += utf8.RuneCountInString(out)
			fmt.Fprint(pp.Output, out)
		}
		if pp.ShowTimeEstimate {
			n, _ := fmt.Fprintf(pp.Output, "%s ", FormatDuration(pp.Ratewatcher.OverallElapsed()))
			drawwidth += n
		}
	} else if pp.Ratewatcher.HasEstimate() {
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

// Clear just clears the line using a \r character and spaces.
func (pp *ProgressPrinter) Clear() {
	pp.Output.Write([]byte("\r"))
	pp.Output.Write([]byte(strings.Repeat(" ", pp.Width+3+pp.lastdrawnwidth)))
	pp.Output.Write([]byte("\r"))
}

// Done reprints the bar using the overall rate and overall elapsed time and terminates the bar
// with a newline.
func (pp *ProgressPrinter) Done() {
	pp.done = true
	if !pp.NonTTY {
		pp.Output.Write([]byte("\r"))
	}
	pp.Reprint()
	pp.Output.Write([]byte("\n"))
}

// Interruptf interrupts the bar enough to print a message and then reprint the bar on the line below.
func (pp *ProgressPrinter) Interruptf(format string, args ...interface{}) {
	if !pp.NonTTY {
		pp.Output.Write([]byte("\r"))
	}
	text := fmt.Sprintf(format, args...)
	if len(text) < pp.lastdrawnwidth {
		text += strings.Repeat(" ", pp.Width+3+pp.lastdrawnwidth-len(text))
	}
	fmt.Fprint(pp.Output, text+"\n")
	if !pp.NonTTY {
		pp.Reprint()
	}
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
