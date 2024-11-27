package main

import (
  "github.com/hpcloud/tail"
  "log"
  "regexp"
  "time"
)

type LogEvent struct {
  dt time.Time
  fields []string
}

func tapLog(fileName string, ch chan LogEvent) {
  t, err := tail.TailFile(fileName, tail.Config{Follow: true})
  if err != nil {
    log.Fatalf("%v", err)
  }
  re := regexp.MustCompile(`(\S+) \S+ \S+ \[(.*?)\] "[^/]*(/.*?)\s`)

  go func() {
    for {
      line := <-t.Lines
      matches := re.FindStringSubmatch(line.Text)
      if len(matches) != 4 {
        log.Fatalf("Invalid line: %s", line.Text)
      }
      layout := "02/Jan/2006:15:00:00-0700"
      dt, err := time.Parse(layout, matches[2])
      if err != nil {
        log.Fatalf("Invalid time: %s", matches[2])
      }
      ch <- LogEvent{dt: dt, fields: matches[1:]}
    }
  }()
}

