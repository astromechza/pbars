# `pbars` `|██████████▏                 |`

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

![progress gif](http://i.imgur.com/tlD6QqI.gifv)
