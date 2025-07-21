# RRogue

A learning-focused roguelike game written in Go, exploring Entity Component System (ECS) architecture and event-driven programming patterns.

## About

RRogue is a traditional turn-based roguelike game built as an educational project to explore modern game architecture patterns. The game features:

- **Entity Component System (ECS)** using `github.com/bytearena/ecs`
- **Event-driven messaging** with publish-subscribe patterns

## Current State

This is a **learning project** currently featuring:
- Single dungeon level with procedural generation
- Turn-based combat between player and monsters  
- Basic inventory and equipment system
- Event-driven UI messaging system
- Game state management

**Note**: The game currently supports only a single level, with plans for multi-level dungeons in future iterations.

## Prerequisites

- [Go](https://golang.org/dl/) 1.19 or later

## Getting Started

### Clone and Build

```bash
# Clone the repository
git clone https://github.com/caustin/rrogue.git
cd rrogue

# Download dependencies
go mod download

# Build the game
go build

# Run the game
./rrogue
```

### Alternative: Run Without Building

```bash
# Run directly with go run (automatically handles dependencies)
go run .
```

### Development

```bash
# Run with go run for development
go run .

# Run tests
go test ./...

# Check for issues
go vet ./...
```

## Controls

- **Arrow Keys**: Move player
- **Mouse**: Alternative movement (click to move)
- **ESC**: Quit game


## Learning Goals

This project explores several key concepts:

1. **Entity Component System (ECS)**: Pure data components with behavior in systems
2. **Event-Driven Architecture**: Loose coupling through message passing
3. **Go Concurrency**: Thread-safe systems with proper synchronization

## Architecture

For detailed architecture information, see [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md).

The project demonstrates a migration from tightly-coupled legacy code to a clean event-driven architecture, showing how to refactor complex game systems incrementally.

## Contributing

This is primarily a learning project, but feedback and suggestions are welcome! Please open an issue to discuss any changes.
