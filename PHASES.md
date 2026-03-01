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

## Phase 9: Monster Drops & Crafting ✅

**Goal**: Material drops and equipment crafting

### Tasks

- [x] Create `entity/material.go` - MaterialType enum and Material struct
- [x] Create `entity/material.go` - MaterialPouch for unlimited material storage
- [x] Define 4 common materials (Scales, Claws, Fangs, Hide)
- [x] Define 5 rare materials (one per boss type)
- [x] Create `entity/equipment.go` - Equipment struct with stat bonuses
- [x] Define 3 equipment slots (Weapon, Armor, Charm)
- [x] Add DropTable to Monster struct
- [x] Regular monsters: 50% chance to drop 0-1 common material
- [x] Boss monsters: guaranteed rare drop + common materials
- [x] Add MaterialPouch and equipment slots to Player
- [x] Add EffectiveAttack(), EffectiveDefense(), EffectiveMaxHP() to Player
- [x] Combat uses effective stats (equipment bonuses apply)
- [x] Create `craft/crafting.go` - Recipe system
- [x] Define 5 basic recipes (common materials only)
- [x] Define 5 boss recipes (require rare materials)
- [x] Crafting menu ('c' key) with recipe browser
- [x] Auto-equip crafted items (replace existing)
- [x] Material pickup on walk-over (auto-add to pouch)
- [x] Materials render on map (common: orange, rare: gold)
- [x] Material pouch display in HUD
- [x] Visible monster HP display in HUD (up to 4 nearby monsters)
- [x] Write tests for all new functionality

### Recipes Implemented

**Basic (Common Materials):**

- Iron Sword: 3 Scales, 2 Claws → +3 ATK
- Bone Knife: 2 Fangs, 1 Scales → +2 ATK
- Leather Armor: 4 Hide, 1 Fangs → +2 DEF, +10 MaxHP
- Hide Vest: 3 Hide → +1 DEF, +5 MaxHP
- Hunter's Charm: 2 Fangs, 2 Claws → +1 ATK, +1 DEF

**Boss (Require Rare Materials):**

- Wyvern Blade: 1 Wyvern Scale, 3 Scales → +6 ATK
- Ogre Armor: 1 Ogre Hide, 4 Hide → +5 DEF, +25 MaxHP
- Troll Gauntlets: 1 Troll Claw, 2 Claws → +4 ATK, +1 DEF
- Cyclops Monocle: 1 Cyclops Eye, 2 Fangs → +2 ATK, +2 DEF
- Minotaur Horn Helm: 1 Minotaur Horn, 3 Hide → +3 DEF, +15 MaxHP

**Deliverable**: Hunt → materials → craft → stronger equipment.

**Commit Message**: `Phase 9: Monster drops and crafting system`

---

## Phase 10: Enhanced AI ✅

**Goal**: Smarter monster behavior

### Tasks

- [x] Implement chase AI (pursue player when in FOV)
- [x] Boss-specific behaviors:
  - [x] Aggressive (Ogre, Troll): Double damage when below 50% HP
  - [x] Teleport (Wyvern): Blink near player with cooldown
  - [x] Summoner (Cyclops): Spawn minion monsters with cooldown
- [x] AI respects FOV (no cheating)
- [x] Different monster archetypes:
  - [x] Wander (Goblin, Bat): Random movement
  - [x] Chase (Wolf): Pursue when visible
  - [x] Aggressive (Spider): Always pursue
  - [x] Defensive (Slime): Retreat when low HP
  - [x] Fleeing (Rat): Run away, fight when cornered
- [x] Write tests for AI behaviors and boss abilities
- [x] Boss defense stat added

**Deliverable**: Dynamic, challenging combat encounters.

**Commit Message**: `Phases 10-12: Enhanced AI, Scoring System, and Hunt Progression`

---

## Phase 11: Scoring & Leaderboard ✅

**Goal**: Score tracking and high score persistence

### Tasks

- [x] Create `score/leaderboard.go` - scoring system
- [x] Score calculation:
  - Regular monster: 10 points
  - Boss monster: 50 points
- [x] Persist top-10 scores to `assets/data/scores.json`
- [x] Initials entry for high scores (3 letters)
- [x] Display leaderboard on splash screen
- [x] Score display in game HUD
- [x] Write tests for scoring logic

**Deliverable**: Competitive replayability with leaderboard.

**Commit Message**: `Phases 10-12: Enhanced AI, Scoring System, and Hunt Progression`

---

## Phase 12: Hunt Progression ✅

**Goal**: Sequential hunts with difficulty scaling

### Tasks

- [x] After boss defeat, offer "Next Hunt" option (N key)
- [x] Persist player equipment across hunts
- [x] Persist material pouch across hunts
- [x] Scale difficulty per hunt:
  - [x] Monster HP: +15% per hunt
  - [x] Monster ATK: +1 per 2 hunts
  - [x] Monsters per room: +1 per 2 hunts
- [x] Track hunt number and display in HUD
- [x] Cumulative score across hunt chain
- [x] Write tests for difficulty scaling and persistence

**Deliverable**: Multi-hunt campaign with escalating challenge.

**Commit Message**: `Phases 10-12: Enhanced AI, Scoring System, and Hunt Progression`

---

## Phase 13: Menus & Polish ⏳

**Goal**: Complete UI flow

### Tasks

- [ ] Main menu: New Game, Leaderboard, Quit
- [ ] Pause menu (ESC during game) with resume option
- [ ] Quit confirmation dialog
- [ ] Equipment display screen
- [ ] Color scheme polish
- [ ] Victory screen with detailed score breakdown
- [ ] Death screen with retry option

**Deliverable**: Professional menu flow and presentation.

---

## Phase 14: Save/Load ⏳

**Goal**: Game state persistence

### Tasks

- [ ] Create `save/save.go` - serialization system
- [ ] Save game state to JSON
- [ ] Auto-save between hunts
- [ ] Load game from main menu
- [ ] Handle corrupted save files gracefully
- [ ] Write tests for serialization round-trip

**Deliverable**: Players can resume interrupted hunts.

---

## Phase 15: Balance & Playtesting ⏳

**Goal**: Tune gameplay feel

### Tasks

- [ ] Reduce regular monster drop rates (currently too high)
- [ ] Balance monster HP/damage values
- [ ] Tune item spawn rates
- [ ] Adjust boss difficulty curve
- [ ] Test crafting progression
- [ ] Performance check
- [ ] Bug fixes from playtesting

**Deliverable**: Fair, fun, polished gameplay.

---

## Phase 16: Final Polish & Documentation ⏳

**Goal**: Ship-ready state

### Tasks

- [ ] Update README.md with final mechanics
- [ ] Code cleanup and documentation
- [ ] Final bug sweep
- [ ] Build instructions (cross-platform)
- [ ] Create release build

**Deliverable**: Complete, documented, playable game.

---

## Notes

- Each phase should result in a **runnable game**
- After phase completion, create specific commit
- Post-phase bugs go in new branch, separate commits per fix
- Merge to main only after approval
- Test-driven development throughout
