package main

import "github.com/caustin/rrogue/level"

// GameMap holds all the level and aggregate information for the entire world.
type GameMap struct {
	Dungeons     []level.Dungeon
	CurrentLevel level.Level
}

// NewGameMap creates a new set of maps for the entire game.
func NewGameMap() GameMap {
	//Return a new game map of a single level for now
	l := level.NewLevel()
	levels := make([]level.Level, 0)
	levels = append(levels, l)
	d := level.Dungeon{Name: "default", Levels: levels}
	dungeons := make([]level.Dungeon, 0)
	dungeons = append(dungeons, d)
	gm := GameMap{Dungeons: dungeons, CurrentLevel: l}
	return gm

}
