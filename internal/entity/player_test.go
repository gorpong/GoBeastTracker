package entity

import (
	"testing"

	"beasttracker/internal/ui"
)

// TestNewPlayer verifies that a new player is created with correct initial values
func TestNewPlayer(t *testing.T) {
	p := NewPlayer(5, 10)

	if p.X != 5 {
		t.Errorf("NewPlayer X = %d, want 5", p.X)
	}
	if p.Y != 10 {
		t.Errorf("NewPlayer Y = %d, want 10", p.Y)
	}
	if p.Glyph != '@' {
		t.Errorf("NewPlayer Glyph = %q, want '@'", p.Glyph)
	}
}

// TestPlayerMove verifies that player moves correctly in each direction
func TestPlayerMove(t *testing.T) {
	tests := []struct {
		name      string
		startX    int
		startY    int
		direction ui.Direction
		wantX     int
		wantY     int
	}{
		{"Move up", 10, 10, ui.DirUp, 10, 9},
		{"Move down", 10, 10, ui.DirDown, 10, 11},
		{"Move left", 10, 10, ui.DirLeft, 9, 10},
		{"Move right", 10, 10, ui.DirRight, 11, 10},
		{"No movement", 10, 10, ui.DirNone, 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPlayer(tt.startX, tt.startY)
			p.Move(tt.direction)

			if p.X != tt.wantX {
				t.Errorf("After Move(%v): X = %d, want %d", tt.direction, p.X, tt.wantX)
			}
			if p.Y != tt.wantY {
				t.Errorf("After Move(%v): Y = %d, want %d", tt.direction, p.Y, tt.wantY)
			}
		})
	}
}

// TestPlayerPosition verifies Position() returns correct coordinates
func TestPlayerPosition(t *testing.T) {
	p := NewPlayer(7, 3)
	x, y := p.Position()

	if x != 7 || y != 3 {
		t.Errorf("Position() = (%d, %d), want (7, 3)", x, y)
	}
}

// TestPlayerSetPosition verifies SetPosition correctly updates coordinates
func TestPlayerSetPosition(t *testing.T) {
	p := NewPlayer(0, 0)
	p.SetPosition(15, 20)

	if p.X != 15 || p.Y != 20 {
		t.Errorf("After SetPosition(15, 20): (%d, %d), want (15, 20)", p.X, p.Y)
	}
}

// TestPlayerMultipleMoves verifies multiple consecutive moves work correctly
func TestPlayerMultipleMoves(t *testing.T) {
	p := NewPlayer(10, 10)

	// Move in a square pattern
	p.Move(ui.DirUp)    // 10, 9
	p.Move(ui.DirRight) // 11, 9
	p.Move(ui.DirDown)  // 11, 10
	p.Move(ui.DirLeft)  // 10, 10

	if p.X != 10 || p.Y != 10 {
		t.Errorf("After square pattern: (%d, %d), want (10, 10)", p.X, p.Y)
	}
}
