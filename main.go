package main

import (
	"flag"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
)

const (
	// https://symbl.cc/en/unicode/blocks/box-drawing/
	HLINE        = '━'
	VLINE        = '┃'
	VLINE_DASHED = '┇'
	RCORNER      = '┣'
	TCORNER      = '┳'
	XCORNER      = '╋'

	CLOCK_FORMAT = "2006-01-02 15:04:05"
)

type Clock interface {
	Tick()
	Resize(width, height int)
	ui.Drawable
}

func main() {
	var mode12 bool
	var stopwatch bool

	flag.BoolVar(&mode12, "12", false, "12h clock")
	flag.BoolVar(&stopwatch, "s", false, "stopwatch mode")
	flag.Parse()

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to start: %v", err)
	}
	defer ui.Close()

	var (
		clock  Clock
		ticker *time.Ticker
	)

	width, height := ui.TerminalDimensions()
	if !stopwatch {
		clock = NewGridClock(height, width, mode12)
		ticker = time.NewTicker(time.Second)
	} else {
		start := time.Now()
		ticker = time.NewTicker(time.Millisecond)
		clock = NewStopWatch(height, width, start, *ticker)
	}

	ui.Render(clock)
	uiEvents := ui.PollEvents()

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "s":
				if stopwatch {
					ticker.Stop()
				}
			case "r":
				if stopwatch {
					ticker.Reset(time.Millisecond)
				}
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				clock.Resize(payload.Width, payload.Height)
			}
		case <-ticker.C:
			clock.Tick()
		}
	}
}
