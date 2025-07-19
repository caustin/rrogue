# Complex Algorithms Documentation

This document provides detailed documentation for the complex algorithms implemented in the RROGUE codebase.

## Table of Contents

1. [A* Pathfinding Algorithm](#a-pathfinding-algorithm)
2. [Dungeon Generation Algorithm](#dungeon-generation-algorithm)
3. [Rectangle Intersection Algorithm](#rectangle-intersection-algorithm)
4. [Manhattan Distance Calculation](#manhattan-distance-calculation)
5. [Combat Resolution Algorithm](#combat-resolution-algorithm)
6. [Random Number Generation](#random-number-generation)
7. [Monster AI Algorithm](#monster-ai-algorithm)
8. [Field of View Integration](#field-of-view-integration)

## A* Pathfinding Algorithm

**File:** `astar.go:63-210`  
**Function:** `AStar.GetPath(level Level, start *Position, end *Position) []Position`

### Purpose
Implements the A* search algorithm for finding optimal paths between two points on the game map, used primarily for monster AI pathfinding.

### Algorithm Overview
A* is a graph traversal and path search algorithm that finds the least-cost path from a given initial node to a goal node. It uses a heuristic function to guide its search.

### Implementation Details

#### Core Data Structures
```go
type node struct {
    Parent   *node      // Previous node in path
    Position *Position  // Current position
    g        int        // Distance from start
    h        int        // Heuristic distance to goal  
    f        int        // Total cost (g + h)
}
```

#### Algorithm Steps
1. **Initialize**: Create start node and add to open list
2. **Main Loop**: Continue until open list is empty
   - Select node with lowest f-score from open list
   - Move current node from open to closed list
   - Check if goal reached - if so, reconstruct path
   - Generate neighbor nodes (4-directional movement)
   - For each neighbor:
     - Skip if in closed list or blocked
     - Calculate g, h, and f scores
     - Add to open list if not present or path is better

#### Heuristic Function
Uses Manhattan distance: `|x1 - x2| + |y1 - y2|`

#### Path Reconstruction
- Follows parent pointers from goal back to start
- Reverses path using reflection-based `reverseSlice()` function

### Time Complexity
- **Best Case**: O(n) where n is the shortest path length
- **Average Case**: O(b^d) where b is branching factor, d is depth
- **Worst Case**: O(V + E) where V is vertices, E is edges

### Space Complexity
O(V) for storing open and closed lists

### Limitations
- Only supports 4-directional movement (cardinal directions)
- No diagonal movement support
- Calls `NewGameData()` every execution (performance impact)

### Optimization Opportunities
1. Cache GameData instance instead of creating new one
2. Add diagonal movement support
3. Implement binary heap for open list instead of linear search
4. Add path caching for frequently requested routes

---

## Dungeon Generation Algorithm

**File:** `level.go:101-149`  
**Function:** `Level.GenerateLevelTiles()`

### Purpose
Procedural generation of dungeon levels using a room-and-corridor approach.

### Algorithm Overview
Generates dungeons by:
1. Creating random rectangular rooms
2. Connecting rooms with L-shaped tunnels
3. Ensuring no room overlaps using intersection testing

### Implementation Details

#### Parameters
```go
MIN_SIZE := 6      // Minimum room dimension
MAX_SIZE := 10     // Maximum room dimension  
MAX_ROOMS := 30    // Maximum number of rooms to attempt
```

#### Generation Steps

1. **Initialize Map**: Create all-wall baseline using `createTiles()`

2. **Room Generation Loop**:
   ```go
   for idx := 0; idx < MAX_ROOMS; idx++ {
       // Generate random room dimensions and position
       w := GetRandomBetween(MIN_SIZE, MAX_SIZE)
       h := GetRandomBetween(MIN_SIZE, MAX_SIZE)
       x := GetDiceRoll(gd.ScreenWidth - w - 1)
       y := GetDiceRoll(levelHeight - h - 1)
   ```

3. **Collision Detection**: Check new room against all existing rooms using `Rect.Intersect()`

4. **Room Placement**: If no collision, carve out room using `createRoom()`

5. **Tunnel Creation**: Connect new room to previous room with L-shaped tunnel:
   ```go
   coinflip := GetDiceRoll(2)
   if coinflip == 2 {
       createHorizontalTunnel(prevX, newX, prevY)
       createVerticalTunnel(prevY, newY, newX)
   } else {
       createHorizontalTunnel(prevX, newX, newY)
       createVerticalTunnel(prevY, newY, prevX)
   }
   ```

#### Room Creation Algorithm
```go
func (level *Level) createRoom(room Rect) {
    for y := room.Y1 + 1; y < room.Y2; y++ {
        for x := room.X1 + 1; x < room.X2; x++ {
            index := level.GetIndexFromXY(x, y)
            level.Tiles[index].Blocked = false
            level.Tiles[index].TileType = FLOOR
            level.Tiles[index].Image = floor
        }
    }
}
```

#### Tunnel Creation Algorithm
- **Horizontal Tunnels**: `createHorizontalTunnel(x1, x2, y)`
- **Vertical Tunnels**: `createVerticalTunnel(y1, y2, x)`
- Both use `min/max` functions to handle bidirectional tunneling

### Time Complexity
O(R × R × A) where:
- R = number of rooms attempted (MAX_ROOMS)
- A = average room area for intersection testing

### Space Complexity
O(W × H) for tile array where W = width, H = height

### Algorithm Characteristics
- **Deterministic Structure**: Always produces connected dungeons
- **Variety**: Random room sizes and positions create unique layouts
- **Guaranteed Connectivity**: All rooms connected via tunnel system

### Potential Improvements
1. **Room Packing**: Better space utilization algorithms
2. **Tunnel Optimization**: Avoid unnecessary long tunnels
3. **Dead-end Removal**: Post-processing to remove unwanted dead ends
4. **Themed Areas**: Different room types or special rooms

---

## Rectangle Intersection Algorithm

**File:** `rect.go:27-29`  
**Function:** `Rect.Intersect(other Rect) bool`

### Purpose
Determines if two axis-aligned rectangles overlap, critical for dungeon room placement collision detection.

### Algorithm Implementation
```go
func (r *Rect) Intersect(other Rect) bool {
    return (r.X1 <= other.X2 && r.X2 >= other.X1 && 
            r.Y1 <= other.Y2 && r.Y2 >= other.Y1)
}
```

### Mathematical Foundation
Two rectangles intersect if they overlap in both X and Y dimensions:
- **X-axis overlap**: `r.X1 ≤ other.X2 AND r.X2 ≥ other.X1`
- **Y-axis overlap**: `r.Y1 ≤ other.Y2 AND r.Y2 ≥ other.Y1`

### Algorithm Analysis
- **Time Complexity**: O(1) - constant time
- **Space Complexity**: O(1) - no additional storage
- **Correctness**: Handles all edge cases including:
  - Complete overlap
  - Partial overlap
  - Edge touching
  - No intersection

### Visual Representation
```
Case 1: Intersection          Case 2: No Intersection
┌─────┐                       ┌─────┐
│  A  │                       │  A  │     ┌─────┐
│  ┌──┼─┐                     │     │     │  B  │
│  │XX│ │ B                   │     │     │     │
└──┼──┘ │                     └─────┘     └─────┘
   │  B │
   └────┘
```

### Testing Coverage
Covered by `rect_test.go` with comprehensive test cases for all intersection scenarios.

---

## Manhattan Distance Calculation

**File:** `components.go:16-20`  
**Function:** `Position.GetManhattanDistance(other *Position) int`

### Purpose
Calculates the Manhattan (taxicab) distance between two positions, used for movement cost estimation and proximity detection.

### Algorithm Implementation
```go
func (p *Position) GetManhattanDistance(other *Position) int {
    xDist := math.Abs(float64(p.X - other.X))
    yDist := math.Abs(float64(p.Y - other.Y))
    return int(xDist) + int(yDist)
}
```

### Mathematical Foundation
Manhattan distance = |x₁ - x₂| + |y₁ - y₂|

### Use Cases in Codebase
1. **A* Heuristic**: Used in `astar.go:186` for pathfinding heuristic
2. **Combat Range**: Used in `monster_systems.go:24` for melee attack detection
3. **AI Decision Making**: Determines if monster should attack or move

### Algorithm Characteristics
- **Time Complexity**: O(1)
- **Space Complexity**: O(1)
- **Admissible Heuristic**: Never overestimates actual path cost (perfect for A*)
- **Grid-Based**: Ideal for tile-based movement systems

### Comparison with Other Distance Metrics
- **Euclidean Distance**: √((x₁-x₂)² + (y₁-y₂)²) - not suitable for grid movement
- **Chebyshev Distance**: max(|x₁-x₂|, |y₁-y₂|) - allows diagonal movement

---

## Combat Resolution Algorithm

**File:** `combat_system.go:9-90`  
**Function:** `AttackSystem(g *Game, attackerPosition *Position, defenderPosition *Position)`

### Purpose
Resolves combat encounters between entities using dice-based mechanics with armor class and damage reduction.

### Algorithm Overview
Implements classic RPG combat mechanics:
1. Entity identification and validation
2. To-hit roll calculation
3. Damage calculation and application
4. Death handling and cleanup

### Implementation Details

#### Entity Resolution Phase
```go
// Find attacker and defender entities by position
for _, playerCombatant := range g.World.Query(g.WorldTags["players"]) {
    pos := playerCombatant.Components[g.Components.Position].(*Position)
    if pos.IsEqual(attackerPosition) {
        attacker = playerCombatant
    }
}
```

#### Combat Mechanics

**To-Hit Calculation**:
```go
toHitRoll := GetDiceRoll(10)  // 1d10
if toHitRoll + attackerWeapon.ToHitBonus > defenderArmor.ArmorClass {
    // Hit successful
}
```

**Damage Calculation**:
```go
damageRoll := GetRandomBetween(attackerWeapon.MinimumDamage, attackerWeapon.MaximumDamage)
damageDone := damageRoll - defenderArmor.Defense
if damageDone < 0 {
    damageDone = 0  // Prevent healing from weak attacks
}
defenderHealth.CurrentHealth -= damageDone
```

#### Death Handling Algorithm
```go
if defenderHealth.CurrentHealth <= 0 {
    if defenderName == "Player" {
        defenderMessage.GameStateMessage = "Game Over!"
        g.Turn = GameOver
    } else {
        // Monster cleanup: unblock tile and dispose entity
        level := g.Map.CurrentLevel
        pos := defender.Components[g.Components.Position].(*Position)
        tile := level.Tiles[level.GetIndexFromXY(pos.X, pos.Y)]
        tile.Blocked = false
        g.World.DisposeEntity(defender.Entity)
    }
}
```

### Combat Statistics
- **Hit Chance**: (1d10 + ToHitBonus) vs Armor Class
- **Damage Range**: WeaponMin to WeaponMax, reduced by Defense
- **Minimum Damage**: 0 (attacks cannot heal)

### Algorithm Complexity
- **Time Complexity**: O(P + M) where P = players, M = monsters
- **Space Complexity**: O(1) for combat calculation

### Design Patterns
- **Component Querying**: ECS pattern for entity resolution
- **Message System**: Combat results communicated via UserMessage components
- **State Management**: Game state transitions on player death

---

## Random Number Generation

**File:** `dice.go`  
**Functions:** `GetDiceRoll()`, `GetRandomBetween()`, `GetRandomInt()`

### Purpose
Provides cryptographically secure random number generation for game mechanics, replacing traditional pseudo-random generators.

### Implementation Details

#### Core Random Function
```go
func GetRandomInt(num int) int {
    x, _ := rand.Int(rand.Reader, big.NewInt(int64(num)))
    return int(x.Int64())
}
```

#### Dice Roll Function
```go
func GetDiceRoll(num int) int {
    x, _ := rand.Int(rand.Reader, big.NewInt(int64(num)))
    return int(x.Int64()) + 1  // Range: 1 to num inclusive
}
```

#### Range Function  
```go
func GetRandomBetween(low int, high int) int {
    return GetDiceRoll(high-low+1) + low - 1
}
```

### Cryptographic Security
- **Source**: `crypto/rand.Reader` - cryptographically secure
- **Quality**: Suitable for security-sensitive applications
- **Performance**: Slower than math/rand but higher quality

### Usage Patterns in Codebase
1. **Combat**: `GetDiceRoll(10)` for to-hit rolls
2. **Damage**: `GetRandomBetween(min, max)` for damage calculation  
3. **Dungeon Generation**: Room placement and tunnel direction
4. **AI Behavior**: Decision randomization

### Algorithm Characteristics
- **Uniform Distribution**: All values equally probable
- **Independence**: Each call independent of previous calls
- **Unpredictability**: Cannot be predicted from previous values

### Performance Considerations
- **Slower**: ~10x slower than `math/rand`
- **Quality**: Cryptographically secure randomness
- **Trade-off**: Security vs. performance (appropriate for game use)

---

## Monster AI Algorithm

**File:** `monster_systems.go:7-46`  
**Function:** `UpdateMonster(game *Game)`

### Purpose
Implements AI behavior for monsters including vision, pathfinding, and combat decisions.

### Algorithm Overview
1. **Player Location**: Identify player position
2. **Monster Processing**: For each monster, determine action based on player visibility
3. **Behavior Selection**: Attack, move toward player, or do nothing

### Implementation Details

#### Player Position Resolution
```go
playerPosition := Position{}
for _, plr := range game.World.Query(game.WorldTags["players"]) {
    pos := plr.Components[game.Components.Position].(*Position)
    playerPosition.X = pos.X
    playerPosition.Y = pos.Y
}
```

#### Monster Behavior Algorithm
```go
for _, result := range game.World.Query(game.WorldTags["monsters"]) {
    pos := result.Components[game.Components.Position].(*Position)
    
    // Calculate monster's field of view
    monsterSees := fov.New()
    monsterSees.Compute(l, pos.X, pos.Y, 8)  // 8-tile vision range
    
    if monsterSees.IsVisible(playerPosition.X, playerPosition.Y) {
        if pos.GetManhattanDistance(&playerPosition) == 1 {
            // Adjacent: Attack player
            AttackSystem(game, pos, &playerPosition)
        } else {
            // Not adjacent: Move toward player using A*
            astar := AStar{}
            path := astar.GetPath(l, pos, &playerPosition)
            if len(path) > 1 {
                // Move to next position in path
                nextTile := l.Tiles[l.GetIndexFromXY(path[1].X, path[1].Y)]
                if !nextTile.Blocked {
                    // Update tile blocking and monster position
                    l.Tiles[l.GetIndexFromXY(pos.X, pos.Y)].Blocked = false
                    pos.X = path[1].X
                    pos.Y = path[1].Y
                    nextTile.Blocked = true
                }
            }
        }
    }
}
```

### Behavioral States
1. **Idle**: Player not visible, no action
2. **Pursue**: Player visible but not adjacent, use A* pathfinding
3. **Attack**: Player adjacent, initiate combat

### AI Characteristics
- **Vision-Based**: Only acts when player is visible
- **Optimal Pathfinding**: Uses A* for intelligent movement
- **Immediate Combat**: Attacks when adjacent
- **Blocking Awareness**: Respects tile blocking rules

### Performance Considerations
- **FOV Calculation**: O(V²) where V is vision range
- **Pathfinding**: O(b^d) for A* algorithm
- **Monster Count**: Linear scaling with number of monsters

### Algorithm Complexity
- **Time Complexity**: O(M × (V² + P)) where M = monsters, V = vision range, P = pathfinding cost
- **Space Complexity**: O(M × V²) for FOV calculations

---

## Field of View Integration

**Library:** `github.com/norendren/go-fov`  
**Usage:** Player visibility and monster AI vision

### Purpose
Provides realistic line-of-sight calculations for both player exploration and monster AI behavior.

### Integration Points

#### Player FOV (`level.go:74`)
```go
isVis := level.PlayerVisible.IsVisible(x, y)
if isVis {
    // Render visible tiles normally
    level.Tiles[idx].IsRevealed = true
} else if tile.IsRevealed == true {
    // Render previously seen tiles with dimming
    op.ColorM.Translate(100, 100, 100, 0.35)
}
```

#### Monster FOV (`monster_systems.go:20-22`)
```go
monsterSees := fov.New()
monsterSees.Compute(l, pos.X, pos.Y, 8)
if monsterSees.IsVisible(playerPosition.X, playerPosition.Y) {
    // Monster can see player, take action
}
```

### Algorithm Features
- **Shadowcasting**: Efficient FOV calculation algorithm
- **Wall Blocking**: Respects opaque tiles (walls)
- **Range Limiting**: Configurable vision distance
- **Optimization**: Fast calculation for real-time use

### Level Integration Interface
```go
func (level Level) InBounds(x, y int) bool
func (level Level) IsOpaque(x, y int) bool
```

### Performance Characteristics
- **Algorithm**: Recursive shadowcasting
- **Time Complexity**: O(R²) where R is vision radius
- **Space Complexity**: O(R²) for visibility map
- **Real-time**: Suitable for per-frame calculation

### Game Mechanics Impact
- **Exploration**: Player discovers map through movement
- **Stealth**: Monsters only react when player is visible
- **Memory**: Previously seen areas remain faintly visible
- **Tactical**: Player can use walls to break line of sight

---

## Performance Analysis Summary

### Critical Path Algorithms
1. **A* Pathfinding**: Most expensive operation, called per monster per turn
2. **FOV Calculation**: Called per monster and per frame for player
3. **Combat Resolution**: Infrequent but complex entity queries

### Optimization Priorities
1. **Cache GameData**: Eliminate repeated `NewGameData()` calls
2. **A* Improvements**: Binary heap for open list, path caching
3. **FOV Optimization**: Cache monster FOV when not moving
4. **Render Pooling**: Already optimized with `render_pool.go`

### Memory Usage
- **Pathfinding**: Temporary node allocation during A* search
- **FOV**: Vision map storage per entity
- **Combat**: Minimal allocation, reuses existing components

### Scalability Considerations
- **Monster Count**: Linear growth in AI processing
- **Map Size**: Quadratic growth in FOV and pathfinding
- **Room Count**: Linear growth in generation time

---

This documentation provides comprehensive coverage of all complex algorithms in the RROGUE codebase, including implementation details, complexity analysis, and optimization opportunities.