package ui

import (
	"github.com/gdamore/tcell/v2"
)

// Direction represents a movement direction
type Direction int

const (
	DirNone Direction = iota
	DirUp
	DirDown
	DirLeft
	DirRight
)

// String returns the string representation of a Direction
func (d Direction) String() string {
	switch d {
	case DirUp:
		return "Up"
	case DirDown:
		return "Down"
	case DirLeft:
		return "Left"
	case DirRight:
		return "Right"
	default:
		return "None"
	}
}

// Delta returns the (dx, dy) coordinate change for this direction
func (d Direction) Delta() (int, int) {
	switch d {
	case DirUp:
		return 0, -1
	case DirDown:
		return 0, 1
	case DirLeft:
		return -1, 0
	case DirRight:
		return 1, 0
	default:
		return 0, 0
	}
}

// Action represents a game action
type Action int

const (
	ActionNone Action = iota
	ActionMove
	ActionQuit
	ActionInventory // Toggle inventory display
	ActionDropMode  // Enter drop mode
	ActionUseItem   // Use item (slot determined by ParseSlotNumber)
	ActionConfirm   // Confirm selection (Enter key)
)

// String returns the string representation of an Action
func (a Action) String() string {
	switch a {
	case ActionNone:
		return "None"
	case ActionMove:
		return "Move"
	case ActionQuit:
		return "Quit"
	case ActionInventory:
		return "Inventory"
	case ActionDropMode:
		return "DropMode"
	case ActionUseItem:
		return "UseItem"
	case ActionConfirm:
		return "Confirm"
	default:
		return "Unknown"
	}
}

// ParseDirection converts a key event into a Direction.
// Returns DirNone if the key is not a movement key.
func ParseDirection(key tcell.Key, r rune) Direction {
	// Arrow keys
	switch key {
	case tcell.KeyUp:
		return DirUp
	case tcell.KeyDown:
		return DirDown
	case tcell.KeyLeft:
		return DirLeft
	case tcell.KeyRight:
		return DirRight
	}

	// Vi keys (hjkl) and WASD
	if key == tcell.KeyRune {
		switch r {
		case 'h', 'a':
			return DirLeft
		case 'j', 's':
			return DirDown
		case 'k', 'w':
			return DirUp
		case 'l', 'd':
			return DirRight
		}
	}

	return DirNone
}

// ParseAction converts a key event into an Action.
// Movement keys return ActionMove, quit keys return ActionQuit.
func ParseAction(key tcell.Key, r rune) Action {
	// Check for quit
	if key == tcell.KeyEscape {
		return ActionQuit
	}
	if key == tcell.KeyRune && (r == 'q' || r == 'Q') {
		return ActionQuit
	}

	// Check for confirm
	if key == tcell.KeyEnter {
		return ActionConfirm
	}

	// Check for inventory
	if key == tcell.KeyRune && (r == 'i' || r == 'I') {
		return ActionInventory
	}

	// Check for drop mode
	if key == tcell.KeyRune && (r == 'x' || r == 'X') {
		return ActionDropMode
	}

	// Check for item use (number keys 1-9)
	if key == tcell.KeyRune && r >= '1' && r <= '9' {
		return ActionUseItem
	}

	// Check for movement
	if ParseDirection(key, r) != DirNone {
		return ActionMove
	}

	return ActionNone
}

// ParseSlotNumber extracts a slot number (1-9) from a rune.
// Returns the slot number and true if valid, or 0 and false if not a valid slot key.
func ParseSlotNumber(r rune) (int, bool) {
	if r >= '1' && r <= '9' {
		return int(r - '0'), true
	}
	return 0, false
}
