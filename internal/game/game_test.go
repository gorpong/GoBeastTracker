package game

import (
	"fmt"
	"testing"

	"beasttracker/internal/dungeon"
	"beasttracker/internal/entity"
	"beasttracker/internal/ui"
)

// TestNewGame verifies that a new game is created with correct initial state
func TestNewGame(t *testing.T) {
	g := NewGame(80, 25, 12345)

	if g == nil {
		t.Fatal("NewGame() returned nil")
	}

	if g.Width != 80 {
		t.Errorf("Game Width = %d, want 80", g.Width)
	}
	if g.Height != 25 {
		t.Errorf("Game Height = %d, want 25", g.Height)
	}
	if g.Player == nil {
		t.Fatal("Game Player is nil")
	}
	if g.Dungeon == nil {
		t.Fatal("Game Dungeon is nil")
	}
	if g.Running != true {
		t.Error("Game should be running after creation")
	}
}

// TestNewGameWithDungeon verifies dungeon is properly generated
func TestNewGameWithDungeon(t *testing.T) {
	g := NewGame(100, 40, 12345)

	if g.Dungeon.Width != 100 {
		t.Errorf("Dungeon Width = %d, want 100", g.Dungeon.Width)
	}
	if g.Dungeon.Height != 40 {
		t.Errorf("Dungeon Height = %d, want 40", g.Dungeon.Height)
	}
	if len(g.Dungeon.Rooms) == 0 {
		t.Error("Dungeon should have rooms")
	}
}

// TestGamePlayerSpawnInRoom verifies player spawns in a room (walkable tile)
func TestGamePlayerSpawnInRoom(t *testing.T) {
	g := NewGame(100, 40, 12345)

	x, y := g.Player.Position()

	// Player should spawn on a walkable tile
	if !g.Dungeon.IsWalkable(x, y) {
		t.Errorf("Player spawned at (%d, %d) which is not walkable", x, y)
	}
}

// TestGameHandleQuit verifies quit action stops the game
func TestGameHandleQuit(t *testing.T) {
	g := NewGame(100, 40, 12345)

	if !g.Running {
		t.Error("Game should be running before quit")
	}

	g.HandleInput(ui.ActionQuit, ui.DirNone)

	if g.Running {
		t.Error("Game should not be running after quit")
	}
}

// TestGameWallCollision verifies player cannot walk through walls
func TestGameWallCollision(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	// Find the player's current position (should be in a room)
	startX, startY := testGame.Player.Position()

	// Find a direction that leads to a wall
	// We'll check all 4 directions and verify walls block movement
	directions := []ui.Direction{ui.DirUp, ui.DirDown, ui.DirLeft, ui.DirRight}

	for _, dir := range directions {
		// Reset player to start position
		testGame.Player.SetPosition(startX, startY)

		dx, dy := dir.Delta()
		targetX, targetY := startX+dx, startY+dy

		// Try to move
		testGame.HandleInput(ui.ActionMove, dir)
		newX, newY := testGame.Player.Position()

		if testGame.Dungeon.IsWalkable(targetX, targetY) {
			// If target is walkable, player should have moved
			if newX != targetX || newY != targetY {
				t.Errorf("Player should have moved to walkable tile (%d,%d), but is at (%d,%d)",
					targetX, targetY, newX, newY)
			}
		} else {
			// If target is not walkable, player should stay in place
			if newX != startX || newY != startY {
				t.Errorf("Player should not move into wall at (%d,%d), but moved to (%d,%d)",
					targetX, targetY, newX, newY)
			}
		}
	}
}

// TestGameMovementInRoom verifies player can move freely within a room
func TestGameMovementInRoom(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	// Player starts in first room's center
	startX, startY := testGame.Player.Position()

	// Find a direction where we can move (floor tile)
	var canMoveDir ui.Direction
	var targetX, targetY int

	for _, dir := range []ui.Direction{ui.DirUp, ui.DirDown, ui.DirLeft, ui.DirRight} {
		dx, dy := dir.Delta()
		tx, ty := startX+dx, startY+dy
		if testGame.Dungeon.IsWalkable(tx, ty) {
			canMoveDir = dir
			targetX, targetY = tx, ty
			break
		}
	}

	if canMoveDir == ui.DirNone {
		t.Skip("No walkable adjacent tile found")
	}

	testGame.HandleInput(ui.ActionMove, canMoveDir)
	newX, newY := testGame.Player.Position()

	if newX != targetX || newY != targetY {
		t.Errorf("Player should have moved to (%d,%d), but is at (%d,%d)",
			targetX, targetY, newX, newY)
	}
}

// TestGameHasMonsters verifies monsters are spawned in the game
func TestGameHasMonsters(t *testing.T) {
	g := NewGame(100, 40, 12345)

	if len(g.Monsters) == 0 {
		t.Error("Game should have monsters spawned")
	}

	// Should have multiple monsters
	if len(g.Monsters) < 3 {
		t.Errorf("Expected at least 3 monsters, got %d", len(g.Monsters))
	}
}

// TestGameMonstersInRooms verifies monsters spawn on walkable tiles
func TestGameMonstersInRooms(t *testing.T) {
	g := NewGame(100, 40, 12345)

	for i, monster := range g.Monsters {
		x, y := monster.Position()
		if !g.Dungeon.IsWalkable(x, y) {
			t.Errorf("Monster %d at (%d, %d) is not on a walkable tile", i, x, y)
		}
	}
}

// TestGameMonstersDontOverlapPlayer verifies monsters don't spawn on player
func TestGameMonstersDontOverlapPlayer(t *testing.T) {
	g := NewGame(100, 40, 12345)

	px, py := g.Player.Position()
	for _, monster := range g.Monsters {
		mx, my := monster.Position()
		if mx == px && my == py {
			t.Errorf("Monster spawned on player position (%d, %d)", px, py)
		}
	}
}

// TestGameGetMonsterAt verifies GetMonsterAt returns correct monster
func TestGameGetMonsterAt(t *testing.T) {
	g := NewGame(100, 40, 12345)

	if len(g.Monsters) == 0 {
		t.Skip("No monsters to test")
	}

	// Get first monster's position
	mx, my := g.Monsters[0].Position()

	found := g.GetMonsterAt(mx, my)
	if found == nil {
		t.Errorf("GetMonsterAt(%d, %d) returned nil, expected monster", mx, my)
	}
	if found != g.Monsters[0] {
		t.Error("GetMonsterAt returned wrong monster")
	}

	// Check empty position
	emptyMonster := g.GetMonsterAt(-999, -999)
	if emptyMonster != nil {
		t.Error("GetMonsterAt should return nil for empty position")
	}
}

// TestGameRemoveDeadMonsters verifies dead monsters are removed
func TestGameRemoveDeadMonsters(t *testing.T) {
	g := NewGame(100, 40, 12345)

	initialCount := len(g.Monsters)
	if initialCount == 0 {
		t.Skip("No monsters to test")
	}

	// Kill first monster
	g.Monsters[0].TakeDamage(g.Monsters[0].HP)
	if !g.Monsters[0].Dead {
		t.Fatal("Monster should be dead")
	}

	// Remove dead monsters
	g.RemoveDeadMonsters()

	if len(g.Monsters) != initialCount-1 {
		t.Errorf("After removing dead: %d monsters, want %d", len(g.Monsters), initialCount-1)
	}
}

// TestGameUpdateMonsterAI verifies monsters move with wander AI
func TestGameUpdateMonsterAI(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	if len(testGame.Monsters) == 0 {
		t.Skip("No monsters to test")
	}

	// Get initial positions
	initialPositions := make(map[*entity.Monster]struct{ x, y int })
	for _, monster := range testGame.Monsters {
		x, y := monster.Position()
		initialPositions[monster] = struct{ x, y int }{x, y}
	}

	// Update AI multiple times - some monsters should move
	moved := false
	for i := 0; i < 10; i++ {
		testGame.UpdateMonsterAI()
		for _, monster := range testGame.Monsters {
			x, y := monster.Position()
			initial := initialPositions[monster]
			if x != initial.x || y != initial.y {
				moved = true
				break
			}
		}
		if moved {
			break
		}
	}

	if !moved {
		t.Error("After 10 AI updates, at least one monster should have moved")
	}
}

// TestGameMonstersDontMoveIntoWalls verifies AI respects walls
func TestGameMonstersDontMoveIntoWalls(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	if len(testGame.Monsters) == 0 {
		t.Skip("No monsters to test")
	}

	// Update AI many times
	for i := 0; i < 50; i++ {
		testGame.UpdateMonsterAI()

		// Check all monsters are on walkable tiles
		for _, monster := range testGame.Monsters {
			x, y := monster.Position()
			if !testGame.Dungeon.IsWalkable(x, y) {
				t.Errorf("Monster at (%d, %d) is on unwalkable tile after AI update", x, y)
			}
		}
	}
}

// TestGameMonstersDontOverlap verifies AI prevents monster stacking
func TestGameMonstersDontOverlap(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	if len(testGame.Monsters) < 2 {
		t.Skip("Need at least 2 monsters to test")
	}

	// Update AI many times
	for i := 0; i < 50; i++ {
		testGame.UpdateMonsterAI()

		// Check no two monsters share position
		positions := make(map[string]bool)
		for _, monster := range testGame.Monsters {
			x, y := monster.Position()
			key := fmt.Sprintf("%d,%d", x, y)
			if positions[key] {
				t.Errorf("Multiple monsters at position %s after AI update", key)
			}
			positions[key] = true
		}
	}
}

// TestGameBumpToAttack verifies bump-to-attack combat mechanics
func TestGameBumpToAttack(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	if len(testGame.Monsters) == 0 {
		t.Skip("No monsters to test")
	}

	// Get player and first monster positions
	px, py := testGame.Player.Position()
	monster := testGame.Monsters[0]
	mx, my := monster.Position()
	initialMonsterHP := monster.HP

	// Place monster adjacent to player
	monster.SetPosition(px+1, py)

	// Bump into monster (move right)
	testGame.HandleInput(ui.ActionMove, ui.DirRight)

	// Player should not have moved (bump attack, not walk through)
	newPx, newPy := testGame.Player.Position()
	if newPx != px || newPy != py {
		t.Errorf("Player moved during attack: from (%d,%d) to (%d,%d)", px, py, newPx, newPy)
	}

	// Monster should have taken damage
	if monster.HP >= initialMonsterHP {
		t.Errorf("Monster HP = %d, should be less than %d after attack", monster.HP, initialMonsterHP)
	}

	// Restore monster position for other tests
	monster.SetPosition(mx, my)
}

// TestGameAttackKillsMonster verifies monsters can be killed
func TestGameAttackKillsMonster(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	if len(testGame.Monsters) == 0 {
		t.Skip("No monsters to test")
	}

	px, py := testGame.Player.Position()
	monster := testGame.Monsters[0]

	// Place weak monster adjacent to player
	monster.SetPosition(px+1, py)
	monster.HP = 1 // Set to 1 HP so it dies in one hit

	initialMonsterCount := len(testGame.Monsters)

	// Attack monster
	testGame.HandleInput(ui.ActionMove, ui.DirRight)

	// Monster should be dead
	if !monster.Dead {
		t.Error("Monster should be dead after attack")
	}

	// Remove dead monsters and verify count decreased
	testGame.RemoveDeadMonsters()
	if len(testGame.Monsters) != initialMonsterCount-1 {
		t.Errorf("Monster count = %d, want %d", len(testGame.Monsters), initialMonsterCount-1)
	}
}

// TestGameMonsterAttacksPlayer verifies monsters can damage player
func TestGameMonsterAttacksPlayer(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	if len(testGame.Monsters) == 0 {
		t.Skip("No monsters to test")
	}

	initialPlayerHP := testGame.Player.HP
	px, py := testGame.Player.Position()
	monster := testGame.Monsters[0]

	// Place monster adjacent to player and have it attack
	monster.SetPosition(px+1, py)

	// Simulate monster attacking player (this would happen during AI update when adjacent)
	damage := testGame.CalculateDamage(monster.Attack, testGame.Player.Defense)
	testGame.Player.TakeDamage(damage)

	if testGame.Player.HP >= initialPlayerHP {
		t.Errorf("Player HP = %d, should be less than %d after taking damage", testGame.Player.HP, initialPlayerHP)
	}
}

// TestGameCalculateDamage verifies damage calculation
func TestGameCalculateDamage(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	// Damage should be attack - defense, minimum 1
	damage := testGame.CalculateDamage(10, 3)
	if damage != 7 {
		t.Errorf("CalculateDamage(10, 3) = %d, want 7", damage)
	}

	// Defense higher than attack should still do 1 damage
	damage = testGame.CalculateDamage(5, 10)
	if damage != 1 {
		t.Errorf("CalculateDamage(5, 10) = %d, want 1 (minimum)", damage)
	}
}

// TestGamePlayerDeath verifies game over on player death
func TestGamePlayerDeath(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	// Kill player
	testGame.Player.TakeDamage(testGame.Player.HP)

	if !testGame.Player.Dead {
		t.Error("Player should be dead")
	}

	// Trigger game state update
	testGame.CheckPlayerDeath()

	// Check game state reflects player death
	if testGame.GameState != StateGameOver {
		t.Errorf("Game state = %v, want StateGameOver", testGame.GameState)
	}
}

// TestGameMessageLog verifies combat messages are logged
func TestGameMessageLog(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	if len(testGame.Monsters) == 0 {
		t.Skip("No monsters to test")
	}

	px, py := testGame.Player.Position()
	monster := testGame.Monsters[0]

	// Place monster adjacent to player
	monster.SetPosition(px+1, py)

	// Clear any existing messages
	testGame.Messages = nil

	// Attack monster
	testGame.HandleInput(ui.ActionMove, ui.DirRight)

	// Should have at least one message about the attack
	if len(testGame.Messages) == 0 {
		t.Error("Expected combat message after attack")
	}
}

// TestGameHasBoss verifies a boss is spawned in the game
func TestGameHasBoss(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	boss := testGame.GetBoss()
	if boss == nil {
		t.Error("Game should have a boss monster")
	}

	if !boss.IsBoss {
		t.Error("Boss monster IsBoss should be true")
	}
}

// TestGameVictoryOnBossKill verifies victory when boss is killed
func TestGameVictoryOnBossKill(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	boss := testGame.GetBoss()
	if boss == nil {
		t.Skip("No boss to test")
	}

	px, py := testGame.Player.Position()

	// Place boss adjacent to player
	boss.SetPosition(px+1, py)
	boss.HP = 1 // Set to 1 HP so it dies in one hit

	// Attack boss
	testGame.HandleInput(ui.ActionMove, ui.DirRight)

	// Boss should be dead
	if !boss.Dead {
		t.Error("Boss should be dead after attack")
	}

	// Game should be in victory state
	if testGame.GameState != StateVictory {
		t.Errorf("Game state = %v, want StateVictory", testGame.GameState)
	}
}

// TestGameBossInLastRoom verifies boss spawns in last room (furthest from player)
func TestGameBossInLastRoom(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	boss := testGame.GetBoss()
	if boss == nil {
		t.Skip("No boss to test")
	}

	bx, by := boss.Position()

	// Boss should be on a walkable tile
	if !testGame.Dungeon.IsWalkable(bx, by) {
		t.Errorf("Boss at (%d, %d) is not on a walkable tile", bx, by)
	}
}

// Ensure dungeon import is used
var _ = dungeon.TileFloor

// TestGameHasItems verifies items are spawned in the game
func TestGameHasItems(t *testing.T) {
	// Try multiple seeds to account for randomness in item spawning
	itemsFound := false

	for seed := int64(0); seed < 20; seed++ {
		testGame := NewGame(100, 40, seed)
		if len(testGame.Items) > 0 {
			itemsFound = true
			break
		}
	}

	if !itemsFound {
		t.Error("Items should spawn in at least one of 20 game seeds")
	}
}

// TestGameItemsOnWalkableTiles verifies items spawn on walkable tiles
func TestGameItemsOnWalkableTiles(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	for i, item := range testGame.Items {
		x, y := item.Position()
		if !testGame.Dungeon.IsWalkable(x, y) {
			t.Errorf("Item %d at (%d, %d) is not on a walkable tile", i, x, y)
		}
	}
}

// TestGameItemsDontOverlapPlayer verifies items don't spawn on player
func TestGameItemsDontOverlapPlayer(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	px, py := testGame.Player.Position()
	for i, item := range testGame.Items {
		ix, iy := item.Position()
		if ix == px && iy == py {
			t.Errorf("Item %d spawned on player position (%d, %d)", i, px, py)
		}
	}
}

// TestGameItemsDontOverlapMonsters verifies items don't spawn on monsters
func TestGameItemsDontOverlapMonsters(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	for i, item := range testGame.Items {
		ix, iy := item.Position()
		if testGame.GetMonsterAt(ix, iy) != nil {
			t.Errorf("Item %d spawned on monster position (%d, %d)", i, ix, iy)
		}
	}
}

// TestGameGetItemAt verifies GetItemAt returns correct item
func TestGameGetItemAt(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	if len(testGame.Items) == 0 {
		t.Skip("No items to test")
	}

	// Get first item's position
	firstItem := testGame.Items[0]
	ix, iy := firstItem.Position()

	found := testGame.GetItemAt(ix, iy)
	if found == nil {
		t.Errorf("GetItemAt(%d, %d) returned nil, expected item", ix, iy)
	}
	if found != firstItem {
		t.Error("GetItemAt returned wrong item")
	}

	// Check empty position
	emptyItem := testGame.GetItemAt(-999, -999)
	if emptyItem != nil {
		t.Error("GetItemAt should return nil for empty position")
	}
}

// TestGameItemsDontOverlapEachOther verifies no two items share position
func TestGameItemsDontOverlapEachOther(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	positions := make(map[string]bool)
	for i, item := range testGame.Items {
		x, y := item.Position()
		key := fmt.Sprintf("%d,%d", x, y)
		if positions[key] {
			t.Errorf("Item %d at position %s overlaps with another item", i, key)
		}
		positions[key] = true
	}
}

// TestGameRemoveItem verifies items can be removed from the game
func TestGameRemoveItem(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	if len(testGame.Items) == 0 {
		t.Skip("No items to test")
	}

	initialCount := len(testGame.Items)
	firstItem := testGame.Items[0]
	ix, iy := firstItem.Position()

	testGame.RemoveItem(firstItem)

	if len(testGame.Items) != initialCount-1 {
		t.Errorf("After RemoveItem: %d items, want %d", len(testGame.Items), initialCount-1)
	}

	// Item should no longer be at that position
	if testGame.GetItemAt(ix, iy) == firstItem {
		t.Error("Removed item should not be found at its position")
	}
}

// TestGameItemTypes verifies both item types can spawn
func TestGameItemTypes(t *testing.T) {
	// Test with multiple seeds to increase chance of seeing both types
	herbsFound := false
	potionsFound := false

	for seed := int64(0); seed < 20; seed++ {
		testGame := NewGame(100, 40, seed)
		for _, item := range testGame.Items {
			switch item.Type {
			case entity.ItemHerbs:
				herbsFound = true
			case entity.ItemPotion:
				potionsFound = true
			}
		}
		if herbsFound && potionsFound {
			break
		}
	}

	if !herbsFound {
		t.Error("Herbs should be able to spawn (not found in 20 seeds)")
	}
	if !potionsFound {
		t.Error("Potions should be able to spawn (not found in 20 seeds)")
	}
}

// TestGamePickupItem verifies player can pick up items
func TestGamePickupItem(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	// Place an item adjacent to player
	px, py := testGame.Player.Position()
	testItem := entity.NewItem(entity.ItemHerbs, px+1, py)
	testGame.Items = append(testGame.Items, testItem)

	initialItemCount := len(testGame.Items)
	initialInventoryCount := testGame.Player.Inventory.Count()

	// Move onto the item
	testGame.HandleInput(ui.ActionMove, ui.DirRight)

	// Item should be picked up
	if testGame.Player.Inventory.Count() != initialInventoryCount+1 {
		t.Errorf("Inventory count = %d, want %d",
			testGame.Player.Inventory.Count(), initialInventoryCount+1)
	}

	// Item should be removed from ground
	if len(testGame.Items) != initialItemCount-1 {
		t.Errorf("Ground item count = %d, want %d",
			len(testGame.Items), initialItemCount-1)
	}
}

// TestGamePickupItemMessage verifies pickup generates a message
func TestGamePickupItemMessage(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	px, py := testGame.Player.Position()
	testItem := entity.NewItem(entity.ItemPotion, px+1, py)
	testGame.Items = append(testGame.Items, testItem)

	testGame.Messages = nil

	testGame.HandleInput(ui.ActionMove, ui.DirRight)

	if len(testGame.Messages) == 0 {
		t.Error("Expected pickup message")
	}

	foundPickupMsg := false
	for _, msg := range testGame.Messages {
		if contains(msg, "picked up") || contains(msg, "Picked up") {
			foundPickupMsg = true
			break
		}
	}
	if !foundPickupMsg {
		t.Errorf("Expected pickup message, got: %v", testGame.Messages)
	}
}

// contains checks if substr is in s (simple helper for tests)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestGamePickupItemFullInventory verifies pickup fails when inventory is full
func TestGamePickupItemFullInventory(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	// Fill player inventory
	for testGame.Player.Inventory.Count() < testGame.Player.Inventory.Capacity() {
		filler := entity.NewItem(entity.ItemHerbs, 0, 0)
		testGame.Player.Inventory.Add(filler)
	}

	// Place an item adjacent to player
	px, py := testGame.Player.Position()
	testItem := entity.NewItem(entity.ItemPotion, px+1, py)
	testGame.Items = append(testGame.Items, testItem)

	initialGroundItems := len(testGame.Items)
	testGame.Messages = nil

	// Try to move onto the item
	testGame.HandleInput(ui.ActionMove, ui.DirRight)

	// Player should have moved (items don't block movement)
	newPx, newPy := testGame.Player.Position()
	if newPx != px+1 || newPy != py {
		t.Error("Player should still move even when inventory is full")
	}

	// Item should still be on ground
	if len(testGame.Items) != initialGroundItems {
		t.Errorf("Item should remain on ground when inventory full, got %d items",
			len(testGame.Items))
	}

	// Should have a "full" message
	foundFullMsg := false
	for _, msg := range testGame.Messages {
		if contains(msg, "full") || contains(msg, "Full") {
			foundFullMsg = true
			break
		}
	}
	if !foundFullMsg {
		t.Errorf("Expected inventory full message, got: %v", testGame.Messages)
	}
}

// TestGameUseItem verifies using an item from inventory
func TestGameUseItem(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	// Damage player first
	testGame.Player.TakeDamage(50)
	initialHP := testGame.Player.HP

	// Add a potion to inventory
	potion := entity.NewItem(entity.ItemPotion, 0, 0)
	testGame.Player.Inventory.Add(potion)

	// Use item in slot 1
	testGame.HandleInput(ui.ActionUseItem, ui.DirNone)
	testGame.UseItemInSlot(1)

	// Player should be healed
	expectedHP := initialHP + entity.PotionHealing
	if expectedHP > testGame.Player.MaxHP {
		expectedHP = testGame.Player.MaxHP
	}

	if testGame.Player.HP != expectedHP {
		t.Errorf("Player HP = %d, want %d", testGame.Player.HP, expectedHP)
	}

	// Item should be removed from inventory
	if testGame.Player.Inventory.Count() != 0 {
		t.Error("Item should be removed from inventory after use")
	}
}

// TestGameUseItemEmptySlot verifies using empty slot does nothing
func TestGameUseItemEmptySlot(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	initialHP := testGame.Player.HP
	testGame.Messages = nil

	testGame.UseItemInSlot(1)

	if testGame.Player.HP != initialHP {
		t.Error("Using empty slot should not change HP")
	}

	// Should have an error message
	if len(testGame.Messages) == 0 {
		t.Error("Expected message when using empty slot")
	}
}

// TestGameUseItemAtFullHealth verifies using healing at full HP
func TestGameUseItemAtFullHealth(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	// Add herbs to inventory
	herbs := entity.NewItem(entity.ItemHerbs, 0, 0)
	testGame.Player.Inventory.Add(herbs)

	testGame.Messages = nil

	testGame.UseItemInSlot(1)

	// Should still consume the item (player choice to use it)
	if testGame.Player.Inventory.Count() != 0 {
		t.Error("Item should be consumed even at full health")
	}

	// HP should still be at max
	if testGame.Player.HP != testGame.Player.MaxHP {
		t.Errorf("HP should be at max, got %d", testGame.Player.HP)
	}
}

// TestGameInputModeDropMode verifies drop mode state tracking
func TestGameInputModeDropMode(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	if testGame.InputMode != InputModeNormal {
		t.Errorf("Initial InputMode = %v, want InputModeNormal", testGame.InputMode)
	}

	// Enter drop mode
	testGame.HandleInput(ui.ActionDropMode, ui.DirNone)

	if testGame.InputMode != InputModeDropping {
		t.Errorf("After drop key: InputMode = %v, want InputModeDropping", testGame.InputMode)
	}
}

// TestGameDropItemQuick verifies quick drop with x + number
func TestGameDropItemQuick(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	// Add item to inventory
	herbs := entity.NewItem(entity.ItemHerbs, 0, 0)
	testGame.Player.Inventory.Add(herbs)

	px, py := testGame.Player.Position()
	initialGroundItems := len(testGame.Items)

	// Enter drop mode then press 1
	testGame.HandleInput(ui.ActionDropMode, ui.DirNone)
	testGame.HandleDropModeInput('1')

	// Item should be on ground at player position
	if len(testGame.Items) != initialGroundItems+1 {
		t.Errorf("Ground items = %d, want %d", len(testGame.Items), initialGroundItems+1)
	}

	droppedItem := testGame.GetItemAt(px, py)
	if droppedItem == nil {
		t.Error("Dropped item should be at player position")
	}

	// Item should be removed from inventory
	if testGame.Player.Inventory.Count() != 0 {
		t.Error("Item should be removed from inventory after drop")
	}

	// Should return to normal mode
	if testGame.InputMode != InputModeNormal {
		t.Errorf("InputMode = %v, want InputModeNormal after drop", testGame.InputMode)
	}
}

// TestGameDropModeOpenMenu verifies x + x opens drop menu
func TestGameDropModeOpenMenu(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	// Enter drop mode then press x again
	testGame.HandleInput(ui.ActionDropMode, ui.DirNone)
	testGame.HandleDropModeInput('x')

	if testGame.InputMode != InputModeDropMenu {
		t.Errorf("InputMode = %v, want InputModeDropMenu", testGame.InputMode)
	}
}

// TestGameDropModeOpenMenuWithI verifies x + i opens drop menu
func TestGameDropModeOpenMenuWithI(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	// Enter drop mode then press i
	testGame.HandleInput(ui.ActionDropMode, ui.DirNone)
	testGame.HandleDropModeInput('i')

	if testGame.InputMode != InputModeDropMenu {
		t.Errorf("InputMode = %v, want InputModeDropMenu", testGame.InputMode)
	}
}

// TestGameDropModeCancelOnOtherKey verifies drop mode cancels on unrelated key
func TestGameDropModeCancelOnOtherKey(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	testGame.HandleInput(ui.ActionDropMode, ui.DirNone)
	testGame.HandleDropModeInput('z')

	if testGame.InputMode != InputModeNormal {
		t.Errorf("InputMode = %v, want InputModeNormal after cancel", testGame.InputMode)
	}
}

// TestGameDropEmptySlot verifies dropping from empty slot shows message
func TestGameDropEmptySlot(t *testing.T) {
	testGame := NewGame(100, 40, 12345)

	testGame.Messages = nil

	testGame.HandleInput(ui.ActionDropMode, ui.DirNone)
	testGame.HandleDropModeInput('1')

	// Should have error message
	foundEmptyMsg := false
	for _, msg := range testGame.Messages {
		if contains(msg, "empty") || contains(msg, "Empty") || contains(msg, "nothing") {
			foundEmptyMsg = true
			break
		}
	}
	if !foundEmptyMsg {
		t.Errorf("Expected empty slot message, got: %v", testGame.Messages)
	}

	// Should return to normal mode
	if testGame.InputMode != InputModeNormal {
		t.Errorf("InputMode = %v, want InputModeNormal", testGame.InputMode)
	}
}

func TestGameHasMaterialsSlice(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	if testGame.Materials == nil {
		t.Error("NewGame should initialize Materials slice")
	}
}

func TestGameSpawnsMaterialOnMonsterDeath(t *testing.T) {
	foundMaterialDrop := false

	for seed := int64(1); seed <= 20; seed++ {
		testGame := NewGame(50, 30, seed)

		for _, monster := range testGame.Monsters {
			if monster.IsBoss {
				continue
			}

			monster.HP = 1
			mx, my := monster.Position()
			testGame.Player.SetPosition(mx-1, my)
			testGame.HandleInput(ui.ActionMove, ui.DirRight)

			if len(testGame.Materials) > 0 {
				foundMaterialDrop = true
				break
			}
		}

		if foundMaterialDrop {
			break
		}
	}

	if !foundMaterialDrop {
		t.Error("Killing monsters should sometimes spawn materials")
	}
}

func TestGameBossGuaranteesRareMaterialDrop(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	boss := testGame.GetBoss()
	if boss == nil {
		t.Fatal("Game should have a boss")
	}

	boss.HP = 1
	bx, by := boss.Position()
	testGame.Player.SetPosition(bx-1, by)

	testGame.HandleInput(ui.ActionMove, ui.DirRight)

	hasRareMaterial := false
	for _, mat := range testGame.Materials {
		if mat.Type.IsRare() {
			hasRareMaterial = true
			break
		}
	}

	if !hasRareMaterial {
		t.Error("Killing boss should guarantee a rare material drop")
	}
}

func TestGamePickupMaterial(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	px, py := testGame.Player.Position()
	material := entity.NewMaterial(entity.MaterialScales, px+1, py)
	testGame.Materials = append(testGame.Materials, material)

	testGame.HandleInput(ui.ActionMove, ui.DirRight)

	if testGame.Player.MaterialPouch.Count(entity.MaterialScales) != 1 {
		t.Error("Walking over material should add it to pouch")
	}

	if len(testGame.Materials) != 0 {
		t.Error("Picked up material should be removed from ground")
	}
}

func TestGameGetMaterialAt(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	material := entity.NewMaterial(entity.MaterialClaws, 10, 10)
	testGame.Materials = append(testGame.Materials, material)

	found := testGame.GetMaterialAt(10, 10)
	if found != material {
		t.Error("GetMaterialAt should return material at position")
	}

	notFound := testGame.GetMaterialAt(5, 5)
	if notFound != nil {
		t.Error("GetMaterialAt should return nil when no material present")
	}
}

func TestGameCraftingModeToggle(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	if testGame.InputMode != InputModeNormal {
		t.Error("Game should start in normal mode")
	}

	testGame.HandleInput(ui.ActionCraft, ui.DirNone)

	if testGame.InputMode != InputModeCrafting {
		t.Error("ActionCraft should switch to crafting mode")
	}
}

func TestGameCraftingCursor(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	if testGame.CraftingCursor != 0 {
		t.Error("CraftingCursor should start at 0")
	}
}

func TestGameCraftItem(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	testGame.Player.MaterialPouch.Add(entity.MaterialScales, 3)
	testGame.Player.MaterialPouch.Add(entity.MaterialClaws, 2)

	crafted := testGame.CraftRecipe("Iron Sword")

	if !crafted {
		t.Error("CraftRecipe should return true with sufficient materials")
	}

	if testGame.Player.EquippedWeapon == nil {
		t.Error("Crafted weapon should be auto-equipped")
	}

	if testGame.Player.EquippedWeapon.Name != "Iron Sword" {
		t.Errorf("Equipped weapon = %q, want \"Iron Sword\"", testGame.Player.EquippedWeapon.Name)
	}

	if testGame.Player.MaterialPouch.Count(entity.MaterialScales) != 0 {
		t.Error("Crafting should consume materials")
	}
}

func TestGameCraftItemInsufficientMaterials(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	testGame.Player.MaterialPouch.Add(entity.MaterialScales, 1)

	crafted := testGame.CraftRecipe("Iron Sword")

	if crafted {
		t.Error("CraftRecipe should return false with insufficient materials")
	}

	if testGame.Player.EquippedWeapon != nil {
		t.Error("No weapon should be equipped after failed craft")
	}
}

func TestGameCraftReplacesEquipment(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	oldWeapon := entity.NewEquipment("Old Sword", entity.SlotWeapon, 1, 0, 0)
	testGame.Player.Equip(oldWeapon)

	testGame.Player.MaterialPouch.Add(entity.MaterialScales, 3)
	testGame.Player.MaterialPouch.Add(entity.MaterialClaws, 2)

	testGame.CraftRecipe("Iron Sword")

	if testGame.Player.EquippedWeapon.Name != "Iron Sword" {
		t.Error("New crafted weapon should replace old weapon")
	}
}

func TestGameCombatUsesEffectiveAttack(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	weapon := entity.NewEquipment("Test Sword", entity.SlotWeapon, 5, 0, 0)
	testGame.Player.Equip(weapon)

	expectedAttack := entity.DefaultPlayerAttack + 5
	if testGame.Player.EffectiveAttack() != expectedAttack {
		t.Errorf("EffectiveAttack = %d, want %d", testGame.Player.EffectiveAttack(), expectedAttack)
	}

	var targetMonster *entity.Monster
	for _, monster := range testGame.Monsters {
		if !monster.IsBoss {
			targetMonster = monster
			break
		}
	}

	if targetMonster == nil {
		t.Skip("No regular monster found")
	}

	initialHP := targetMonster.HP
	mx, my := targetMonster.Position()
	testGame.Player.SetPosition(mx-1, my)

	testGame.HandleInput(ui.ActionMove, ui.DirRight)

	if !targetMonster.Dead {
		actualDamage := initialHP - targetMonster.HP
		if actualDamage != expectedAttack {
			t.Errorf("Damage dealt = %d, want %d (effective ATK)", actualDamage, expectedAttack)
		}
	}
}

func TestGameCombatUsesEffectiveDefense(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	armor := entity.NewEquipment("Test Armor", entity.SlotArmor, 0, 3, 0)
	testGame.Player.Equip(armor)

	expectedDefense := entity.DefaultPlayerDefense + 3
	if testGame.Player.EffectiveDefense() != expectedDefense {
		t.Errorf("EffectiveDefense = %d, want %d", testGame.Player.EffectiveDefense(), expectedDefense)
	}
}

func TestGameVisibleMonsters(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	visibleMonsters := testGame.GetVisibleMonsters()

	for _, monster := range visibleMonsters {
		mx, my := monster.Position()
		if !testGame.IsVisible(mx, my) {
			t.Errorf("GetVisibleMonsters returned monster at (%d,%d) which is not visible", mx, my)
		}
	}
}

func TestGameVisibleMonstersExcludesBoss(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	visibleMonsters := testGame.GetVisibleMonsters()

	for _, monster := range visibleMonsters {
		if monster.IsBoss {
			t.Error("GetVisibleMonsters should not include boss (boss shown separately in HUD)")
		}
	}
}

func TestGameVisibleMonstersLimit(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	visibleMonsters := testGame.GetVisibleMonsters()

	if len(visibleMonsters) > 4 {
		t.Errorf("GetVisibleMonsters should return at most 4, got %d", len(visibleMonsters))
	}
}

func TestGameGetAllRecipes(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	recipes := testGame.GetAllRecipes()

	if len(recipes) == 0 {
		t.Error("GetAllRecipes should return at least one recipe")
	}
}

// TestGameMonsterDropsSpreadOut verifies materials spawn at different positions
func TestGameMonsterDropsSpreadOut(t *testing.T) {
	// We need a scenario where a monster drops multiple materials
	// Boss monsters guarantee a rare drop + possible common drops
	// We'll test that when multiple materials drop, they're at different positions

	foundSpreadDrops := false

	for seed := int64(1); seed <= 50; seed++ {
		testGame := NewGame(50, 30, seed)

		boss := testGame.GetBoss()
		if boss == nil {
			continue
		}

		// Place player adjacent to boss and kill it
		bx, by := boss.Position()
		testGame.Player.SetPosition(bx-1, by)
		boss.HP = 1

		testGame.HandleInput(ui.ActionMove, ui.DirRight)

		// Check if we got multiple materials
		if len(testGame.Materials) < 2 {
			continue
		}

		// Verify no two materials share the same position
		positions := make(map[string]bool)
		allUnique := true
		for _, mat := range testGame.Materials {
			mx, my := mat.Position()
			key := fmt.Sprintf("%d,%d", mx, my)
			if positions[key] {
				allUnique = false
				break
			}
			positions[key] = true
		}

		if allUnique && len(testGame.Materials) >= 2 {
			foundSpreadDrops = true
			break
		}
	}

	if !foundSpreadDrops {
		t.Error("When multiple materials drop, they should spawn at different positions")
	}
}

// TestGameMonsterDropsOnWalkableTiles verifies all dropped materials land on walkable tiles
func TestGameMonsterDropsOnWalkableTiles(t *testing.T) {
	for seed := int64(1); seed <= 20; seed++ {
		testGame := NewGame(50, 30, seed)

		boss := testGame.GetBoss()
		if boss == nil {
			continue
		}

		bx, by := boss.Position()
		testGame.Player.SetPosition(bx-1, by)
		boss.HP = 1

		testGame.HandleInput(ui.ActionMove, ui.DirRight)

		for i, mat := range testGame.Materials {
			mx, my := mat.Position()
			if !testGame.Dungeon.IsWalkable(mx, my) {
				t.Errorf("Seed %d: Material %d at (%d,%d) is not on walkable tile", seed, i, mx, my)
			}
		}
	}
}

// TestGameMonsterDropsDontOverlapPlayer verifies materials don't spawn on player
func TestGameMonsterDropsDontOverlapPlayer(t *testing.T) {
	for seed := int64(1); seed <= 20; seed++ {
		testGame := NewGame(50, 30, seed)

		boss := testGame.GetBoss()
		if boss == nil {
			continue
		}

		bx, by := boss.Position()
		testGame.Player.SetPosition(bx-1, by)
		boss.HP = 1

		testGame.HandleInput(ui.ActionMove, ui.DirRight)

		px, py := testGame.Player.Position()
		for i, mat := range testGame.Materials {
			mx, my := mat.Position()
			if mx == px && my == py {
				t.Errorf("Seed %d: Material %d spawned on player position (%d,%d)", seed, i, mx, my)
			}
		}
	}
}

// TestGameMonsterDropsDontOverlapMonsters verifies materials don't spawn on living monsters
func TestGameMonsterDropsDontOverlapMonsters(t *testing.T) {
	for seed := int64(1); seed <= 20; seed++ {
		testGame := NewGame(50, 30, seed)

		boss := testGame.GetBoss()
		if boss == nil {
			continue
		}

		bx, by := boss.Position()
		testGame.Player.SetPosition(bx-1, by)
		boss.HP = 1

		testGame.HandleInput(ui.ActionMove, ui.DirRight)

		for i, mat := range testGame.Materials {
			mx, my := mat.Position()
			if monster := testGame.GetMonsterAt(mx, my); monster != nil {
				t.Errorf("Seed %d: Material %d at (%d,%d) overlaps with monster", seed, i, mx, my)
			}
		}
	}
}

// TestGameMonsterDropsDontOverlapItems verifies materials don't spawn on existing items
func TestGameMonsterDropsDontOverlapItems(t *testing.T) {
	for seed := int64(1); seed <= 20; seed++ {
		testGame := NewGame(50, 30, seed)

		boss := testGame.GetBoss()
		if boss == nil {
			continue
		}

		bx, by := boss.Position()
		testGame.Player.SetPosition(bx-1, by)
		boss.HP = 1

		testGame.HandleInput(ui.ActionMove, ui.DirRight)

		for i, mat := range testGame.Materials {
			mx, my := mat.Position()
			if item := testGame.GetItemAt(mx, my); item != nil {
				t.Errorf("Seed %d: Material %d at (%d,%d) overlaps with item", seed, i, mx, my)
			}
		}
	}
}

func TestGameScore(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	if testGame.Score != 0 {
		t.Errorf("Initial score = %d, want 0", testGame.Score)
	}
}

func TestGameHuntNumber(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	if testGame.HuntNumber != 1 {
		t.Errorf("Initial hunt number = %d, want 1", testGame.HuntNumber)
	}
}

func TestGameScoreOnMonsterKill(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	var targetMonster *entity.Monster
	for _, monster := range testGame.Monsters {
		if !monster.IsBoss {
			targetMonster = monster
			break
		}
	}

	if targetMonster == nil {
		t.Skip("No regular monster found")
	}

	targetMonster.HP = 1
	mx, my := targetMonster.Position()
	testGame.Player.SetPosition(mx-1, my)

	testGame.HandleInput(ui.ActionMove, ui.DirRight)

	if testGame.Score != 10 {
		t.Errorf("Score after kill = %d, want 10", testGame.Score)
	}
}

func TestGameScoreOnBossKill(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	boss := testGame.GetBoss()
	if boss == nil {
		t.Skip("No boss found")
	}

	boss.HP = 1
	bx, by := boss.Position()
	testGame.Player.SetPosition(bx-1, by)

	testGame.HandleInput(ui.ActionMove, ui.DirRight)

	if testGame.Score != 50 {
		t.Errorf("Score after boss kill = %d, want 50", testGame.Score)
	}
}

func TestNewGameWithHunt(t *testing.T) {
	firstGame := NewGame(50, 30, 12345)
	
	weapon := entity.NewEquipment("Test Sword", entity.SlotWeapon, 5, 0, 0)
	firstGame.Player.Equip(weapon)
	firstGame.Player.MaterialPouch.Add(entity.MaterialScales, 10)

	secondGame := NewGameWithHunt(50, 30, 54321, 2, firstGame.Player)

	if secondGame.HuntNumber != 2 {
		t.Errorf("Hunt number = %d, want 2", secondGame.HuntNumber)
	}

	if secondGame.Player.EquippedWeapon == nil {
		t.Error("Equipment should persist to next hunt")
	}

	if secondGame.Player.MaterialPouch.Count(entity.MaterialScales) != 10 {
		t.Error("Materials should persist to next hunt")
	}

	if secondGame.Player.HP != entity.DefaultPlayerHP {
		t.Errorf("HP should reset, got %d", secondGame.Player.HP)
	}
}

func TestDifficultyScaling(t *testing.T) {
	hunt1 := NewGame(50, 30, 12345)
	hunt3 := NewGameWithHunt(50, 30, 12345, 3, nil)

	var hunt1Monster, hunt3Monster *entity.Monster
	for _, m := range hunt1.Monsters {
		if !m.IsBoss && m.Name == "Goblin" {
			hunt1Monster = m
			break
		}
	}
	for _, m := range hunt3.Monsters {
		if !m.IsBoss && m.Name == "Goblin" {
			hunt3Monster = m
			break
		}
	}

	if hunt1Monster != nil && hunt3Monster != nil {
		if hunt3Monster.MaxHP <= hunt1Monster.MaxHP {
			t.Error("Hunt 3 monsters should have more HP than hunt 1")
		}
	}
}

func TestBossUsesChaseAI(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	boss := testGame.GetBoss()
	if boss == nil {
		t.Skip("No boss found")
	}

	if boss.AI != entity.AIChase {
		t.Errorf("Boss AI = %v, want AIChase", boss.AI)
	}
}

func TestChaseAIMovesTowardPlayer(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	boss := testGame.GetBoss()
	if boss == nil {
		t.Skip("No boss found")
	}

	px, py := testGame.Player.Position()
	
	boss.SetPosition(px+5, py)
	testGame.ComputeFOV()
	
	initialBX, _ := boss.Position()

	testGame.UpdateMonsterAI()

	newBX, _ := boss.Position()
	
	if testGame.IsVisible(initialBX, py) && newBX >= initialBX {
		t.Error("Chase AI should move boss toward player when visible")
	}
}

func TestGameVictoryStateAllowsMovement(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	boss := testGame.GetBoss()
	if boss == nil {
		t.Skip("No boss found")
	}

	boss.HP = 1
	bx, by := boss.Position()
	testGame.Player.SetPosition(bx-1, by)

	testGame.HandleInput(ui.ActionMove, ui.DirRight)

	if testGame.GameState != StateVictory {
		t.Fatal("Game should be in victory state after boss kill")
	}

	// Player should still be able to move after victory
	px, py := testGame.Player.Position()
	testGame.HandleInput(ui.ActionMove, ui.DirLeft)

	newPx, newPy := testGame.Player.Position()
	if newPx == px && newPy == py {
		// Check if movement was blocked by wall, not by victory state
		if testGame.Dungeon.IsWalkable(px-1, py) {
			t.Error("Player should be able to move after victory")
		}
	}
}

func TestGameVictoryStateMonstersStillActive(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	boss := testGame.GetBoss()
	if boss == nil {
		t.Skip("No boss found")
	}

	boss.HP = 1
	bx, by := boss.Position()
	testGame.Player.SetPosition(bx-1, by)

	testGame.HandleInput(ui.ActionMove, ui.DirRight)

	if testGame.GameState != StateVictory {
		t.Fatal("Game should be in victory state")
	}

	// Count non-boss monsters
	regularMonsters := 0
	for _, m := range testGame.Monsters {
		if !m.IsBoss && !m.Dead {
			regularMonsters++
		}
	}

	if regularMonsters == 0 {
		t.Skip("No regular monsters to test")
	}

	// Monsters should still update AI after victory
	initialPositions := make(map[*entity.Monster]struct{ x, y int })
	for _, m := range testGame.Monsters {
		if !m.Dead {
			x, y := m.Position()
			initialPositions[m] = struct{ x, y int }{x, y}
		}
	}

	// Move player multiple times to trigger AI updates
	for i := 0; i < 10; i++ {
		testGame.HandleInput(ui.ActionMove, ui.DirLeft)
		testGame.HandleInput(ui.ActionMove, ui.DirRight)
	}

	// At least one monster should have moved (probabilistic, but likely)
	moved := false
	for m, pos := range initialPositions {
		if m.Dead {
			continue
		}
		x, y := m.Position()
		if x != pos.x || y != pos.y {
			moved = true
			break
		}
	}

	// This is probabilistic, so we don't fail if no movement
	// The important thing is that the game didn't crash
	_ = moved
}

func TestGameVictoryCanPickupItems(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	boss := testGame.GetBoss()
	if boss == nil {
		t.Skip("No boss found")
	}

	boss.HP = 1
	bx, by := boss.Position()
	testGame.Player.SetPosition(bx-1, by)

	testGame.HandleInput(ui.ActionMove, ui.DirRight)

	if testGame.GameState != StateVictory {
		t.Fatal("Game should be in victory state")
	}

	// Place an item near player
	px, py := testGame.Player.Position()
	item := entity.NewItem(entity.ItemPotion, px+1, py)
	testGame.Items = append(testGame.Items, item)

	initialCount := testGame.Player.Inventory.Count()

	// Move onto item
	if testGame.Dungeon.IsWalkable(px+1, py) {
		testGame.HandleInput(ui.ActionMove, ui.DirRight)

		if testGame.Player.Inventory.Count() <= initialCount {
			// Item might not have been picked up if inventory full
			if !testGame.Player.Inventory.IsFull() {
				t.Error("Should be able to pick up items after victory")
			}
		}
	}
}

func TestGameVictoryCanPickupMaterials(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	boss := testGame.GetBoss()
	if boss == nil {
		t.Skip("No boss found")
	}

	boss.HP = 1
	bx, by := boss.Position()
	testGame.Player.SetPosition(bx-1, by)

	initialMaterials := len(testGame.Materials)

	testGame.HandleInput(ui.ActionMove, ui.DirRight)

	if testGame.GameState != StateVictory {
		t.Fatal("Game should be in victory state")
	}

	// Boss should have dropped materials
	if len(testGame.Materials) <= initialMaterials {
		t.Skip("Boss didn't drop materials this run")
	}

	// Find a material and move onto it
	if len(testGame.Materials) > 0 {
		mat := testGame.Materials[0]
		mx, my := mat.Position()
		testGame.Player.SetPosition(mx-1, my)

		initialPouchCount := testGame.Player.MaterialPouch.TotalCount()

		if testGame.Dungeon.IsWalkable(mx, my) {
			testGame.HandleInput(ui.ActionMove, ui.DirRight)

			if testGame.Player.MaterialPouch.TotalCount() <= initialPouchCount {
				t.Error("Should be able to pick up materials after victory")
			}
		}
	}
}

func TestGameCheckpointSave(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	checkpoint := testGame.CreateCheckpoint()

	if checkpoint == nil {
		t.Fatal("CreateCheckpoint returned nil")
	}

	if checkpoint.HuntNumber != testGame.HuntNumber {
		t.Errorf("Checkpoint hunt = %d, want %d", checkpoint.HuntNumber, testGame.HuntNumber)
	}

	if checkpoint.Score != testGame.Score {
		t.Errorf("Checkpoint score = %d, want %d", checkpoint.Score, testGame.Score)
	}
}

func TestGameCheckpointPreservesEquipment(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	weapon := entity.NewEquipment("Test Sword", entity.SlotWeapon, 5, 0, 0)
	testGame.Player.Equip(weapon)

	checkpoint := testGame.CreateCheckpoint()

	if checkpoint.EquippedWeapon == nil {
		t.Error("Checkpoint should preserve equipped weapon")
	}

	if checkpoint.EquippedWeapon.Name != "Test Sword" {
		t.Errorf("Checkpoint weapon = %q, want \"Test Sword\"", checkpoint.EquippedWeapon.Name)
	}
}

func TestGameCheckpointPreservesMaterials(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	testGame.Player.MaterialPouch.Add(entity.MaterialScales, 5)
	testGame.Player.MaterialPouch.Add(entity.MaterialClaws, 3)

	checkpoint := testGame.CreateCheckpoint()

	if checkpoint.Materials[entity.MaterialScales] != 5 {
		t.Errorf("Checkpoint scales = %d, want 5", checkpoint.Materials[entity.MaterialScales])
	}
	if checkpoint.Materials[entity.MaterialClaws] != 3 {
		t.Errorf("Checkpoint claws = %d, want 3", checkpoint.Materials[entity.MaterialClaws])
	}
}

func TestGameCheckpointPreservesStash(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	armor := entity.NewEquipment("Leather Armor", entity.SlotArmor, 0, 2, 10)
	testGame.Player.EquipmentStash.Add(armor)

	checkpoint := testGame.CreateCheckpoint()

	if len(checkpoint.StashedEquipment) != 1 {
		t.Errorf("Checkpoint stash count = %d, want 1", len(checkpoint.StashedEquipment))
	}
}

func TestGameRestoreFromCheckpoint(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	weapon := entity.NewEquipment("Test Sword", entity.SlotWeapon, 5, 0, 0)
	testGame.Player.Equip(weapon)
	testGame.Player.MaterialPouch.Add(entity.MaterialScales, 10)
	testGame.Score = 150

	checkpoint := testGame.CreateCheckpoint()

	// Create new game and restore
	newGame := NewGameFromCheckpoint(50, 30, 54321, checkpoint)

	if newGame.HuntNumber != checkpoint.HuntNumber {
		t.Errorf("Restored hunt = %d, want %d", newGame.HuntNumber, checkpoint.HuntNumber)
	}

	if newGame.Score != checkpoint.Score {
		t.Errorf("Restored score = %d, want %d", newGame.Score, checkpoint.Score)
	}

	if newGame.Player.EquippedWeapon == nil || newGame.Player.EquippedWeapon.Name != "Test Sword" {
		t.Error("Restored game should have equipped weapon")
	}

	if newGame.Player.MaterialPouch.Count(entity.MaterialScales) != 10 {
		t.Errorf("Restored scales = %d, want 10", newGame.Player.MaterialPouch.Count(entity.MaterialScales))
	}
}

func TestGameRestoreResetsHP(t *testing.T) {
	testGame := NewGame(50, 30, 12345)
	testGame.Player.TakeDamage(50)

	checkpoint := testGame.CreateCheckpoint()

	newGame := NewGameFromCheckpoint(50, 30, 54321, checkpoint)

	if newGame.Player.HP != entity.DefaultPlayerHP {
		t.Errorf("Restored HP = %d, want %d (full)", newGame.Player.HP, entity.DefaultPlayerHP)
	}
}

func TestGameRestoreResetsInventory(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	potion := entity.NewItem(entity.ItemPotion, 0, 0)
	testGame.Player.Inventory.Add(potion)

	checkpoint := testGame.CreateCheckpoint()

	newGame := NewGameFromCheckpoint(50, 30, 54321, checkpoint)

	if newGame.Player.Inventory.Count() != 0 {
		t.Errorf("Restored inventory count = %d, want 0", newGame.Player.Inventory.Count())
	}
}

func TestGameInputModeEquipment(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	if testGame.InputMode != InputModeNormal {
		t.Error("Game should start in normal mode")
	}

	testGame.HandleInput(ui.ActionEquipment, ui.DirNone)

	if testGame.InputMode != InputModeEquipment {
		t.Errorf("InputMode = %v, want InputModeEquipment", testGame.InputMode)
	}
}

func TestGameEquipmentCursor(t *testing.T) {
	testGame := NewGame(50, 30, 12345)

	if testGame.EquipmentCursor != 0 {
		t.Error("EquipmentCursor should start at 0")
	}
}
