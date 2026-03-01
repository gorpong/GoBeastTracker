package ui

import (
	"github.com/gdamore/tcell/v2"
)

type Direction int

const (
	DirNone Direction = iota
	DirUp
	DirDown
	DirLeft
	DirRight
)

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

type Action int

const (
	ActionNone Action = iota
	ActionMove
	ActionQuit
	ActionInventory
	ActionDropMode
	ActionUseItem
	ActionConfirm
	ActionCraft
	ActionNextHunt
	ActionBackspace
)

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
	case ActionCraft:
		return "Craft"
	case ActionNextHunt:
		return "NextHunt"
	case ActionBackspace:
		return "Backspace"
	default:
		return "Unknown"
	}
}

func ParseDirection(key tcell.Key, r rune) Direction {
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

func ParseAction(key tcell.Key, r rune) Action {
	if key == tcell.KeyEscape {
		return ActionQuit
	}
	if key == tcell.KeyRune && (r == 'q' || r == 'Q') {
		return ActionQuit
	}

	if key == tcell.KeyEnter {
		return ActionConfirm
	}

	if key == tcell.KeyBackspace || key == tcell.KeyBackspace2 {
		return ActionBackspace
	}

	if key == tcell.KeyRune && (r == 'i' || r == 'I') {
		return ActionInventory
	}

	if key == tcell.KeyRune && (r == 'x' || r == 'X') {
		return ActionDropMode
	}

	if key == tcell.KeyRune && (r == 'c' || r == 'C') {
		return ActionCraft
	}

	if key == tcell.KeyRune && (r == 'n' || r == 'N') {
		return ActionNextHunt
	}

	if key == tcell.KeyRune && r >= '1' && r <= '9' {
		return ActionUseItem
	}

	if ParseDirection(key, r) != DirNone {
		return ActionMove
	}

	return ActionNone
}

func ParseSlotNumber(r rune) (int, bool) {
	if r >= '1' && r <= '9' {
		return int(r - '0'), true
	}
	return 0, false
}

// IsLetter returns true if the rune is A-Z (for initials entry)
func IsLetter(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

// ToUpper converts a letter to uppercase
func ToUpper(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - 32
	}
	return r
}
