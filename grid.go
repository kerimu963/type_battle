package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CellState int

const (
	CellEmpty CellState = iota
	CellActive
)

type Cell struct {
	State CellState
}

type Grid struct {
	Columns int
	Rows    int
	Cells   [][]Cell
	x       int
	y       int
	width   int
	height  int
}

func NewGrid(columns, rows, x, y, width, height int) *Grid {
	cells := make([][]Cell, rows)
	for row := range cells {
		cells[row] = make([]Cell, columns)
	}

	return &Grid{
		Columns: columns,
		Rows:    rows,
		Cells:   cells,
		x:       x,
		y:       y,
		width:   width,
		height:  height,
	}
}

func (g *Grid) ToggleAt(x, y int) {
	if x < g.x || x >= g.x+g.width || y < g.y || y >= g.y+g.height {
		return
	}

	column := (x - g.x) * g.Columns / g.width
	row := (y - g.y) * g.Rows / g.height
	cell := &g.Cells[row][column]

	if cell.State == CellEmpty {
		cell.State = CellActive
	} else {
		cell.State = CellEmpty
	}
}

func (g *Grid) Draw(screen *ebiten.Image) {
	cellWidth := float32(g.width) / float32(g.Columns)
	cellHeight := float32(g.height) / float32(g.Rows)
	gridX := float32(g.x)
	gridY := float32(g.y)
	vector.FillRect(screen, gridX, gridY, float32(g.width), float32(g.height), color.White, false)

	for row := range g.Rows {
		for column := range g.Columns {
			if g.Cells[row][column].State == CellActive {
				x := gridX + float32(column)*cellWidth
				y := gridY + float32(row)*cellHeight
				vector.FillRect(screen, x, y, cellWidth, cellHeight, color.RGBA{R: 80, G: 160, B: 255, A: 255}, false)
			}
		}
	}

	gridColor := color.RGBA{R: 80, G: 80, B: 80, A: 255}
	for column := 0; column <= g.Columns; column++ {
		x := gridX + float32(column)*cellWidth
		vector.StrokeLine(screen, x, gridY, x, gridY+float32(g.height), 1, gridColor, false)
	}
	for row := 0; row <= g.Rows; row++ {
		y := gridY + float32(row)*cellHeight
		vector.StrokeLine(screen, gridX, y, gridX+float32(g.width), y, 1, gridColor, false)
	}
}

func (g *Grid) ActiveCellCount() int {
	count := 0
	for row := range g.Rows {
		for column := range g.Columns {
			if g.Cells[row][column].State == CellActive {
				count++
			}
		}
	}
	return count
}
