# `pbars` `|██████████▏                 |`

[![GoDoc](https://godoc.org/github.com/AstromechZA/pbars?status.svg)](https://godoc.org/github.com/AstromechZA/pbars)
[![Build](https://travis-ci.org/AstromechZA/pbars.svg?branch=master)](https://travis-ci.org/AstromechZA/pbars/)

```golang
// setup the printer
pp := pbars.NewProgressPrinter("My Title", 50, true)

// go!
for i := 0; i < 400; i++ {
    pp.Update(int64(i+1), 400)
    time.Sleep(16 * time.Millisecond)
}

// new line required afterwards :)
pp.Done()
```

![progress gif](http://i.imgur.com/tlD6QqI.gif?raw=true)

## Features

- UTF8 blocks vs pure ASCII (sometimes a terminal doesn't support utf8 chars)

```golang
pbars.NewProgressPrinter("My Title", 50, true)   // utf8 mode
pbars.NewProgressPrinter("My Title", 50, false)   // ascii mode
```

- Customisable unit formats

By default the progress bar rate is formatted as 'units' per second. But often you'll want a measure of bytes or bits.
The ProgressPrinter struct allows you to set the `UnitFunc` to be whatever you want as long as it looks like
`func(v float64) string`. An example is the `pbars.ByteFormatFunc` that will convert to `B`, `KB`, `MB` etc.
