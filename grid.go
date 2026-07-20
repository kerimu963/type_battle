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
	width   int
	height  int
}

func NewGrid(columns, rows, width, height int) *Grid {
	cells := make([][]Cell, rows)
	for row := range cells {
		cells[row] = make([]Cell, columns)
	}

	return &Grid{
		Columns: columns,
		Rows:    rows,
		Cells:   cells,
		width:   width,
		height:  height,
	}
}

func (g *Grid) ToggleAt(x, y int) {
	if x < 0 || x >= g.width || y < 0 || y >= g.height {
		return
	}

	column := x * g.Columns / g.width
	row := y * g.Rows / g.height
	cell := &g.Cells[row][column]

	if cell.State == CellEmpty {
		cell.State = CellActive
	} else {
		cell.State = CellEmpty
	}
}

func (g *Grid) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	cellWidth := float32(g.width) / float32(g.Columns)
	cellHeight := float32(g.height) / float32(g.Rows)

	for row := range g.Rows {
		for column := range g.Columns {
			if g.Cells[row][column].State == CellActive {
				x := float32(column) * cellWidth
				y := float32(row) * cellHeight
				vector.FillRect(screen, x, y, cellWidth, cellHeight, color.RGBA{R: 80, G: 160, B: 255, A: 255}, false)
			}
		}
	}

	gridColor := color.RGBA{R: 80, G: 80, B: 80, A: 255}
	for column := 0; column <= g.Columns; column++ {
		x := float32(column) * cellWidth
		vector.StrokeLine(screen, x, 0, x, float32(g.height), 1, gridColor, false)
	}
	for row := 0; row <= g.Rows; row++ {
		y := float32(row) * cellHeight
		vector.StrokeLine(screen, 0, y, float32(g.width), y, 1, gridColor, false)
	}
}
