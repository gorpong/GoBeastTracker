package ui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

// TestParseDirection verifies that keyboard events map to correct directions
func TestParseDirection(t *testing.T) {
	tests := []struct {
		name     string
		key      tcell.Key
		rune     rune
		expected Direction
	}{
		// Arrow keys
		{"Arrow Up", tcell.KeyUp, 0, DirUp},
		{"Arrow Down", tcell.KeyDown, 0, DirDown},
		{"Arrow Left", tcell.KeyLeft, 0, DirLeft},
		{"Arrow Right", tcell.KeyRight, 0, DirRight},

		// Vi keys (hjkl)
		{"Vi h (left)", tcell.KeyRune, 'h', DirLeft},
		{"Vi j (down)", tcell.KeyRune, 'j', DirDown},
		{"Vi k (up)", tcell.KeyRune, 'k', DirUp},
		{"Vi l (right)", tcell.KeyRune, 'l', DirRight},

		// WASD keys
		{"WASD w (up)", tcell.KeyRune, 'w', DirUp},
		{"WASD a (left)", tcell.KeyRune, 'a', DirLeft},
		{"WASD s (down)", tcell.KeyRune, 's', DirDown},
		{"WASD d (right)", tcell.KeyRune, 'd', DirRight},

		// Non-movement keys should return DirNone
		{"Non-movement q", tcell.KeyRune, 'q', DirNone},
		{"Non-movement x", tcell.KeyRune, 'x', DirNone},
		{"Non-movement Enter", tcell.KeyEnter, 0, DirNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseDirection(tt.key, tt.rune)
			if result != tt.expected {
				t.Errorf("ParseDirection(%v, %q) = %v, want %v", tt.key, tt.rune, result, tt.expected)
			}
		})
	}
}

// TestParseAction verifies that keyboard events map to correct actions
func TestParseAction(t *testing.T) {
	tests := []struct {
		name     string
		key      tcell.Key
		rune     rune
		expected Action
	}{
		// Quit actions
		{"Quit q", tcell.KeyRune, 'q', ActionQuit},
		{"Quit Q", tcell.KeyRune, 'Q', ActionQuit},
		{"Quit ESC", tcell.KeyEscape, 0, ActionQuit},

		// Movement actions
		{"Move Up", tcell.KeyUp, 0, ActionMove},
		{"Move h", tcell.KeyRune, 'h', ActionMove},
		{"Move w", tcell.KeyRune, 'w', ActionMove},

		// No action for unbound keys
		{"Unbound x", tcell.KeyRune, 'x', ActionNone},
		{"Unbound z", tcell.KeyRune, 'z', ActionNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAction(tt.key, tt.rune)
			if result != tt.expected {
				t.Errorf("ParseAction(%v, %q) = %v, want %v", tt.key, tt.rune, result, tt.expected)
			}
		})
	}
}

// TestDirectionDelta verifies that directions produce correct coordinate deltas
func TestDirectionDelta(t *testing.T) {
	tests := []struct {
		dir    Direction
		wantDX int
		wantDY int
	}{
		{DirNone, 0, 0},
		{DirUp, 0, -1},
		{DirDown, 0, 1},
		{DirLeft, -1, 0},
		{DirRight, 1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.dir.String(), func(t *testing.T) {
			dx, dy := tt.dir.Delta()
			if dx != tt.wantDX || dy != tt.wantDY {
				t.Errorf("Direction(%v).Delta() = (%d, %d), want (%d, %d)",
					tt.dir, dx, dy, tt.wantDX, tt.wantDY)
			}
		})
	}
}

// TestDirectionString verifies string representation of directions
func TestDirectionString(t *testing.T) {
	tests := []struct {
		dir  Direction
		want string
	}{
		{DirNone, "None"},
		{DirUp, "Up"},
		{DirDown, "Down"},
		{DirLeft, "Left"},
		{DirRight, "Right"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.dir.String(); got != tt.want {
				t.Errorf("Direction.String() = %q, want %q", got, tt.want)
			}
		})
	}
}
