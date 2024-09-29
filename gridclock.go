package main

import (
	"fmt"
	"image"
	"strconv"
	"time"

	ui "github.com/gizak/termui/v3"
)

type GridClock struct {
	ui.Block
	Mode12 bool
}

func NewGridClock(height, width int, mode12 bool) *GridClock {
	c := &GridClock{Block: *ui.NewBlock(), Mode12: mode12}
	c.SetRect(0, 0, width, height)
	return c
}

func (c *GridClock) Resize(width, height int) {
	c.SetRect(0, 0, width, height)
	ui.Clear()
	ui.Render(c)
}

func (c *GridClock) Tick() {
	ui.Render(c)
}

func drawVLine(buf *ui.Buffer, rune rune, y1, y2, x int) {
	for i := y1; i < y2-1; i++ {
		// skip drawing on border
		if x == 0 {
			continue
		}
		buf.SetCell(ui.NewCell(rune), image.Point{x, i})
	}
}

func drawHLine(buf *ui.Buffer, rune rune, x1, x2, y int) {
	for j := x1; j < x2-1; j++ {
		// skip drawing on border
		if y == 0 {
			continue
		}
		buf.SetCell(ui.NewCell(rune), image.Point{j, y})
	}
}

func (c *GridClock) Draw(buf *ui.Buffer) {
	var left, top int
	now := time.Now()

	c.Block.Draw(buf)

	title := fmt.Sprintf("%s", now.Format(CLOCK_FORMAT))
	buf.SetString(title, ui.NewStyle(ui.ColorWhite), image.Point{c.Min.X + 2, c.Max.Y - 1})

	// hour line
	hour := now.Hour()
	hourY := c.Max.Y * ((hour % 12) + 1) / (12 + 1)
	drawHLine(buf, HLINE, left+3, c.Max.X, hourY)

	var hourStr string
	if hour > 12 && !c.Mode12 {
		hourStr = strconv.Itoa(hour)
	} else {
		if (hour%12) < 10 && (hour%12) > 0 {
			hourStr = fmt.Sprintf("0%d", hour)
		} else {
			hourStr = strconv.Itoa(hour)
		}
	}

	buf.SetString(hourStr, ui.NewStyle(ui.ColorWhite), image.Point{left, hourY})
	top += hourY

	// minute lines
	minute := now.Minute()
	minX := (c.Max.X - left) * ((minute / 5) + 1) / (12 + 1)
	minY := (c.Max.Y - top) * ((minute % 5) + 1) / (5 + 1)

	drawHLine(buf, HLINE, left+minX, c.Max.X, top+minY)
	drawVLine(buf, VLINE, top, c.Max.Y, left+minX)
	drawVLine(buf, VLINE_DASHED, 1, top+1, left+minX)

	buf.SetCell(ui.NewCell(RCORNER), image.Point{left + minX, top + minY})
	buf.SetCell(ui.NewCell(XCORNER), image.Point{left + minX, top})

	var minStr string
	if minute < 10 {
		minStr = fmt.Sprintf("0%d", minute)
	} else {
		minStr = strconv.Itoa(minute)
	}
	buf.SetString(minStr, ui.NewStyle(ui.ColorWhite), image.Point{left + minX, 0})

	left += minX
	top += minY

	// second lines
	second := now.Second()
	secX := (c.Max.X - left) * ((second / 5) + 1) / (12 + 1)
	secY := (c.Max.Y - top) * ((second % 5) + 1) / (5 + 1)

	drawHLine(buf, HLINE, left+secX, c.Max.X, top+secY)
	drawVLine(buf, VLINE, top, c.Max.Y, left+secX)

	buf.SetCell(ui.NewCell(RCORNER), image.Point{left + secX, top + secY})
	buf.SetCell(ui.NewCell(TCORNER), image.Point{left + secX, top})

	var secStr string
	if second < 10 {
		secStr = fmt.Sprintf("0%d", second)
	} else {
		secStr = strconv.Itoa(second)
	}
	buf.SetString(secStr, ui.NewStyle(ui.ColorWhite), image.Point{left - 2, top})
}
