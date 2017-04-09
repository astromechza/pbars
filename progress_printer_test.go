package pbars

import (
	"io/ioutil"
	"testing"
	"time"

	"bytes"

	. "gopkg.in/go-playground/assert.v1"
)

func BenchmarkProgressPrinterUpdate(b *testing.B) {
	pp := NewProgressPrinter("Some title", 100, true)
	pp.Output = ioutil.Discard
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pp.Update(int64(i), int64(b.N))
	}
}

func TestProgressPrinterUTF8Stuff(t *testing.T) {
	pp := NewProgressPrinter("Special λϴ", 40, true)
	buff := bytes.NewBufferString("")
	pp.Output = buff

	customRateWatcher := NewUpdateAveragingRateWatcher(4)
	customRateWatcher.timefunc = func() time.Time { return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC) }
	pp.Ratewatcher = customRateWatcher

	pp.Update(0, 100)
	Equal(t, buff.String(), "\rSpecial λϴ |                                      | 0.00% ")
	buff.Reset()

	customRateWatcher.timefunc = func() time.Time { return time.Date(2000, 1, 1, 0, 10, 0, 0, time.UTC) }
	pp.Ratewatcher = customRateWatcher

	pp.Update(53, 100)
	pp.Reprint()

	Equal(t, buff.String(), "Special λϴ |████████████████████▏                 | 53.00% 0.09/s 8m52.07s ")
}
