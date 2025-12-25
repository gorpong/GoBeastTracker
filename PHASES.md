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

**Deliverable**: Running program that displays text in terminal and responds to quit command.

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

## Phase 2: Dungeon Generation ⏳

**Goal**: Procedurally generate dungeon layouts

### Tasks
- [ ] Create `dungeon/tile.go` - tile types (floor, wall, door)
- [ ] Create `dungeon/room.go` - room structure
- [ ] Create `dungeon/generator.go` - BSP or room-corridor algorithm
- [ ] Implement seed-based RNG for determinism
- [ ] Render dungeon to screen (walls '#', floors '.')
- [ ] Write tests for generation (fixed seed = fixed layout)

**Deliverable**: Procedurally generated 100×40 dungeon displayed on screen.

**Commit Message**: `Phase 2: Procedural dungeon generation`

---

## Phase 3: Player Movement & Collision ⏳

**Goal**: Physics and camera system

### Tasks
- [ ] Create `game/world.go` - dungeon container
- [ ] Implement collision detection (walls, boundaries)
- [ ] Player spawns in valid starting room
- [ ] Implement scrolling camera (viewport tracks player)
- [ ] Write tests for collision logic

**Deliverable**: Player navigates dungeon with proper collision, camera follows.

**Commit Message**: `Phase 3: Player collision and scrolling camera`

---

## Phase 4: Monsters & Basic AI ⏳

**Goal**: Populate dungeon with monsters

### Tasks
- [ ] Create `entity/monster.go` - monster struct
- [ ] Define 3-4 regular monster types (stats, glyphs)
- [ ] Implement spawn system (place monsters in rooms)
- [ ] Basic AI: wander randomly
- [ ] Monsters block player movement
- [ ] Write tests for AI behavior

**Deliverable**: Dungeon has wandering monsters that block movement.

**Commit Message**: `Phase 4: Monster spawning and basic AI`

---

## Phase 5: Field of View ⏳

**Goal**: Implement fog of war

### Tasks
- [ ] Create `fov/fov.go` - shadowcasting algorithm
- [ ] Compute FOV on player position change
- [ ] Render visible tiles normally
- [ ] Render "memory" tiles dimly (previously seen)
- [ ] Hide monsters outside FOV
- [ ] Write tests for FOV calculations

**Deliverable**: Fog of war with exploration memory.

**Commit Message**: `Phase 5: Field of view and fog of war`

---

## Phase 6: Combat System ⏳

**Goal**: Bump-to-attack combat

### Tasks
- [ ] Create `combat/combat.go` - damage calculation
- [ ] Create `entity/player.go` - player stats (HP, ATK, DEF)
- [ ] Implement bump-to-attack mechanic
- [ ] HP tracking and death handling
- [ ] Monster removal on death
- [ ] Player death → game over screen
- [ ] Combat message log
- [ ] Write tests for damage calculations

**Deliverable**: Functional combat with death consequences.

**Commit Message**: `Phase 6: Combat system and HP tracking`

---

## Phase 7: Target Monster & Win Condition ⏳

**Goal**: Boss monster and hunt victory

### Tasks
- [ ] Define 3-5 boss monster types (unique stats, glyphs)
- [ ] Spawn one boss per dungeon (target monster)
- [ ] Boss has 3-5x HP of regular monsters
- [ ] Win condition: defeat target monster
- [ ] Victory screen with hunt completion
- [ ] Write tests for boss spawning

**Deliverable**: Complete hunt with clear objective and victory.

**Commit Message**: `Phase 7: Boss monsters and win condition`

---

## Phase 8: Items & Healing ⏳

**Goal**: Item system with healing

### Tasks
- [ ] Create `entity/item.go` - item struct
- [ ] Define healing item types (potion, herbs)
- [ ] Spawn items in dungeon
- [ ] Inventory system (list, pickup on walk-over)
- [ ] Use item command (keybind)
- [ ] Create `ui/hud.go` - display HP and inventory
- [ ] Write tests for inventory operations

**Deliverable**: Player can find and use healing items.

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
- [ ] Write README.md (how to play, controls, gameplay loop)
- [ ] Code cleanup and documentation
- [ ] Final bug sweep
- [ ] Build instructions (cross-platform)
- [ ] Add LICENSE file
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
