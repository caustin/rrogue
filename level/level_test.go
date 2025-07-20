package level

import (
	"github.com/caustin/rrogue/config"
	"testing"
)

func TestGetIndexFromXY(t *testing.T) {
	level := Level{}

	tests := []struct {
		name     string
		x, y     int
		expected int
	}{
		{
			name:     "origin (0,0)",
			x:        0,
			y:        0,
			expected: 0,
		},
		{
			name:     "first row position (5,0)",
			x:        5,
			y:        0,
			expected: 5,
		},
		{
			name:     "second row start (0,1)",
			x:        0,
			y:        1,
			expected: 80, // y * ScreenWidth = 1 * 80
		},
		{
			name:     "middle position (10,5)",
			x:        10,
			y:        5,
			expected: 410, // (5 * 80) + 10
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := level.GetIndexFromXY(tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("GetIndexFromXY(%d, %d) = %d, expected %d", tt.x, tt.y, result, tt.expected)
			}
		})
	}
}

func TestInBounds(t *testing.T) {
	level := Level{}
	gd := config.NewGameData()
	levelHeight = gd.ScreenHeight - gd.UIHeight // Set global variable

	tests := []struct {
		name     string
		x, y     int
		expected bool
	}{
		{
			name:     "valid origin",
			x:        0,
			y:        0,
			expected: true,
		},
		{
			name:     "valid position near edge",
			x:        gd.ScreenWidth - 1,
			y:        levelHeight - 1,
			expected: true,
		},
		{
			name:     "negative x",
			x:        -1,
			y:        10,
			expected: false,
		},
		{
			name:     "negative y",
			x:        10,
			y:        -1,
			expected: false,
		},
		{
			name:     "x at boundary (valid)",
			x:        gd.ScreenWidth,
			y:        10,
			expected: true,
		},
		{
			name:     "y at boundary (valid)",
			x:        10,
			y:        levelHeight,
			expected: true,
		},
		{
			name:     "x too large",
			x:        gd.ScreenWidth + 1,
			y:        10,
			expected: false,
		},
		{
			name:     "y too large",
			x:        10,
			y:        levelHeight + 1,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := level.InBounds(tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("InBounds(%d, %d) = %v, expected %v", tt.x, tt.y, result, tt.expected)
			}
		})
	}
}

func TestIsOpaque(t *testing.T) {
	level := Level{}
	gd := config.NewGameData()
	levelHeight := gd.ScreenHeight - gd.UIHeight

	// Create tiles for testing
	tiles := make([]*MapTile, levelHeight*gd.ScreenWidth)
	for i := range tiles {
		tiles[i] = &MapTile{
			TileType: WALL, // Default to wall (opaque)
		}
	}
	level.Tiles = tiles

	// Set some tiles to floor (transparent)
	floorIndex := level.GetIndexFromXY(5, 5)
	level.Tiles[floorIndex].TileType = FLOOR

	tests := []struct {
		name     string
		x, y     int
		expected bool
	}{
		{
			name:     "wall tile is opaque",
			x:        0,
			y:        0,
			expected: true,
		},
		{
			name:     "floor tile is not opaque",
			x:        5,
			y:        5,
			expected: false,
		},
		{
			name:     "another wall tile is opaque",
			x:        10,
			y:        10,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := level.IsOpaque(tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("IsOpaque(%d, %d) = %v, expected %v", tt.x, tt.y, result, tt.expected)
			}
		})
	}
}

func TestMinMax(t *testing.T) {
	tests := []struct {
		name   string
		x, y   int
		minExp int
		maxExp int
	}{
		{
			name:   "x less than y",
			x:      5,
			y:      10,
			minExp: 5,
			maxExp: 10,
		},
		{
			name:   "x greater than y",
			x:      15,
			y:      8,
			minExp: 8,
			maxExp: 15,
		},
		{
			name:   "x equals y",
			x:      7,
			y:      7,
			minExp: 7,
			maxExp: 7,
		},
		{
			name:   "negative numbers",
			x:      -3,
			y:      -8,
			minExp: -8,
			maxExp: -3,
		},
		{
			name:   "mixed positive and negative",
			x:      -5,
			y:      3,
			minExp: -5,
			maxExp: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			minResult := min(tt.x, tt.y)
			maxResult := max(tt.x, tt.y)

			if minResult != tt.minExp {
				t.Errorf("min(%d, %d) = %d, expected %d", tt.x, tt.y, minResult, tt.minExp)
			}

			if maxResult != tt.maxExp {
				t.Errorf("max(%d, %d) = %d, expected %d", tt.x, tt.y, maxResult, tt.maxExp)
			}
		})
	}
}

func TestCreateRoom(t *testing.T) {
	level := Level{}
	gd := config.NewGameData()
	levelHeight := gd.ScreenHeight - gd.UIHeight

	// Initialize tiles as walls
	tiles := make([]*MapTile, levelHeight*gd.ScreenWidth)
	for i := range tiles {
		tiles[i] = &MapTile{
			Blocked:  true,
			TileType: WALL,
		}
	}
	level.Tiles = tiles

	// Create a 3x3 room at position (5,5)
	room := NewRect(5, 5, 3, 3)
	level.createRoom(room)

	// Check that interior tiles are floors and not blocked
	for y := room.Y1 + 1; y < room.Y2; y++ {
		for x := room.X1 + 1; x < room.X2; x++ {
			index := level.GetIndexFromXY(x, y)
			tile := level.Tiles[index]

			if tile.Blocked {
				t.Errorf("Tile at (%d, %d) should not be blocked after createRoom", x, y)
			}

			if tile.TileType != FLOOR {
				t.Errorf("Tile at (%d, %d) should be FLOOR, got %d", x, y, tile.TileType)
			}
		}
	}

	// Check that border tiles are still walls (room borders, not interior)
	borderTiles := []struct{ x, y int }{
		{room.X1, room.Y1},     // top-left corner
		{room.X1, room.Y1 + 1}, // left edge
		{room.X1 + 1, room.Y1}, // top edge
	}

	for _, pos := range borderTiles {
		index := level.GetIndexFromXY(pos.x, pos.y)
		tile := level.Tiles[index]

		if !tile.Blocked {
			t.Errorf("Border tile at (%d, %d) should still be blocked", pos.x, pos.y)
		}

		if tile.TileType != WALL {
			t.Errorf("Border tile at (%d, %d) should still be WALL, got %d", pos.x, pos.y, tile.TileType)
		}
	}
}
