# Codebase Reorganization Plan

## Current State Analysis

### Complexity Issues
- **25 Go files** in root directory
- **Mixed concerns** - systems, components, utilities all together
- **No clear separation** between game logic layers
- **Difficult navigation** for new contributors
- **Testing scattered** throughout codebase

### Current File Categories
1. **Entry Point**: `main.go`
2. **ECS Core**: `components.go`, `world.go`
3. **Game Systems**: `*_system.go` files (6 files)
4. **World Management**: `level.go`, `dungeon.go`, `map.go`
5. **Algorithms**: `astar.go`
6. **Utilities**: `dice.go`, `rect.go`, `gamedata.go`, `render_pool.go`, `turnstate.go`
7. **Tests**: `*_test.go` files (6 files)
8. **Assets**: `assets/`, `fonts/`
9. **Documentation**: `plans/`, `*.md` files

## Proposed Directory Structure

### Option A: Standard Go Project Layout (Recommended)

```
rrogue/
├── cmd/
│   └── rrogue/
│       └── main.go                 # Application entry point
├── internal/                       # Private application code
│   ├── components/
│   │   ├── components.go           # ECS component definitions
│   │   └── components_test.go
│   ├── systems/
│   │   ├── combat.go               # combat_system.go
│   │   ├── combat_test.go
│   │   ├── hud.go                  # hud_system.go
│   │   ├── monster.go              # monster_systems.go
│   │   ├── player.go               # player_systems.go
│   │   ├── render.go               # render_system.go
│   │   ├── userlog.go              # userlog_system.go
│   │   └── systems_test.go
│   ├── world/
│   │   ├── world.go                # ECS world initialization
│   │   ├── level.go                # Level generation and management
│   │   ├── level_test.go
│   │   ├── dungeon.go              # Dungeon container
│   │   ├── map.go                  # Game map management
│   │   └── turnstate.go            # Game state management
│   ├── engine/
│   │   ├── gamedata.go             # Game configuration
│   │   ├── render_pool.go          # Memory pool optimization
│   │   └── pathfinding.go          # astar.go
│   ├── ui/
│   │   ├── screens.go              # Future: menu screens
│   │   └── input.go                # Future: input handling
│   └── utils/
│       ├── dice.go                 # Random number generation
│       ├── dice_test.go
│       ├── rect.go                 # Rectangle utilities
│       └── rect_test.go
├── assets/                         # Game assets (unchanged)
│   ├── *.png
│   └── fonts/
├── docs/                           # Documentation
│   ├── ALGORITHMS.md
│   └── README.md
├── plans/                          # Development plans
│   └── *.md
├── go.mod
├── go.sum
├── CLAUDE.md
└── LICENSE
```

### Option B: Simplified Structure (Alternative)

```
rrogue/
├── main.go                         # Keep in root for simplicity
├── components/
│   ├── components.go
│   └── components_test.go
├── systems/
│   ├── *.go (all system files)
│   └── *_test.go
├── world/
│   ├── world.go
│   ├── level.go
│   ├── dungeon.go
│   ├── map.go
│   └── turnstate.go
├── engine/
│   ├── gamedata.go
│   ├── render_pool.go
│   └── pathfinding.go
├── utils/
│   ├── dice.go
│   ├── rect.go
│   └── *_test.go
├── assets/
├── docs/
└── plans/
```

## Detailed Migration Plan

### Phase 1: Create Directory Structure

**Week 1: Foundation Setup**

1. **Create directories**:
   ```bash
   mkdir -p cmd/rrogue
   mkdir -p internal/{components,systems,world,engine,utils,ui}
   mkdir -p docs
   ```

2. **Move documentation**:
   ```bash
   mv ALGORITHMS.md docs/
   mv README.md docs/
   ```

### Phase 2: Move and Refactor Core Components

**Week 1-2: ECS Core**

1. **Move entry point**:
   ```bash
   mv main.go cmd/rrogue/main.go
   ```

2. **Move components**:
   ```bash
   mv components.go internal/components/
   mv components_test.go internal/components/
   ```

3. **Move world management**:
   ```bash
   mv world.go internal/world/
   mv level.go internal/world/
   mv level_test.go internal/world/
   mv dungeon.go internal/world/
   mv map.go internal/world/
   mv turnstate.go internal/world/
   mv turnstate_test.go internal/world/
   ```

### Phase 3: Reorganize Systems

**Week 2: Game Systems**

1. **Move and rename systems**:
   ```bash
   mv combat_system.go internal/systems/combat.go
   mv combat_test.go internal/systems/combat_test.go
   mv hud_system.go internal/systems/hud.go
   mv monster_systems.go internal/systems/monster.go
   mv player_systems.go internal/systems/player.go
   mv render_system.go internal/systems/render.go
   mv userlog_system.go internal/systems/userlog.go
   ```

2. **Update package declarations** in each system file:
   ```go
   package systems
   ```

### Phase 4: Engine and Utilities

**Week 2-3: Supporting Code**

1. **Move engine code**:
   ```bash
   mv gamedata.go internal/engine/
   mv render_pool.go internal/engine/
   mv astar.go internal/engine/pathfinding.go
   ```

2. **Move utilities**:
   ```bash
   mv dice.go internal/utils/
   mv dice_test.go internal/utils/
   mv rect.go internal/utils/
   mv rect_test.go internal/utils/
   ```

### Phase 5: Update Imports and Dependencies

**Week 3: Import Refactoring**

1. **Update main.go imports**:
   ```go
   package main
   
   import (
       "github.com/caustin/rrogue/internal/components"
       "github.com/caustin/rrogue/internal/systems"
       "github.com/caustin/rrogue/internal/world"
       "github.com/caustin/rrogue/internal/engine"
   )
   ```

2. **Update all internal imports** across packages

3. **Update go.mod module path** if needed

## Benefits of Reorganization

### For New Contributors
1. **Clear structure** - Easy to find related code
2. **Logical grouping** - Systems, components, utilities separated
3. **Standard Go layout** - Familiar to Go developers
4. **Reduced cognitive load** - Smaller files per directory

### For Maintenance
1. **Better testing** - Test files grouped with implementation
2. **Easier refactoring** - Related code grouped together
3. **Clearer dependencies** - Import structure shows relationships
4. **Package-level documentation** - Each package can have clear purpose

### For Development
1. **Feature isolation** - Changes contained within packages
2. **Parallel development** - Multiple developers can work on different packages
3. **Code reuse** - Utilities and engine code clearly separated
4. **Future expansion** - Easy to add new systems or components

## Go-Specific Best Practices Applied

### Package Naming
- **Short, descriptive names**: `systems`, `utils`, `engine`
- **No stuttering**: `systems.Combat` not `systems.CombatSystem`
- **Clear responsibility**: Each package has single purpose

### Import Organization
- **Internal packages**: Use `internal/` to prevent external imports
- **Logical grouping**: Related functionality in same package
- **Dependency direction**: Higher-level packages import lower-level ones

### File Organization
- **One concept per file**: `combat.go` contains only combat logic
- **Tests alongside code**: `combat_test.go` next to `combat.go`
- **Package documentation**: `doc.go` files for complex packages

### Interface Design
- **Small interfaces**: Following Go's "accept interfaces, return structs"
- **Package boundaries**: Well-defined APIs between packages
- **Dependency injection**: Clear component relationships

## Migration Checklist

### Pre-Migration
- [ ] Backup current codebase
- [ ] Ensure all tests pass
- [ ] Document current dependencies
- [ ] Plan import path updates

### During Migration
- [ ] Create directory structure
- [ ] Move files systematically (one package at a time)
- [ ] Update package declarations
- [ ] Fix import statements
- [ ] Run tests after each package migration
- [ ] Update documentation

### Post-Migration
- [ ] Run full test suite
- [ ] Update build scripts/commands
- [ ] Update CLAUDE.md with new structure
- [ ] Create package documentation
- [ ] Update contribution guidelines

## Recommended Timeline

### Week 1: Planning and Foundation
- Finalize directory structure
- Create directories
- Move documentation and assets

### Week 2: Core Migration
- Move ECS components and world
- Move game systems
- Update basic imports

### Week 3: Engine and Utilities
- Move engine code and utilities
- Complete import updates
- Full testing and validation

### Week 4: Documentation and Polish
- Update all documentation
- Create package docs
- Performance testing
- Final cleanup

## Alternative: Gradual Migration

If full reorganization seems too disruptive, consider **gradual migration**:

1. **Start with utilities** - Move `dice.go`, `rect.go` to `utils/`
2. **Group systems** - Move all `*_system.go` to `systems/`
3. **Separate tests** - Move tests to appropriate packages
4. **Final core migration** - Move components and world last

## Success Criteria

1. **All tests pass** after reorganization
2. **Clear package boundaries** with minimal circular dependencies
3. **Improved readability** for new contributors
4. **Maintained performance** - no regression in game performance
5. **Updated documentation** reflecting new structure

This reorganization will transform the codebase from a flat structure to a well-organized, maintainable Go project that follows community standards and best practices.