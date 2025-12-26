package entity

import (
	"beasttracker/internal/ui"
)

// AIType represents the behavior pattern of a monster
type AIType int

const (
	AIWander AIType = iota // Wanders randomly
	AIChase                // Chases player when in sight
)

// String returns the string representation of an AIType
func (a AIType) String() string {
	switch a {
	case AIWander:
		return "Wander"
	case AIChase:
		return "Chase"
	default:
		return "Unknown"
	}
}

// Monster represents an enemy in the game
type Monster struct {
	Name   string
	Glyph  rune
	X      int
	Y      int
	HP     int
	MaxHP  int
	Attack int
	AI     AIType
	Dead   bool
}

// NewMonster creates a new monster with the specified attributes
func NewMonster(name string, glyph rune, x, y, hp, attack int) *Monster {
	return &Monster{
		Name:   name,
		Glyph:  glyph,
		X:      x,
		Y:      y,
		HP:     hp,
		MaxHP:  hp,
		Attack: attack,
		AI:     AIWander,
		Dead:   false,
	}
}

// Position returns the monster's current coordinates
func (m *Monster) Position() (int, int) {
	return m.X, m.Y
}

// SetPosition sets the monster's position to the specified coordinates
func (m *Monster) SetPosition(x, y int) {
	m.X = x
	m.Y = y
}

// Move moves the monster in the specified direction
func (m *Monster) Move(dir ui.Direction) {
	dx, dy := dir.Delta()
	m.X += dx
	m.Y += dy
}

// TakeDamage reduces the monster's HP by the specified amount
func (m *Monster) TakeDamage(damage int) {
	m.HP -= damage
	if m.HP <= 0 {
		m.HP = 0
		m.Dead = true
	}
}

// IsAlive returns true if the monster is still alive
func (m *Monster) IsAlive() bool {
	return !m.Dead
}
