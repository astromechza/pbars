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

## Example

There's an [example in the /example directory](example/main.go) which does the following:

![progress gif](http://i.imgur.com/GFJTmhE.gif?raw=true)

## Features

- UTF8 blocks vs pure ASCII (sometimes a terminal doesn't support utf8 chars)

```golang
pbars.NewProgressPrinter("My Title", 50, true)   // utf8 mode
pbars.NewProgressPrinter("My Title", 50, false)   // ascii mode
```

- Uses the overall units per second and elapsed time once you call `Done`

- Customisable unit formats

By default the progress bar rate is formatted as 'units' per second. But often you'll want a measure of bytes or bits.
The ProgressPrinter struct allows you to set the `UnitFunc` to be whatever you want as long as it looks like
`func(v float64) string`. An example is the `pbars.ByteFormatFunc` that will convert to `B`, `KB`, `MB` etc.

- `Interruptf` method for printing messages while the progress bar continues (see the example)

- `Clear` method for clearing and removing the progress bar once you no longer need it

- MIT licensed, use it, abuse it.
