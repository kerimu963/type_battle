package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	canvasWidth  = 800
	canvasHeight = 600
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return canvasWidth, canvasHeight
}

func main() {
	ebiten.SetWindowSize(canvasWidth, canvasHeight)
	ebiten.SetWindowTitle("Type Battle")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
