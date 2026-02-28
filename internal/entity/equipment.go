package entity

import (
	"fmt"
	"strings"
)

// EquipmentSlot represents where equipment can be equipped
type EquipmentSlot int

const (
	SlotWeapon EquipmentSlot = iota
	SlotArmor
	SlotCharm
)

// String returns the display name of the equipment slot
func (es EquipmentSlot) String() string {
	switch es {
	case SlotWeapon:
		return "Weapon"
	case SlotArmor:
		return "Armor"
	case SlotCharm:
		return "Charm"
	default:
		return "Unknown"
	}
}

// Equipment represents a craftable piece of equipment
type Equipment struct {
	Name       string
	Slot       EquipmentSlot
	ATKBonus   int
	DEFBonus   int
	MaxHPBonus int
}

// NewEquipment creates a new piece of equipment
func NewEquipment(name string, slot EquipmentSlot, atkBonus, defBonus, maxHPBonus int) *Equipment {
	return &Equipment{
		Name:       name,
		Slot:       slot,
		ATKBonus:   atkBonus,
		DEFBonus:   defBonus,
		MaxHPBonus: maxHPBonus,
	}
}

// StatsString returns a compact string showing the stat bonuses
func (e *Equipment) StatsString() string {
	var parts []string

	if e.ATKBonus > 0 {
		parts = append(parts, fmt.Sprintf("+%d ATK", e.ATKBonus))
	}
	if e.DEFBonus > 0 {
		parts = append(parts, fmt.Sprintf("+%d DEF", e.DEFBonus))
	}
	if e.MaxHPBonus > 0 {
		parts = append(parts, fmt.Sprintf("+%d MaxHP", e.MaxHPBonus))
	}

	if len(parts) == 0 {
		return "No bonuses"
	}

	return strings.Join(parts, ", ")
}

// Description returns a full description of the equipment
func (e *Equipment) Description() string {
	return fmt.Sprintf("%s (%s): %s", e.Name, e.Slot.String(), e.StatsString())
}
