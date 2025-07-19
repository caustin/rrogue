# Centralized Logging and Debug Overlay Implementation Plan

## Current State Analysis

### Existing Logging
- Only `log.Fatal()` calls in 5 files for initialization errors
- No structured logging or debug output capabilities
- Commented-out `fmt.Printf()` statements in `userlog_system.go`
- No centralized logging mechanism
- No development/debugging visibility into system operations

## Phase 1: Core Centralized Logger

### 1.1 Create Logger Infrastructure

**New File: `logger.go`**
```go
type Logger struct {
    level           LogLevel
    output          io.Writer
    enableTimestamp bool
    enableSource    bool
    mu             sync.Mutex
}

type LogLevel int
const (
    DEBUG LogLevel = iota
    INFO
    WARN
    ERROR
    FATAL
    OFF
)
```

**Key Features:**
- Thread-safe logging with mutex
- Configurable log levels
- Multiple output destinations (file, console, memory buffer)
- Structured log entry format
- Source system identification

### 1.2 Logger Methods

**Core Logging Methods:**
```go
func (l *Logger) Debug(system string, msg string, args ...interface{})
func (l *Logger) Info(system string, msg string, args ...interface{})
func (l *Logger) Warn(system string, msg string, args ...interface{})
func (l *Logger) Error(system string, msg string, args ...interface{})
func (l *Logger) Fatal(system string, msg string, args ...interface{})
```

**Configuration Methods:**
```go
func (l *Logger) SetLevel(level LogLevel)
func (l *Logger) SetOutput(w io.Writer)
func (l *Logger) EnableTimestamp(enable bool)
func (l *Logger) EnableSource(enable bool)
```

### 1.3 Integration with ECS

**Modify `components.go`:**
- Add Logger field to Components struct
- Initialize in NewComponents()

**Modify `world.go`:**
- Create logger instance during world initialization
- Pass logger to Components struct

**Update all system files:**
- Replace `log.Fatal()` calls with `g.Components.Logger.Fatal()`
- Add debug/info logging for key operations

## Phase 2: Logger Configuration and Management

### 2.1 Configuration System

**Environment Variable Support:**
```go
func NewLoggerFromEnv() *Logger {
    level := parseLogLevel(os.Getenv("RROGUE_LOG_LEVEL"))
    output := getLogOutput(os.Getenv("RROGUE_LOG_OUTPUT"))
    return NewLogger(level, output)
}
```

**Configuration Options:**
- `RROGUE_LOG_LEVEL`: DEBUG, INFO, WARN, ERROR, FATAL, OFF
- `RROGUE_LOG_OUTPUT`: stdout, stderr, file path
- `RROGUE_LOG_FORMAT`: timestamp, source system inclusion

### 2.2 Multiple Output Destinations

**File Logging:**
```go
func (l *Logger) AddFileOutput(filename string) error
func (l *Logger) SetRotatingFileOutput(filename string, maxSize int) error
```

**Memory Buffer (for overlay):**
```go
type MemoryBuffer struct {
    entries    []LogEntry
    maxEntries int
    mu         sync.RWMutex
}
```

### 2.3 Performance Optimizations

- Lazy string formatting (only format if log level allows)
- String pooling for frequent log messages
- Buffered I/O for file output
- Minimal allocation during logging calls

## Phase 3: System Integration

### 3.1 Replace Existing Logging

**Files to Update:**
- `main.go`: Replace log.Fatal with centralized logger
- `world.go`: Replace log.Fatal calls
- `level.go`: Replace log.Fatal calls
- `hud_system.go`: Replace log.Fatal calls
- `userlog_system.go`: Replace log.Fatal calls, uncomment debug prints

### 3.2 Add Comprehensive Logging

**Combat System Logging:**
```go
g.Components.Logger.Debug("combat_system", "Player attacks %s for %d damage", monsterName, damage)
g.Components.Logger.Info("combat_system", "Monster %s defeated", monsterName)
```

**AI System Logging:**
```go
g.Components.Logger.Debug("monster_ai", "Monster %s pathfinding to player", monsterID)
g.Components.Logger.Debug("monster_ai", "Monster %s taking turn", monsterID)
```

**Player System Logging:**
```go
g.Components.Logger.Debug("player_system", "Player moved to position (%d, %d)", x, y)
g.Components.Logger.Info("player_system", "Player picked up item: %s", itemName)
```

### 3.3 Performance Monitoring

**Render System Logging:**
```go
g.Components.Logger.Debug("render_system", "Frame rendered in %v", duration)
g.Components.Logger.Warn("render_system", "Frame time exceeded target: %v", duration)
```

**Memory Management:**
```go
g.Components.Logger.Debug("render_pool", "DrawOptions pool size: %d", poolSize)
g.Components.Logger.Debug("memory", "GC triggered, heap size: %d MB", heapMB)
```

## Phase 4: Testing and Validation

### 4.1 Unit Tests

**Test Files to Create:**
- `logger_test.go`: Core logger functionality
- `logger_integration_test.go`: ECS integration tests

**Test Coverage:**
- Log level filtering
- Thread safety
- Output formatting
- Performance benchmarks
- Memory usage validation

### 4.2 Integration Testing

- Verify all systems can log without errors
- Test logger performance under high load
- Validate log output formatting
- Test configuration from environment variables

## Phase 5: Debug Overlay System (Final Phase)

### 5.1 Overlay Infrastructure

**New File: `debug_overlay.go`**
```go
type DebugOverlay struct {
    enabled       bool
    visible       bool
    logBuffer     *MemoryBuffer
    scrollOffset  int
    maxLines      int
    font          font.Face
    bgColor       color.RGBA
    textColor     map[LogLevel]color.RGBA
}
```

### 5.2 Overlay Features

**Visual Components:**
- Semi-transparent background
- Color-coded log levels
- Timestamp and system name display
- Auto-scroll to latest entries
- Scrollable history with keyboard controls

**Integration Points:**
- Hook into logger's memory buffer
- Render in main draw loop after game content
- Handle input in player input system

### 5.3 Overlay Controls

**Keyboard Shortcuts:**
- **F3** or **`**: Toggle overlay visibility
- **Page Up/Down**: Scroll through log history
- **F4**: Cycle log level filter (DEBUG → INFO → WARN → ERROR)
- **Ctrl+L**: Clear log buffer
- **Ctrl+S**: Save current log to file

### 5.4 Overlay Rendering

**Render Integration:**
```go
func ProcessDebugOverlay(g *Game, screen *ebiten.Image) {
    if !g.Components.Logger.DebugOverlay.visible {
        return
    }
    
    // Render semi-transparent background
    // Render log entries with color coding
    // Render scroll indicators
    // Render help text
}
```

### 5.5 Advanced Overlay Features

**Log Filtering:**
- Filter by system name
- Filter by log level
- Search functionality
- Highlight recent entries

**Performance Metrics:**
- FPS counter
- Memory usage display
- Entity count display
- System timing information

## Implementation Timeline

### Week 1: Core Logger (Phase 1-2)
- Implement basic Logger struct and methods
- Add configuration system
- Create unit tests

### Week 2: System Integration (Phase 3)
- Replace all existing log.Fatal calls
- Add comprehensive logging to all systems
- Performance testing and optimization

### Week 3: Testing and Validation (Phase 4)
- Complete test coverage
- Integration testing
- Performance benchmarking
- Bug fixes and refinements

### Week 4: Debug Overlay (Phase 5)
- Implement overlay infrastructure
- Add rendering and input handling
- Advanced features and polish
- Final testing and documentation

## Success Criteria

1. **Centralized Logging**: All systems use centralized logger
2. **Performance**: No measurable impact on game performance
3. **Configurability**: Easy to enable/disable logging levels
4. **Development Experience**: Clear visibility into system operations
5. **Debug Overlay**: Real-time log viewing in-game
6. **Maintainability**: Clean, well-tested logging infrastructure

## Backward Compatibility

All changes will maintain existing game functionality. The logging system is purely additive and will not affect game logic or user experience when disabled.

## Dependencies

- No new external dependencies required
- Uses existing Go standard library (io, sync, time, os)
- Integrates with existing ECS and rendering systems