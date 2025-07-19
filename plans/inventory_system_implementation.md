# Inventory System Implementation Plan

## Current State Analysis

### Existing Equipment System
- **MeleeWeapon Component**: Name, MinimumDamage, MaximumDamage, ToHitBonus
- **Armor Component**: Name, Defense, ArmorClass
- **Fixed Equipment**: Players and monsters have hardcoded equipment in `world.go`
- **No Inventory**: Items are directly attached to entities as components
- **No Item Management**: Cannot equip/unequip or carry multiple items

### Current Equipment in Use
**Player Equipment:**
- Battle Axe (10-20 damage, +3 to hit)
- Plate Armor (15 defense, 18 AC)

**Monster Equipment:**
- Orcs: Rusty Sword (5-10 damage, +1 to hit), Leather Armor (3 defense, 6 AC)
- Skeletons: Bone Club (2-8 damage, +0 to hit), No Armor (1 defense, 4 AC)

## Proposed Inventory System Architecture

### Core Components

#### Item System
```go
type Item struct {
    ID          string      // Unique identifier
    Name        string      // Display name
    Description string      // Item description
    ItemType    ItemType    // Weapon, Armor, Consumable, etc.
    Rarity      Rarity      // Common, Uncommon, Rare, Epic
    Value       int         // Gold value
    Weight      int         // Encumbrance (future)
    Properties  ItemProperties // Type-specific data
}

type ItemType int
const (
    ItemTypeWeapon ItemType = iota
    ItemTypeArmor
    ItemTypePotion
    ItemTypeScroll
    ItemTypeMisc
)

type Rarity int
const (
    RarityCommon Rarity = iota
    RarityUncommon
    RarityRare
    RarityEpic
    RarityLegendary
)
```

#### Equipment Slots
```go
type EquipmentSlot int
const (
    SlotMainHand EquipmentSlot = iota
    SlotOffHand  // Future: shields, dual wielding
    SlotArmor
    SlotRing1    // Future: rings
    SlotRing2
    SlotAmulet   // Future: amulets
)

type Equipment struct {
    Slots map[EquipmentSlot]*Item
}
```

#### Inventory Component
```go
type Inventory struct {
    Items       []*Item
    MaxCapacity int
    Gold        int
}

// Core inventory operations
func (inv *Inventory) AddItem(item *Item) bool
func (inv *Inventory) RemoveItem(itemID string) *Item
func (inv *Inventory) HasItem(itemID string) bool
func (inv *Inventory) GetItemsByType(itemType ItemType) []*Item
func (inv *Inventory) GetCapacity() (used, max int)
```

### Item Type Implementations

#### Weapon Properties
```go
type WeaponProperties struct {
    MinDamage   int
    MaxDamage   int
    ToHitBonus  int
    WeaponType  WeaponType // Sword, Axe, Club, etc.
    TwoHanded   bool
}

type WeaponType int
const (
    WeaponTypeSword WeaponType = iota
    WeaponTypeAxe
    WeaponTypeClub
    WeaponTypeDagger
    WeaponTypeBow // Future: ranged weapons
)
```

#### Armor Properties
```go
type ArmorProperties struct {
    Defense     int
    ArmorClass  int
    ArmorType   ArmorType
    Slot        EquipmentSlot
}

type ArmorType int
const (
    ArmorTypeCloth ArmorType = iota
    ArmorTypeLeather
    ArmorTypeChainMail
    ArmorTypePlateMail
)
```

#### Consumable Properties
```go
type ConsumableProperties struct {
    Effect      ConsumableEffect
    Magnitude   int
    Duration    int // 0 for instant effects
    Charges     int // For scrolls
}

type ConsumableEffect int
const (
    EffectHealHealth ConsumableEffect = iota
    EffectRestoreMana // Future: magic system
    EffectBoostStrength
    EffectBoostDexterity
    EffectCurePoison
    EffectTeleport
    EffectIdentify
)
```

## Implementation Phases

### Phase 1: Core Item and Inventory System

#### 1.1 Create Item Framework
**New Files:**
- `items.go` - Item struct, types, and core functionality
- `item_definitions.go` - Predefined items and item creation functions
- `inventory.go` - Inventory component and operations

**Item Creation System:**
```go
func CreateWeapon(name string, minDamage, maxDamage, toHit int, weaponType WeaponType) *Item
func CreateArmor(name string, defense, ac int, armorType ArmorType) *Item  
func CreatePotion(name string, effect ConsumableEffect, magnitude int) *Item
func CreateScroll(name string, effect ConsumableEffect, charges int) *Item
```

#### 1.2 Update Component System
**Modify `components.go`:**
```go
type Inventory struct {
    Items       []*Item
    MaxCapacity int
    Gold        int
}

type Equipment struct {
    Slots map[EquipmentSlot]*Item
}

// Deprecate direct weapon/armor components in favor of Equipment
```

**Update `world.go`:**
- Replace direct MeleeWeapon/Armor components with Equipment + Inventory
- Add inventory initialization for players and monsters

### Phase 2: Equipment Management

#### 2.1 Equipment System
**Equipment Operations:**
```go
func (eq *Equipment) EquipItem(item *Item, slot EquipmentSlot) *Item
func (eq *Equipment) UnequipItem(slot EquipmentSlot) *Item
func (eq *Equipment) GetEquippedItem(slot EquipmentSlot) *Item
func (eq *Equipment) GetWeaponStats() (minDmg, maxDmg, toHit int)
func (eq *Equipment) GetArmorStats() (defense, ac int)
```

#### 2.2 Combat System Integration
**Modify `combat_system.go`:**
- Replace direct component access with equipment-based stat calculation
- Add equipment-based damage and armor calculations

```go
func GetCombatStats(entity *ecs.QueryResult, components *Components) (weapon, armor stats)
```

### Phase 3: Item Interaction and UI

#### 3.1 Player Actions
**New Input Commands:**
- **I**: Open inventory screen
- **E**: Equip/use selected item
- **D**: Drop selected item
- **G**: Pick up item from ground

**Item Pickup System:**
```go
type GroundItems struct {
    Items []*Item
}

func PickupItem(g *Game, playerPos *Position) bool
func DropItem(g *Game, playerPos *Position, item *Item)
```

#### 3.2 Inventory UI
**Inventory Screen Layout:**
```
┌─────────────────────────────────────┐
│              INVENTORY              │
├─────────────────┬───────────────────┤
│ Equipment       │ Backpack (15/20)  │
│                 │                   │
│ Weapon:         │ [P] Health Potion │
│ Battle Axe      │ [S] Scroll of..   │
│                 │ [W] Steel Sword   │
│ Armor:          │ [A] Chain Mail    │
│ Plate Armor     │                   │
│                 │                   │
│ [E]quip [D]rop [U]se [ESC]lose     │
└─────────────────────────────────────┘
```

### Phase 4: Item Generation and Drops

#### 4.1 Item Database
**Predefined Items:**
```go
var WeaponTemplates = map[string]*Item{
    "rusty_sword": CreateWeapon("Rusty Sword", 3, 6, 0, WeaponTypeSword),
    "steel_sword": CreateWeapon("Steel Sword", 6, 12, 2, WeaponTypeSword),
    "battle_axe":  CreateWeapon("Battle Axe", 10, 20, 3, WeaponTypeAxe),
    // ... more weapons
}

var ArmorTemplates = map[string]*Item{
    "cloth_robe":   CreateArmor("Cloth Robe", 1, 10, ArmorTypeCloth),
    "leather_armor": CreateArmor("Leather Armor", 3, 12, ArmorTypeLeather),
    "plate_armor":  CreateArmor("Plate Armor", 15, 18, ArmorTypePlateMail),
    // ... more armor
}

var ConsumableTemplates = map[string]*Item{
    "health_potion": CreatePotion("Health Potion", EffectHealHealth, 20),
    "teleport_scroll": CreateScroll("Scroll of Teleport", EffectTeleport, 1),
    // ... more consumables
}
```

#### 4.2 Monster Drop System
**Drop Tables:**
```go
type DropTable struct {
    Items []DropEntry
}

type DropEntry struct {
    ItemID     string
    DropChance float64 // 0.0 to 1.0
    MinCount   int
    MaxCount   int
}

var MonsterDropTables = map[string]DropTable{
    "orc": {
        Items: []DropEntry{
            {"rusty_sword", 0.3, 1, 1},
            {"leather_armor", 0.2, 1, 1},
            {"health_potion", 0.1, 1, 1},
        },
    },
    "skeleton": {
        Items: []DropEntry{
            {"bone_club", 0.4, 1, 1},
            {"teleport_scroll", 0.05, 1, 1},
        },
    },
}
```

**Drop Generation:**
```go
func GenerateDrops(monsterType string) []*Item
func OnMonsterDeath(g *Game, monster *ecs.QueryResult, position *Position)
```

### Phase 5: Advanced Features

#### 5.1 Item Enhancement
- **Item Quality**: Poor, Normal, Superior, Exceptional
- **Prefixes/Suffixes**: "+1 Sword of Slaying", "Blessed Armor"
- **Randomized Stats**: Base item with random bonuses

#### 5.2 Consumable Effects
**Potion System:**
```go
func UsePotion(g *Game, entity *ecs.QueryResult, potion *Item)
func ApplyHealthPotion(health *Health, magnitude int)
func ApplyBuffPotion(entity *ecs.QueryResult, effect ConsumableEffect, magnitude, duration int)
```

**Scroll System:**
```go
func UseScroll(g *Game, caster *ecs.QueryResult, scroll *Item)
func CastTeleport(g *Game, caster *ecs.QueryResult)
func CastIdentify(g *Game, caster *ecs.QueryResult, target *Item)
```

#### 5.3 Item Management QoL
- **Auto-pickup**: Automatically pick up certain item types
- **Item stacking**: Stack identical consumables
- **Quick-use**: Hotkeys for potions/scrolls
- **Item comparison**: Compare stats when hovering over equipment

## Database Schema (Future: Save/Load)

### Item Persistence
```go
type SavedItem struct {
    ID          string                 `json:"id"`
    TemplateID  string                 `json:"template_id"`
    Properties  map[string]interface{} `json:"properties"`
    Charges     int                    `json:"charges,omitempty"`
}

type SavedInventory struct {
    Items    []SavedItem `json:"items"`
    Gold     int         `json:"gold"`
    Capacity int         `json:"capacity"`
}

type SavedEquipment struct {
    Slots map[string]SavedItem `json:"slots"`
}
```

## Integration Points

### Combat System Updates
```go
// Replace direct component access
attackerWeapon := attacker.Components[g.Components.MeleeWeapon].(*MeleeWeapon)
defenderArmor := defender.Components[g.Components.Armor].(*Armor)

// With equipment-based access  
attackerEquipment := attacker.Components[g.Components.Equipment].(*Equipment)
defenderEquipment := defender.Components[g.Components.Equipment].(*Equipment)

weaponStats := attackerEquipment.GetWeaponStats()
armorStats := defenderEquipment.GetArmorStats()
```

### HUD System Updates
```go
// Display equipped items in HUD
equipment := player.Components[g.Components.Equipment].(*Equipment)
weapon := equipment.GetEquippedItem(SlotMainHand)
armor := equipment.GetEquippedItem(SlotArmor)

if weapon != nil {
    weaponText := fmt.Sprintf("Weapon: %s", weapon.Name)
    // ... render weapon info
}
```

### World Generation Updates
```go
// Replace hardcoded equipment in entity creation
// Old:
.AddComponent(components.MeleeWeapon, &MeleeWeapon{...})
.AddComponent(components.Armor, &Armor{...})

// New:  
.AddComponent(components.Equipment, &Equipment{Slots: make(map[EquipmentSlot]*Item)})
.AddComponent(components.Inventory, &Inventory{MaxCapacity: 20})

// Then equip starting items
startingWeapon := CreateWeapon("Battle Axe", 10, 20, 3, WeaponTypeAxe)
equipment.EquipItem(startingWeapon, SlotMainHand)
```

## Implementation Timeline

### Week 1-2: Core Foundation (Phase 1)
- Implement Item, ItemType, and Inventory structs
- Create item definition system
- Add Inventory and Equipment components
- Basic item creation and management functions

### Week 3: Equipment System (Phase 2)  
- Implement equipment slot management
- Update combat system to use equipment stats
- Replace hardcoded weapon/armor components
- Update HUD to show equipped items

### Week 4: UI and Interaction (Phase 3)
- Create inventory screen
- Implement pickup/drop mechanics
- Add inventory navigation and item usage
- Ground item system

### Week 5-6: Drops and Generation (Phase 4)
- Create comprehensive item database
- Implement monster drop tables
- Add drop generation on monster death
- Test item progression and balance

### Week 7-8: Polish and Advanced Features (Phase 5)
- Add consumable effects (potions, scrolls)
- Implement item quality and randomization
- Add quality-of-life improvements
- Performance optimization and testing

## Success Criteria

1. **Complete Item System**: Full item creation, management, and persistence
2. **Equipment Management**: Equip/unequip items with proper stat calculation
3. **Inventory UI**: Intuitive inventory screen with full functionality
4. **Monster Drops**: Random item drops with configurable drop tables
5. **Consumable System**: Working potions and scrolls with various effects
6. **Performance**: No significant impact on game performance
7. **Save Compatibility**: Items properly serialize for save/load system

## Testing Strategy

### Unit Tests
- Item creation and property validation
- Inventory operations (add, remove, capacity)
- Equipment slot management
- Drop table generation
- Consumable effect application

### Integration Tests
- Combat system with equipment-based stats
- Monster death and item generation
- Player item pickup and usage
- Save/load compatibility with items

### Balance Testing
- Item progression curves
- Drop rates and item availability
- Combat impact of different equipment
- Economic balance (item values)

This inventory system will transform the game from fixed equipment to a full RPG-style item management system while maintaining compatibility with existing save systems and performance requirements.