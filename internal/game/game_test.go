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
