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
	X              int
	Y              int
	Glyph          rune
	HP             int
	MaxHP          int
	Attack         int
	Defense        int
	Dead           bool
	Inventory      *Inventory
	MaterialPouch  *MaterialPouch
	EquippedWeapon *Equipment
	EquippedArmor  *Equipment
	EquippedCharm  *Equipment
}

// NewPlayer creates a new player at the specified position
func NewPlayer(x, y int) *Player {
	return &Player{
		X:              x,
		Y:              y,
		Glyph:          '@',
		HP:             DefaultPlayerHP,
		MaxHP:          DefaultPlayerHP,
		Attack:         DefaultPlayerAttack,
		Defense:        DefaultPlayerDefense,
		Dead:           false,
		Inventory:      NewInventory(DefaultInventoryCapacity),
		MaterialPouch:  NewMaterialPouch(),
		EquippedWeapon: nil,
		EquippedArmor:  nil,
		EquippedCharm:  nil,
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

// Heal restores HP to the player, not exceeding EffectiveMaxHP
func (p *Player) Heal(amount int) {
	p.HP += amount
	effectiveMax := p.EffectiveMaxHP()
	if p.HP > effectiveMax {
		p.HP = effectiveMax
	}
}

// IsAlive returns true if the player is still alive
func (p *Player) IsAlive() bool {
	return !p.Dead
}

// Equip equips the given equipment to the appropriate slot.
// Returns the previously equipped item (if any) which was replaced.
func (p *Player) Equip(equipment *Equipment) *Equipment {
	var old *Equipment

	switch equipment.Slot {
	case SlotWeapon:
		old = p.EquippedWeapon
		p.EquippedWeapon = equipment
	case SlotArmor:
		old = p.EquippedArmor
		p.EquippedArmor = equipment
	case SlotCharm:
		old = p.EquippedCharm
		p.EquippedCharm = equipment
	}

	return old
}

// Unequip removes equipment from the specified slot.
// Returns the removed equipment, or nil if slot was empty.
func (p *Player) Unequip(slot EquipmentSlot) *Equipment {
	var removed *Equipment

	switch slot {
	case SlotWeapon:
		removed = p.EquippedWeapon
		p.EquippedWeapon = nil
	case SlotArmor:
		removed = p.EquippedArmor
		p.EquippedArmor = nil
	case SlotCharm:
		removed = p.EquippedCharm
		p.EquippedCharm = nil
	}

	return removed
}

// EffectiveAttack returns the player's total attack including equipment bonuses
func (p *Player) EffectiveAttack() int {
	total := p.Attack

	if p.EquippedWeapon != nil {
		total += p.EquippedWeapon.ATKBonus
	}
	if p.EquippedArmor != nil {
		total += p.EquippedArmor.ATKBonus
	}
	if p.EquippedCharm != nil {
		total += p.EquippedCharm.ATKBonus
	}

	return total
}

// EffectiveDefense returns the player's total defense including equipment bonuses
func (p *Player) EffectiveDefense() int {
	total := p.Defense

	if p.EquippedWeapon != nil {
		total += p.EquippedWeapon.DEFBonus
	}
	if p.EquippedArmor != nil {
		total += p.EquippedArmor.DEFBonus
	}
	if p.EquippedCharm != nil {
		total += p.EquippedCharm.DEFBonus
	}

	return total
}

// EffectiveMaxHP returns the player's total max HP including equipment bonuses
func (p *Player) EffectiveMaxHP() int {
	total := p.MaxHP

	if p.EquippedWeapon != nil {
		total += p.EquippedWeapon.MaxHPBonus
	}
	if p.EquippedArmor != nil {
		total += p.EquippedArmor.MaxHPBonus
	}
	if p.EquippedCharm != nil {
		total += p.EquippedCharm.MaxHPBonus
	}

	return total
}
