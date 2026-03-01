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

		// Inventory and drop actions
		{"Inventory i", tcell.KeyRune, 'i', ActionInventory},
		{"Drop mode x", tcell.KeyRune, 'x', ActionDropMode},

		// No action for unbound keys
		{"Unbound z", tcell.KeyRune, 'z', ActionNone},
		{"Unbound y", tcell.KeyRune, 'y', ActionNone},
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

// TestParseActionInventory verifies inventory key mapping
func TestParseActionInventory(t *testing.T) {
	result := ParseAction(tcell.KeyRune, 'i')
	if result != ActionInventory {
		t.Errorf("ParseAction('i') = %v, want ActionInventory", result)
	}
}

// TestParseActionDropMode verifies drop mode key mapping
func TestParseActionDropMode(t *testing.T) {
	result := ParseAction(tcell.KeyRune, 'x')
	if result != ActionDropMode {
		t.Errorf("ParseAction('x') = %v, want ActionDropMode", result)
	}
}

// TestParseActionUseItem verifies number keys trigger item use
func TestParseActionUseItem(t *testing.T) {
	tests := []struct {
		name string
		rune rune
		want Action
	}{
		{"Key 1", '1', ActionUseItem},
		{"Key 2", '2', ActionUseItem},
		{"Key 3", '3', ActionUseItem},
		{"Key 4", '4', ActionUseItem},
		{"Key 5", '5', ActionUseItem},
		{"Key 6", '6', ActionUseItem},
		{"Key 7", '7', ActionUseItem},
		{"Key 8", '8', ActionUseItem},
		{"Key 9", '9', ActionUseItem},
		{"Key 0 is not item", '0', ActionNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseAction(tcell.KeyRune, tt.rune)
			if result != tt.want {
				t.Errorf("ParseAction('%c') = %v, want %v", tt.rune, result, tt.want)
			}
		})
	}
}

// TestParseSlotNumber verifies slot number extraction from key
func TestParseSlotNumber(t *testing.T) {
	tests := []struct {
		name     string
		rune     rune
		wantSlot int
		wantOK   bool
	}{
		{"Key 1", '1', 1, true},
		{"Key 2", '2', 2, true},
		{"Key 3", '3', 3, true},
		{"Key 4", '4', 4, true},
		{"Key 5", '5', 5, true},
		{"Key 6", '6', 6, true},
		{"Key 7", '7', 7, true},
		{"Key 8", '8', 8, true},
		{"Key 9", '9', 9, true},
		{"Key 0 invalid", '0', 0, false},
		{"Letter invalid", 'a', 0, false},
		{"Symbol invalid", '!', 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slot, ok := ParseSlotNumber(tt.rune)
			if ok != tt.wantOK {
				t.Errorf("ParseSlotNumber('%c') ok = %v, want %v", tt.rune, ok, tt.wantOK)
			}
			if slot != tt.wantSlot {
				t.Errorf("ParseSlotNumber('%c') slot = %d, want %d", tt.rune, slot, tt.wantSlot)
			}
		})
	}
}

// TestParseActionConfirm verifies enter key mapping
func TestParseActionConfirm(t *testing.T) {
	result := ParseAction(tcell.KeyEnter, 0)
	if result != ActionConfirm {
		t.Errorf("ParseAction(KeyEnter) = %v, want ActionConfirm", result)
	}
}

// TestParseActionCancel verifies escape cancels without quitting in menus
func TestParseActionCancel(t *testing.T) {
	// Note: Escape is both Quit and Cancel - context determines behavior
	// For now, ActionQuit takes precedence; game logic will handle menu context
	result := ParseAction(tcell.KeyEscape, 0)
	if result != ActionQuit {
		t.Errorf("ParseAction(KeyEscape) = %v, want ActionQuit", result)
	}
}

func TestParseActionCraft(t *testing.T) {
	action := ParseAction(tcell.KeyRune, 'c')

	if action != ActionCraft {
		t.Errorf("ParseAction for 'c' = %v, want ActionCraft", action)
	}
}

func TestParseActionCraftUppercase(t *testing.T) {
	action := ParseAction(tcell.KeyRune, 'C')

	if action != ActionCraft {
		t.Errorf("ParseAction for 'C' = %v, want ActionCraft", action)
	}
}

func TestActionCraftString(t *testing.T) {
	if ActionCraft.String() != "Craft" {
		t.Errorf("ActionCraft.String() = %q, want \"Craft\"", ActionCraft.String())
	}
}

func TestParseActionEquipment(t *testing.T) {
	action := ParseAction(tcell.KeyRune, 'e')

	if action != ActionEquipment {
		t.Errorf("ParseAction for 'e' = %v, want ActionEquipment", action)
	}
}

func TestParseActionEquipmentUppercase(t *testing.T) {
	action := ParseAction(tcell.KeyRune, 'E')

	if action != ActionEquipment {
		t.Errorf("ParseAction for 'E' = %v, want ActionEquipment", action)
	}
}

func TestActionEquipmentString(t *testing.T) {
	if ActionEquipment.String() != "Equipment" {
		t.Errorf("ActionEquipment.String() = %q, want \"Equipment\"", ActionEquipment.String())
	}
}

func TestParseActionRestart(t *testing.T) {
	action := ParseAction(tcell.KeyRune, 'r')

	if action != ActionRestart {
		t.Errorf("ParseAction for 'r' = %v, want ActionRestart", action)
	}
}

func TestParseActionRestartUppercase(t *testing.T) {
	action := ParseAction(tcell.KeyRune, 'R')

	if action != ActionRestart {
		t.Errorf("ParseAction for 'R' = %v, want ActionRestart", action)
	}
}

func TestActionRestartString(t *testing.T) {
	if ActionRestart.String() != "Restart" {
		t.Errorf("ActionRestart.String() = %q, want \"Restart\"", ActionRestart.String())
	}
}

func TestParseActionNextHunt(t *testing.T) {
	action := ParseAction(tcell.KeyRune, 'n')

	if action != ActionNextHunt {
		t.Errorf("ParseAction for 'n' = %v, want ActionNextHunt", action)
	}
}

func TestActionNextHuntString(t *testing.T) {
	if ActionNextHunt.String() != "NextHunt" {
		t.Errorf("ActionNextHunt.String() = %q, want \"NextHunt\"", ActionNextHunt.String())
	}
}

func TestParseActionBackspace(t *testing.T) {
	action := ParseAction(tcell.KeyBackspace, 0)

	if action != ActionBackspace {
		t.Errorf("ParseAction for Backspace = %v, want ActionBackspace", action)
	}
}

func TestParseActionBackspace2(t *testing.T) {
	action := ParseAction(tcell.KeyBackspace2, 0)

	if action != ActionBackspace {
		t.Errorf("ParseAction for Backspace2 = %v, want ActionBackspace", action)
	}
}

func TestActionBackspaceString(t *testing.T) {
	if ActionBackspace.String() != "Backspace" {
		t.Errorf("ActionBackspace.String() = %q, want \"Backspace\"", ActionBackspace.String())
	}
}

func TestIsLetter(t *testing.T) {
	tests := []struct {
		r    rune
		want bool
	}{
		{'A', true},
		{'Z', true},
		{'a', true},
		{'z', true},
		{'M', true},
		{'1', false},
		{'!', false},
		{' ', false},
		{'@', false},
	}

	for _, tc := range tests {
		got := IsLetter(tc.r)
		if got != tc.want {
			t.Errorf("IsLetter(%q) = %v, want %v", tc.r, got, tc.want)
		}
	}
}

func TestToUpper(t *testing.T) {
	tests := []struct {
		r    rune
		want rune
	}{
		{'a', 'A'},
		{'z', 'Z'},
		{'m', 'M'},
		{'A', 'A'},
		{'Z', 'Z'},
		{'1', '1'},
	}

	for _, tc := range tests {
		got := ToUpper(tc.r)
		if got != tc.want {
			t.Errorf("ToUpper(%q) = %q, want %q", tc.r, got, tc.want)
		}
	}
}
