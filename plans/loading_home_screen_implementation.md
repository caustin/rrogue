# Loading/Home Screen Implementation Plan

## Current State Analysis

### Existing Game Flow
- Game starts directly into gameplay via `NewGame()` in `main.go:77`
- No menu system or state management beyond turn states
- Turn states: `WaitingForPlayerInput`, `ProcessingPlayerAction`, `ProcessingMonsterTurn`, `GameOver`
- No save/load functionality currently exists
- Game initialization creates new world, map, and entities immediately

### Current Architecture
- `Game` struct contains all game state (World, Map, Components, Turn, etc.)
- Single entry point through `ebiten.RunGame(g)` 
- No concept of "screens" or "scenes" beyond game states

## Proposed Architecture

### Game State Management

**New Game States (Expanding TurnState):**
```go
type GameState int

const (
    // Menu States
    MainMenu GameState = iota
    LoadGame
    Settings
    
    // Game States  
    Playing
    Paused
    GameOver
    
    // Loading States
    LoadingNewGame
    LoadingSavedGame
)

type TurnState int // Keep existing for in-game turns
```

**Game State Manager:**
```go
type GameStateManager struct {
    currentState GameState
    previousState GameState
    stateData map[GameState]interface{}
    transitions map[GameState][]GameState
}
```

### Screen System Architecture

**Base Screen Interface:**
```go
type Screen interface {
    Update(g *Game) error
    Draw(g *Game, screen *ebiten.Image)
    OnEnter(g *Game)
    OnExit(g *Game)
}
```

**Screen Implementations:**
- `MainMenuScreen` - Home screen with New Game/Load Game options
- `LoadGameScreen` - Save file selection and loading
- `SettingsScreen` - Game configuration options  
- `GameplayScreen` - Current game logic (refactored)
- `PauseScreen` - In-game pause menu
- `GameOverScreen` - Death/victory screen

## Phase 1: Core State Management

### 1.1 Create State System

**New Files:**
- `gamestate.go` - GameState constants and management
- `screen.go` - Screen interface and base functionality
- `screen_manager.go` - Screen switching and management

**State Manager Features:**
- State transition validation
- State history for back navigation
- State-specific data storage
- Event-driven state changes

### 1.2 Refactor Main Game Loop

**Modify `main.go`:**
```go
type Game struct {
    // Existing fields...
    StateManager *GameStateManager
    CurrentScreen Screen
    Screens map[GameState]Screen
}

func (g *Game) Update() error {
    return g.CurrentScreen.Update(g)
}

func (g *Game) Draw(screen *ebiten.Image) {
    g.CurrentScreen.Draw(g, screen)
}
```

**Screen Registration:**
```go
func (g *Game) initializeScreens() {
    g.Screens[MainMenu] = &MainMenuScreen{}
    g.Screens[LoadGame] = &LoadGameScreen{}
    g.Screens[Playing] = &GameplayScreen{}
    g.Screens[GameOver] = &GameOverScreen{}
}
```

## Phase 2: Main Menu Screen

### 2.1 Main Menu UI Design

**Visual Layout:**
```
    ╔════════════════════════════════╗
    ║           RROGUE               ║
    ║                                ║
    ║         [New Game]             ║
    ║         [Load Game]            ║
    ║         [Settings]             ║
    ║         [Exit]                 ║
    ║                                ║
    ║    Version: 0.0.19             ║
    ╚════════════════════════════════╝
```

**Menu Implementation:**
```go
type MainMenuScreen struct {
    selectedOption int
    menuOptions []MenuOption
    titleFont font.Face
    menuFont font.Face
    backgroundImage *ebiten.Image
}

type MenuOption struct {
    text string
    action func(*Game)
    enabled bool
}
```

### 2.2 Menu Navigation

**Input Handling:**
- Arrow keys/WASD for menu navigation
- Enter/Space to select option
- Escape to exit (from main menu)
- Mouse support for menu selection

**Menu Actions:**
```go
func (m *MainMenuScreen) handleNewGame(g *Game) {
    g.StateManager.TransitionTo(LoadingNewGame)
}

func (m *MainMenuScreen) handleLoadGame(g *Game) {
    g.StateManager.TransitionTo(LoadGame)
}
```

### 2.3 Visual Polish

**Background Options:**
- Static background image
- Animated particle effects
- Rotating/moving background elements
- Game logo/title treatment

**Audio Integration (Future):**
- Background music
- Menu selection sounds
- Ambient effects

## Phase 3: Save/Load System Foundation

### 3.1 Save Data Structure

**Save File Format (JSON):**
```go
type SaveData struct {
    Version string `json:"version"`
    Timestamp time.Time `json:"timestamp"`
    PlayerData PlayerSaveData `json:"player"`
    WorldData WorldSaveData `json:"world"`
    MapData MapSaveData `json:"map"`
    GameStats GameStatistics `json:"stats"`
}

type PlayerSaveData struct {
    Position Position `json:"position"`
    Health Health `json:"health"`
    Armor Armor `json:"armor"`
    Weapon MeleeWeapon `json:"weapon"`
    Level int `json:"level"`
}
```

### 3.2 Serialization System

**Component Serialization:**
```go
type Serializable interface {
    ToJSON() ([]byte, error)
    FromJSON([]byte) error
}

// Implement for each component type
func (h *Health) ToJSON() ([]byte, error) {
    return json.Marshal(h)
}
```

**World Serialization:**
```go
func SerializeWorld(world *ecs.Manager) (*WorldSaveData, error) {
    // Serialize all entities and components
    // Handle entity relationships
    // Compress data if needed
}
```

### 3.3 Save File Management

**Save Directory Structure:**
```
saves/
├── save_001.json
├── save_002.json
├── save_003.json
└── autosave.json
```

**Save File Operations:**
```go
func ListSaveFiles() ([]SaveFileInfo, error)
func LoadSaveFile(filename string) (*SaveData, error) 
func WriteSaveFile(filename string, data *SaveData) error
func DeleteSaveFile(filename string) error
```

## Phase 4: Load Game Screen

### 4.1 Save File Browser

**UI Layout:**
```
    ╔════════════════════════════════╗
    ║          Load Game             ║
    ║                                ║
    ║ ┌────────────────────────────┐ ║
    ║ │ Save 1 - 2025-07-19 14:30  │ ║
    ║ │ Level 3, Health: 45/60     │ ║
    ║ ├────────────────────────────┤ ║
    ║ │ Save 2 - 2025-07-18 20:15  │ ║
    ║ │ Level 1, Health: 60/60     │ ║
    ║ └────────────────────────────┘ ║
    ║                                ║
    ║      [Load]  [Delete]  [Back]  ║
    ╚════════════════════════════════╝
```

**Save File Display:**
- Save file name and timestamp
- Character stats preview
- Game progress indicators
- Save file thumbnails (future enhancement)

### 4.2 Loading Process

**Loading States:**
```go
type LoadingScreen struct {
    progress float64
    message string
    currentStep string
}
```

**Loading Steps:**
1. Validate save file
2. Deserialize world data
3. Recreate entities and components
4. Initialize game systems
5. Transition to gameplay

## Phase 5: Integration and Polish

### 5.1 Game Integration

**Gameplay Screen Refactor:**
- Move current game logic to `GameplayScreen`
- Implement pause functionality
- Add in-game menu access
- Handle save game triggers

**Pause Menu:**
```go
type PauseScreen struct {
    backgroundGame *ebiten.Image // Screenshot of game
    menuOptions []MenuOption
}
```

### 5.2 Settings Screen

**Configuration Options:**
- Graphics settings (resolution, fullscreen)
- Audio settings (volume, effects)
- Gameplay settings (difficulty, auto-save)
- Keybinding configuration

**Settings Persistence:**
```go
type GameSettings struct {
    Graphics GraphicsSettings `json:"graphics"`
    Audio AudioSettings `json:"audio"`  
    Gameplay GameplaySettings `json:"gameplay"`
    Controls ControlSettings `json:"controls"`
}
```

### 5.3 Error Handling and Validation

**Save File Validation:**
- Version compatibility checks
- Corruption detection
- Graceful error recovery
- User-friendly error messages

**Loading Error Recovery:**
- Backup save file creation
- Partial load recovery
- Error reporting and logging

## Implementation Timeline

### Week 1: Core Infrastructure (Phase 1)
- Implement GameState system
- Create Screen interface and manager
- Refactor main game loop
- Basic state transitions

### Week 2: Main Menu (Phase 2)  
- Design and implement MainMenuScreen
- Add menu navigation and input handling
- Create visual assets and styling
- Basic menu functionality

### Week 3: Save System Foundation (Phase 3)
- Design save data structures
- Implement serialization system
- Create save file management
- Basic save/load functionality

### Week 4: Load Screen and Integration (Phase 4-5)
- Implement LoadGameScreen
- Create loading process
- Integrate with existing game systems
- Polish and error handling

## Success Criteria

1. **Professional Menu System**: Clean, navigable main menu
2. **Reliable Save/Load**: Persistent game state across sessions
3. **Seamless Integration**: No disruption to existing gameplay
4. **Error Resilience**: Graceful handling of save file issues
5. **Extensibility**: Easy to add new menu options and features

## Technical Considerations

### Performance
- Lazy loading of save file previews
- Efficient serialization for large worlds
- Memory management during state transitions

### Compatibility
- Save file version management
- Backward compatibility strategy
- Migration tools for save format changes

### User Experience
- Fast loading times
- Clear progress indicators
- Intuitive navigation
- Helpful error messages

## Future Enhancements

- **Cloud Save Support**: Save files in cloud storage
- **Save File Screenshots**: Visual previews of save states
- **Multiple Save Slots**: Organized save file management
- **Quick Save/Load**: Hotkey-based save/load functionality
- **Save File Import/Export**: Share save files between players
- **Achievement System**: Track player progress and milestones