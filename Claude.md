We are building a **terminal-based ASCII roguelike game in Go** with a **Monster Hunter–inspired gameplay loop**, designed to be **playable and fun** but not a full commercial game.

### Core goals

* **Turn-based roguelike**, rendered entirely in a terminal (ASCII / glyph grid).
* Written in **Go**, prioritizing **clarity, explicit control flow, and readable loops**.
* Use **simple, idiomatic Go**—no clever tricks, no premature abstractions.
* Prefer **composition over inheritance**.
* Avoid ECS unless it becomes obviously necessary because this is intended to be a small game

### Gameplay loop (Monster Hunter–like)

* The game is about **hunting specific monsters**, not clearing endless floors.
* Each dungeon run has:

 * One **primary target monster** (a “hunt”).
 * Optional smaller monsters.
 * Limited healing/resources.
 * The player wins by **tracking down and defeating the target monster**, not by grinding.

### Minimal feature set (MVP but fun)

Implement these fully and cleanly:

1. **Procedural dungeon generation**

   * Rooms + corridors are fine.
   * Deterministic via RNG seed.

2. **Player movement & collision**

   * Walls, monsters, bounds.

3. **Monsters**

   * Basic AI (wander, chase when in sight).
   * One special **boss-like target monster** with more HP and distinct behavior.

4. **Combat**

   * Bump-to-attack.
   * HP, damage, death.

5. **Field of View (FOV)**

   * Player has limited vision.
   * Previously seen tiles remain dim.

6. **Items**

   * Healing items.
   * Optional buffs (attack/defense).

7. **Monster drops & crafting**

   * Target monster drops a unique material.
   * Crafting or upgrading equipment is simple and menu-based.

8. **Win/Loss conditions**

   * Win: defeat the target monster.
   * Loss: player dies.
   * Simple Scoring System:
        * Each monster defeated: 10 points
        * Each boss monster defeated: 20 points
        * Each level cleared: 50 points
    * Leaderboard:
        * Top-10 high scores
        * If new score would go into top-10, ask for player's initials
        * Display leaderboard after game is over and on game splash screen while waiting for user to begin

### Scope constraints (important)

* This is **not** a full roguelike framework.
* No procedural lore.
* No complex skill trees.
* No animation beyond terminal redraws.
* No async or goroutines unless strictly necessary.
* Allow save/resume (if not too difficult).

### Technical constraints

* Use a **terminal UI library suitable for games** (e.g. `tcell`), but make a recommendation if `tcell` isn't your preferred library.
* Rendering should be grid-based, not widget-based.
* Separate **game state**, **logic**, and **rendering**, but keep it simple.
* All state should be explicit and inspectable.

### Code quality expectations

* Small, understandable structs.
* Clear naming - **IMPORTANT**: Avoid single-character variable names except:
  * Idiomatic Go receiver names (e.g., `func (g *Game)` is fine)
  * Very short functions (<20 lines) where context is obvious
  * Loop counters in tight loops (i, j, k)
  * In functions >20 lines, use descriptive names like `testGame`, `generatedDungeon`, `monster` instead of `g`, `d`, `m`
* Functions that do one thing.
* No hidden global state except where justified (e.g., RNG).
* Comments explaining *why*, not *what*.

### Development style

* Build iteratively.
* Each step should result in a **runnable game**.
* Favor “working and fun” over “perfect.”
* When tradeoffs exist, choose the simpler solution.
* Use test-driven development to create tests for each separate component with appropriate mocking as necessary.

### Test-Driven Development (TDD) Guidelines

**TDD is mandatory for this project. Follow this cycle strictly:**

1. **RED**: Write tests FIRST that define expected behavior. Tests are the specification.
2. **GREEN**: Write the minimal code necessary to make tests pass.
3. **REFACTOR**: Clean up code while keeping tests green.

**Key principles:**
- Tests define what the code SHOULD do, not what it currently does.
- Never write implementation code before its corresponding tests exist.
- Only modify a test if the test itself is incorrectly written—never to match flawed implementation.
- Tests should be aspirational: they describe the correct behavior, and the code must conform to them.
- When a test fails, fix the implementation, not the test (unless the test has a bug).

### Output expectations

* Start by proposing:

  1. A **high-level architecture** (files/modules).
  2. A **development plan broken into milestones**

* Then proceed step by step, writing real Go code.
* As each phase is complete, create a specific commit message
* If there are post-phase fixes, create a new branch and do each fix as a separate commit. Don't merge into main until I give final approval that all bugs are fixed and approve moving to the next phase.
* Ask clarifying questions up front as the start of the project and the start of each phase if you need clarification, but only ask if you think there's a specific decision that has a significant impact on the code creation for that phase.

The goal is **polished, playable terminal roguelike with a clear Monster Hunter-style hunt loop**.

---

## Lessons Learned (Phases 0-7)

### User Preferences
* **Combined Phases**: User prefers combining related phases (e.g., Phase 6 & 7: Combat + Boss) for efficiency when they're cohesive
* **Testing Approach**: User is comfortable testing multiple completed phases together rather than testing incrementally
* **Bug Reports**: User provides clear, actionable feedback about missing functionality (e.g., "monsters don't attack")
* **Communication Style**: User appreciates concise updates and summaries without unnecessary verbosity

### Implementation Patterns Established

#### Project Structure (Working Well)
```
internal/
  ├── dungeon/      # Tile-based map generation, rooms, corridors
  ├── entity/       # Player and Monster structs with combat stats
  ├── fov/          # Shadowcasting field of view
  ├── game/         # Game state, combat logic, AI updates
  └── ui/           # Input parsing, screen abstraction
main.go             # Rendering loop, HUD, victory/game over screens
```

#### Key Design Decisions
1. **Game State Enum**: Use `GameStateType` (StatePlaying, StateGameOver, StateVictory) for clear state transitions
2. **Message System**: Keep last 5 messages in a slice, display most recent in HUD
3. **Combat Flow**:
   - Player attacks on movement into monster (bump-to-attack)
   - Monsters attack during AI update when adjacent to player
   - **CRITICAL**: Must call attack function when adjacent, not just skip movement
4. **Boss Design**:
   - Single `IsBoss` boolean field on Monster struct
   - Spawn in last room (furthest from player)
   - 5-7x HP of regular monsters
   - Different visual style (purple vs red)
5. **FOV Integration**:
   - Compute on game start and after each player movement
   - Separate "visible" vs "explored" tracking
   - Hide monsters outside FOV, render explored tiles dimly

#### Testing Patterns
* Use `testGame` variable name in longer test functions for clarity
* Place monster at `px+1, py` to test adjacent combat
* Set monster HP to 1 for one-hit kill tests
* Use `CheckPlayerDeath()` explicitly in tests to trigger game state changes
* Mock map implements interface with `GetWidth()`, `GetHeight()`, `IsTransparent()`

#### Common Pitfalls to Avoid
1. **Monster AI**: Don't just make monsters avoid the player - make them attack when adjacent!
2. **Variable Naming**: In functions >20 lines, use `generatedDungeon` not `d`, `testGame` not `g`, `monster` not `m`
3. **Interface Methods**: When implementing interfaces like `fov.Map`, use methods not fields (e.g., `GetWidth()` not `Width`)
4. **Commit Messages**: Can combine related phases in one commit with clear sections

#### Effective TDD Workflow
1. Write test in `*_test.go` file (RED)
2. Run `go test ./internal/[package]/...` to confirm failure
3. Implement in corresponding `.go` file (GREEN)
4. Run tests again to confirm pass
5. Commit when phase complete, not after each test

#### UI/Rendering Learnings
* Reserve 2 rows at top for HUD (title, stats, boss info)
* Reserve 2 rows at bottom (messages, instructions)
* Color-code HP: green (>60%), yellow (>30%), red (≤30%)
* Boss in purple (`tcell.ColorPurple`), regular monsters in red
* ASCII art for game over and victory screens adds polish

### Phase Development Velocity
* **Phase 0-2**: Setup and basic rendering (slower, foundational)
* **Phase 3**: Merged into Phase 2 (wall collision)
* **Phase 4**: Monsters and AI (moderate pace)
* **Phase 5**: FOV system (moderate pace, complex algorithm)
* **Phase 6-7**: Combat + Boss (fast when combined, ~90 minutes)
* **Bug Fix**: Monster attack implementation (15 minutes)

### Future Phase Considerations
* **Phase 8** (Items): Will need inventory UI, item struct, pickup/use mechanics
* **Phase 9** (Crafting): Material drops on monster death, crafting menu, equipment system
* **Phase 10** (Enhanced AI): Chase behavior when player visible, boss special attacks
* **Scoring** (Phase 11): Points system, leaderboard JSON persistence
* **Save/Load** (Phase 14): Game state serialization - consider early if user wants to pause development

### Development Flow Preferences
* Start with TODO list for complex phases
* Mark tasks complete immediately after finishing (not in batches)
* Commit with detailed messages including implementation notes
* User testing happens after phase completion, not during
* Bug fixes get separate commits with "Fix:" prefix
