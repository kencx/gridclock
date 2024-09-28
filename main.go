package main

import (
	"fmt"
	"image"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
)

const (
	// https://symbl.cc/en/unicode/blocks/box-drawing/
	HLINE   = '━'
	VLINE   = '┃'
	rcorner = '┏'
	lcorner = '┫'
	tcorner = '┳'
	bcorner = '┻'
	xcorner = '╋'

	CLOCK_FORMAT = "15:04:05"
)

type Clock struct {
	ui.Block
	internalClock *Clock

	Height, Width int

	// no of partitions
	xpart, ypart int

	// starting offset
	xstart, ystart int
	// running offset
	xoffset, yoffset int

	// size of each step
	xstep, ystep int

	// time interval between steps
	interval time.Duration
	counter  int
}

func NewClock(height, width, xpart, ypart int, interval time.Duration) *Clock {
	x := width / xpart
	y := height / ypart

	c := &Clock{
		Block: *ui.NewBlock(),

		Height: height,
		Width:  width,

		xpart: xpart,
		ypart: ypart,

		// spaces out the partitions more evenly
		xstart: x/2 + 1,
		ystart: y/2 + 1,
		// xstart: int(math.Round(float64(x)/2) + 0.5),
		// ystart: int(math.Round(float64(y)/2) + 0.5),

		xstep: x,
		ystep: y,

		interval: interval,
		counter:  0,
	}

	c.xoffset = c.xstart
	c.yoffset = c.ystart
	c.Title = fmt.Sprintf("%s width: %d, height: %d; xoffset: %d; yoffset: %d; elapsed: %d",
		time.Now().Format(CLOCK_FORMAT),
		c.Width,
		c.Height,
		c.xoffset,
		c.yoffset,
		c.counter,
	)

	c.SetRect(0, 0, width, height)
	return c
}

func (c *Clock) Resize(width, height int) {
	x := width / c.xpart
	y := height / c.ypart

	c.Height = height
	c.Width = width

	c.xstart = x/2 + 1
	c.ystart = y/2 + 1

	c.xstep = x
	c.ystep = y

	c.SetRect(0, 0, width, height)
	ui.Clear()
	ui.Render(c)
}

func (c *Clock) Tick() {
	if c.ystep < c.Height {
		c.yoffset += c.ystep
	}

	if c.yoffset >= c.Height-1 {
		c.yoffset = c.ystart

		if c.xstep < c.Width {
			c.xoffset += c.xstep
		}

		if c.xoffset >= c.Width-1 {
			c.xoffset = c.xstart
		}
	}

	c.counter += 1
	c.Title = fmt.Sprintf("%s width: %d, height: %d; xoffset: %d; yoffset: %d; elapsed: %d",
		time.Now().Format(CLOCK_FORMAT),
		c.Width,
		c.Height,
		c.xoffset,
		c.yoffset,
		c.counter,
	)
	ui.Render(c)
}

func (c *Clock) Draw(buf *ui.Buffer) {
	c.Block.Draw(buf)

	// draw vertical
	for i := 1; i < c.Height-1; i++ {
		// skip drawing on border
		if c.xoffset == 0 {
			continue
		}
		buf.SetCell(ui.NewCell(VLINE), image.Point{c.xoffset, i})
	}

	// draw horizontal
	for j := 1; j < c.Width-1; j++ {
		// skip drawing on border
		if c.yoffset == 0 {
			continue
		}
		buf.SetCell(ui.NewCell(HLINE), image.Point{j, c.yoffset})
	}

	if c.xoffset != 0 && c.yoffset != 0 {
		buf.SetCell(ui.NewCell(xcorner), image.Point{c.xoffset, c.yoffset})
	}
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to start: %v", err)
	}
	defer ui.Close()

	width, height := ui.TerminalDimensions()
	clock := NewClock(height, width, 60, 12, time.Hour)
	ui.Render(clock)

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				clock.Resize(payload.Width, payload.Height)
			}
		case <-ticker:
			clock.Tick()
		}
	}
}
