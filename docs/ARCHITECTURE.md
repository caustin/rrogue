# RRogue Game Architecture Documentation

## Table of Contents
- [Overview](#overview)
- [Package Structure](#package-structure)
- [Core Components](#core-components)
- [Event System](#event-system)
- [World Service](#world-service)
- [Systems Architecture](#systems-architecture)
- [Game Loop](#game-loop)
- [Data Flow](#data-flow)
- [Adding New Features](#adding-new-features)
- [Testing Strategy](#testing-strategy)
- [Future Considerations](#future-considerations)

## Overview

RRogue is a roguelike game built with Go using an Entity Component System (ECS) architecture combined with event-driven messaging. The design emphasizes clean separation of concerns, dependency injection, and scalable extensibility.

### Key Architectural Principles
- **Composition over Inheritance**: ECS pattern for flexible entity composition
- **Dependency Injection**: Systems receive only what they need via interfaces
- **Event-Driven Communication**: Systems communicate through events, not direct calls
- **Single Responsibility**: Each package has a clear, focused purpose
- **Loose Coupling**: Components can be modified independently

## Package Structure

```
/rrogue-0.0.19/
├── main.go                     # Application entry point
├── components/                 # ECS component definitions
│   └── components.go
├── config/                     # Configuration data
│   └── gamedata.go
├── events/                     # Event system implementation
│   ├── event.go               # Base event types and interfaces
│   ├── bus.go                 # Event bus implementation
│   ├── combat_events.go       # Combat-specific events
│   ├── game_events.go         # Game state events
│   └── bus_test.go           # Event system tests
├── game/                       # Core game logic and systems
│   ├── game.go               # Main game struct and loop
│   ├── combat_system.go      # Combat logic (legacy, being refactored)
│   ├── hud_system.go         # UI rendering
│   ├── map.go               # Map data structures
│   ├── monster_systems.go   # Monster AI and behavior
│   ├── player_systems.go    # Player input and movement
│   ├── render_system.go     # Rendering pipeline
│   ├── turnstate.go         # Turn state management
│   └── userlog_system.go    # User message logging
├── level/                      # Level generation and management
│   ├── astar.go             # Pathfinding algorithms
│   ├── dungeon.go           # Dungeon generation
│   └── level.go             # Level data structures
├── systems/                    # Event-driven system implementations
│   ├── combat.go            # Event-driven combat system
│   ├── gamestate.go         # Game state management system
│   ├── ui.go                # User interface and message system
│   ├── registry.go          # System registry and lifecycle management
│   ├── gamebridge.go        # Temporary bridge for game state access
│   └── mapbridge.go         # Map tile management bridge
├── utils/                      # Utility functions
│   ├── dice.go              # Random number generation
│   ├── rect.go              # Rectangle utilities
│   └── render_pool.go       # Rendering optimizations
├── world/                      # ECS world management
│   ├── service.go           # WorldService interface
│   └── gameworld.go         # WorldService implementation
└── docs/                       # Documentation
    ├── CLAUDE.md            # Development context
    └── ARCHITECTURE.md      # This file
```

## Core Components

### Entity Component System (ECS)

The game uses the `github.com/bytearena/ecs` library to implement a pure ECS architecture:

- **Entities**: Game objects identified by unique IDs
- **Components**: Pure data structures (Position, Health, Armor, etc.)
- **Systems**: Logic that operates on entities with specific component combinations

#### Component Types (components/components.go)
```go
type Position struct {
    X, Y int
}

type Health struct {
    MaxHealth, CurrentHealth int
}

type Armor struct {
    Name string
    Defense, ArmorClass int
}

type MeleeWeapon struct {
    Name string
    MinimumDamage, MaximumDamage, ToHitBonus int
}

type Name struct {
    Label string
}

type UserMessage struct {
    AttackMessage, DeadMessage, GameStateMessage string
}

type Renderable struct {
    Image *ebiten.Image
}

type Player struct{}
type Monster struct{}
type Movable struct{}
```

### Game Struct (game/game.go)

The main Game struct orchestrates the entire application:

```go
type Game struct {
    Map           GameMap                 // Level and map data
    World         world.WorldService      // ECS world interface
    EventBus      *events.EventBus        // Event messaging system
    Systems       *systems.SystemRegistry // Centralized system management
    GameData      config.GameData         // Display configuration
    Turn          TurnState               // Current turn state
    TurnCounter   int                     // Turn tracking
    AutoMoveState *AutoMoveState          // Player auto-movement state
}
```

## Event System

The event system enables loose coupling between game systems through publish-subscribe messaging.

### Event Bus (events/bus.go)

Thread-safe event dispatcher that manages subscriptions and publishing:

```go
type EventBus struct {
    subscribers map[EventType][]EventHandler
    mutex       sync.RWMutex
}

type EventHandler func(event Event)
```

#### Key Methods
- `Subscribe(eventType EventType, handler EventHandler)`: Register event handler
- `Publish(event Event)`: Send event to all subscribers
- `PublishMany(events []Event)`: Send multiple events in sequence

### Event Types

#### Base Event Interface (events/event.go)
```go
type Event interface {
    Type() EventType
    Timestamp() time.Time
}

type EventType string
```

#### Combat Events (events/combat_events.go)
- **AttackEvent**: Represents an attack between entities
- **DamageEvent**: Damage being applied to an entity
- **DeathEvent**: Entity death notification
- **MoveEvent**: Entity movement tracking

#### Game State Events (events/game_events.go)
- **TurnStartEvent**: Beginning of player/monster turn
- **TurnEndEvent**: End of turn
- **TurnChangeEvent**: Turn state transitions between player/monster/game over
- **TurnCounterEvent**: Turn counter updates and increments
- **GameOverEvent**: Game termination (enhanced with FinalTurn field)
- **TileBlockedEvent/TileUnblockedEvent**: Map state changes

### Event Flow Example

#### Combat and Death Sequence
```
Player Attack Sequence:
1. TakePlayerAction() → CombatSystem.ProcessAttack() → AttackEvent
2. CombatSystem.HandleAttack() → MessageEvent + DamageEvent  
3. UISystem.HandleMessage() → adds message to queue
4. CombatSystem.HandleDamage() → DeathEvent (if fatal)
5. GameStateSystem.HandleDeath() → GameOverEvent (if player death)
6. MapBridge.HandleEntityDeath() → unblocks tile (if monster death)
```

#### Turn Transition Sequence  
```
Player Turn Complete:
1. TakePlayerAction() returns true → GameStateSystem.ChangeTurn()
2. GameStateSystem.ChangeTurn() → TurnChangeEvent + IncrementTurn()
3. GameStateSystem.HandleTurnChange() → syncs g.Turn to ProcessingMonsterTurn
4. UpdateMonster() processes AI → GameStateSystem.ChangeTurn() 
5. GameStateSystem.HandleTurnChange() → syncs g.Turn to WaitingForPlayerInput
```

## World Service

The WorldService interface abstracts ECS operations and provides dependency injection for systems.

### Interface Definition (world/service.go)

```go
type WorldService interface {
    // Entity queries
    QueryPlayers() []*ecs.QueryResult
    QueryMonsters() []*ecs.QueryResult
    QueryRenderables() []*ecs.QueryResult
    QueryMessengers() []*ecs.QueryResult
    
    // Component access
    GetPosition(entity *ecs.QueryResult) *components.Position
    GetHealth(entity *ecs.QueryResult) *components.Health
    GetArmor(entity *ecs.QueryResult) *components.Armor
    GetMeleeWeapon(entity *ecs.QueryResult) *components.MeleeWeapon
    GetName(entity *ecs.QueryResult) *components.Name
    GetUserMessage(entity *ecs.QueryResult) *components.UserMessage
    GetRenderable(entity *ecs.QueryResult) *components.Renderable
    
    // Entity lifecycle
    DisposeEntity(entity *ecs.QueryResult)
    
    // Raw access for advanced use cases
    GetManager() *ecs.Manager
}
```

### Implementation (world/gameworld.go)

```go
type GameWorld struct {
    manager    *ecs.Manager
    tags       map[string]ecs.Tag
    components *ComponentReferences
}
```

The GameWorld implements WorldService and manages:
- ECS manager instance
- Entity tag queries (players, monsters, renderables)
- Component type registrations
- Entity creation and initialization

## Systems Architecture

### Legacy Systems (game/ package)

Current systems that directly access the Game struct:
- **TakePlayerAction**: Handles player input and movement
- **UpdateMonster**: Monster AI and pathfinding
- **ProcessRenderables**: Entity rendering
- **ProcessHUD**: UI display
- **ProcessUserLog**: Message logging
- **AttackSystem**: Combat logic (being migrated)

### Event-Driven Systems (systems/ package)

New systems that use dependency injection and event messaging:

#### SystemRegistry (systems/registry.go)
```go
type SystemRegistry struct {
    Combat     *CombatSystem
    GameState  *GameStateSystem
    UI         *UISystem
    GameBridge *GameBridge
    MapBridge  *MapBridge
    
    world    world.WorldService
    eventBus *events.EventBus
}
```

**Responsibilities:**
- Manage system lifecycle and dependencies
- Initialize all systems with proper dependency injection
- Register event handlers for all systems
- Provide centralized access to systems

#### CombatSystem (systems/combat.go)
```go
type CombatSystem struct {
    world    world.WorldService
    eventBus *events.EventBus
}
```

**Responsibilities:**
- Subscribe to AttackEvent and DamageEvent
- Process hit/miss calculations
- Apply damage to entities
- Publish damage and death events

**Key Methods:**
- `HandleAttack(event events.Event)`: Process attack attempts
- `HandleDamage(event events.Event)`: Apply damage and check for death
- `ProcessAttack(attackerPos, defenderPos)`: Initiate attack sequence

#### GameStateSystem (systems/gamestate.go)
```go
type GameStateSystem struct {
    world    world.WorldService
    eventBus *events.EventBus
    
    currentState TurnState
    turnCounter  int
    mutex        sync.RWMutex
}
```

**Responsibilities:**
- Centralized turn state management (WaitingForPlayerInput, ProcessingMonsterTurn, GameOver)
- Handle player death and trigger game over events
- Manage turn transitions through events
- Thread-safe state synchronization with Game struct

**Key Methods:**
- `HandleDeath(event events.Event)`: Process entity deaths and trigger game over for players
- `ChangeTurn(toState TurnState)`: Publish turn change events  
- `TriggerGameOver(reason string)`: Publish game over events
- `IncrementTurn()`: Increment and publish turn counter events

#### UISystem (systems/ui.go)
```go
type UISystem struct {
    world    world.WorldService
    eventBus *events.EventBus
    
    messages []UIMessage
    mutex    sync.RWMutex
}
```

**Responsibilities:**
- Event-driven user interface message management
- Thread-safe message queue with persistence across frames
- Replace direct UserMessage component manipulation
- Latest-first message ordering for better UX

**Key Methods:**
- `HandleMessage(event events.Event)`: Process MessageEvent and add to queue
- `GetMessageTexts()`: Return current messages for display (latest first)
- `AddMessage(text, messageType string)`: Direct message addition

#### MapBridge (systems/mapbridge.go)
```go
type MapBridge struct {
    eventBus *events.EventBus
    gameRef  interface{}
}
```

**Responsibilities:**
- Handle entity death events for tile cleanup
- Unblock map tiles when monsters die
- Bridge between event system and map operations during migration

## Game Loop

### Main Loop (game/game.go)

```go
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
    }
    return nil
}
```

### Turn-Based Flow

1. **Player Turn**:
   - Process input (movement, attacks, special actions)
   - Validate actions against game rules
   - Apply immediate effects
   - Transition to monster turn

2. **Monster Turn**:
   - Update each monster's AI
   - Process monster actions (movement, attacks)
   - Apply effects
   - Return to player turn

3. **Event Processing**:
   - Events are published during turn processing
   - Subscribers handle events immediately
   - State changes propagate through event chain

## Data Flow

### Input Processing
```
Player Input → TakePlayerAction() → Game State Changes → Event Publishing
```

### Combat Flow
```
Attack Request → AttackEvent → CombatSystem → DamageEvent → 
DeathEvent → GameStateSystem/MapSystem → State Updates
```

### Rendering Flow
```
Game State → Query Entities → Component Data → Render Pipeline → Display
```

### Event Propagation
```
System Action → Event Publication → Event Bus → Subscribed Systems → 
Side Effects → Additional Events → Cascading Updates
```

## Adding New Features

### Adding a New Component

1. Define component in `components/components.go`:
```go
type NewComponent struct {
    Field1 int
    Field2 string
}
```

2. Register in `world/gameworld.go`:
```go
components := &ComponentReferences{
    // ... existing components
    NewComponent: manager.NewComponent(),
}
```

3. Add WorldService method:
```go
GetNewComponent(entity *ecs.QueryResult) *components.NewComponent
```

### Adding a New Event

1. Define event in appropriate events file:
```go
type NewEvent struct {
    BaseEvent
    Data interface{}
}

func NewNewEvent(data interface{}) *NewEvent {
    return &NewEvent{
        BaseEvent: NewBaseEvent(NewEventType),
        Data:      data,
    }
}
```

2. Add event type constant:
```go
const NewEventType EventType = "new_event"
```

### Adding a New System

1. Create system struct with dependencies:
```go
type NewSystem struct {
    world    world.WorldService
    eventBus *events.EventBus
    // other dependencies
}
```

2. Implement constructor with dependency injection:
```go
func NewNewSystem(world world.WorldService, eventBus *events.EventBus) *NewSystem {
    return &NewSystem{world: world, eventBus: eventBus}
}
```

3. Register event handlers:
```go
func (s *NewSystem) RegisterHandlers() {
    s.eventBus.Subscribe(EventType, s.HandleEvent)
}
```

4. Add to SystemRegistry:
```go
// In NewSystemRegistry
registry.NewSystem = NewNewSystem(world, eventBus)

// In RegisterAllHandlers
if r.NewSystem != nil {
    r.NewSystem.RegisterHandlers()
}
```

## Testing Strategy

### Unit Testing

**Event System**: Comprehensive tests in `events/bus_test.go`
- Event publication and subscription
- Multiple subscriber handling
- Event creation and validation

**Systems**: Test each system in isolation
- Mock WorldService dependencies
- Verify event publishing behavior
- Test edge cases and error conditions

### Integration Testing

**Game Loop**: Test complete turn cycles
- Player action processing
- Monster AI behavior
- Event propagation chains

**Combat**: Test attack sequences end-to-end
- Hit/miss calculations
- Damage application
- Death handling

### Testing Example

```go
func TestCombatSystem_HandleAttack(t *testing.T) {
    // Setup
    mockWorld := &MockWorldService{}
    eventBus := events.NewEventBus()
    combat := systems.NewCombatSystem(mockWorld, eventBus)
    
    // Configure mocks
    mockWorld.On("GetArmor", mock.Anything).Return(&components.Armor{ArmorClass: 10})
    
    // Test
    attackEvent := events.NewAttackEvent(attacker, defender, pos1, pos2, 15, true)
    combat.HandleAttack(attackEvent)
    
    // Verify
    mockWorld.AssertCalled(t, "GetArmor", defender)
}
```

## Future Considerations

### Current State

**Migration Status**: Core event-driven systems are fully implemented and operational

#### ✅ **Completed Migrations**
- ✅ **CombatSystem**: Event-driven combat with dependency injection
- ✅ **GameStateSystem**: Centralized turn state management and game over handling  
- ✅ **UISystem**: Event-driven user interface message management
- ✅ **SystemRegistry**: Lifecycle management for all event-driven systems
- ✅ **Event Integration**: Combat messages flow through UISystem instead of direct component access
- ✅ **Temporary Handler Removal**: Game struct no longer handles events directly

#### 🔄 **Current Architecture**
- SystemRegistry manages CombatSystem, GameStateSystem, UISystem, and bridge systems
- GameStateSystem handles all turn transitions and game over conditions
- UISystem provides thread-safe message queue with latest-first ordering
- MapBridge handles tile cleanup during transition period
- WorldService interface provides clean ECS abstraction

### Planned Improvements

1. **Complete System Migration**:
   - Implement proper MapSystem to replace MapBridge
   - Remove legacy AttackSystem function (marked low priority)  
   - Migrate remaining legacy systems to event-driven architecture

2. **Enhanced Event System**:
   - Event queuing and delayed execution
   - Event priority and ordering
   - Event replay and undo capabilities

3. **Save/Load System**:
   - Event sourcing for game state persistence
   - Replay functionality for debugging
   - Checkpoint and restore capabilities

4. **Networking Support**:
   - Event synchronization between clients
   - Deterministic event ordering
   - Network-aware event batching

5. **Modding Framework**:
   - Plugin system for custom systems
   - Event-based mod API
   - Dynamic component registration

6. **Performance Optimizations**:
   - Event batching and coalescing
   - System execution scheduling
   - Memory pool management

### Design Patterns for Extension

**Command Pattern**: Implement undoable actions
**Observer Pattern**: Enhanced event subscription with filters
**Strategy Pattern**: Pluggable AI behaviors
**Factory Pattern**: Dynamic entity creation
**State Machine**: Complex game state management

### Scalability Considerations

- **System Prioritization**: Critical systems run first
- **Event Throttling**: Prevent event flood scenarios
- **Resource Management**: Component pooling and reuse
- **Parallel Processing**: Concurrent system execution
- **Memory Optimization**: Efficient component storage

## Conclusion

The RRogue architecture provides a solid foundation for roguelike game development with clear separation of concerns, testable components, and extensible design. The event-driven messaging system enables complex interactions while maintaining loose coupling between systems.

The architecture supports both immediate game development needs and future expansion into multiplayer, modding, and advanced features. The clean package structure and dependency injection patterns make the codebase maintainable and approachable for new developers.

---

*This documentation reflects the current state of the architecture as of the event system integration. It should be updated as the system evolves and new patterns emerge.*