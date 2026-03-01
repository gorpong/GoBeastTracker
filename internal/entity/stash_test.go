package entity

import (
	"testing"
)

func TestNewEquipmentStash(t *testing.T) {
	stash := NewEquipmentStash()

	if stash == nil {
		t.Fatal("NewEquipmentStash returned nil")
	}

	if stash.Count() != 0 {
		t.Errorf("New stash count = %d, want 0", stash.Count())
	}
}

func TestEquipmentStashAdd(t *testing.T) {
	stash := NewEquipmentStash()

	sword := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)
	stash.Add(sword)

	if stash.Count() != 1 {
		t.Errorf("Stash count = %d, want 1", stash.Count())
	}
}

func TestEquipmentStashAddMultiple(t *testing.T) {
	stash := NewEquipmentStash()

	sword := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)
	armor := NewEquipment("Leather Armor", SlotArmor, 0, 2, 10)
	charm := NewEquipment("Hunter's Charm", SlotCharm, 1, 1, 0)

	stash.Add(sword)
	stash.Add(armor)
	stash.Add(charm)

	if stash.Count() != 3 {
		t.Errorf("Stash count = %d, want 3", stash.Count())
	}
}

func TestEquipmentStashGetBySlot(t *testing.T) {
	stash := NewEquipmentStash()

	sword1 := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)
	sword2 := NewEquipment("Wyvern Blade", SlotWeapon, 6, 0, 0)
	armor := NewEquipment("Leather Armor", SlotArmor, 0, 2, 10)

	stash.Add(sword1)
	stash.Add(sword2)
	stash.Add(armor)

	weapons := stash.GetBySlot(SlotWeapon)
	if len(weapons) != 2 {
		t.Errorf("Weapons count = %d, want 2", len(weapons))
	}

	armors := stash.GetBySlot(SlotArmor)
	if len(armors) != 1 {
		t.Errorf("Armors count = %d, want 1", len(armors))
	}

	charms := stash.GetBySlot(SlotCharm)
	if len(charms) != 0 {
		t.Errorf("Charms count = %d, want 0", len(charms))
	}
}

func TestEquipmentStashRemove(t *testing.T) {
	stash := NewEquipmentStash()

	sword := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)
	stash.Add(sword)

	removed := stash.Remove(sword)
	if !removed {
		t.Error("Remove should return true for existing item")
	}

	if stash.Count() != 0 {
		t.Errorf("Stash count after remove = %d, want 0", stash.Count())
	}
}

func TestEquipmentStashRemoveNonexistent(t *testing.T) {
	stash := NewEquipmentStash()

	sword := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)

	removed := stash.Remove(sword)
	if removed {
		t.Error("Remove should return false for nonexistent item")
	}
}

func TestEquipmentStashGetAll(t *testing.T) {
	stash := NewEquipmentStash()

	sword := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)
	armor := NewEquipment("Leather Armor", SlotArmor, 0, 2, 10)

	stash.Add(sword)
	stash.Add(armor)

	all := stash.GetAll()
	if len(all) != 2 {
		t.Errorf("GetAll length = %d, want 2", len(all))
	}
}

func TestEquipmentStashContains(t *testing.T) {
	stash := NewEquipmentStash()

	sword := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)
	otherSword := NewEquipment("Bone Knife", SlotWeapon, 2, 0, 0)

	stash.Add(sword)

	if !stash.Contains(sword) {
		t.Error("Stash should contain added sword")
	}

	if stash.Contains(otherSword) {
		t.Error("Stash should not contain unadded sword")
	}
}

func TestEquipmentStashFindByName(t *testing.T) {
	stash := NewEquipmentStash()

	sword := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)
	armor := NewEquipment("Leather Armor", SlotArmor, 0, 2, 10)

	stash.Add(sword)
	stash.Add(armor)

	found := stash.FindByName("Iron Sword")
	if found == nil {
		t.Fatal("FindByName should find Iron Sword")
	}
	if found.Name != "Iron Sword" {
		t.Errorf("Found name = %q, want \"Iron Sword\"", found.Name)
	}

	notFound := stash.FindByName("Wyvern Blade")
	if notFound != nil {
		t.Error("FindByName should return nil for nonexistent item")
	}
}

func TestEquipmentStashClear(t *testing.T) {
	stash := NewEquipmentStash()

	stash.Add(NewEquipment("Sword", SlotWeapon, 1, 0, 0))
	stash.Add(NewEquipment("Armor", SlotArmor, 0, 1, 0))

	stash.Clear()

	if stash.Count() != 0 {
		t.Errorf("Stash count after clear = %d, want 0", stash.Count())
	}
}
