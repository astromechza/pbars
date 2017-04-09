package pbars

import (
	"container/ring"
	"fmt"
	"math"
	"time"
)

type UpdateAveragingRateWatcher struct {
	startTime           time.Time
	lastUpdateTime      time.Time
	lastUpdatePosition  int64
	previousRatesBuffer *ring.Ring
	estimatedPercentage float32
	estimatedRemaining  time.Duration
	estimatedRate       float32
	hasEstimate         bool
	timefunc            func() time.Time
}

func (nw *UpdateAveragingRateWatcher) PercentageComplete() float32 {
	return nw.estimatedPercentage
}

func (nw *UpdateAveragingRateWatcher) EstimatedUnitsPerSecond() float32 {
	return nw.estimatedRate
}

func (nw *UpdateAveragingRateWatcher) EstimatedRemaining() time.Duration {
	return nw.estimatedRemaining
}

func (nw *UpdateAveragingRateWatcher) OverallUnitsPerSecond() float32 {
	return float32(nw.lastUpdatePosition) * float32(time.Second) / float32(nw.OverallElapsed())
}

func (nw *UpdateAveragingRateWatcher) OverallElapsed() time.Duration {
	return nw.lastUpdateTime.Sub(nw.startTime)
}

func (nw *UpdateAveragingRateWatcher) Update(position, length int64) {
	now := nw.timefunc()
	if position > length {
		position = length
	}
	if nw.startTime.IsZero() {
		nw.startTime = now
		nw.lastUpdatePosition = position
		nw.lastUpdateTime = now
		return
	}
	nw.hasEstimate = false
	if position > 0 && length > 0 {
		nw.estimatedPercentage = float32(position) / float32(length)
		positionDelta := position - nw.lastUpdatePosition
		if positionDelta >= 0 {
			elapsed := now.Sub(nw.lastUpdateTime)
			if elapsed <= 0 {
				return
			}
			currentRate := float32(positionDelta) / (float32(elapsed) / float32(time.Second))
			nw.previousRatesBuffer.Value = currentRate
			nw.previousRatesBuffer = nw.previousRatesBuffer.Next()

			total := float32(0)
			count := float32(0)
			nw.previousRatesBuffer.Do(func(x interface{}) {
				if x != nil {
					v := x.(float32)
					total += v
					count += 1.0
				}
			})
			if count > 0 {
				nw.estimatedRate = total / count
				amountLeft := float64(length - position)
				timeAtRate := (amountLeft / float64(nw.estimatedRate)) * float64(time.Second)
				nw.estimatedRemaining = time.Duration(math.Abs(timeAtRate))
				nw.hasEstimate = true
			}
		}
		nw.lastUpdatePosition = position
		nw.lastUpdateTime = now
	}
}

func (nw *UpdateAveragingRateWatcher) HasEstimate() bool {
	return nw.estimatedRemaining < (1<<63 - 1)
}

func (nw *UpdateAveragingRateWatcher) String() string {
	if nw.hasEstimate {
		return fmt.Sprintf("%.2f%% %.4f/s %s", 100.0*nw.estimatedPercentage, nw.estimatedRate, nw.estimatedRemaining)
	}
	return fmt.Sprintf("%.2f%%", 100.0*nw.estimatedPercentage)
}

func NewUpdateAveragingRateWatcher(window int) *UpdateAveragingRateWatcher {
	if window < 1 {
		window = 1
	}
	return &UpdateAveragingRateWatcher{
		previousRatesBuffer: ring.New(window),
		estimatedPercentage: 0.0,
		estimatedRemaining:  1<<63 - 1,
		estimatedRate:       0,
		timefunc:            time.Now,
	}
}

var _ RateWatcher = (*UpdateAveragingRateWatcher)(nil)
