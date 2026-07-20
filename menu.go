package main

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	autoPlayButtonX      = 20
	autoPlayButtonY      = 530
	autoPlayButtonWidth  = 160
	autoPlayButtonHeight = 28
	fastPlayButtonX      = 20
	fastPlayButtonY      = 565
	fastPlayButtonWidth  = 160
	fastPlayButtonHeight = 28
)

var menuFontFace = func() *text.GoTextFace {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		panic(err)
	}
	return &text.GoTextFace{Source: source, Size: 12}
}()

func drawMenu(screen *ebiten.Image, grid *Grid, currentTime int, autoPlay, fastPlay bool) {
	background := color.RGBA{R: 32, G: 37, B: 48, A: 255}
	separator := color.RGBA{R: 90, G: 100, B: 120, A: 255}
	screenHeight := float32(screen.Bounds().Dy())

	vector.FillRect(screen, 0, 0, menuWidth, screenHeight, background, false)
	vector.StrokeLine(screen, menuWidth-1, 0, menuWidth-1, screenHeight, 2, separator, false)

	menuText := fmt.Sprintf(
		"TYPE BATTLE\n\nTIME: %d\n\nGrid: %d x %d",
		currentTime,
		grid.Columns,
		grid.Rows,
	)
	ebitenutil.DebugPrintAt(screen, menuText, 20, 20)

	counts := grid.StateCounts()
	for state := CellNormal; state < cellStateCount; state++ {
		drawStateLegend(screen, 108+float64(state)*22, state, cellStateName(state), counts[state])
	}
	ebitenutil.DebugPrintAt(screen, "Play: 10/s  Fast: 60/s", 20, 508)

	buttonColor := color.RGBA{R: 75, G: 85, B: 105, A: 255}
	buttonText := "AUTO PLAY: OFF"
	if autoPlay {
		buttonColor = color.RGBA{R: 45, G: 150, B: 85, A: 255}
		buttonText = "AUTO PLAY: ON"
	}
	vector.FillRect(screen, autoPlayButtonX, autoPlayButtonY, autoPlayButtonWidth, autoPlayButtonHeight, buttonColor, false)
	vector.StrokeRect(screen, autoPlayButtonX, autoPlayButtonY, autoPlayButtonWidth, autoPlayButtonHeight, 1, separator, false)
	ebitenutil.DebugPrintAt(screen, buttonText, autoPlayButtonX+25, autoPlayButtonY+9)

	fastButtonColor := color.RGBA{R: 75, G: 85, B: 105, A: 255}
	fastButtonText := "FAST PLAY: OFF"
	if fastPlay {
		fastButtonColor = color.RGBA{R: 205, G: 90, B: 50, A: 255}
		fastButtonText = "FAST PLAY: ON"
	}
	vector.FillRect(screen, fastPlayButtonX, fastPlayButtonY, fastPlayButtonWidth, fastPlayButtonHeight, fastButtonColor, false)
	vector.StrokeRect(screen, fastPlayButtonX, fastPlayButtonY, fastPlayButtonWidth, fastPlayButtonHeight, 1, separator, false)
	ebitenutil.DebugPrintAt(screen, fastButtonText, fastPlayButtonX+23, fastPlayButtonY+9)
}

func drawStateLegend(screen *ebiten.Image, y float64, state CellState, name string, count int) {
	vector.FillRect(screen, 20, float32(y), 16, 16, cellColor(state), false)
	vector.StrokeRect(screen, 20, float32(y), 16, 16, 1, color.White, false)

	op := &text.DrawOptions{}
	op.GeoM.Translate(44, y)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, fmt.Sprintf("%s: %d", name, count), menuFontFace, op)
}

func isAutoPlayButtonAt(x, y int) bool {
	return x >= autoPlayButtonX && x < autoPlayButtonX+autoPlayButtonWidth &&
		y >= autoPlayButtonY && y < autoPlayButtonY+autoPlayButtonHeight
}

func isFastPlayButtonAt(x, y int) bool {
	return x >= fastPlayButtonX && x < fastPlayButtonX+fastPlayButtonWidth &&
		y >= fastPlayButtonY && y < fastPlayButtonY+fastPlayButtonHeight
}
