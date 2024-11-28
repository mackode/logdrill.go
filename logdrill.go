package main

import (
  "flag"
  "fmt"
  "os"
)

func main() {
  flag.Parse()
  files := flag.Args()
  if len(files) == 0 {
    fmt.Printf("No input file\n")
    os.Exit(1)
  }
  logCh := make(chan LogEvent)
  for _, file := range files {
    tapLog(file, logCh)
  }
  rowCh := make(chan string)
  limitInCh := make(chan string)
  limitOutCh := make(chan string)
  geoCh := make(chan string)
  geoCached := Memoize(ipLookup)
  go limiter(limitInCh, limitOutCh)
  go func() {
    for {
      select {
        case ip := <-limitOutCh:
          geoCh <- geoCached(ip)
      }
    }
  }()
  go func() {
    for ev := range logCh {
      limitInCh <- ev.fields[0]
      rowCh <- fmt.Sprintf("%s %s %s", ev.dt.Format("15:00:00"), ev.fields[0], ev.fields[2])
    }
  }()
  uiStart(rowCh, geoCh)
}
