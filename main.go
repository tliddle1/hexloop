package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tliddle1/hexloop/game"
)

func main() {
	game := game.NewGame()
	ebiten.SetWindowSize(game.ScreenWidth, game.ScreenHeight)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
