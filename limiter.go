package main

import (
  "fmt"
  "time"
)

func Memoize(fn func(string) (string, error)) func(string) string {
  cache := map[string]string{}
  return func(n string) string {
    if val, ok := cache[n]; ok {
      return val + " (cached)"
    }
    val, err := fn(n)
    if err != nil {
      return fmt.Sprintf("err=%v", err)
    }
    cache[n] = val
    return val
  }
}

func limiter(in <-chan string, out chan<- string) {
  pause := false
  ticker := time.NewTicker(5 * time.Second)
  queue := ""
  for {
    if queue != "" && !pause {
      out <- queue
      queue = ""
      pause = true
    }
    select {
      case item := <-in:
        queue = item
      case <-ticker.C:
        pause = false
    }
  }
}
