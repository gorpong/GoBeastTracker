package entity

import (
	"beasttracker/internal/ui"
)

// Player represents the player character in the game
type Player struct {
	X     int
	Y     int
	Glyph rune
}

// NewPlayer creates a new player at the specified position
func NewPlayer(x, y int) *Player {
	return &Player{
		X:     x,
		Y:     y,
		Glyph: '@',
	}
}

// Move moves the player in the specified direction
func (p *Player) Move(dir ui.Direction) {
	dx, dy := dir.Delta()
	p.X += dx
	p.Y += dy
}

// Position returns the player's current coordinates
func (p *Player) Position() (int, int) {
	return p.X, p.Y
}

// SetPosition sets the player's position to the specified coordinates
func (p *Player) SetPosition(x, y int) {
	p.X = x
	p.Y = y
}
