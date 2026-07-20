package main

import (
	"image/color"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CellState int

const (
	CellFire CellState = iota
	CellWater
	CellGrass
)

type Cell struct {
	State CellState
}

type cellPosition struct {
	row    int
	column int
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
	initialCells := generateInitialCells(columns * rows)
	cells := make([][]Cell, rows)
	for row := range cells {
		cells[row] = make([]Cell, columns)
		for column := range cells[row] {
			cells[row][column] = initialCells[row*columns+column]
		}
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

func generateInitialCells(cellCount int) []Cell {
	states := []CellState{CellFire, CellWater, CellGrass}
	rand.Shuffle(len(states), func(i, j int) {
		states[i], states[j] = states[j], states[i]
	})

	counts := make(map[CellState]int, len(states))
	baseCount := cellCount / len(states)
	for _, state := range states {
		counts[state] = baseCount
	}
	for i := 0; i < cellCount%len(states); i++ {
		counts[states[i]]++
	}

	cells := make([]Cell, 0, cellCount)
	for _, state := range states {
		for range counts[state] {
			cells = append(cells, Cell{State: state})
		}
	}

	rand.Shuffle(len(cells), func(i, j int) {
		cells[i], cells[j] = cells[j], cells[i]
	})

	return cells
}

func (g *Grid) ResolveAttacks() {
	previous := g.copyCells()
	next := g.copyCells()

	for row := range g.Rows {
		for column := range g.Columns {
			targets := g.neighborsOf(row, column)
			if len(targets) == 0 {
				continue
			}
			target := targets[rand.IntN(len(targets))]
			attackerState := previous[row][column].State
			targetState := previous[target.row][target.column].State

			if isSuperEffective(attackerState, targetState) {
				next[target.row][target.column].State = attackerState
			}
		}
	}

	g.Cells = next
}

func (g *Grid) copyCells() [][]Cell {
	cells := make([][]Cell, g.Rows)
	for row := range g.Rows {
		cells[row] = make([]Cell, g.Columns)
		copy(cells[row], g.Cells[row])
	}
	return cells
}

func (g *Grid) neighborsOf(row, column int) []cellPosition {
	neighbors := make([]cellPosition, 0, 4)
	if row > 0 {
		neighbors = append(neighbors, cellPosition{row: row - 1, column: column})
	}
	if row+1 < g.Rows {
		neighbors = append(neighbors, cellPosition{row: row + 1, column: column})
	}
	if column > 0 {
		neighbors = append(neighbors, cellPosition{row: row, column: column - 1})
	}
	if column+1 < g.Columns {
		neighbors = append(neighbors, cellPosition{row: row, column: column + 1})
	}
	return neighbors
}

func isSuperEffective(attacker, target CellState) bool {
	return attacker == CellFire && target == CellGrass ||
		attacker == CellWater && target == CellFire ||
		attacker == CellGrass && target == CellWater
}

func (g *Grid) Draw(screen *ebiten.Image) {
	cellWidth := float32(g.width) / float32(g.Columns)
	cellHeight := float32(g.height) / float32(g.Rows)
	gridX := float32(g.x)
	gridY := float32(g.y)
	vector.FillRect(screen, gridX, gridY, float32(g.width), float32(g.height), color.White, false)

	for row := range g.Rows {
		for column := range g.Columns {
			x := gridX + float32(column)*cellWidth
			y := gridY + float32(row)*cellHeight
			vector.FillRect(screen, x, y, cellWidth, cellHeight, cellColor(g.Cells[row][column].State), false)
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

func (g *Grid) CountState(state CellState) int {
	count := 0
	for row := range g.Rows {
		for column := range g.Columns {
			if g.Cells[row][column].State == state {
				count++
			}
		}
	}
	return count
}

func cellColor(state CellState) color.Color {
	switch state {
	case CellFire:
		return color.RGBA{R: 230, G: 70, B: 60, A: 255}
	case CellWater:
		return color.RGBA{R: 60, G: 130, B: 230, A: 255}
	case CellGrass:
		return color.RGBA{R: 70, G: 180, B: 90, A: 255}
	default:
		return color.White
	}
}
