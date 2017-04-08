package pbars

import (
	"testing"

	"time"

	. "gopkg.in/go-playground/assert.v1"
)

func TestNaiveRateWatcher(t *testing.T) {
	n := NewNaiveRateWatcher()
	Equal(t, n.String(), "0.00%")
	n.timefunc = func() time.Time { return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC) }
	n.Update(0, 100)
	Equal(t, n.String(), "0.00%")
	n.timefunc = func() time.Time { return time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC) }
	n.Update(1, 100)
	Equal(t, n.String(), "1.00% 1.0000/s 1m39s")
	n.timefunc = func() time.Time { return time.Date(2000, 1, 1, 0, 0, 2, 0, time.UTC) }
	n.Update(50, 100)
	Equal(t, n.String(), "50.00% 25.0000/s 2s")
	n.timefunc = func() time.Time { return time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC) }
	n.Update(100, 100)
	Equal(t, n.String(), "100.00% 10.0000/s 0s")
}

func TestNaiveRateWatcherEdgeCases(t *testing.T) {
	n := NewNaiveRateWatcher()
	n.timefunc = func() time.Time { return time.Date(2000, 1, 1, 0, 1, 0, 0, time.UTC) }
	n.Update(101, 100)
	Equal(t, n.String(), "0.00%")
	n.timefunc = func() time.Time { return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC) }
	n.Update(101, 100)
	Equal(t, n.String(), "0.00%")
	n.timefunc = func() time.Time { return time.Date(2000, 1, 1, 0, 10, 0, 0, time.UTC) }
	n.Update(101, 100)
	Equal(t, n.String(), "100.00% 0.1852/s 0s")
}

func BenchmarkNaiveRateWatcherUpdateNTimes(b *testing.B) {
	n := NewNaiveRateWatcher()
	n.Update(0, int64(b.N))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n.Update(int64(i), int64(b.N))
	}
}
