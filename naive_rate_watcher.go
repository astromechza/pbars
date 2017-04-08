package pbars

import (
	"fmt"
	"math"
	"time"
)

type NaiveRateWatcher struct {
	startTime           time.Time
	estimateTime        time.Time
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

func (nw *NaiveRateWatcher) Update(position, length int64) {
	if nw.startTime.IsZero() {
		nw.startTime = nw.timefunc()
		return
	}
	nw.estimateTime = nw.timefunc()
	elapsed := nw.estimateTime.Sub(nw.startTime)
	if position > length {
		position = length
	}
	nw.hasEstimate = false
	if elapsed > 0 && position > 0 && length > 0 {
		nw.estimatedPercentage = float32(position) / float32(length)
		amountLeft := float64(length - position)
		nw.estimatedRate = float32(position) * float32(time.Second) / float32(elapsed)
		timeAtRate := (amountLeft * float64(time.Second) / float64(nw.estimatedRate))
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
