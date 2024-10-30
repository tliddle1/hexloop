package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	squareMode := false
	game := NewGame(squareMode)
	ebiten.SetWindowSize(game.screenWidth, game.screenHeight)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
