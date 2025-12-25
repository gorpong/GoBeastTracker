package game

import (
	"testing"

	"beasttracker/internal/ui"
)

// TestNewGame verifies that a new game is created with correct initial state
func TestNewGame(t *testing.T) {
	g := NewGame(80, 25)

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
	if g.Running != true {
		t.Error("Game should be running after creation")
	}
}

// TestGamePlayerSpawn verifies player spawns within bounds
func TestGamePlayerSpawn(t *testing.T) {
	g := NewGame(80, 25)

	x, y := g.Player.Position()

	if x < 0 || x >= 80 {
		t.Errorf("Player X = %d, should be in [0, 80)", x)
	}
	if y < 0 || y >= 25 {
		t.Errorf("Player Y = %d, should be in [0, 25)", y)
	}
}

// TestGameHandleInput verifies input handling updates game state correctly
func TestGameHandleInput(t *testing.T) {
	g := NewGame(80, 25)
	startX, startY := g.Player.Position()

	// Move right
	g.HandleInput(ui.ActionMove, ui.DirRight)
	x, y := g.Player.Position()

	if x != startX+1 {
		t.Errorf("After move right: X = %d, want %d", x, startX+1)
	}
	if y != startY {
		t.Errorf("After move right: Y = %d, want %d (unchanged)", y, startY)
	}
}

// TestGameHandleQuit verifies quit action stops the game
func TestGameHandleQuit(t *testing.T) {
	g := NewGame(80, 25)

	if !g.Running {
		t.Error("Game should be running before quit")
	}

	g.HandleInput(ui.ActionQuit, ui.DirNone)

	if g.Running {
		t.Error("Game should not be running after quit")
	}
}

// TestGameBoundaryCheck verifies player cannot move outside bounds
func TestGameBoundaryCheck(t *testing.T) {
	g := NewGame(80, 25)

	// Position player at top-left corner
	g.Player.SetPosition(0, 0)

	// Try to move up (should be blocked)
	g.HandleInput(ui.ActionMove, ui.DirUp)
	x, y := g.Player.Position()
	if y != 0 {
		t.Errorf("Player moved up past boundary: Y = %d, want 0", y)
	}

	// Try to move left (should be blocked)
	g.HandleInput(ui.ActionMove, ui.DirLeft)
	x, y = g.Player.Position()
	if x != 0 {
		t.Errorf("Player moved left past boundary: X = %d, want 0", x)
	}

	// Position player at bottom-right corner
	g.Player.SetPosition(79, 24)

	// Try to move down (should be blocked)
	g.HandleInput(ui.ActionMove, ui.DirDown)
	x, y = g.Player.Position()
	if y != 24 {
		t.Errorf("Player moved down past boundary: Y = %d, want 24", y)
	}

	// Try to move right (should be blocked)
	g.HandleInput(ui.ActionMove, ui.DirRight)
	x, y = g.Player.Position()
	if x != 79 {
		t.Errorf("Player moved right past boundary: X = %d, want 79", x)
	}
}

// TestGameValidMove verifies player can move within bounds
func TestGameValidMove(t *testing.T) {
	g := NewGame(80, 25)

	// Position player in center
	g.Player.SetPosition(40, 12)

	// Move in all directions
	directions := []struct {
		dir    ui.Direction
		wantX  int
		wantY  int
	}{
		{ui.DirRight, 41, 12},
		{ui.DirDown, 41, 13},
		{ui.DirLeft, 40, 13},
		{ui.DirUp, 40, 12},
	}

	for _, d := range directions {
		g.HandleInput(ui.ActionMove, d.dir)
		x, y := g.Player.Position()
		if x != d.wantX || y != d.wantY {
			t.Errorf("After move %v: got (%d, %d), want (%d, %d)",
				d.dir, x, y, d.wantX, d.wantY)
		}
	}
}
