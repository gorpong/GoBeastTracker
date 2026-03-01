package entity

import (
	"beasttracker/internal/ui"
)

const (
	DefaultPlayerHP      = 100
	DefaultPlayerAttack  = 10
	DefaultPlayerDefense = 2
)

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
	EquipmentStash *EquipmentStash
	EquippedWeapon *Equipment
	EquippedArmor  *Equipment
	EquippedCharm  *Equipment
}

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
		EquipmentStash: NewEquipmentStash(),
		EquippedWeapon: nil,
		EquippedArmor:  nil,
		EquippedCharm:  nil,
	}
}

func (p *Player) Move(dir ui.Direction) {
	dx, dy := dir.Delta()
	p.X += dx
	p.Y += dy
}

func (p *Player) Position() (int, int) {
	return p.X, p.Y
}

func (p *Player) SetPosition(x, y int) {
	p.X = x
	p.Y = y
}

func (p *Player) TakeDamage(damage int) {
	p.HP -= damage
	if p.HP <= 0 {
		p.HP = 0
		p.Dead = true
	}
}

func (p *Player) Heal(amount int) {
	p.HP += amount
	effectiveMax := p.EffectiveMaxHP()
	if p.HP > effectiveMax {
		p.HP = effectiveMax
	}
}

func (p *Player) IsAlive() bool {
	return !p.Dead
}

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

// EquipFromStash equips an item from the stash, returning old item to stash
// Returns false if the equipment is not in the stash
func (p *Player) EquipFromStash(equipment *Equipment) bool {
	if !p.EquipmentStash.Contains(equipment) {
		return false
	}

	p.EquipmentStash.Remove(equipment)

	old := p.Equip(equipment)
	if old != nil {
		p.EquipmentStash.Add(old)
	}

	return true
}

// UnequipToStash removes equipped item and places it in the stash
func (p *Player) UnequipToStash(slot EquipmentSlot) {
	removed := p.Unequip(slot)
	if removed != nil {
		p.EquipmentStash.Add(removed)
	}
}

// GetAllEquipped returns all currently equipped items
func (p *Player) GetAllEquipped() []*Equipment {
	result := make([]*Equipment, 0, 3)

	if p.EquippedWeapon != nil {
		result = append(result, p.EquippedWeapon)
	}
	if p.EquippedArmor != nil {
		result = append(result, p.EquippedArmor)
	}
	if p.EquippedCharm != nil {
		result = append(result, p.EquippedCharm)
	}

	return result
}

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
