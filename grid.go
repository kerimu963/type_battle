package main

import (
	"image/color"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type CellState int

const (
	superEffectiveChangeChance = 1.0
	normalChangeChance         = 1.0 / 4.0
	notVeryEffectiveChance     = 1.0 / 8.0
	noEffectChangeChance       = 0.0
)

const (
	CellNormal CellState = iota
	CellFire
	CellWater
	CellElectric
	CellGrass
	CellIce
	CellFighting
	CellPoison
	CellGround
	CellFlying
	CellPsychic
	CellBug
	CellRock
	CellGhost
	CellDragon
	CellDark
	CellSteel
	CellFairy
	cellStateCount
)

type Cell struct {
	State CellState
}

type attackEffect int

const (
	effectNone attackEffect = iota
	effectNotVeryEffective
	effectNormal
	effectSuperEffective
)

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

func (g *Grid) SetBounds(x, y, width, height int) {
	g.x = x
	g.y = y
	g.width = width
	g.height = height
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
	states := make([]CellState, cellStateCount)
	for state := range cellStateCount {
		states[state] = CellState(state)
	}
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
	successfulAttacks := make([][]CellState, g.Rows*g.Columns)

	for row := range g.Rows {
		for column := range g.Columns {
			targets := g.neighborsOf(row, column)
			if len(targets) == 0 {
				continue
			}
			target := targets[rand.IntN(len(targets))]
			attackerState := previous[row][column].State
			targetState := previous[target.row][target.column].State

			chance := changeChance(attackEffectiveness(attackerState, targetState))
			if rand.Float64() < chance {
				targetIndex := target.row*g.Columns + target.column
				successfulAttacks[targetIndex] = append(successfulAttacks[targetIndex], attackerState)
			}
		}
	}

	for targetIndex, attacks := range successfulAttacks {
		if len(attacks) == 0 {
			continue
		}
		row := targetIndex / g.Columns
		column := targetIndex % g.Columns
		next[row][column].State = attacks[rand.IntN(len(attacks))]
	}

	g.Cells = next
}

func attackEffectiveness(attacker, target CellState) attackEffect {
	if hasNoEffect(attacker, target) {
		return effectNone
	}
	if isSuperEffective(attacker, target) {
		return effectSuperEffective
	}
	if isNotVeryEffective(attacker, target) {
		return effectNotVeryEffective
	}
	return effectNormal
}

func changeChance(effect attackEffect) float64 {
	switch effect {
	case effectSuperEffective:
		return superEffectiveChangeChance
	case effectNormal:
		return normalChangeChance
	case effectNotVeryEffective:
		return notVeryEffectiveChance
	case effectNone:
		return noEffectChangeChance
	default:
		return 0
	}
}

func hasNoEffect(attacker, target CellState) bool {
	return attacker == CellNormal && target == CellGhost ||
		attacker == CellElectric && target == CellGround ||
		attacker == CellFighting && target == CellGhost ||
		attacker == CellPoison && target == CellSteel ||
		attacker == CellGround && target == CellFlying ||
		attacker == CellPsychic && target == CellDark ||
		attacker == CellGhost && target == CellNormal ||
		attacker == CellDragon && target == CellFairy
}

func isNotVeryEffective(attacker, target CellState) bool {
	switch attacker {
	case CellNormal:
		return target == CellRock || target == CellSteel
	case CellFire:
		return target == CellFire || target == CellWater || target == CellRock || target == CellDragon
	case CellWater:
		return target == CellWater || target == CellGrass || target == CellDragon
	case CellElectric:
		return target == CellElectric || target == CellGrass || target == CellDragon
	case CellGrass:
		return target == CellFire || target == CellGrass || target == CellPoison || target == CellFlying ||
			target == CellBug || target == CellDragon || target == CellSteel
	case CellIce:
		return target == CellFire || target == CellWater || target == CellIce || target == CellSteel
	case CellFighting:
		return target == CellPoison || target == CellFlying || target == CellPsychic || target == CellBug || target == CellFairy
	case CellPoison:
		return target == CellPoison || target == CellGround || target == CellRock || target == CellGhost
	case CellGround:
		return target == CellGrass || target == CellBug
	case CellFlying:
		return target == CellElectric || target == CellRock || target == CellSteel
	case CellPsychic:
		return target == CellPsychic || target == CellSteel
	case CellBug:
		return target == CellFire || target == CellFighting || target == CellPoison || target == CellFlying ||
			target == CellGhost || target == CellSteel || target == CellFairy
	case CellRock:
		return target == CellFighting || target == CellGround || target == CellSteel
	case CellGhost:
		return target == CellDark
	case CellDragon:
		return target == CellSteel
	case CellDark:
		return target == CellFighting || target == CellDark || target == CellFairy
	case CellSteel:
		return target == CellFire || target == CellWater || target == CellElectric || target == CellSteel
	case CellFairy:
		return target == CellFire || target == CellPoison || target == CellSteel
	default:
		return false
	}
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
	switch attacker {
	case CellNormal:
		return false
	case CellFire:
		return target == CellGrass || target == CellIce || target == CellBug || target == CellSteel
	case CellWater:
		return target == CellFire || target == CellGround || target == CellRock
	case CellElectric:
		return target == CellWater || target == CellFlying
	case CellGrass:
		return target == CellWater || target == CellGround || target == CellRock
	case CellIce:
		return target == CellGrass || target == CellGround || target == CellFlying || target == CellDragon
	case CellFighting:
		return target == CellNormal || target == CellIce || target == CellRock || target == CellDark || target == CellSteel
	case CellPoison:
		return target == CellGrass || target == CellFairy
	case CellGround:
		return target == CellFire || target == CellElectric || target == CellPoison || target == CellRock || target == CellSteel
	case CellFlying:
		return target == CellGrass || target == CellFighting || target == CellBug
	case CellPsychic:
		return target == CellFighting || target == CellPoison
	case CellBug:
		return target == CellGrass || target == CellPsychic || target == CellDark
	case CellRock:
		return target == CellFire || target == CellIce || target == CellFlying || target == CellBug
	case CellGhost:
		return target == CellPsychic || target == CellGhost
	case CellDragon:
		return target == CellDragon
	case CellDark:
		return target == CellPsychic || target == CellGhost
	case CellSteel:
		return target == CellIce || target == CellRock || target == CellFairy
	case CellFairy:
		return target == CellFighting || target == CellDragon || target == CellDark
	default:
		return false
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

func (g *Grid) StateCounts() [cellStateCount]int {
	var counts [cellStateCount]int
	for row := range g.Rows {
		for column := range g.Columns {
			counts[g.Cells[row][column].State]++
		}
	}
	return counts
}

func cellColor(state CellState) color.Color {
	switch state {
	case CellNormal:
		return color.RGBA{R: 0xA8, G: 0xA7, B: 0x7A, A: 255}
	case CellFire:
		return color.RGBA{R: 0xEE, G: 0x81, B: 0x30, A: 255}
	case CellWater:
		return color.RGBA{R: 0x63, G: 0x90, B: 0xF0, A: 255}
	case CellElectric:
		return color.RGBA{R: 0xF7, G: 0xD0, B: 0x2C, A: 255}
	case CellGrass:
		return color.RGBA{R: 0x7A, G: 0xC7, B: 0x4C, A: 255}
	case CellIce:
		return color.RGBA{R: 0x96, G: 0xD9, B: 0xD6, A: 255}
	case CellFighting:
		return color.RGBA{R: 0xC2, G: 0x2E, B: 0x28, A: 255}
	case CellPoison:
		return color.RGBA{R: 0xA3, G: 0x3E, B: 0xA1, A: 255}
	case CellGround:
		return color.RGBA{R: 0xE2, G: 0xBF, B: 0x65, A: 255}
	case CellFlying:
		return color.RGBA{R: 0xA9, G: 0x8F, B: 0xF3, A: 255}
	case CellPsychic:
		return color.RGBA{R: 0xF9, G: 0x55, B: 0x87, A: 255}
	case CellBug:
		return color.RGBA{R: 0xA6, G: 0xB9, B: 0x1A, A: 255}
	case CellRock:
		return color.RGBA{R: 0xB6, G: 0xA1, B: 0x36, A: 255}
	case CellGhost:
		return color.RGBA{R: 0x73, G: 0x57, B: 0x97, A: 255}
	case CellDragon:
		return color.RGBA{R: 0x6F, G: 0x35, B: 0xFC, A: 255}
	case CellDark:
		return color.RGBA{R: 0x70, G: 0x57, B: 0x46, A: 255}
	case CellSteel:
		return color.RGBA{R: 0xB7, G: 0xB7, B: 0xCE, A: 255}
	case CellFairy:
		return color.RGBA{R: 0xD6, G: 0x85, B: 0xAD, A: 255}
	default:
		return color.White
	}
}

func cellStateName(state CellState) string {
	names := [...]string{
		"ノーマル", "ほのお", "みず", "でんき", "くさ", "こおり",
		"かくとう", "どく", "じめん", "ひこう", "エスパー", "むし",
		"いわ", "ゴースト", "ドラゴン", "あく", "はがね", "フェアリー",
	}
	if state < 0 || state >= cellStateCount {
		return "不明"
	}
	return names[state]
}
