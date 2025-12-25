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
* Clear naming.
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
