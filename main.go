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
	gridColumns           = 200
	gridRows              = 100
	autoPlayIntervalTicks = 6
	fastPlayIntervalTicks = 1
	chartHistoryLimit     = 300
)

type Game struct {
	grid            *Grid
	secondaryCanvas *Canvas
	time            int
	autoPlay        bool
	fastPlay        bool
	autoPlayCounter int
}

func newGame() *Game {
	game := &Game{
		grid:            NewGrid(gridColumns, gridRows, menuWidth, 0, gridWidth, canvasHeight),
		secondaryCanvas: NewCanvas(menuWidth, canvasHeight*3/4, gridWidth, canvasHeight/4),
		time:            1,
	}
	game.secondaryCanvas.AddSnapshot(game.time, game.grid.StateCounts())
	return game
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
			g.fastPlay = false
			g.autoPlayCounter = 0
		} else if isFastPlayButtonAt(x, y) {
			g.fastPlay = !g.fastPlay
			g.autoPlay = false
			g.autoPlayCounter = 0
		}
	}

	if g.autoPlay || g.fastPlay {
		g.autoPlayCounter++
		interval := autoPlayIntervalTicks
		if g.fastPlay {
			interval = fastPlayIntervalTicks
		}
		if g.autoPlayCounter >= interval {
			g.advanceTime()
			g.autoPlayCounter = 0
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.grid.Draw(screen)
	g.secondaryCanvas.Draw(screen)
	drawMenu(screen, g.grid, g.time, g.autoPlay, g.fastPlay)
}

func (g *Game) advanceTime() {
	g.grid.ResolveAttacks()
	g.time++
	g.secondaryCanvas.AddSnapshot(g.time, g.grid.StateCounts())
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if outsideWidth <= menuWidth || outsideHeight <= 0 {
		return canvasWidth, canvasHeight
	}
	secondaryHeight := outsideHeight / 4
	gridHeight := outsideHeight - secondaryHeight
	g.grid.SetBounds(menuWidth, 0, outsideWidth-menuWidth, gridHeight)
	g.secondaryCanvas.SetBounds(menuWidth, gridHeight, outsideWidth-menuWidth, secondaryHeight)
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowTitle("Type Battle")
	ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(newGame()); err != nil {
		log.Fatal(err)
	}
}
