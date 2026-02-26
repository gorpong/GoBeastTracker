# BeastTracker

A terminal-based ASCII roguelike game with Monster Hunter-inspired gameplay.
Hunt down powerful boss monsters in procedurally generated dungeons.

## Current Status

**Phases Complete: 0-8** (out of 16 planned phases)

The game is currently playable with core mechanics implemented:

- ✅ Procedural dungeon generation with rooms and corridors
- ✅ Player movement with fog of war
- ✅ Combat system (bump-to-attack)
- ✅ Regular monsters with wander AI
- ✅ Boss monsters with victory condition
- ✅ Game over and victory screens
- ✅ Item system with healing (herbs and potions)
- ✅ Inventory management with pickup/drop mechanics

## Game Mechanics

### Objective

Hunt down and defeat the **target boss monster** lurking in the deepest room of
the dungeon. The boss is displayed in purple and has significantly more HP than
regular monsters.

### Controls

- **Movement/Attack**: Arrow keys, HJKL (vi-style), or WASD
- **Use Item**: 1-9 (uses item in that inventory slot)
- **Inventory**: I (toggle inventory view)
- **Drop Item**: X + 1-9 (quick drop) or X + X (drop menu)
- **Quit**: Q or ESC

### Items

Items spawn randomly in dungeon rooms and are automatically picked up when walked over (if inventory has space).

| Item   | Glyph | Healing | Spawn Rate |
|--------|-------|---------|------------|
| Herbs  |  `"`  |  25 HP  |    75%     |
| Potion |  `!`  |  60 HP  |    25%     |

### Inventory

- **Capacity**: 6 slots
- **Display**: Shown in HUD as `[1:"] [2:!] [3:-]` etc.
- **Full Inventory**: Items remain on ground with warning message

### Combat

- **Bump-to-Attack**: Walk into an enemy to attack them
- **Damage Calculation**: `Damage = Attacker's ATK - Defender's DEF` (minimum 1 damage)
- **Monster Retaliation**: Monsters adjacent to the player will attack on their turn

### Player Stats

- **HP**: 100 (health points - game over when reaching 0)
- **ATK**: 10 (attack power)
- **DEF**: 2 (defense - reduces incoming damage)

### Monsters

#### Regular Monsters (Red)

Spawn 1-3 per room with random wander AI:

- **Goblin** (g): 15 HP, 3 ATK
- **Rat** (r): 8 HP, 2 ATK
- **Spider** (s): 12 HP, 2 ATK
- **Bat** (b): 10 HP, 2 ATK

#### Boss Monsters (Purple)

One boss spawns in the furthest room from the player:

- **Wyvern** (W): 80 HP, 12 ATK
- **Ogre** (O): 100 HP, 10 ATK
- **Troll** (T): 90 HP, 11 ATK
- **Cyclops** (C): 85 HP, 13 ATK
- **Minotaur** (M): 95 HP, 12 ATK

### Field of View (FOV)

- Player has 8-tile vision radius using shadowcasting algorithm
- Unexplored areas are completely hidden
- Previously explored areas are dimmed
- Monsters only visible when in line of sight

### HUD Information

- **Row 1**: Title, HP (color-coded: green/yellow/red), Boss target info
- **Row 2**: ATK/DEF stats, current position
- **Row 3**: Inventory bar showing items in slots 1-6
- **Bottom**: Latest message and context-sensitive controls

## Installation & Running

### Prerequisites

- Go 1.16 or higher
- Terminal with color support

### Build and Run

```bash
# Clone or navigate to the project directory
cd /path/to/beasttracker

# Install dependencies
go mod download

# Run the game
go run .

# Or build an executable
go build -o beasttracker
./beasttracker
```

## Development

### Project Structure

```text
beasttracker/
├── main.go                 # Entry point, game loop, rendering
├── internal/
│   ├── dungeon/           # Dungeon generation, tiles, rooms
│   │   ├── generator.go
│   │   ├── room.go
│   │   └── tile.go
│   ├── entity/            # Game entities (player, monsters)
│   │   ├── player.go
│   │   └── monster.go
│   ├── fov/               # Field of view (shadowcasting)
│   │   └── fov.go
│   ├── game/              # Game state, logic, combat
│   │   └── game.go
│   └── ui/                # Input handling, screen abstraction
│       ├── input.go
│       └── screen.go
├── PHASES.md              # Development roadmap
├── Claude.md              # AI assistant guidelines
└── README.md              # This file
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./internal/game/...

# Run tests with coverage
go test -cover ./...
```

### Test-Driven Development (TDD)

This project follows strict TDD:

1. **RED**: Write failing tests that define expected behavior
2. **GREEN**: Write minimal code to make tests pass
3. **REFACTOR**: Clean up code while keeping tests green

### Code Quality

- **Simple, idiomatic Go**: No clever tricks, clear control flow
- **Descriptive naming**: Variables >20 lines should not be single-character
- **Composition over inheritance**: Keep structs simple and composable
- **No premature optimization**: Implement what's needed, nothing more

## Gameplay Tips

1. **Explore Carefully**: Use fog of war to your advantage - monsters can't see you if you can't see them
2. **Pick Your Battles**: Regular monsters give you combat experience before facing the boss
3. **Watch Your HP**: The color-coded HP display warns you when you're in danger (green > yellow > red)
4. **Track the Boss**: The HUD shows boss HP so you can plan your approach
5. **Mind the Corridors**: Narrow passages can protect you from being surrounded
6. **Conserve Potions**: Herbs are common but weak; save potions for emergencies
7. **Manage Inventory**: Drop herbs to make room for potions when inventory is full

## Technical Details

### Dependencies

- **tcell/v2**: Terminal cell-based UI library for rendering and input

### Dungeon Generation

- Room-and-corridor algorithm
- 100×40 tile map with scrolling camera
- Deterministic seed-based generation for reproducibility

### FOV Algorithm

- Recursive shadowcasting in 8 octants
- Efficient line-of-sight with proper wall blocking
- Separate tracking of visible vs. explored tiles

### Combat System

- Turn-based: player moves, then all monsters move
- Simple damage formula with minimum damage guarantee
- Message log shows last 5 combat events

## Contributing

This is a learning project following a structured development plan. If you'd like to contribute:

1. Review [PHASES.md](PHASES.md) for current phase
2. Review [Claude.md](Claude.md) for development guidelines
3. Follow TDD principles
4. Keep changes focused and well-tested

## License

MIT (see [LICENSE.md](LICENSE.md))

## Credits

Developed as a Monster Hunter-inspired roguelike learning project.

Built with:

- [Go](https://golang.org/)
- [tcell](https://github.com/gdamore/tcell)
