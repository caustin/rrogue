package main

import (
	"github.com/caustin/rrogue/game"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	g := game.NewGame()
	ebiten.SetWindowResizable(true)

	ebiten.SetWindowTitle("Tower")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
