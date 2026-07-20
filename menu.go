package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func drawMenu(screen *ebiten.Image, grid *Grid, currentTime int) {
	background := color.RGBA{R: 32, G: 37, B: 48, A: 255}
	separator := color.RGBA{R: 90, G: 100, B: 120, A: 255}

	vector.FillRect(screen, 0, 0, menuWidth, canvasHeight, background, false)
	vector.StrokeLine(screen, menuWidth-1, 0, menuWidth-1, canvasHeight, 2, separator, false)

	menuText := fmt.Sprintf(
		"TYPE BATTLE\n\nTIME: %d\n\nGrid: %d x %d\nActive: %d\n\nSpace:\nNext TIME\n\nLeft click:\nToggle cell",
		currentTime,
		grid.Columns,
		grid.Rows,
		grid.ActiveCellCount(),
	)
	ebitenutil.DebugPrintAt(screen, menuText, 20, 20)
}
