package game

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"beasttracker/internal/craft"
	"beasttracker/internal/dungeon"
	"beasttracker/internal/entity"
	"beasttracker/internal/fov"
	"beasttracker/internal/ui"
)

const (
	monstersPerRoom    = 2
	playerFOVRadius    = 8
	maxMessages        = 5
	maxItemsPerRoom    = 2
	maxVisibleMonsters = 4
)

type GameStateType int

const (
	StatePlaying GameStateType = iota
	StateGameOver
	StateVictory
)

type InputMode int

const (
	InputModeNormal InputMode = iota
	InputModeDropping
	InputModeDropMenu
	InputModeInventory
	InputModeCrafting
)

type Game struct {
	Width          int
	Height         int
	Player         *entity.Player
	Dungeon        *dungeon.Dungeon
	Monsters       []*entity.Monster
	Items          []*entity.Item
	Materials      []*entity.Material
	FOV            *fov.FOVMap
	Running        bool
	Seed           int64
	GameState      GameStateType
	Messages       []string
	InputMode      InputMode
	CraftingCursor int
}

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
		Width:          width,
		Height:         height,
		Player:         entity.NewPlayer(playerX, playerY),
		Dungeon:        generatedDungeon,
		Monsters:       make([]*entity.Monster, 0),
		Items:          make([]*entity.Item, 0),
		Materials:      make([]*entity.Material, 0),
		FOV:            fov.NewFOVMap(width, height),
		Running:        true,
		Seed:           seed,
		GameState:      StatePlaying,
		Messages:       make([]string, 0),
		InputMode:      InputModeNormal,
		CraftingCursor: 0,
	}

	newGame.spawnMonsters(rng)
	newGame.spawnBoss(rng)
	newGame.spawnItems(rng)
	newGame.ComputeFOV()

	return newGame
}

func (g *Game) ComputeFOV() {
	px, py := g.Player.Position()
	fov.Compute(g.FOV, g.Dungeon, px, py, playerFOVRadius)
}

func (g *Game) IsVisible(x, y int) bool {
	return g.FOV.IsVisible(x, y)
}

func (g *Game) IsExplored(x, y int) bool {
	return g.FOV.IsExplored(x, y)
}

func (g *Game) spawnMonsters(rng *rand.Rand) {
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

	startRoom := 1
	if len(g.Dungeon.Rooms) <= 1 {
		startRoom = 0
	}

	for i := startRoom; i < len(g.Dungeon.Rooms); i++ {
		room := g.Dungeon.Rooms[i]
		numMonsters := rng.Intn(3) + 1

		for j := 0; j < numMonsters; j++ {
			x := room.X + rng.Intn(room.Width)
			y := room.Y + rng.Intn(room.Height)

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

			mType := monsterTypes[rng.Intn(len(monsterTypes))]
			monster := entity.NewMonster(mType.name, mType.glyph, x, y, mType.hp, mType.attack)
			monster.DropTable = entity.GetRegularMonsterDropTable()
			g.Monsters = append(g.Monsters, monster)
		}
	}
}

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

func (g *Game) GetBoss() *entity.Monster {
	for _, monster := range g.Monsters {
		if monster.IsBoss {
			return monster
		}
	}
	return nil
}

func (g *Game) spawnBoss(rng *rand.Rand) {
	if len(g.Dungeon.Rooms) < 2 {
		return
	}

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

	lastRoom := g.Dungeon.Rooms[len(g.Dungeon.Rooms)-1]
	cx, cy := lastRoom.Center()

	bossType := bossTypes[rng.Intn(len(bossTypes))]
	boss := entity.NewBossMonster(bossType.name, bossType.glyph, cx, cy, bossType.hp, bossType.attack)
	boss.DropTable = entity.GetBossDropTable(bossType.name)
	g.Monsters = append(g.Monsters, boss)
}

func (g *Game) spawnItems(rng *rand.Rand) {
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
			for attempt := 0; attempt < 10; attempt++ {
				x := room.X + rng.Intn(room.Width)
				y := room.Y + rng.Intn(room.Height)

				if !g.isValidItemPosition(x, y) {
					continue
				}

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

func (g *Game) GetItemAt(x, y int) *entity.Item {
	for _, item := range g.Items {
		ix, iy := item.Position()
		if ix == x && iy == y {
			return item
		}
	}
	return nil
}

func (g *Game) GetMaterialAt(x, y int) *entity.Material {
	for _, material := range g.Materials {
		mx, my := material.Position()
		if mx == x && my == y {
			return material
		}
	}
	return nil
}

func (g *Game) RemoveItem(item *entity.Item) {
	for i, it := range g.Items {
		if it == item {
			g.Items = append(g.Items[:i], g.Items[i+1:]...)
			return
		}
	}
}

func (g *Game) RemoveMaterial(material *entity.Material) {
	for i, mat := range g.Materials {
		if mat == material {
			g.Materials = append(g.Materials[:i], g.Materials[i+1:]...)
			return
		}
	}
}

func (g *Game) RemoveDeadMonsters() {
	alive := make([]*entity.Monster, 0, len(g.Monsters))
	for _, monster := range g.Monsters {
		if !monster.Dead {
			alive = append(alive, monster)
		}
	}
	g.Monsters = alive
}

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

func (g *Game) updateWanderAI(monster *entity.Monster, rng *rand.Rand) {
	px, py := g.Player.Position()
	mx, my := monster.Position()

	dx := px - mx
	dy := py - my
	if (dx == 1 || dx == -1) && dy == 0 || (dy == 1 || dy == -1) && dx == 0 {
		g.monsterAttack(monster)
		return
	}

	if rng.Intn(2) == 0 {
		return
	}

	directions := []ui.Direction{ui.DirUp, ui.DirDown, ui.DirLeft, ui.DirRight}
	dir := directions[rng.Intn(len(directions))]

	ddx, ddy := dir.Delta()
	newX := monster.X + ddx
	newY := monster.Y + ddy

	if !g.Dungeon.IsWalkable(newX, newY) {
		return
	}

	if newX == px && newY == py {
		return
	}

	if g.GetMonsterAt(newX, newY) != nil {
		return
	}

	monster.SetPosition(newX, newY)
}

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
	case ui.ActionCraft:
		g.InputMode = InputModeCrafting
		g.CraftingCursor = 0
	}
}

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

	g.tryPickupItem(newX, newY)
	g.tryPickupMaterial(newX, newY)

	g.ComputeFOV()
}

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

func (g *Game) tryPickupMaterial(x, y int) {
	material := g.GetMaterialAt(x, y)
	if material == nil {
		return
	}

	g.Player.MaterialPouch.Add(material.Type, 1)
	g.RemoveMaterial(material)
	g.AddMessage(fmt.Sprintf("Picked up %s.", material.Name()))
}

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

func (g *Game) HandleDropModeInput(r rune) {
	if r == 'x' || r == 'X' || r == 'i' || r == 'I' {
		g.InputMode = InputModeDropMenu
		return
	}

	slot, ok := ui.ParseSlotNumber(r)
	if ok {
		g.dropItemFromSlot(slot)
		g.InputMode = InputModeNormal
		return
	}

	g.InputMode = InputModeNormal
}

func (g *Game) dropItemFromSlot(slot int) {
	item := g.Player.Inventory.GetSlot(slot)
	if item == nil {
		g.AddMessage(fmt.Sprintf("Slot %d is empty.", slot))
		return
	}

	g.Player.Inventory.Remove(slot)

	px, py := g.Player.Position()
	item.X = px
	item.Y = py
	g.Items = append(g.Items, item)

	g.AddMessage(fmt.Sprintf("Dropped %s.", item.Name()))
}

func (g *Game) AddMessage(msg string) {
	g.Messages = append(g.Messages, msg)
	if len(g.Messages) > maxMessages {
		g.Messages = g.Messages[len(g.Messages)-maxMessages:]
	}
}

func (g *Game) CalculateDamage(attack, defense int) int {
	damage := attack - defense
	if damage < 1 {
		damage = 1
	}
	return damage
}

func (g *Game) playerAttack(monster *entity.Monster) {
	damage := g.CalculateDamage(g.Player.EffectiveAttack(), 0)
	monster.TakeDamage(damage)

	if monster.Dead {
		g.spawnMonsterDrops(monster)

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

func (g *Game) spawnMonsterDrops(monster *entity.Monster) {
	if monster.DropTable == nil {
		return
	}

	drops := monster.DropTable.GenerateDrops()
	if len(drops) == 0 {
		return
	}

	mx, my := monster.Position()

	// Find valid positions around the monster for drops
	// Start with the monster's position, then spiral outward
	validPositions := g.findDropPositions(mx, my, len(drops))

	for i, materialType := range drops {
		var posX, posY int
		if i < len(validPositions) {
			posX, posY = validPositions[i].x, validPositions[i].y
		} else {
			// Fallback: if we couldn't find enough positions, stack on monster position
			// This should be rare in practice
			posX, posY = mx, my
		}

		material := entity.NewMaterial(materialType, posX, posY)
		g.Materials = append(g.Materials, material)
	}
}

// position is a simple coordinate pair for drop placement
type position struct {
	x, y int
}

// findDropPositions finds valid positions for material drops, starting from center
// and spiraling outward. Returns up to count positions.
func (g *Game) findDropPositions(centerX, centerY, count int) []position {
	positions := make([]position, 0, count)

	// Check center first (where monster died)
	if g.isValidDropPosition(centerX, centerY) {
		positions = append(positions, position{centerX, centerY})
		if len(positions) >= count {
			return positions
		}
	}

	// Spiral outward checking adjacent tiles
	// Order: immediate neighbors first, then expand radius
	maxRadius := 5 // Don't spread too far

	for radius := 1; radius <= maxRadius && len(positions) < count; radius++ {
		// Check all tiles at this radius (square ring around center)
		for dx := -radius; dx <= radius; dx++ {
			for dy := -radius; dy <= radius; dy++ {
				// Only check tiles on the edge of this radius
				if abs(dx) != radius && abs(dy) != radius {
					continue
				}

				checkX, checkY := centerX+dx, centerY+dy
				if g.isValidDropPosition(checkX, checkY) {
					// Also check this position isn't already in our list
					alreadyUsed := false
					for _, p := range positions {
						if p.x == checkX && p.y == checkY {
							alreadyUsed = true
							break
						}
					}
					if !alreadyUsed {
						positions = append(positions, position{checkX, checkY})
						if len(positions) >= count {
							return positions
						}
					}
				}
			}
		}
	}

	return positions
}

// isValidDropPosition checks if a position is suitable for dropping a material
func (g *Game) isValidDropPosition(x, y int) bool {
	// Must be walkable
	if !g.Dungeon.IsWalkable(x, y) {
		return false
	}

	// Can't be on player
	px, py := g.Player.Position()
	if x == px && y == py {
		return false
	}

	// Can't be on a living monster
	if g.GetMonsterAt(x, y) != nil {
		return false
	}

	// Can't be on an existing item
	if g.GetItemAt(x, y) != nil {
		return false
	}

	// Can't be on an existing material
	if g.GetMaterialAt(x, y) != nil {
		return false
	}

	return true
}

func (g *Game) monsterAttack(monster *entity.Monster) {
	damage := g.CalculateDamage(monster.Attack, g.Player.EffectiveDefense())
	g.Player.TakeDamage(damage)

	g.AddMessage(fmt.Sprintf("The %s hits you for %d damage!", monster.Name, damage))

	if g.Player.Dead {
		g.GameState = StateGameOver
		g.AddMessage("You have been slain!")
	}
}

func (g *Game) CheckPlayerDeath() {
	if g.Player.Dead && g.GameState != StateGameOver {
		g.GameState = StateGameOver
	}
}

func (g *Game) GetVisibleMonsters() []*entity.Monster {
	visible := make([]*entity.Monster, 0)

	for _, monster := range g.Monsters {
		if monster.Dead || monster.IsBoss {
			continue
		}

		mx, my := monster.Position()
		if g.IsVisible(mx, my) {
			visible = append(visible, monster)
		}
	}

	if len(visible) > maxVisibleMonsters {
		px, py := g.Player.Position()
		sort.Slice(visible, func(i, j int) bool {
			mx1, my1 := visible[i].Position()
			mx2, my2 := visible[j].Position()
			dist1 := abs(mx1-px) + abs(my1-py)
			dist2 := abs(mx2-px) + abs(my2-py)
			return dist1 < dist2
		})
		visible = visible[:maxVisibleMonsters]
	}

	return visible
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (g *Game) CraftRecipe(recipeName string) bool {
	recipe := craft.GetRecipeByName(recipeName)
	if recipe == nil {
		g.AddMessage(fmt.Sprintf("Unknown recipe: %s", recipeName))
		return false
	}

	if !recipe.CanCraft(g.Player.MaterialPouch) {
		g.AddMessage("Not enough materials!")
		return false
	}

	equipment, ok := recipe.Craft(g.Player.MaterialPouch)
	if !ok {
		return false
	}

	oldEquip := g.Player.Equip(equipment)
	if oldEquip != nil {
		g.AddMessage(fmt.Sprintf("Crafted %s! (Replaced %s)", equipment.Name, oldEquip.Name))
	} else {
		g.AddMessage(fmt.Sprintf("Crafted and equipped %s!", equipment.Name))
	}

	return true
}

func (g *Game) GetAllRecipes() []*craft.Recipe {
	return craft.GetAllRecipes()
}
