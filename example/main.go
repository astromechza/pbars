package main

import (
	"time"

	"github.com/AstromechZA/pbars"
)

func main() {

	// setup the printer
	pp := pbars.NewProgressPrinter("My λ Title", 50, true)
	pp.UnitFunc = pbars.ByteFormatFunc

	// go!
	for i := 0; i < 400; i++ {
		pp.Update(int64(i+1), 400)
		time.Sleep(16 * time.Millisecond)
	}

	// new line required afterwards :)
	pp.Done()

	// setup the printer
	pp = pbars.NewProgressPrinter("My Title", 50, false)

	// interrupt 1
	pp.Interruptf("Beginning progress..")

	// go!
	for i := 0; i < 200; i++ {
		pp.Update(int64(i+1), 400)
		time.Sleep(16 * time.Millisecond)
	}

	// interrupt 2
	pp.Interruptf("part 2..")

	// go!
	for i := 200; i < 400; i++ {
		pp.Update(int64(i+1), 400)
		time.Sleep(16 * time.Millisecond)
	}

	// new line required afterwards :)
	pp.Done()
}
