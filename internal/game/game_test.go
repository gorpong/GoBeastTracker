package game

import (
	"testing"

	"beasttracker/internal/dungeon"
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
	g := NewGame(100, 40, 12345)

	// Find the player's current position (should be in a room)
	startX, startY := g.Player.Position()

	// Find a direction that leads to a wall
	// We'll check all 4 directions and verify walls block movement
	directions := []ui.Direction{ui.DirUp, ui.DirDown, ui.DirLeft, ui.DirRight}

	for _, dir := range directions {
		// Reset player to start position
		g.Player.SetPosition(startX, startY)

		dx, dy := dir.Delta()
		targetX, targetY := startX+dx, startY+dy

		// Try to move
		g.HandleInput(ui.ActionMove, dir)
		newX, newY := g.Player.Position()

		if g.Dungeon.IsWalkable(targetX, targetY) {
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
	g := NewGame(100, 40, 12345)

	// Player starts in first room's center
	startX, startY := g.Player.Position()

	// Find a direction where we can move (floor tile)
	var canMoveDir ui.Direction
	var targetX, targetY int

	for _, dir := range []ui.Direction{ui.DirUp, ui.DirDown, ui.DirLeft, ui.DirRight} {
		dx, dy := dir.Delta()
		tx, ty := startX+dx, startY+dy
		if g.Dungeon.IsWalkable(tx, ty) {
			canMoveDir = dir
			targetX, targetY = tx, ty
			break
		}
	}

	if canMoveDir == ui.DirNone {
		t.Skip("No walkable adjacent tile found")
	}

	g.HandleInput(ui.ActionMove, canMoveDir)
	newX, newY := g.Player.Position()

	if newX != targetX || newY != targetY {
		t.Errorf("Player should have moved to (%d,%d), but is at (%d,%d)",
			targetX, targetY, newX, newY)
	}
}

// Ensure dungeon import is used
var _ = dungeon.TileFloor
