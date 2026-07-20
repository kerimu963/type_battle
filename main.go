package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	menuWidth    = 200
	gridWidth    = 800
	canvasWidth  = menuWidth + gridWidth
	canvasHeight = 600
	gridColumns  = 10
	gridRows     = 8
)

type Game struct {
	grid *Grid
	time int
}

func newGame() *Game {
	return &Game{
		grid: NewGrid(gridColumns, gridRows, menuWidth, 0, gridWidth, canvasHeight),
		time: 1,
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.advanceTime()
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.grid.ToggleAt(x, y)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.grid.Draw(screen)
	drawMenu(screen, g.grid, g.time)
}

func (g *Game) advanceTime() {
	g.time++
}

func (g *Game) Layout(_, _ int) (int, int) {
	return canvasWidth, canvasHeight
}

func main() {
	ebiten.SetWindowSize(canvasWidth, canvasHeight)
	ebiten.SetWindowTitle("Type Battle")

	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
