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

// TestPlayerStats verifies player has combat stats
func TestPlayerStats(t *testing.T) {
	p := NewPlayer(0, 0)

	if p.HP <= 0 {
		t.Errorf("Player HP = %d, want > 0", p.HP)
	}
	if p.MaxHP <= 0 {
		t.Errorf("Player MaxHP = %d, want > 0", p.MaxHP)
	}
	if p.Attack <= 0 {
		t.Errorf("Player Attack = %d, want > 0", p.Attack)
	}
	if p.Defense < 0 {
		t.Errorf("Player Defense = %d, want >= 0", p.Defense)
	}
}

// TestPlayerTakeDamage verifies damage reduces HP
func TestPlayerTakeDamage(t *testing.T) {
	p := NewPlayer(0, 0)
	initialHP := p.HP

	p.TakeDamage(5)

	if p.HP != initialHP-5 {
		t.Errorf("After taking 5 damage: HP = %d, want %d", p.HP, initialHP-5)
	}
	if p.Dead {
		t.Error("Player should not be dead after minor damage")
	}
}

// TestPlayerDeath verifies player dies when HP reaches 0
func TestPlayerDeath(t *testing.T) {
	p := NewPlayer(0, 0)

	p.TakeDamage(p.HP)

	if p.HP != 0 {
		t.Errorf("After fatal damage: HP = %d, want 0", p.HP)
	}
	if !p.Dead {
		t.Error("Player should be dead at 0 HP")
	}
}

// TestPlayerOverkillDamage verifies HP doesn't go negative
func TestPlayerOverkillDamage(t *testing.T) {
	p := NewPlayer(0, 0)

	p.TakeDamage(p.HP + 100)

	if p.HP < 0 {
		t.Errorf("HP should not be negative: HP = %d", p.HP)
	}
	if p.HP != 0 {
		t.Errorf("HP should be 0 after overkill: HP = %d", p.HP)
	}
}

// TestPlayerIsAlive verifies IsAlive returns correct status
func TestPlayerIsAlive(t *testing.T) {
	p := NewPlayer(0, 0)

	if !p.IsAlive() {
		t.Error("New player should be alive")
	}

	p.TakeDamage(p.HP)

	if p.IsAlive() {
		t.Error("Player should not be alive after fatal damage")
	}
}
