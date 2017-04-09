package pbars

import (
	"testing"
	"time"

	. "gopkg.in/go-playground/assert.v1"
)

func TestUpdateAveragingRateWatcher(t *testing.T) {
	n := NewUpdateAveragingRateWatcher(4)
	n.timefunc = func() time.Time { return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC) }
	n.Update(0, 100)
	Equal(t, n.String(), "0.00%")
	n.timefunc = func() time.Time { return time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC) }
	n.Update(1, 100)
	Equal(t, n.String(), "1.00% 1.0000/s 1m39s")
	n.timefunc = func() time.Time { return time.Date(2000, 1, 1, 0, 0, 2, 0, time.UTC) }
	n.Update(2, 100)
	Equal(t, n.String(), "2.00% 1.0000/s 1m38s")
	n.timefunc = func() time.Time { return time.Date(2000, 1, 1, 0, 0, 3, 0, time.UTC) }
	n.Update(3, 100)
	Equal(t, n.String(), "3.00% 1.0000/s 1m37s")
	n.timefunc = func() time.Time { return time.Date(2000, 1, 1, 0, 0, 4, 0, time.UTC) }
	n.Update(4, 100)
	Equal(t, n.String(), "4.00% 1.0000/s 1m36s")
	n.timefunc = func() time.Time { return time.Date(2000, 1, 1, 0, 0, 5, 0, time.UTC) }
	n.Update(7, 100)
	Equal(t, n.String(), "7.00% 1.5000/s 1m2s")
	n.timefunc = func() time.Time { return time.Date(2000, 1, 1, 0, 0, 6, 0, time.UTC) }
	n.Update(9, 100)
	Equal(t, n.String(), "9.00% 1.7500/s 52s")
	n.timefunc = func() time.Time { return time.Date(2000, 1, 1, 0, 0, 7, 0, time.UTC) }
	n.Update(11, 100)
	Equal(t, n.String(), "11.00% 2.0000/s 44.5s")
	n.timefunc = func() time.Time { return time.Date(2000, 1, 1, 0, 0, 8, 0, time.UTC) }
	n.Update(11, 100)
	Equal(t, n.String(), "11.00% 1.7500/s 50.857142857s")
}

func BenchmarkUpdateAveragingRateWatcherUpdateNTimes(b *testing.B) {
	n := NewUpdateAveragingRateWatcher(4)
	n.Update(0, int64(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n.Update(int64(i), int64(b.N))
	}
}
