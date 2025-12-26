package game

import (
	"math/rand"
	"time"

	"beasttracker/internal/dungeon"
	"beasttracker/internal/entity"
	"beasttracker/internal/ui"
)

const (
	monstersPerRoom = 2 // Average monsters per room
)

// Game holds all game state
type Game struct {
	Width    int
	Height   int
	Player   *entity.Player
	Dungeon  *dungeon.Dungeon
	Monsters []*entity.Monster
	Running  bool
	Seed     int64
}

// NewGame creates a new game with the specified dimensions and RNG seed
func NewGame(width, height int, seed int64) *Game {
	rng := rand.New(rand.NewSource(seed))

	// Generate dungeon
	generatedDungeon := dungeon.GenerateDungeon(width, height, seed)

	// Spawn player in the center of the first room
	var playerX, playerY int
	if len(generatedDungeon.Rooms) > 0 {
		playerX, playerY = generatedDungeon.Rooms[0].Center()
	} else {
		// Fallback to center if no rooms (shouldn't happen)
		playerX = width / 2
		playerY = height / 2
	}

	newGame := &Game{
		Width:    width,
		Height:   height,
		Player:   entity.NewPlayer(playerX, playerY),
		Dungeon:  generatedDungeon,
		Monsters: make([]*entity.Monster, 0),
		Running:  true,
		Seed:     seed,
	}

	// Spawn monsters in rooms (skip first room where player spawns)
	newGame.spawnMonsters(rng)

	return newGame
}

// spawnMonsters populates the dungeon with monsters
func (g *Game) spawnMonsters(rng *rand.Rand) {
	// Define monster types
	monsterTypes := []struct {
		name   string
		glyph  rune
		hp     int
		attack int
	}{
		{"Goblin", 'g', 15, 3},
		{"Rat", 'r', 8, 2},
		{"Spider", 's', 12, 2},
		{"Bat", 'b', 10, 2},
	}

	// Skip first room (player starts there)
	startRoom := 1
	if len(g.Dungeon.Rooms) <= 1 {
		startRoom = 0
	}

	for i := startRoom; i < len(g.Dungeon.Rooms); i++ {
		room := g.Dungeon.Rooms[i]

		// Spawn 1-3 monsters per room
		numMonsters := rng.Intn(3) + 1

		for j := 0; j < numMonsters; j++ {
			// Pick random position in room
			x := room.X + rng.Intn(room.Width)
			y := room.Y + rng.Intn(room.Height)

			// Ensure position is walkable and not occupied
			if !g.Dungeon.IsWalkable(x, y) {
				continue
			}
			if g.GetMonsterAt(x, y) != nil {
				continue
			}

			px, py := g.Player.Position()
			if x == px && y == py {
				continue
			}

			// Pick random monster type
			mType := monsterTypes[rng.Intn(len(monsterTypes))]
			monster := entity.NewMonster(mType.name, mType.glyph, x, y, mType.hp, mType.attack)
			g.Monsters = append(g.Monsters, monster)
		}
	}
}

// GetMonsterAt returns the monster at the specified position, or nil
func (g *Game) GetMonsterAt(x, y int) *entity.Monster {
	for _, monster := range g.Monsters {
		if !monster.Dead {
			mx, my := monster.Position()
			if mx == x && my == y {
				return monster
			}
		}
	}
	return nil
}

// RemoveDeadMonsters removes all dead monsters from the game
func (g *Game) RemoveDeadMonsters() {
	alive := make([]*entity.Monster, 0, len(g.Monsters))
	for _, monster := range g.Monsters {
		if !monster.Dead {
			alive = append(alive, monster)
		}
	}
	g.Monsters = alive
}

// UpdateMonsterAI updates all monster AI behaviors
func (g *Game) UpdateMonsterAI() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, monster := range g.Monsters {
		if monster.Dead {
			continue
		}

		switch monster.AI {
		case entity.AIWander:
			g.updateWanderAI(monster, rng)
		}
	}
}

// updateWanderAI makes a monster wander randomly
func (g *Game) updateWanderAI(monster *entity.Monster, rng *rand.Rand) {
	// 50% chance to move on any given turn
	if rng.Intn(2) == 0 {
		return
	}

	// Pick a random direction
	directions := []ui.Direction{ui.DirUp, ui.DirDown, ui.DirLeft, ui.DirRight}
	dir := directions[rng.Intn(len(directions))]

	dx, dy := dir.Delta()
	newX := monster.X + dx
	newY := monster.Y + dy

	// Check if target position is valid
	if !g.Dungeon.IsWalkable(newX, newY) {
		return
	}

	// Don't move into player
	px, py := g.Player.Position()
	if newX == px && newY == py {
		return
	}

	// Don't move into another monster
	if g.GetMonsterAt(newX, newY) != nil {
		return
	}

	monster.SetPosition(newX, newY)
}

// HandleInput processes player input and updates game state
func (g *Game) HandleInput(action ui.Action, dir ui.Direction) {
	switch action {
	case ui.ActionQuit:
		g.Running = false
	case ui.ActionMove:
		g.tryMovePlayer(dir)
		// After player moves, update monster AI
		g.UpdateMonsterAI()
	}
}

// tryMovePlayer attempts to move the player in the given direction.
// Movement is blocked if it would hit a wall, monster, or go out of bounds.
func (g *Game) tryMovePlayer(dir ui.Direction) {
	dx, dy := dir.Delta()
	newX := g.Player.X + dx
	newY := g.Player.Y + dy

	// Check if target position is walkable (includes bounds check)
	if !g.Dungeon.IsWalkable(newX, newY) {
		return
	}

	// Check for monster at target position
	if g.GetMonsterAt(newX, newY) != nil {
		return
	}

	g.Player.SetPosition(newX, newY)
}
