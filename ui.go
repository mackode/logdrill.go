package main

import (
  t "github.com/gizak/termui/v3"
  "github.com/gizak/termui/v3/widgets"
  "log"
)

func uiStart(rowCh chan string, geoCh chan string) {
  err := t.Init()
  if err != nil {
    log.Fatalln("Termui init failed")
  }
  defer t.Close()
  lb := widgets.NewList()
  lb.Title = "Logdrill"
  lb.TextStyle.Fg = t.ColorBlack
  geo := widgets.NewParagraph()
  geo.TextStyle.Fg = t.ColorGreen
  geo.Text = ""

  var width, height int
  listSize := func() int {
    return height - 3
  }
  resize := func() {
    width, height = t.TerminalDimensions()
    lb.SetRect(0, 0, width, listSize())
    geo.SetRect(0, listSize(), width, height)
    t.Render(lb, geo)
  }
  resize()
  uiEvents := t.PollEvents()
  for {
    select {
      case e := <-uiEvents:
        switch e.ID {
          case "<Resize>":
            resize()
          case "q", "<C-c>":
            return
        }
      case line := <-geoCh:
        geo.Text = line
        t.Render(geo)
      case line := <-rowCh:
        if len(lb.Rows) >= listSize() {
          lb.Rows = lb.Rows[0:listSize()]
        }
        lb.Rows = append([]string{line}, lb.Rows...)
        t.Render(lb)
    }
  }
}
