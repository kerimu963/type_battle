package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	menuWidth             = 200
	gridWidth             = 600
	canvasWidth           = menuWidth + gridWidth
	canvasHeight          = 600
	gridColumns           = 250
	gridRows              = 200
	autoPlayIntervalTicks = 6
)

type Game struct {
	grid            *Grid
	time            int
	autoPlay        bool
	autoPlayCounter int
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
		g.autoPlayCounter = 0
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if isAutoPlayButtonAt(x, y) {
			g.autoPlay = !g.autoPlay
			g.autoPlayCounter = 0
		}
	}

	if g.autoPlay {
		g.autoPlayCounter++
		if g.autoPlayCounter >= autoPlayIntervalTicks {
			g.advanceTime()
			g.autoPlayCounter = 0
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.grid.Draw(screen)
	drawMenu(screen, g.grid, g.time, g.autoPlay)
}

func (g *Game) advanceTime() {
	g.grid.ResolveAttacks()
	g.time++
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth <= menuWidth || outsideHeight <= 0 {
		return canvasWidth, canvasHeight
	}
	g.grid.SetBounds(menuWidth, 0, outsideWidth-menuWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowTitle("Type Battle")
	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
