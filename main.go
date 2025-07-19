package main

import (
	_ "image/png"
	"log"

	"github.com/bytearena/ecs"
	"github.com/hajimehoshi/ebiten/v2"
)

// Game holds all data the entire game will need.
type Game struct {
	Map           GameMap
	World         *ecs.Manager
	WorldTags     map[string]ecs.Tag
	Components    *Components
	GameData      GameData
	Turn          TurnState
	TurnCounter   int
	AutoMoveState *AutoMoveState
}

// NewGame creates a new Game Object and initializes the data
// This is a pretty solid refactor candidate for later
func NewGame() *Game {
	g := &Game{}
	g.Map = NewGameMap()
	g.GameData = NewGameData()
	world, tags, components := InitializeWorld(g.Map.CurrentLevel)

	g.WorldTags = tags
	g.World = world
	g.Components = components
	g.Turn = WaitingForPlayerInput
	g.TurnCounter = 0
	return g

}

// Update is called each tic.
func (g *Game) Update() error {
	switch g.Turn {
	case WaitingForPlayerInput:
		if TakePlayerAction(g) {
			g.Turn = ProcessingMonsterTurn
			g.TurnCounter++
		}
	case ProcessingMonsterTurn:
		UpdateMonster(g)
		g.Turn = WaitingForPlayerInput
	default:
		panic("unhandled default case")
	}

	return nil

}

// Draw is called each draw cycle and is where we will blit.
func (g *Game) Draw(screen *ebiten.Image) {
	//Draw the Map
	level := g.Map.CurrentLevel
	level.DrawLevel(screen, g.GameData)
	ProcessRenderables(g, level, screen)
	ProcessUserLog(g, screen)
	ProcessHUD(g, screen)

}

// Layout will return the screen dimensions.
func (g *Game) Layout(w, h int) (int, int) {
	return g.GameData.TileWidth * g.GameData.ScreenWidth, g.GameData.TileHeight * g.GameData.ScreenHeight

}

func main() {

	g := NewGame()
	ebiten.SetWindowResizable(true)

	ebiten.SetWindowTitle("Tower")

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
