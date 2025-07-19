# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a roguelike game written in Go using the Ebiten game engine. The project implements a traditional roguelike
with dungeon generation, turn-based combat, and Entity Component System (ECS) architecture. The codebase has undergone
significant refactoring for improved performance, code quality, and maintainability.

## Development Commands

### Building and Running
- `go build` - Build the executable
- `go run .` - Build and run the game directly
- `go mod tidy` - Clean up dependencies

### Testing and Quality
- `go test ./...` - Run all tests (comprehensive test coverage implemented)
- `go fmt ./...` - Format all Go files
- `go vet ./...` - Run Go vet for static analysis

## Architecture

### Core Components
- **Entity Component System**: Uses `github.com/bytearena/ecs` for game entity management
- **Game Engine**: Built on `github.com/hajimehoshi/ebiten/v2` for 2D graphics and input
- **Field of View**: Uses `github.com/norendren/go-fov` for vision calculations

### Key Files and Systems
- `main.go` - Entry point, Game struct with cached GameData, main game loop (Update/Draw/Layout)
- `world.go` - ECS world initialization with dependency injection, Components struct
- `components.go` - All ECS component definitions (Position, Health, Monster, etc.)
- `level.go` - Level/map structure and optimized rendering with sync.Pool
- `dungeon.go` - Dungeon container for multiple levels
- `map.go` - Game map management and current level tracking
- `render_pool.go` - Memory pool for DrawImageOptions to reduce allocations
- `*_system.go` files - Game systems using dependency injection:
  - `player_systems.go` - Player input and actions
  - `monster_systems.go` - AI and monster behavior
  - `combat_system.go` - Combat resolution
  - `render_system.go` - Entity rendering with pooled memory management
  - `hud_system.go` - UI/HUD rendering
  - `userlog_system.go` - Message log system
- `*_test.go` files - Comprehensive test coverage:
  - `rect_test.go` - Rectangle intersection and geometric tests
  - `turnstate_test.go` - Turn state management tests
  - `level_test.go` - Level generation and utility function tests
  - `components_test.go` - ECS component validation tests
  - `dice_test.go` - Random number generation and dice rolling tests

### Game Flow
- Turn-based system with WaitingForPlayerInput/ProcessingMonsterTurn states
- Player input triggers actions, then monster AI runs
- Rendering happens every frame but game logic is turn-based
- Uses A* pathfinding for monster movement
- Performance optimized render loop with 99%+ allocation reduction

### Player Controls
- **Arrow Keys**: Single-step movement and combat
- **Period (.) + Arrow Keys**: Timed auto-movement in direction (~8 moves/second) until:
  - A monster becomes visible in player's field of view
  - Player reaches a corridor junction (>2 walkable directions)
  - Player reaches a room or open area
  - Movement is blocked by walls or obstacles
  - Any key is pressed to interrupt movement
- **Q**: Wait/skip turn
- **Escape**: Stop auto-movement

### Assets
- PNG sprites in `assets/` directory (player, monsters, tiles)
- Font files in `fonts/` directory

## Dependencies
- Go 1.17+
- Ebiten v2 game engine
- ECS library for entity management
- Go-fov for field of view calculations

## Recent Major Improvements

### Performance Optimizations
- **Render Pool**: Implemented sync.Pool for DrawImageOptions, reducing allocations from ~240,000/sec to <1,000/sec
- **GameData Caching**: Eliminated repeated struct allocations by caching GameData in Game struct
- **Memory Management**: 99%+ reduction in garbage collection pressure during gameplay

### Code Quality Enhancements  
- **ECS Refactoring**: Eliminated global ECS components, introduced dependency injection via Components struct
- **Turn System**: Refactored turn states from BeforePlayerAction/PlayerTurn/MonsterTurn to WaitingForPlayerInput/ProcessingMonsterTurn
- **Bug Fixes**: Fixed critical rectangle intersection logic, dice rolling algorithms, and entity cleanup issues

### Test Coverage
- **Comprehensive Testing**: Added extensive test suites covering rectangle operations, turn states, level generation, components, and dice systems
- **Quality Assurance**: All core game logic now has proper unit test coverage with edge case validation

### Architecture Improvements
- **Dependency Injection**: ECS components now properly injected rather than using global variables
- **Clean Code**: Reduced tight coupling and improved maintainability across the codebase
- **Documentation**: Added detailed performance analysis and optimization plans in `/plans/` directory

## Performance Notes
- The render loop has been heavily optimized - avoid creating new DrawImageOptions manually
- Use GetDrawOptions()/PutDrawOptions() from render_pool.go for all rendering operations
- GameData is cached in the Game struct - pass the cached instance rather than calling NewGameData()
- All major performance bottlenecks have been identified and resolved