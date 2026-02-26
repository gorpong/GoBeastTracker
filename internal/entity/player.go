package entity

import (
	"beasttracker/internal/ui"
)

// Default player stats
const (
	DefaultPlayerHP      = 100
	DefaultPlayerAttack  = 10
	DefaultPlayerDefense = 2
)

// Player represents the player character in the game
type Player struct {
	X         int
	Y         int
	Glyph     rune
	HP        int
	MaxHP     int
	Attack    int
	Defense   int
	Dead      bool
	Inventory *Inventory
}

// NewPlayer creates a new player at the specified position
func NewPlayer(x, y int) *Player {
	return &Player{
		X:         x,
		Y:         y,
		Glyph:     '@',
		HP:        DefaultPlayerHP,
		MaxHP:     DefaultPlayerHP,
		Attack:    DefaultPlayerAttack,
		Defense:   DefaultPlayerDefense,
		Dead:      false,
		Inventory: NewInventory(DefaultInventoryCapacity),
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

// TakeDamage reduces the player's HP by the specified amount
func (p *Player) TakeDamage(damage int) {
	p.HP -= damage
	if p.HP <= 0 {
		p.HP = 0
		p.Dead = true
	}
}

// Heal restores HP to the player, not exceeding MaxHP
func (p *Player) Heal(amount int) {
	p.HP += amount
	if p.HP > p.MaxHP {
		p.HP = p.MaxHP
	}
}

// IsAlive returns true if the player is still alive
func (p *Player) IsAlive() bool {
	return !p.Dead
}
