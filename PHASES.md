# BeastTracker - Development Phases

## Project Decisions

- **Dungeon Size**: 100×40 tiles with scrolling camera
- **Hunt Structure**: Sequential hunts with escalating difficulty
- **Progression**: Crafted equipment persists across hunts
- **Boss Variety**: 3-5 unique boss types with distinct behaviors
- **Boost Items**: Attack, Defense, Speed (extra actions)

---

## Phase Status Legend

- ⏳ Not Started
- 🚧 In Progress
- ✅ Complete
- 🐛 Bug Fixes

---

## Phase 0: Project Setup ✅

**Goal**: Initialize project structure and verify tooling

### Tasks

- [x] Initialize Go module (`go mod init`)
- [x] Add tcell dependency
- [x] Create directory structure (internal/, assets/)
- [x] Write tests for screen abstraction (TDD - RED)
- [x] Implement screen.go to pass tests (TDD - GREEN)
- [x] Create basic main.go with tcell splash screen
- [x] Verify terminal rendering works (manual test by user)
- [x] Create .gitignore

**Deliverable**: Running program that displays text in terminal and responds to 
quit command.

**Commit Message**: `Phase 0: Project setup and tcell initialization`

---

## Phase 1: Core Rendering & Input ✅

**Goal**: Basic grid rendering and player movement (with boundary collision)

### Tasks

- [x] Write tests for input handling - Direction, Action types (TDD - RED)
- [x] Create `ui/input.go` - keyboard event handling (TDD - GREEN)
- [x] Write tests for player position and movement (TDD - RED)
- [x] Create `entity/player.go` - player struct (TDD - GREEN)
- [x] Write tests for game state management (TDD - RED)
- [x] Create `game/game.go` - game loop and state (TDD - GREEN)
- [x] Handle movement input (arrow keys + vi keys: hjkl + WASD)
- [x] Handle quit (q/ESC)
- [x] Implement boundary collision (player can't leave screen)
- [x] Update main.go with game loop
- [x] Manual testing by user

**Deliverable**: Player '@' moves on screen with boundary collision.

**Commit Message**: `Phase 1: Core rendering and input system`

---

## Phase 2: Dungeon Generation ✅

**Goal**: Procedurally generate dungeon layouts

### Tasks

- [x] Write tests for tile types (TDD - RED)
- [x] Create `dungeon/tile.go` - tile types with walkability/transparency (TDD - GREEN)
- [x] Write tests for room structure (TDD - RED)
- [x] Create `dungeon/room.go` - room bounds, center, intersection (TDD - GREEN)
- [x] Write tests for dungeon generator (TDD - RED)
- [x] Create `dungeon/generator.go` - room-corridor algorithm (TDD - GREEN)
- [x] Implement seed-based RNG for determinism
- [x] Integrate dungeon into Game struct with wall collision
- [x] Implement scrolling camera (center on player)
- [x] Render dungeon to screen (walls '#', floors '.')
- [x] Player spawns in first room
- [x] Manual testing by user

**Deliverable**: Procedurally generated 100×40 dungeon with scrolling camera.

**Commit Message**: `Phase 2: Procedural dungeon generation`

---

## Phase 3: Player Movement & Collision ✅ (Merged into Phase 2)

**Note**: Wall collision and camera scrolling were implemented as part of Phase 2.

All tasks completed in Phase 2.

---

## Phase 4: Monsters & Basic AI ✅

**Goal**: Populate dungeon with monsters

### Tasks

- [x] Write tests for Monster structure (TDD - RED)
- [x] Create `entity/monster.go` - monster struct with HP, attack, AI type (TDD - GREEN)
- [x] Define 4 regular monster types (Goblin, Rat, Spider, Bat)
- [x] Write tests for monster spawning in game (TDD - RED)
- [x] Implement spawn system - 1-3 monsters per room (TDD - GREEN)
- [x] Monsters block player movement
- [x] Monsters render in red on screen
- [x] Write tests for wander AI behavior (TDD - RED)
- [x] Implement wander AI - random movement (TDD - GREEN)
- [x] Integrate AI updates into game loop (monsters move after player)
- [x] Manual testing by user

**Deliverable**: Dungeon populated with wandering monsters that block movement.

**Commit Message**: `Phase 4: Monster spawning and basic AI`

---

## Phase 5: Field of View ✅

**Goal**: Implement fog of war

### Tasks

- [x] Write tests for FOV calculations (TDD - RED)
- [x] Create `fov/fov.go` - shadowcasting algorithm (TDD - GREEN)
- [x] Compute FOV on player position change
- [x] Render visible tiles normally
- [x] Render "memory" tiles dimly (previously seen)
- [x] Hide monsters outside FOV
- [x] Manual testing by user

**Deliverable**: Fog of war with exploration memory.

**Commit Message**: `Phase 5: Field of view and fog of war`

---

## Phase 6: Combat System ✅

**Goal**: Bump-to-attack combat

### Tasks

- [x] Write tests for player stats (TDD - RED)
- [x] Add player stats (HP, ATK, DEF) to entity/player.go (TDD - GREEN)
- [x] Write tests for combat mechanics (TDD - RED)
- [x] Implement bump-to-attack mechanic (TDD - GREEN)
- [x] Implement damage calculation (attack - defense, min 1)
- [x] HP tracking and death handling
- [x] Monster removal on death
- [x] Player death triggers game over state
- [x] Combat message log
- [x] HP display with color coding (green/yellow/red)
- [x] Manual testing by user

**Deliverable**: Functional combat with death consequences.

**Commit Message**: `Phase 6: Combat system and HP tracking`

---

## Phase 7: Target Monster & Win Condition ✅

**Goal**: Boss monster and hunt victory

### Tasks

- [x] Write tests for boss monsters (TDD - RED)
- [x] Define 5 boss monster types (Wyvern, Ogre, Troll, Cyclops, Minotaur)
- [x] Add IsBoss field to Monster struct (TDD - GREEN)
- [x] Write tests for boss spawning and victory (TDD - RED)
- [x] Spawn one boss per dungeon in last room (TDD - GREEN)
- [x] Boss has 5-7x HP of regular monsters
- [x] Win condition: defeat target monster
- [x] Victory screen with hunt completion
- [x] Game over screen on player death
- [x] Boss displayed in purple, regular monsters in red
- [x] HUD shows boss target info
- [x] Manual testing by user

**Deliverable**: Complete hunt with clear objective and victory.

**Commit Message**: `Phase 7: Boss monsters and win condition`

---

## Phase 8: Items & Healing ✅

**Goal**: Item system with healing

### Tasks

- [x] Create `entity/item.go` - item struct with ItemType enum
- [x] Define healing item types (Herbs 25HP, Potion 60HP)
- [x] Create `entity/inventory.go` - inventory management with capacity
- [x] Add Inventory field to Player with Heal() method
- [x] Spawn items in dungeon (0-2 per room, weighted random)
- [x] Implement auto-pickup on walk-over (if space available)
- [x] Inventory full message when no space
- [x] Use item command (number keys 1-9)
- [x] Drop mode ('x' + number or 'x' + 'x' for menu)
- [x] Update `ui/input.go` - new actions for inventory/drop/use
- [x] Render items on map (herbs: `"`, potions: `!`)
- [x] Display inventory bar in HUD
- [x] Context-sensitive instructions based on input mode
- [x] Write tests for all new functionality

**Deliverable**: Player can find, pick up, use, and drop healing items.

**Commit Message**: `Phase 8: Item system and healing`

---

## Phase 9: Monster Drops & Crafting ⏳

**Goal**: Material drops and equipment crafting

### Tasks

- [ ] Monsters drop materials on death (scales, claws, etc.)
- [ ] Boss drops unique rare material
- [ ] Create `craft/crafting.go` - recipe system
- [ ] Crafting menu UI (select recipe, craft item)
- [ ] Equipment types: Weapon (+ATK), Armor (+DEF), Charm (+SPD)
- [ ] Equipment affects player stats
- [ ] Write tests for crafting logic

**Deliverable**: Hunt → materials → craft → stronger equipment.

**Commit Message**: `Phase 9: Monster drops and crafting system`

---

## Phase 10: Enhanced AI ⏳

**Goal**: Smarter monster behavior

### Tasks

- [ ] Implement chase AI (pursue player when in FOV)
- [ ] Boss-specific behaviors (aggressive, teleport, summon)
- [ ] AI respects FOV (no cheating)
- [ ] Different monster archetypes (aggressive, defensive, fleeing)
- [ ] Write tests for AI state transitions

**Deliverable**: Dynamic, challenging combat encounters.

**Commit Message**: `Phase 10: Enhanced monster AI`

---

## Phase 11: Scoring & Leaderboard ⏳

**Goal**: Score tracking and high score persistence

### Tasks

- [ ] Create `score/leaderboard.go` - scoring system
- [ ] Score calculation:
  - Regular monster: 10 points
  - Boss monster: 20 points
  - Hunt completion: 50 points
- [ ] Persist top-10 scores to `assets/data/scores.json`
- [ ] Initials entry for high scores
- [ ] Display leaderboard on splash screen and game over
- [ ] Write tests for scoring logic

**Deliverable**: Competitive replayability with leaderboard.

**Commit Message**: `Phase 11: Scoring and leaderboard system`

---

## Phase 12: Hunt Progression ⏳

**Goal**: Sequential hunts with difficulty scaling

### Tasks

- [ ] After boss defeat, offer "Next Hunt" option
- [ ] Persist player equipment across hunts
- [ ] Scale difficulty: more/stronger monsters, tougher bosses
- [ ] Track hunt number and display progress
- [ ] Create meta-progression screen (equipment, hunt history)
- [ ] Write tests for difficulty scaling

**Deliverable**: Multi-hunt campaign with escalating challenge.

**Commit Message**: `Phase 12: Hunt progression and difficulty scaling`

---

## Phase 13: Menus & Polish ⏳

**Goal**: Complete UI flow

### Tasks

- [ ] Create `ui/menu.go` - menu system
- [ ] Main menu: New Game, Continue, Leaderboard, Quit
- [ ] Pause menu (ESC during game)
- [ ] Quit confirmation dialog (protect against accidental 'q' press)
- [ ] Crafting menu improvements (show stats, materials)
- [ ] Color scheme and visual polish (tcell colors)
- [ ] Victory screen with score breakdown
- [ ] Death screen with score and retry option

**Deliverable**: Professional menu flow and presentation.

**Commit Message**: `Phase 13: Menu system and UI polish`

---

## Phase 14: Save/Load ⏳

**Goal**: Game state persistence

### Tasks

- [ ] Create `save/save.go` - serialization system
- [ ] Save game state to JSON
- [ ] Save on quit or manual save command
- [ ] Load game from main menu
- [ ] Handle corrupted save files gracefully
- [ ] Write tests for serialization round-trip

**Deliverable**: Players can resume interrupted hunts.

**Commit Message**: `Phase 14: Save and load system`

---

## Phase 15: Balance & Playtesting ⏳

**Goal**: Tune gameplay feel

### Tasks

- [ ] Balance monster HP/damage values
- [ ] Tune item spawn rates
- [ ] Adjust boss difficulty curve
- [ ] Test crafting progression (is grind reasonable?)
- [ ] Performance check (large dungeons, many monsters)
- [ ] Bug fixes from playtesting

**Deliverable**: Fair, fun, polished gameplay.

**Commit Message**: `Phase 15: Balance tuning and bug fixes`

---

## Phase 16: Final Polish & Documentation ⏳

**Goal**: Ship-ready state

### Tasks

- [ ] Update README.md (how to play, controls, gameplay loop)
- [ ] Code cleanup and documentation
- [ ] Final bug sweep
- [ ] Build instructions (cross-platform)
- [ ] Create release build

**Deliverable**: Complete, documented, playable game.

**Commit Message**: `Phase 16: Final polish and documentation`

---

## Notes

- Each phase should result in a **runnable game**
- After phase completion, create specific commit
- Post-phase bugs go in new branch, separate commits per fix
- Merge to main only after approval
- Test-driven development throughout
