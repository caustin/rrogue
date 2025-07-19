# User Log System and HUD Fixes Plan

## Current Issues Identified

### User Log System (`userlog_system.go`)

1. **Message History Loss**: The `lastText` slice gets completely replaced each frame (line 79), causing previous messages to be lost
2. **Single Frame Display**: Messages only appear when new ones arrive, then disappear the next frame  
3. **No Message Persistence**: Messages are cleared immediately after reading, preventing them from staying visible
4. **Limited Capacity**: Hard-coded limit of 5 messages with no scrolling capability
5. **Inefficient Message Clearing**: AttackMessage and DeadMessage are cleared immediately, preventing persistence across frames

### HUD System (`hud_system.go`)

1. **Performance Issue**: Calls `NewGameData()` every frame instead of using cached GameData from Game struct
2. **Font Loading**: Creates font face every frame instead of reusing existing font
3. **Inconsistent Font Usage**: Uses `hudFont` but draws with `mplusNormalFont` from userlog system
4. **Missing Log Integration**: HUD doesn't properly integrate with the user log display

## Proposed Solutions

### Phase 1: Fix Message Persistence

1. **Implement Message History Buffer**
   - Create a circular buffer or slice to store message history (configurable size, default 20-50 messages)
   - Replace immediate replacement with append operation
   - Add timestamp tracking for message aging

2. **Fix Message Display Logic**
   - Keep messages visible across multiple frames
   - Only clear messages after they've been displayed for sufficient time
   - Implement proper message lifecycle management

### Phase 2: Improve Message Management

1. **Enhanced UserMessage Component**
   - Add message priority levels
   - Add message timestamps
   - Add message categories (combat, movement, system, etc.)
   - Add message expiration settings

2. **Message Queue System**
   - Implement proper message queuing instead of immediate clearing
   - Add message deduplication to prevent spam
   - Support for temporary vs permanent messages

### Phase 3: UI/UX Improvements

1. **Scrolling Support**
   - Add keyboard controls for scrolling through message history
   - Visual indicators for available scroll content
   - Auto-scroll to latest messages

2. **Visual Enhancements**
   - Color coding for different message types
   - Message fade-out effects for old messages
   - Better spacing and formatting

### Phase 4: Performance and Integration Fixes

1. **HUD System Optimization**
   - Use cached GameData from Game struct instead of calling `NewGameData()`
   - Consolidate font loading between HUD and userlog systems
   - Fix font consistency issues

2. **Memory Management**
   - Implement message pooling to reduce allocations
   - Efficient string handling for message storage
   - Proper cleanup of old messages

## Implementation Steps

### Step 1: Message History Buffer
- Modify `userlog_system.go` to use persistent message storage
- Replace `lastText = tmpMessages` with append operation
- Add configurable message history size

### Step 2: Message Lifecycle
- Prevent immediate message clearing
- Add message timestamp tracking
- Implement proper message expiration

### Step 3: HUD Integration
- Fix performance issues in `hud_system.go`
- Consolidate font management
- Ensure proper GameData usage

### Step 4: Enhanced Features
- Add scrolling controls
- Implement message categorization
- Add visual improvements

## Testing Plan

1. **Unit Tests**
   - Message buffer operations
   - Message lifecycle management
   - Font loading and caching

2. **Integration Tests**
   - Message display across multiple frames
   - HUD and log system interaction
   - Performance impact measurement

3. **User Experience Testing**
   - Message visibility during gameplay
   - Scrolling functionality
   - Performance during high message volume

## Expected Outcomes

- **Persistent Message Log**: Players will see a complete history of game events
- **Better Performance**: Reduced allocations and improved rendering efficiency
- **Enhanced UX**: Scrollable message history with proper categorization
- **Consistent Display**: Reliable message visibility across game frames
- **Maintainable Code**: Cleaner separation of concerns between HUD and log systems

## Backward Compatibility

All changes will maintain existing component interfaces and game logic. The UserMessage component structure may be extended but will remain backward compatible with existing message creation code.