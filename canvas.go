package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type canvasSnapshot struct {
	time   int
	counts [cellStateCount]int
}

type Canvas struct {
	x       int
	y       int
	width   int
	height  int
	history []canvasSnapshot
}

func NewCanvas(x, y, width, height int) *Canvas {
	return &Canvas{x: x, y: y, width: width, height: height}
}

func (c *Canvas) SetBounds(x, y, width, height int) {
	c.x = x
	c.y = y
	c.width = width
	c.height = height
}

func (c *Canvas) AddSnapshot(currentTime int, counts [cellStateCount]int) {
	c.history = append(c.history, canvasSnapshot{time: currentTime, counts: counts})
	if len(c.history) > chartHistoryLimit {
		c.history = c.history[len(c.history)-chartHistoryLimit:]
	}
}

func (c *Canvas) Draw(screen *ebiten.Image) {
	x := float32(c.x)
	y := float32(c.y)
	width := float32(c.width)
	height := float32(c.height)
	vector.FillRect(screen, x, y, width, height, color.White, false)
	vector.StrokeRect(screen, x, y, width, height, 2, color.Black, false)

	const (
		leftMargin   = float32(50)
		rightMargin  = float32(12)
		topMargin    = float32(12)
		bottomMargin = float32(30)
	)
	plotX := x + leftMargin
	plotY := y + topMargin
	plotWidth := width - leftMargin - rightMargin
	plotHeight := height - topMargin - bottomMargin
	if plotWidth <= 0 || plotHeight <= 0 {
		return
	}

	axisColor := color.RGBA{R: 70, G: 70, B: 70, A: 255}
	vector.StrokeLine(screen, plotX, plotY, plotX, plotY+plotHeight, 1, axisColor, false)
	vector.StrokeLine(screen, plotX, plotY+plotHeight, plotX+plotWidth, plotY+plotHeight, 1, axisColor, false)
	ebitenutil.DebugPrintAt(screen, "Count", c.x+5, c.y+5)
	ebitenutil.DebugPrintAt(screen, "TIME", c.x+c.width/2, c.y+c.height-18)

	if len(c.history) == 0 {
		return
	}

	maxCount := 1
	for _, snapshot := range c.history {
		for _, count := range snapshot.counts {
			if count > maxCount {
				maxCount = count
			}
		}
	}
	firstTime := c.history[0].time
	lastTime := c.history[len(c.history)-1].time
	ebitenutil.DebugPrintAt(screen, "0", c.x+30, int(plotY+plotHeight)-6)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", maxCount), c.x+5, int(plotY)-4)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", firstTime), int(plotX), int(plotY+plotHeight)+4)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", lastTime), int(plotX+plotWidth)-30, int(plotY+plotHeight)+4)

	if len(c.history) < 2 {
		return
	}

	maxSegments := max(1, int(plotWidth))
	step := max(1, (len(c.history)-1+maxSegments-1)/maxSegments)
	for state := CellNormal; state < cellStateCount; state++ {
		previousIndex := 0
		for index := step; index < len(c.history); index += step {
			c.drawHistorySegment(screen, state, previousIndex, index, plotX, plotY, plotWidth, plotHeight, maxCount)
			previousIndex = index
		}
		lastIndex := len(c.history) - 1
		if previousIndex != lastIndex {
			c.drawHistorySegment(screen, state, previousIndex, lastIndex, plotX, plotY, plotWidth, plotHeight, maxCount)
		}
	}
}

func (c *Canvas) drawHistorySegment(screen *ebiten.Image, state CellState, from, to int, x, y, width, height float32, maxCount int) {
	lastIndex := len(c.history) - 1
	x1 := x + float32(from)/float32(lastIndex)*width
	x2 := x + float32(to)/float32(lastIndex)*width
	y1 := y + height - float32(c.history[from].counts[state])/float32(maxCount)*height
	y2 := y + height - float32(c.history[to].counts[state])/float32(maxCount)*height
	vector.StrokeLine(screen, x1, y1, x2, y2, 1, cellColor(state), true)
}
