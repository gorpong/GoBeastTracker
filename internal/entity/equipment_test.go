package entity

import (
	"testing"
)

// TestEquipmentSlotString verifies slot type string representation
func TestEquipmentSlotString(t *testing.T) {
	tests := []struct {
		slot EquipmentSlot
		want string
	}{
		{SlotWeapon, "Weapon"},
		{SlotArmor, "Armor"},
		{SlotCharm, "Charm"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.slot.String(); got != tt.want {
				t.Errorf("EquipmentSlot.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestNewEquipment verifies equipment creation with correct values
func TestNewEquipment(t *testing.T) {
	equipment := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)

	if equipment.Name != "Iron Sword" {
		t.Errorf("Equipment Name = %q, want \"Iron Sword\"", equipment.Name)
	}
	if equipment.Slot != SlotWeapon {
		t.Errorf("Equipment Slot = %v, want SlotWeapon", equipment.Slot)
	}
	if equipment.ATKBonus != 3 {
		t.Errorf("Equipment ATKBonus = %d, want 3", equipment.ATKBonus)
	}
	if equipment.DEFBonus != 0 {
		t.Errorf("Equipment DEFBonus = %d, want 0", equipment.DEFBonus)
	}
	if equipment.MaxHPBonus != 0 {
		t.Errorf("Equipment MaxHPBonus = %d, want 0", equipment.MaxHPBonus)
	}
}

// TestEquipmentWeapon verifies weapon equipment stats
func TestEquipmentWeapon(t *testing.T) {
	weapon := NewEquipment("Wyvern Blade", SlotWeapon, 6, 0, 0)

	if weapon.Slot != SlotWeapon {
		t.Errorf("Weapon Slot = %v, want SlotWeapon", weapon.Slot)
	}
	if weapon.ATKBonus != 6 {
		t.Errorf("Weapon ATKBonus = %d, want 6", weapon.ATKBonus)
	}
}

// TestEquipmentArmor verifies armor equipment stats
func TestEquipmentArmor(t *testing.T) {
	armor := NewEquipment("Leather Armor", SlotArmor, 0, 2, 10)

	if armor.Slot != SlotArmor {
		t.Errorf("Armor Slot = %v, want SlotArmor", armor.Slot)
	}
	if armor.DEFBonus != 2 {
		t.Errorf("Armor DEFBonus = %d, want 2", armor.DEFBonus)
	}
	if armor.MaxHPBonus != 10 {
		t.Errorf("Armor MaxHPBonus = %d, want 10", armor.MaxHPBonus)
	}
}

// TestEquipmentCharm verifies charm equipment stats
func TestEquipmentCharm(t *testing.T) {
	charm := NewEquipment("Hunter's Charm", SlotCharm, 1, 1, 0)

	if charm.Slot != SlotCharm {
		t.Errorf("Charm Slot = %v, want SlotCharm", charm.Slot)
	}
	if charm.ATKBonus != 1 {
		t.Errorf("Charm ATKBonus = %d, want 1", charm.ATKBonus)
	}
	if charm.DEFBonus != 1 {
		t.Errorf("Charm DEFBonus = %d, want 1", charm.DEFBonus)
	}
}

// TestEquipmentMixedStats verifies equipment can have multiple stat bonuses
func TestEquipmentMixedStats(t *testing.T) {
	equipment := NewEquipment("Balanced Gear", SlotCharm, 2, 3, 15)

	if equipment.ATKBonus != 2 {
		t.Errorf("ATKBonus = %d, want 2", equipment.ATKBonus)
	}
	if equipment.DEFBonus != 3 {
		t.Errorf("DEFBonus = %d, want 3", equipment.DEFBonus)
	}
	if equipment.MaxHPBonus != 15 {
		t.Errorf("MaxHPBonus = %d, want 15", equipment.MaxHPBonus)
	}
}

// TestEquipmentTotalATK verifies total ATK calculation helper
func TestEquipmentTotalATK(t *testing.T) {
	weapon := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)
	charm := NewEquipment("Attack Charm", SlotCharm, 2, 0, 0)

	total := weapon.ATKBonus + charm.ATKBonus

	if total != 5 {
		t.Errorf("Total ATK from equipment = %d, want 5", total)
	}
}

// TestEquipmentTotalDEF verifies total DEF calculation helper
func TestEquipmentTotalDEF(t *testing.T) {
	armor := NewEquipment("Leather Armor", SlotArmor, 0, 2, 0)
	charm := NewEquipment("Defense Charm", SlotCharm, 0, 1, 0)

	total := armor.DEFBonus + charm.DEFBonus

	if total != 3 {
		t.Errorf("Total DEF from equipment = %d, want 3", total)
	}
}

// TestEquipmentDescription verifies equipment description generation
func TestEquipmentDescription(t *testing.T) {
	tests := []struct {
		name     string
		equip    *Equipment
		wantDesc string
	}{
		{
			name:     "Weapon with ATK only",
			equip:    NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0),
			wantDesc: "Iron Sword (Weapon): +3 ATK",
		},
		{
			name:     "Armor with DEF and HP",
			equip:    NewEquipment("Leather Armor", SlotArmor, 0, 2, 10),
			wantDesc: "Leather Armor (Armor): +2 DEF, +10 MaxHP",
		},
		{
			name:     "Charm with mixed stats",
			equip:    NewEquipment("Hunter's Charm", SlotCharm, 1, 1, 0),
			wantDesc: "Hunter's Charm (Charm): +1 ATK, +1 DEF",
		},
		{
			name:     "Equipment with all stats",
			equip:    NewEquipment("Ultimate Gear", SlotCharm, 2, 3, 15),
			wantDesc: "Ultimate Gear (Charm): +2 ATK, +3 DEF, +15 MaxHP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.equip.Description(); got != tt.wantDesc {
				t.Errorf("Description() = %q, want %q", got, tt.wantDesc)
			}
		})
	}
}

// TestEquipmentStatsString verifies short stats string
func TestEquipmentStatsString(t *testing.T) {
	tests := []struct {
		name      string
		equip     *Equipment
		wantStats string
	}{
		{
			name:      "ATK only",
			equip:     NewEquipment("Sword", SlotWeapon, 5, 0, 0),
			wantStats: "+5 ATK",
		},
		{
			name:      "DEF only",
			equip:     NewEquipment("Shield", SlotArmor, 0, 3, 0),
			wantStats: "+3 DEF",
		},
		{
			name:      "HP only",
			equip:     NewEquipment("Vitality Charm", SlotCharm, 0, 0, 20),
			wantStats: "+20 MaxHP",
		},
		{
			name:      "Multiple stats",
			equip:     NewEquipment("Mixed", SlotCharm, 2, 1, 5),
			wantStats: "+2 ATK, +1 DEF, +5 MaxHP",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.equip.StatsString(); got != tt.wantStats {
				t.Errorf("StatsString() = %q, want %q", got, tt.wantStats)
			}
		})
	}
}
