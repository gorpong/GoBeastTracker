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
	"beasttracker/internal/save"
	"beasttracker/internal/score"
	"beasttracker/internal/ui"
)

const (
	baseMonstersPerRoom = 2
	playerFOVRadius     = 8
	maxMessages         = 5
	maxItemsPerRoom     = 2
	maxVisibleMonsters  = 4
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
	InputModeEquipment
)

type Game struct {
	Width           int
	Height          int
	Player          *entity.Player
	Dungeon         *dungeon.Dungeon
	Monsters        []*entity.Monster
	Items           []*entity.Item
	Materials       []*entity.Material
	FOV             *fov.FOVMap
	Running         bool
	Seed            int64
	GameState       GameStateType
	Messages        []string
	InputMode       InputMode
	CraftingCursor  int
	EquipmentCursor int
	Score           int
	HuntNumber      int
}

func NewGame(width, height int, seed int64) *Game {
	return NewGameWithHunt(width, height, seed, 1, nil)
}

func NewGameWithHunt(width, height int, seed int64, huntNumber int, previousPlayer *entity.Player) *Game {
	rng := rand.New(rand.NewSource(seed))

	generatedDungeon := dungeon.GenerateDungeon(width, height, seed)

	var playerX, playerY int
	if len(generatedDungeon.Rooms) > 0 {
		playerX, playerY = generatedDungeon.Rooms[0].Center()
	} else {
		playerX = width / 2
		playerY = height / 2
	}

	var player *entity.Player
	if previousPlayer != nil {
		player = entity.NewPlayer(playerX, playerY)
		player.EquippedWeapon = previousPlayer.EquippedWeapon
		player.EquippedArmor = previousPlayer.EquippedArmor
		player.EquippedCharm = previousPlayer.EquippedCharm
		player.MaterialPouch = previousPlayer.MaterialPouch
		player.EquipmentStash = previousPlayer.EquipmentStash
	} else {
		player = entity.NewPlayer(playerX, playerY)
	}

	newGame := &Game{
		Width:           width,
		Height:          height,
		Player:          player,
		Dungeon:         generatedDungeon,
		Monsters:        make([]*entity.Monster, 0),
		Items:           make([]*entity.Item, 0),
		Materials:       make([]*entity.Material, 0),
		FOV:             fov.NewFOVMap(width, height),
		Running:         true,
		Seed:            seed,
		GameState:       StatePlaying,
		Messages:        make([]string, 0),
		InputMode:       InputModeNormal,
		CraftingCursor:  0,
		EquipmentCursor: 0,
		Score:           0,
		HuntNumber:      huntNumber,
	}

	newGame.spawnMonsters(rng)
	newGame.spawnBoss(rng)
	newGame.spawnItems(rng)
	newGame.ComputeFOV()

	return newGame
}

func NewGameFromCheckpoint(width, height int, seed int64, checkpoint *save.SaveData) *Game {
	rng := rand.New(rand.NewSource(seed))

	generatedDungeon := dungeon.GenerateDungeon(width, height, seed)

	var playerX, playerY int
	if len(generatedDungeon.Rooms) > 0 {
		playerX, playerY = generatedDungeon.Rooms[0].Center()
	} else {
		playerX = width / 2
		playerY = height / 2
	}

	player := entity.NewPlayer(playerX, playerY)

	// Restore equipment
	if checkpoint.EquippedWeapon != nil {
		player.EquippedWeapon = checkpoint.EquippedWeapon
	}
	if checkpoint.EquippedArmor != nil {
		player.EquippedArmor = checkpoint.EquippedArmor
	}
	if checkpoint.EquippedCharm != nil {
		player.EquippedCharm = checkpoint.EquippedCharm
	}

	// Restore stash
	for _, equip := range checkpoint.StashedEquipment {
		player.EquipmentStash.Add(equip)
	}

	// Restore materials
	for matType, count := range checkpoint.Materials {
		player.MaterialPouch.Add(matType, count)
	}

	newGame := &Game{
		Width:           width,
		Height:          height,
		Player:          player,
		Dungeon:         generatedDungeon,
		Monsters:        make([]*entity.Monster, 0),
		Items:           make([]*entity.Item, 0),
		Materials:       make([]*entity.Material, 0),
		FOV:             fov.NewFOVMap(width, height),
		Running:         true,
		Seed:            seed,
		GameState:       StatePlaying,
		Messages:        make([]string, 0),
		InputMode:       InputModeNormal,
		CraftingCursor:  0,
		EquipmentCursor: 0,
		Score:           checkpoint.Score,
		HuntNumber:      checkpoint.HuntNumber,
	}

	newGame.spawnMonsters(rng)
	newGame.spawnBoss(rng)
	newGame.spawnItems(rng)
	newGame.ComputeFOV()

	return newGame
}

// CreateCheckpoint creates save data from current game state
func (g *Game) CreateCheckpoint() *save.SaveData {
	checkpoint := save.NewSaveData("", g.HuntNumber, g.Score)

	checkpoint.EquippedWeapon = g.Player.EquippedWeapon
	checkpoint.EquippedArmor = g.Player.EquippedArmor
	checkpoint.EquippedCharm = g.Player.EquippedCharm

	checkpoint.StashedEquipment = g.Player.EquipmentStash.GetAll()

	for _, matType := range g.Player.MaterialPouch.AllMaterials() {
		checkpoint.Materials[matType] = g.Player.MaterialPouch.Count(matType)
	}

	return checkpoint
}

func (g *Game) getMonstersPerRoom() int {
	return baseMonstersPerRoom + (g.HuntNumber-1)/2
}

func (g *Game) getMonsterHPMultiplier() float64 {
	return 1.0 + float64(g.HuntNumber-1)*0.15
}

func (g *Game) getMonsterAttackBonus() int {
	return (g.HuntNumber - 1) / 2
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
		ai     entity.AIType
	}{
		{"Goblin", 'g', 15, 3, entity.AIWander},
		{"Rat", 'r', 8, 2, entity.AIFleeing},
		{"Spider", 's', 12, 2, entity.AIAggressive},
		{"Bat", 'b', 10, 2, entity.AIWander},
		{"Wolf", 'w', 18, 4, entity.AIChase},
		{"Slime", 'S', 20, 2, entity.AIDefensive},
	}

	startRoom := 1
	if len(g.Dungeon.Rooms) <= 1 {
		startRoom = 0
	}

	monstersPerRoom := g.getMonstersPerRoom()
	hpMultiplier := g.getMonsterHPMultiplier()
	attackBonus := g.getMonsterAttackBonus()

	for i := startRoom; i < len(g.Dungeon.Rooms); i++ {
		room := g.Dungeon.Rooms[i]
		numMonsters := rng.Intn(monstersPerRoom) + 1

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
			scaledHP := int(float64(mType.hp) * hpMultiplier)
			scaledAttack := mType.attack + attackBonus

			monster := entity.NewMonsterWithAI(mType.name, mType.glyph, x, y, scaledHP, scaledAttack, mType.ai)
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
		if monster.IsBoss && !monster.Dead {
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
		name     string
		glyph    rune
		hp       int
		attack   int
		behavior entity.BossBehavior
	}{
		{"Wyvern", 'W', 80, 12, entity.BossTeleport},
		{"Ogre", 'O', 100, 10, entity.BossAggressive},
		{"Troll", 'T', 90, 11, entity.BossAggressive},
		{"Cyclops", 'C', 85, 13, entity.BossSummoner},
		{"Minotaur", 'M', 95, 12, entity.BossNormal},
	}

	lastRoom := g.Dungeon.Rooms[len(g.Dungeon.Rooms)-1]
	cx, cy := lastRoom.Center()

	bossType := bossTypes[rng.Intn(len(bossTypes))]

	hpMultiplier := g.getMonsterHPMultiplier()
	attackBonus := g.getMonsterAttackBonus()
	scaledHP := int(float64(bossType.hp) * hpMultiplier)
	scaledAttack := bossType.attack + attackBonus

	boss := entity.NewBossMonsterWithBehavior(bossType.name, bossType.glyph, cx, cy, scaledHP, scaledAttack, bossType.behavior)
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

		monster.TickCooldowns()

		if monster.IsBoss {
			g.updateBossAI(monster, rng)
			continue
		}

		switch monster.AI {
		case entity.AIWander:
			g.updateWanderAI(monster, rng)
		case entity.AIChase:
			g.updateChaseAI(monster, rng)
		case entity.AIAggressive:
			g.updateAggressiveAI(monster, rng)
		case entity.AIDefensive:
			g.updateDefensiveAI(monster, rng)
		case entity.AIFleeing:
			g.updateFleeingAI(monster, rng)
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

func (g *Game) updateChaseAI(monster *entity.Monster, rng *rand.Rand) {
	px, py := g.Player.Position()
	mx, my := monster.Position()

	dx := px - mx
	dy := py - my
	if (dx == 1 || dx == -1) && dy == 0 || (dy == 1 || dy == -1) && dx == 0 {
		g.monsterAttack(monster)
		return
	}

	if !g.IsVisible(mx, my) {
		g.updateWanderAI(monster, rng)
		return
	}

	g.moveMonsterToward(monster, px, py)
}

func (g *Game) updateBossAI(boss *entity.Monster, rng *rand.Rand) {
	px, py := g.Player.Position()
	bx, by := boss.Position()

	dx := px - bx
	dy := py - by
	if (dx == 1 || dx == -1) && dy == 0 || (dy == 1 || dy == -1) && dx == 0 {
		g.bossAttack(boss)
		return
	}

	switch boss.BossBehavior {
	case entity.BossTeleport:
		if boss.CanTeleport() && rng.Intn(4) == 0 {
			if g.tryBossTeleport(boss, rng) {
				boss.TeleportCooldown = 5
				g.AddMessage(fmt.Sprintf("The %s vanishes and reappears nearby!", boss.Name))
				return
			}
		}
	case entity.BossSummoner:
		if boss.CanSummon() && rng.Intn(6) == 0 {
			if g.tryBossSummon(boss, rng) {
				boss.SummonCooldown = 8
				return
			}
		}
	}

	g.updateChaseAI(boss, rng)
}

func (g *Game) bossAttack(boss *entity.Monster) {
	damage := g.CalculateDamage(boss.GetEffectiveAttack(), g.Player.EffectiveDefense())
	g.Player.TakeDamage(damage)

	if boss.IsEnraged() && boss.BossBehavior == entity.BossAggressive {
		g.AddMessage(fmt.Sprintf("The enraged %s hits you for %d damage!", boss.Name, damage))
	} else {
		g.AddMessage(fmt.Sprintf("The %s hits you for %d damage!", boss.Name, damage))
	}

	if g.Player.Dead {
		g.GameState = StateGameOver
		g.AddMessage("You have been slain!")
	}
}

func (g *Game) tryBossTeleport(boss *entity.Monster, rng *rand.Rand) bool {
	px, py := g.Player.Position()

	var validPositions []position
	for radius := 2; radius <= 4; radius++ {
		for dx := -radius; dx <= radius; dx++ {
			for dy := -radius; dy <= radius; dy++ {
				if abs(dx) != radius && abs(dy) != radius {
					continue
				}

				checkX, checkY := px+dx, py+dy
				if g.isValidTeleportPosition(checkX, checkY) {
					validPositions = append(validPositions, position{checkX, checkY})
				}
			}
		}
	}

	if len(validPositions) == 0 {
		return false
	}

	chosen := validPositions[rng.Intn(len(validPositions))]
	boss.SetPosition(chosen.x, chosen.y)
	return true
}

func (g *Game) isValidTeleportPosition(x, y int) bool {
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

	return true
}

func (g *Game) tryBossSummon(boss *entity.Monster, rng *rand.Rand) bool {
	bx, by := boss.Position()

	minionTypes := []struct {
		name   string
		glyph  rune
		hp     int
		attack int
	}{
		{"Imp", 'i', 8, 2},
		{"Shade", 'h', 6, 3},
	}

	var validPositions []position
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			checkX, checkY := bx+dx, by+dy
			if g.Dungeon.IsWalkable(checkX, checkY) && g.GetMonsterAt(checkX, checkY) == nil {
				px, py := g.Player.Position()
				if checkX != px || checkY != py {
					validPositions = append(validPositions, position{checkX, checkY})
				}
			}
		}
	}

	if len(validPositions) == 0 {
		return false
	}

	chosen := validPositions[rng.Intn(len(validPositions))]
	mType := minionTypes[rng.Intn(len(minionTypes))]

	minion := entity.NewMonsterWithAI(mType.name, mType.glyph, chosen.x, chosen.y, mType.hp, mType.attack, entity.AIAggressive)
	minion.DropTable = entity.GetRegularMonsterDropTable()
	g.Monsters = append(g.Monsters, minion)

	g.AddMessage(fmt.Sprintf("The %s summons a %s!", boss.Name, mType.name))
	return true
}

func (g *Game) updateAggressiveAI(monster *entity.Monster, rng *rand.Rand) {
	px, py := g.Player.Position()
	mx, my := monster.Position()

	dx := px - mx
	dy := py - my
	if (dx == 1 || dx == -1) && dy == 0 || (dy == 1 || dy == -1) && dx == 0 {
		g.monsterAttack(monster)
		return
	}

	g.moveMonsterToward(monster, px, py)
}

func (g *Game) updateDefensiveAI(monster *entity.Monster, rng *rand.Rand) {
	px, py := g.Player.Position()
	mx, my := monster.Position()

	dx := px - mx
	dy := py - my
	adjacent := (dx == 1 || dx == -1) && dy == 0 || (dy == 1 || dy == -1) && dx == 0

	if monster.IsLowHP() {
		if adjacent {
			g.monsterAttack(monster)
		}
		g.moveMonsterAway(monster, px, py)
		return
	}

	if adjacent {
		g.monsterAttack(monster)
		return
	}

	if g.IsVisible(mx, my) {
		g.moveMonsterToward(monster, px, py)
	} else {
		g.updateWanderAI(monster, rng)
	}
}

func (g *Game) updateFleeingAI(monster *entity.Monster, rng *rand.Rand) {
	px, py := g.Player.Position()
	mx, my := monster.Position()

	dx := px - mx
	dy := py - my
	if (dx == 1 || dx == -1) && dy == 0 || (dy == 1 || dy == -1) && dx == 0 {
		if !g.canMonsterFlee(monster, px, py) {
			g.monsterAttack(monster)
		}
		return
	}

	g.moveMonsterAway(monster, px, py)
}

func (g *Game) moveMonsterToward(monster *entity.Monster, targetX, targetY int) {
	mx, my := monster.Position()
	dir := monster.GetChaseDirection(targetX, targetY)
	if dir == ui.DirNone {
		return
	}

	ddx, ddy := dir.Delta()
	newX, newY := mx+ddx, my+ddy

	if !g.Dungeon.IsWalkable(newX, newY) || g.GetMonsterAt(newX, newY) != nil {
		altDir := g.findAlternateChaseDirection(monster, targetX, targetY, dir)
		if altDir == ui.DirNone {
			return
		}
		ddx, ddy = altDir.Delta()
		newX, newY = mx+ddx, my+ddy
	}

	if !g.Dungeon.IsWalkable(newX, newY) || g.GetMonsterAt(newX, newY) != nil {
		return
	}

	px, py := g.Player.Position()
	if newX == px && newY == py {
		return
	}

	monster.SetPosition(newX, newY)
}

func (g *Game) moveMonsterAway(monster *entity.Monster, targetX, targetY int) {
	mx, my := monster.Position()
	dir := monster.GetFleeDirection(targetX, targetY)
	if dir == ui.DirNone {
		return
	}

	ddx, ddy := dir.Delta()
	newX, newY := mx+ddx, my+ddy

	if !g.Dungeon.IsWalkable(newX, newY) || g.GetMonsterAt(newX, newY) != nil {
		alternatives := g.getPerpendicularDirections(dir)
		moved := false
		for _, alt := range alternatives {
			adx, ady := alt.Delta()
			newX, newY = mx+adx, my+ady
			if g.Dungeon.IsWalkable(newX, newY) && g.GetMonsterAt(newX, newY) == nil {
				moved = true
				break
			}
		}
		if !moved {
			return
		}
	}

	px, py := g.Player.Position()
	if newX == px && newY == py {
		return
	}

	monster.SetPosition(newX, newY)
}

func (g *Game) canMonsterFlee(monster *entity.Monster, playerX, playerY int) bool {
	mx, my := monster.Position()
	dir := monster.GetFleeDirection(playerX, playerY)
	if dir == ui.DirNone {
		return false
	}

	ddx, ddy := dir.Delta()
	newX, newY := mx+ddx, my+ddy

	return g.Dungeon.IsWalkable(newX, newY) && g.GetMonsterAt(newX, newY) == nil
}

func (g *Game) getPerpendicularDirections(dir ui.Direction) []ui.Direction {
	switch dir {
	case ui.DirUp, ui.DirDown:
		return []ui.Direction{ui.DirLeft, ui.DirRight}
	case ui.DirLeft, ui.DirRight:
		return []ui.Direction{ui.DirUp, ui.DirDown}
	default:
		return nil
	}
}

func (g *Game) findAlternateChaseDirection(monster *entity.Monster, targetX, targetY int, blockedDir ui.Direction) ui.Direction {
	mx, my := monster.Position()
	dx := targetX - mx
	dy := targetY - my

	var alternatives []ui.Direction

	if blockedDir == ui.DirLeft || blockedDir == ui.DirRight {
		if dy > 0 {
			alternatives = append(alternatives, ui.DirDown)
		}
		if dy < 0 {
			alternatives = append(alternatives, ui.DirUp)
		}
	} else {
		if dx > 0 {
			alternatives = append(alternatives, ui.DirRight)
		}
		if dx < 0 {
			alternatives = append(alternatives, ui.DirLeft)
		}
	}

	for _, alt := range alternatives {
		adx, ady := alt.Delta()
		newX, newY := mx+adx, my+ady
		if g.Dungeon.IsWalkable(newX, newY) && g.GetMonsterAt(newX, newY) == nil {
			return alt
		}
	}

	return ui.DirNone
}

func (g *Game) HandleInput(action ui.Action, dir ui.Direction) {
	// Allow quitting from any state
	if action == ui.ActionQuit {
		g.Running = false
		return
	}

	// Game over blocks all input except quit
	if g.GameState == StateGameOver {
		return
	}

	// Victory state allows continued play
	switch action {
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
	case ui.ActionEquipment:
		g.InputMode = InputModeEquipment
		g.EquipmentCursor = 0
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
	g.tryPickupMaterials(newX, newY)

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

func (g *Game) tryPickupMaterials(x, y int) {
	for {
		material := g.GetMaterialAt(x, y)
		if material == nil {
			return
		}

		g.Player.MaterialPouch.Add(material.Type, 1)
		g.RemoveMaterial(material)
		g.AddMessage(fmt.Sprintf("Picked up %s.", material.Name()))
	}
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
	damage := g.CalculateDamage(g.Player.EffectiveAttack(), monster.Defense)
	monster.TakeDamage(damage)

	if monster.Dead {
		g.spawnMonsterDrops(monster)

		if monster.IsBoss {
			g.Score += score.PointsPerBoss
			g.AddMessage(fmt.Sprintf("You have slain the %s! VICTORY!", monster.Name))
			g.GameState = StateVictory
		} else {
			g.Score += score.PointsPerMonster
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

	validPositions := g.findDropPositions(mx, my, len(drops))

	for i, materialType := range drops {
		var posX, posY int
		if i < len(validPositions) {
			posX, posY = validPositions[i].x, validPositions[i].y
		} else {
			posX, posY = mx, my
		}

		material := entity.NewMaterial(materialType, posX, posY)
		g.Materials = append(g.Materials, material)
	}
}

type position struct {
	x, y int
}

func (g *Game) findDropPositions(centerX, centerY, count int) []position {
	positions := make([]position, 0, count)

	if g.isValidDropPosition(centerX, centerY) {
		positions = append(positions, position{centerX, centerY})
		if len(positions) >= count {
			return positions
		}
	}

	maxRadius := 5

	for radius := 1; radius <= maxRadius && len(positions) < count; radius++ {
		for dx := -radius; dx <= radius; dx++ {
			for dy := -radius; dy <= radius; dy++ {
				if abs(dx) != radius && abs(dy) != radius {
					continue
				}

				checkX, checkY := centerX+dx, centerY+dy
				if g.isValidDropPosition(checkX, checkY) {
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

func (g *Game) isValidDropPosition(x, y int) bool {
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
		g.Player.EquipmentStash.Add(oldEquip)
		g.AddMessage(fmt.Sprintf("Crafted %s! (Old %s moved to stash)", equipment.Name, oldEquip.Name))
	} else {
		g.AddMessage(fmt.Sprintf("Crafted and equipped %s!", equipment.Name))
	}

	return true
}

func (g *Game) GetAllRecipes() []*craft.Recipe {
	return craft.GetAllRecipes()
}

// GetEquipmentList returns all equipment available for a slot (equipped + stashed)
func (g *Game) GetEquipmentList(slot entity.EquipmentSlot) []*entity.Equipment {
	result := make([]*entity.Equipment, 0)

	// Add currently equipped item first
	var equipped *entity.Equipment
	switch slot {
	case entity.SlotWeapon:
		equipped = g.Player.EquippedWeapon
	case entity.SlotArmor:
		equipped = g.Player.EquippedArmor
	case entity.SlotCharm:
		equipped = g.Player.EquippedCharm
	}

	if equipped != nil {
		result = append(result, equipped)
	}

	// Add stashed items
	stashed := g.Player.EquipmentStash.GetBySlot(slot)
	result = append(result, stashed...)

	return result
}

// IsEquipped returns true if the equipment is currently equipped
func (g *Game) IsEquipped(equipment *entity.Equipment) bool {
	return equipment == g.Player.EquippedWeapon ||
		equipment == g.Player.EquippedArmor ||
		equipment == g.Player.EquippedCharm
}
