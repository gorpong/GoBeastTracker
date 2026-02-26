package game

import (
	"fmt"
	"math/rand"
	"time"

	"beasttracker/internal/dungeon"
	"beasttracker/internal/entity"
	"beasttracker/internal/fov"
	"beasttracker/internal/ui"
)

const (
	monstersPerRoom = 2
	playerFOVRadius = 8
	maxMessages     = 5
	maxItemsPerRoom = 2
)

// GameStateType represents the current state of the game
type GameStateType int

const (
	StatePlaying GameStateType = iota
	StateGameOver
	StateVictory
)

// InputMode represents the current input context
type InputMode int

const (
	InputModeNormal    InputMode = iota
	InputModeDropping            // Waiting for slot number or menu trigger
	InputModeDropMenu            // Showing drop menu
	InputModeInventory           // Showing inventory (future use)
)

// Game holds all game state
type Game struct {
	Width     int
	Height    int
	Player    *entity.Player
	Dungeon   *dungeon.Dungeon
	Monsters  []*entity.Monster
	Items     []*entity.Item
	FOV       *fov.FOVMap
	Running   bool
	Seed      int64
	GameState GameStateType
	Messages  []string
	InputMode InputMode
}

// NewGame creates a new game with the specified dimensions and RNG seed
func NewGame(width, height int, seed int64) *Game {
	rng := rand.New(rand.NewSource(seed))

	generatedDungeon := dungeon.GenerateDungeon(width, height, seed)

	var playerX, playerY int
	if len(generatedDungeon.Rooms) > 0 {
		playerX, playerY = generatedDungeon.Rooms[0].Center()
	} else {
		playerX = width / 2
		playerY = height / 2
	}

	newGame := &Game{
		Width:     width,
		Height:    height,
		Player:    entity.NewPlayer(playerX, playerY),
		Dungeon:   generatedDungeon,
		Monsters:  make([]*entity.Monster, 0),
		Items:     make([]*entity.Item, 0),
		FOV:       fov.NewFOVMap(width, height),
		Running:   true,
		Seed:      seed,
		GameState: StatePlaying,
		Messages:  make([]string, 0),
		InputMode: InputModeNormal,
	}

	newGame.spawnMonsters(rng)
	newGame.spawnBoss(rng)
	newGame.spawnItems(rng)
	newGame.ComputeFOV()

	return newGame
}

// ComputeFOV calculates the field of view from the player's current position
func (g *Game) ComputeFOV() {
	px, py := g.Player.Position()
	fov.Compute(g.FOV, g.Dungeon, px, py, playerFOVRadius)
}

// IsVisible returns true if the tile at (x, y) is currently visible to the player
func (g *Game) IsVisible(x, y int) bool {
	return g.FOV.IsVisible(x, y)
}

// IsExplored returns true if the tile at (x, y) has been explored by the player
func (g *Game) IsExplored(x, y int) bool {
	return g.FOV.IsExplored(x, y)
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

// GetBoss returns the boss monster, or nil if none exists
func (g *Game) GetBoss() *entity.Monster {
	for _, monster := range g.Monsters {
		if monster.IsBoss {
			return monster
		}
	}
	return nil
}

// spawnBoss spawns a boss monster in the last (furthest) room
func (g *Game) spawnBoss(rng *rand.Rand) {
	if len(g.Dungeon.Rooms) < 2 {
		return // Need at least 2 rooms
	}

	// Define boss types
	bossTypes := []struct {
		name   string
		glyph  rune
		hp     int
		attack int
	}{
		{"Wyvern", 'W', 80, 12},
		{"Ogre", 'O', 100, 10},
		{"Troll", 'T', 90, 11},
		{"Cyclops", 'C', 85, 13},
		{"Minotaur", 'M', 95, 12},
	}

	// Spawn in the last room (furthest from player)
	lastRoom := g.Dungeon.Rooms[len(g.Dungeon.Rooms)-1]
	cx, cy := lastRoom.Center()

	// Pick a random boss type
	bossType := bossTypes[rng.Intn(len(bossTypes))]
	boss := entity.NewBossMonster(bossType.name, bossType.glyph, cx, cy, bossType.hp, bossType.attack)
	g.Monsters = append(g.Monsters, boss)
}

// spawnItems populates the dungeon with items
func (g *Game) spawnItems(rng *rand.Rand) {
	// Herbs are more common than potions to encourage exploration
	// while keeping strong healing rare
	itemTypes := []struct {
		itemType entity.ItemType
		weight   int
	}{
		{entity.ItemHerbs, 3},
		{entity.ItemPotion, 1},
	}

	totalWeight := 0
	for _, it := range itemTypes {
		totalWeight += it.weight
	}

	for _, room := range g.Dungeon.Rooms {
		numItems := rng.Intn(maxItemsPerRoom + 1)

		for j := 0; j < numItems; j++ {
			// Limit attempts to avoid infinite loops in small or crowded rooms
			for attempt := 0; attempt < 10; attempt++ {
				x := room.X + rng.Intn(room.Width)
				y := room.Y + rng.Intn(room.Height)

				if !g.isValidItemPosition(x, y) {
					continue
				}

				// Weighted random selection for item type
				roll := rng.Intn(totalWeight)
				var selectedType entity.ItemType
				for _, it := range itemTypes {
					roll -= it.weight
					if roll < 0 {
						selectedType = it.itemType
						break
					}
				}

				item := entity.NewItem(selectedType, x, y)
				g.Items = append(g.Items, item)
				break
			}
		}
	}
}

// isValidItemPosition ensures items don't block movement or stack
func (g *Game) isValidItemPosition(x, y int) bool {
	if !g.Dungeon.IsWalkable(x, y) {
		return false
	}

	px, py := g.Player.Position()
	if x == px && y == py {
		return false
	}

	if g.GetMonsterAt(x, y) != nil {
		return false
	}

	if g.GetItemAt(x, y) != nil {
		return false
	}

	return true
}

// GetItemAt returns the item at the specified position, or nil
func (g *Game) GetItemAt(x, y int) *entity.Item {
	for _, item := range g.Items {
		ix, iy := item.Position()
		if ix == x && iy == y {
			return item
		}
	}
	return nil
}

// RemoveItem removes an item from the game
func (g *Game) RemoveItem(item *entity.Item) {
	for i, it := range g.Items {
		if it == item {
			g.Items = append(g.Items[:i], g.Items[i+1:]...)
			return
		}
	}
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
	px, py := g.Player.Position()
	mx, my := monster.Position()

	// Check if adjacent to player - if so, attack instead of moving
	dx := px - mx
	dy := py - my
	if (dx == 1 || dx == -1) && dy == 0 || (dy == 1 || dy == -1) && dx == 0 {
		g.monsterAttack(monster)
		return
	}

	// 50% chance to move on any given turn
	if rng.Intn(2) == 0 {
		return
	}

	// Pick a random direction
	directions := []ui.Direction{ui.DirUp, ui.DirDown, ui.DirLeft, ui.DirRight}
	dir := directions[rng.Intn(len(directions))]

	ddx, ddy := dir.Delta()
	newX := monster.X + ddx
	newY := monster.Y + ddy

	// Check if target position is valid
	if !g.Dungeon.IsWalkable(newX, newY) {
		return
	}

	// Don't move into player
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
		g.UpdateMonsterAI()
	case ui.ActionDropMode:
		g.InputMode = InputModeDropping
	case ui.ActionInventory:
		g.InputMode = InputModeInventory
	}
}

// tryMovePlayer attempts to move the player in the given direction.
// Movement is blocked if it would hit a wall. If a monster is present, attack it.
func (g *Game) tryMovePlayer(dir ui.Direction) {
	dx, dy := dir.Delta()
	newX := g.Player.X + dx
	newY := g.Player.Y + dy

	if !g.Dungeon.IsWalkable(newX, newY) {
		return
	}

	if monster := g.GetMonsterAt(newX, newY); monster != nil {
		g.playerAttack(monster)
		return
	}

	g.Player.SetPosition(newX, newY)

	// Check for item pickup after moving
	g.tryPickupItem(newX, newY)

	g.ComputeFOV()
}

// tryPickupItem attempts to pick up an item at the given position
func (g *Game) tryPickupItem(x, y int) {
	item := g.GetItemAt(x, y)
	if item == nil {
		return
	}

	if g.Player.Inventory.IsFull() {
		g.AddMessage("Inventory full! Press 'x' to drop an item first.")
		return
	}

	g.Player.Inventory.Add(item)
	g.RemoveItem(item)
	g.AddMessage(fmt.Sprintf("Picked up %s.", item.Name()))
}

// UseItemInSlot uses the item in the specified inventory slot (1-indexed)
func (g *Game) UseItemInSlot(slot int) {
	item := g.Player.Inventory.GetSlot(slot)
	if item == nil {
		g.AddMessage(fmt.Sprintf("No item in slot %d.", slot))
		return
	}

	healAmount := item.HealingValue()
	if healAmount > 0 {
		g.Player.Heal(healAmount)
		g.AddMessage(fmt.Sprintf("Used %s. Restored %d HP.", item.Name(), healAmount))
	}

	g.Player.Inventory.Remove(slot)
}

// HandleDropModeInput processes input while in drop mode
func (g *Game) HandleDropModeInput(r rune) {
	// Check for menu triggers
	if r == 'x' || r == 'X' || r == 'i' || r == 'I' {
		g.InputMode = InputModeDropMenu
		return
	}

	// Check for slot number
	slot, ok := ui.ParseSlotNumber(r)
	if ok {
		g.dropItemFromSlot(slot)
		g.InputMode = InputModeNormal
		return
	}

	// Any other key cancels drop mode
	g.InputMode = InputModeNormal
}

// dropItemFromSlot drops an item from inventory to the ground at player position
func (g *Game) dropItemFromSlot(slot int) {
	item := g.Player.Inventory.GetSlot(slot)
	if item == nil {
		g.AddMessage(fmt.Sprintf("Slot %d is empty.", slot))
		return
	}

	g.Player.Inventory.Remove(slot)

	// Place item at player's feet
	px, py := g.Player.Position()
	item.X = px
	item.Y = py
	g.Items = append(g.Items, item)

	g.AddMessage(fmt.Sprintf("Dropped %s.", item.Name()))
}

// AddMessage adds a message to the message log
func (g *Game) AddMessage(msg string) {
	g.Messages = append(g.Messages, msg)
	// Keep only the last maxMessages
	if len(g.Messages) > maxMessages {
		g.Messages = g.Messages[len(g.Messages)-maxMessages:]
	}
}

// CalculateDamage calculates damage dealt based on attack and defense
func (g *Game) CalculateDamage(attack, defense int) int {
	damage := attack - defense
	if damage < 1 {
		damage = 1 // Minimum 1 damage
	}
	return damage
}

// playerAttack handles the player attacking a monster
func (g *Game) playerAttack(monster *entity.Monster) {
	damage := g.CalculateDamage(g.Player.Attack, 0) // Monsters have no defense for now
	monster.TakeDamage(damage)

	if monster.Dead {
		if monster.IsBoss {
			g.AddMessage(fmt.Sprintf("You have slain the %s! VICTORY!", monster.Name))
			g.GameState = StateVictory
		} else {
			g.AddMessage(fmt.Sprintf("You killed the %s!", monster.Name))
		}
		g.RemoveDeadMonsters()
	} else {
		g.AddMessage(fmt.Sprintf("You hit the %s for %d damage.", monster.Name, damage))
	}
}

// monsterAttack handles a monster attacking the player
func (g *Game) monsterAttack(monster *entity.Monster) {
	damage := g.CalculateDamage(monster.Attack, g.Player.Defense)
	g.Player.TakeDamage(damage)

	g.AddMessage(fmt.Sprintf("The %s hits you for %d damage!", monster.Name, damage))

	if g.Player.Dead {
		g.GameState = StateGameOver
		g.AddMessage("You have been slain!")
	}
}

// CheckPlayerDeath checks if the player is dead and updates game state
func (g *Game) CheckPlayerDeath() {
	if g.Player.Dead && g.GameState != StateGameOver {
		g.GameState = StateGameOver
	}
}
