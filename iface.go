package pbars

import (
	"fmt"
	"time"
)

// ProgressReceiver is the interface implemented by many of the objects in this package.
// It is simply a thing that can receive progress updates.
type ProgressReceiver interface {

	// Update performs some work based on the current progres.
	Update(progress, length int64)
}

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

	// OverallUnitsPerSecond should return the overall rate seen by the watcher. Pretty much length / elapsed.
	OverallUnitsPerSecond() float32

	// OverallElapsed should return the time since the bar was first seen by the user
	OverallElapsed() time.Duration
}
