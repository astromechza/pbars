package pbars

import (
	"fmt"
	"math"
	"time"
)

type NaiveRateWatcher struct {
	firstupdatetime     time.Time
	lastseenposition    int64
	lastupdatetime      time.Time
	estimatedPercentage float32
	estimatedRemaining  time.Duration
	estimatedRate       float32
	hasEstimate         bool
	timefunc            func() time.Time
}

func (nw *NaiveRateWatcher) PercentageComplete() float32 {
	return nw.estimatedPercentage
}

func (nw *NaiveRateWatcher) EstimatedUnitsPerSecond() float32 {
	return nw.estimatedRate
}

func (nw *NaiveRateWatcher) EstimatedRemaining() time.Duration {
	return nw.estimatedRemaining
}

func (nw *NaiveRateWatcher) OverallUnitsPerSecond() float32 {
	return float32(nw.lastseenposition) * float32(time.Second) / float32(nw.OverallElapsed())
}

func (nw *NaiveRateWatcher) OverallElapsed() time.Duration {
	return nw.lastupdatetime.Sub(nw.firstupdatetime)
}

func (nw *NaiveRateWatcher) Update(position, length int64) {
	nw.lastupdatetime = nw.timefunc()
	if position > length {
		position = length
	}
	nw.lastseenposition = position
	if nw.firstupdatetime.IsZero() {
		nw.firstupdatetime = nw.lastupdatetime
		return
	}
	elapsed := nw.OverallElapsed()
	nw.hasEstimate = false
	if elapsed > 0 && position > 0 && length > 0 {
		nw.estimatedPercentage = float32(position) / float32(length)
		nw.estimatedRate = nw.OverallUnitsPerSecond()
		timeAtRate := (float64(length-position) * float64(time.Second) / float64(nw.estimatedRate))
		nw.estimatedRemaining = time.Duration(math.Abs(timeAtRate))
		nw.hasEstimate = true
	}
}

func (nw *NaiveRateWatcher) HasEstimate() bool {
	return nw.hasEstimate
}

func (nw *NaiveRateWatcher) String() string {
	if nw.hasEstimate {
		return fmt.Sprintf("%.2f%% %.4f/s %s", 100.0*nw.estimatedPercentage, nw.estimatedRate, nw.estimatedRemaining)
	}
	return fmt.Sprintf("%.2f%%", 100.0*nw.estimatedPercentage)
}

func NewNaiveRateWatcher() *NaiveRateWatcher {
	return &NaiveRateWatcher{
		estimatedPercentage: 0.0,
		estimatedRemaining:  1<<63 - 1,
		estimatedRate:       0,
		timefunc:            time.Now,
	}
}

var _ RateWatcher = (*NaiveRateWatcher)(nil)
