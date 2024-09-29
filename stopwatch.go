package main

import (
	"fmt"
	"image"
	"strconv"
	"time"

	ui "github.com/gizak/termui/v3"
)

type StopWatch struct {
	ui.Block
	start  time.Time
	ticker time.Ticker
}

func NewStopWatch(height, width int, start time.Time, ticker time.Ticker) *StopWatch {
	c := &StopWatch{Block: *ui.NewBlock(), start: start, ticker: ticker}
	c.SetRect(0, 0, width, height)
	return c
}

func (c *StopWatch) Resize(width, height int) {
	c.SetRect(0, 0, width, height)
	ui.Clear()
	ui.Render(c)
}

func (c *StopWatch) Tick() {
	ui.Render(c)
}

func (c *StopWatch) Draw(buf *ui.Buffer) {
	var left, top int
	c.Block.Draw(buf)

	t := <-c.ticker.C
	elapsed := t.Sub(c.start)

	// minute line
	minute := int(elapsed.Minutes())
	minY := c.Max.Y * (minute + 1) / (10 + 1)
	if minY > c.Max.Y {
		minY -= c.Max.Y
	}
	drawHLine(buf, HLINE, left+3, c.Max.X, minY)

	var minStr string
	if minute < 10 {
		minStr = fmt.Sprintf("0%d", minute)
	} else {
		minStr = strconv.Itoa(minute)
	}

	buf.SetString(minStr, ui.NewStyle(ui.ColorWhite), image.Point{left, minY})
	top += minY

	// sec lines
	sec := int(elapsed.Seconds())
	// reset sec after every min
	if sec > 60 {
		sec -= ((sec / 60) * 60)
	}

	secX := (c.Max.X - left) * ((sec / 5) + 1) / (12 + 1)
	secY := (c.Max.Y - top) * ((sec % 5) + 1) / (5 + 1)

	if secX > c.Max.X {
		secX -= c.Max.X
	}
	if secY > c.Max.Y {
		secY -= c.Max.Y
	}

	drawHLine(buf, HLINE, left+secX, c.Max.X, top+secY)
	drawVLine(buf, VLINE, top, c.Max.Y, left+secX)
	drawVLine(buf, VLINE_DASHED, 1, top+1, left+secX)

	buf.SetCell(ui.NewCell(RCORNER), image.Point{left + secX, top + secY})
	buf.SetCell(ui.NewCell(XCORNER), image.Point{left + secX, top})

	var secStr string
	if sec < 10 {
		secStr = fmt.Sprintf("0%d", sec)
	} else {
		secStr = strconv.Itoa(sec)
	}
	buf.SetString(secStr, ui.NewStyle(ui.ColorWhite), image.Point{left + secX, 0})

	left += secX
	top += secY

	// millisecond lines
	millisecond := int(elapsed.Milliseconds())
	msecX := (c.Max.X - left) * ((millisecond % 100) + 1) / (50 + 1)
	msecY := (c.Max.Y - top) * ((millisecond % 100) + 1) / (20 + 1)

	if msecX > c.Max.X {
		msecX -= c.Max.X
	}
	if msecY > c.Max.Y {
		msecY -= c.Max.Y
	}

	drawHLine(buf, HLINE, left+msecX, c.Max.X, top+msecY)
	drawVLine(buf, VLINE, top, c.Max.Y, left+msecX)

	buf.SetCell(ui.NewCell(RCORNER), image.Point{left + msecX, top + msecY})
	buf.SetCell(ui.NewCell(TCORNER), image.Point{left + msecX, top})

	title := fmt.Sprintf("%d m %d s %d ms", minute, sec%60, millisecond%1000)
	buf.SetString(title, ui.NewStyle(ui.ColorWhite), image.Point{c.Min.X + 2, c.Max.Y - 1})
}
