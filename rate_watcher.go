package pbars

import (
	"fmt"
	"time"
)

// RateWatcher presents an object that can watch the accumulation of progress updates and report
// the estimated rate of completion.
type RateWatcher interface {
	ProgressReceiver
	fmt.Stringer

	// PercentageComplete should return a float between 0.0 and 1.0
	PercentageComplete() float32

	// HasEstimate
	HasEstimate() bool

	// EstimatedUnitsPerSecond should return a value representing the number of units progressing per second
	EstimatedUnitsPerSecond() float32

	// EstimatedRemaining should return an estimated time remaining
	EstimatedRemaining() time.Duration
}
