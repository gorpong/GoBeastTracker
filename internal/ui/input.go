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
)

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

	// Check for movement
	if ParseDirection(key, r) != DirNone {
		return ActionMove
	}

	return ActionNone
}
