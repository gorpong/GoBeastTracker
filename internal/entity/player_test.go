package entity

import (
	"testing"

	"beasttracker/internal/ui"
)

// TestNewPlayer verifies that a new player is created with correct initial values
func TestNewPlayer(t *testing.T) {
	p := NewPlayer(5, 10)

	if p.X != 5 {
		t.Errorf("NewPlayer X = %d, want 5", p.X)
	}
	if p.Y != 10 {
		t.Errorf("NewPlayer Y = %d, want 10", p.Y)
	}
	if p.Glyph != '@' {
		t.Errorf("NewPlayer Glyph = %q, want '@'", p.Glyph)
	}
}

// TestPlayerMove verifies that player moves correctly in each direction
func TestPlayerMove(t *testing.T) {
	tests := []struct {
		name      string
		startX    int
		startY    int
		direction ui.Direction
		wantX     int
		wantY     int
	}{
		{"Move up", 10, 10, ui.DirUp, 10, 9},
		{"Move down", 10, 10, ui.DirDown, 10, 11},
		{"Move left", 10, 10, ui.DirLeft, 9, 10},
		{"Move right", 10, 10, ui.DirRight, 11, 10},
		{"No movement", 10, 10, ui.DirNone, 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPlayer(tt.startX, tt.startY)
			p.Move(tt.direction)

			if p.X != tt.wantX {
				t.Errorf("After Move(%v): X = %d, want %d", tt.direction, p.X, tt.wantX)
			}
			if p.Y != tt.wantY {
				t.Errorf("After Move(%v): Y = %d, want %d", tt.direction, p.Y, tt.wantY)
			}
		})
	}
}

// TestPlayerPosition verifies Position() returns correct coordinates
func TestPlayerPosition(t *testing.T) {
	p := NewPlayer(7, 3)
	x, y := p.Position()

	if x != 7 || y != 3 {
		t.Errorf("Position() = (%d, %d), want (7, 3)", x, y)
	}
}

// TestPlayerSetPosition verifies SetPosition correctly updates coordinates
func TestPlayerSetPosition(t *testing.T) {
	p := NewPlayer(0, 0)
	p.SetPosition(15, 20)

	if p.X != 15 || p.Y != 20 {
		t.Errorf("After SetPosition(15, 20): (%d, %d), want (15, 20)", p.X, p.Y)
	}
}

// TestPlayerMultipleMoves verifies multiple consecutive moves work correctly
func TestPlayerMultipleMoves(t *testing.T) {
	p := NewPlayer(10, 10)

	// Move in a square pattern
	p.Move(ui.DirUp)    // 10, 9
	p.Move(ui.DirRight) // 11, 9
	p.Move(ui.DirDown)  // 11, 10
	p.Move(ui.DirLeft)  // 10, 10

	if p.X != 10 || p.Y != 10 {
		t.Errorf("After square pattern: (%d, %d), want (10, 10)", p.X, p.Y)
	}
}

// TestPlayerStats verifies player has combat stats
func TestPlayerStats(t *testing.T) {
	p := NewPlayer(0, 0)

	if p.HP <= 0 {
		t.Errorf("Player HP = %d, want > 0", p.HP)
	}
	if p.MaxHP <= 0 {
		t.Errorf("Player MaxHP = %d, want > 0", p.MaxHP)
	}
	if p.Attack <= 0 {
		t.Errorf("Player Attack = %d, want > 0", p.Attack)
	}
	if p.Defense < 0 {
		t.Errorf("Player Defense = %d, want >= 0", p.Defense)
	}
}

// TestPlayerTakeDamage verifies damage reduces HP
func TestPlayerTakeDamage(t *testing.T) {
	p := NewPlayer(0, 0)
	initialHP := p.HP

	p.TakeDamage(5)

	if p.HP != initialHP-5 {
		t.Errorf("After taking 5 damage: HP = %d, want %d", p.HP, initialHP-5)
	}
	if p.Dead {
		t.Error("Player should not be dead after minor damage")
	}
}

// TestPlayerDeath verifies player dies when HP reaches 0
func TestPlayerDeath(t *testing.T) {
	p := NewPlayer(0, 0)

	p.TakeDamage(p.HP)

	if p.HP != 0 {
		t.Errorf("After fatal damage: HP = %d, want 0", p.HP)
	}
	if !p.Dead {
		t.Error("Player should be dead at 0 HP")
	}
}

// TestPlayerOverkillDamage verifies HP doesn't go negative
func TestPlayerOverkillDamage(t *testing.T) {
	p := NewPlayer(0, 0)

	p.TakeDamage(p.HP + 100)

	if p.HP < 0 {
		t.Errorf("HP should not be negative: HP = %d", p.HP)
	}
	if p.HP != 0 {
		t.Errorf("HP should be 0 after overkill: HP = %d", p.HP)
	}
}

// TestPlayerIsAlive verifies IsAlive returns correct status
func TestPlayerIsAlive(t *testing.T) {
	p := NewPlayer(0, 0)

	if !p.IsAlive() {
		t.Error("New player should be alive")
	}

	p.TakeDamage(p.HP)

	if p.IsAlive() {
		t.Error("Player should not be alive after fatal damage")
	}
}

// TestPlayerHeal verifies healing restores HP correctly
func TestPlayerHeal(t *testing.T) {
	tests := []struct {
		name       string
		currentHP  int
		maxHP      int
		healAmount int
		wantHP     int
	}{
		{
			name:       "Heal partial damage",
			currentHP:  50,
			maxHP:      100,
			healAmount: 25,
			wantHP:     75,
		},
		{
			name:       "Heal does not exceed max HP",
			currentHP:  90,
			maxHP:      100,
			healAmount: 25,
			wantHP:     100,
		},
		{
			name:       "Heal from critical HP",
			currentHP:  5,
			maxHP:      100,
			healAmount: 60,
			wantHP:     65,
		},
		{
			name:       "Heal at full HP does nothing",
			currentHP:  100,
			maxHP:      100,
			healAmount: 25,
			wantHP:     100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := NewPlayer(0, 0)
			player.HP = tt.currentHP
			player.MaxHP = tt.maxHP

			player.Heal(tt.healAmount)

			if player.HP != tt.wantHP {
				t.Errorf("After Heal(%d): HP = %d, want %d", tt.healAmount, player.HP, tt.wantHP)
			}
		})
	}
}

// TestPlayerHasInventory verifies player has an inventory
func TestPlayerHasInventory(t *testing.T) {
	player := NewPlayer(0, 0)

	if player.Inventory == nil {
		t.Fatal("Player should have an inventory")
	}

	if player.Inventory.Capacity() != DefaultInventoryCapacity {
		t.Errorf("Player inventory capacity = %d, want %d",
			player.Inventory.Capacity(), DefaultInventoryCapacity)
	}
}

// TestPlayerMaterialPouch verifies player has material pouch
func TestPlayerMaterialPouch(t *testing.T) {
	player := NewPlayer(0, 0)

	if player.MaterialPouch == nil {
		t.Error("NewPlayer should initialize MaterialPouch")
	}
	if player.MaterialPouch.TotalCount() != 0 {
		t.Error("New player's MaterialPouch should be empty")
	}
}

// TestPlayerEquipmentSlots verifies player has equipment slots
func TestPlayerEquipmentSlots(t *testing.T) {
	player := NewPlayer(0, 0)

	if player.EquippedWeapon != nil {
		t.Error("New player should have no weapon equipped")
	}
	if player.EquippedArmor != nil {
		t.Error("New player should have no armor equipped")
	}
	if player.EquippedCharm != nil {
		t.Error("New player should have no charm equipped")
	}
}

// TestPlayerEquipWeapon verifies equipping a weapon
func TestPlayerEquipWeapon(t *testing.T) {
	player := NewPlayer(0, 0)
	weapon := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)

	player.Equip(weapon)

	if player.EquippedWeapon != weapon {
		t.Error("Equip() should set EquippedWeapon for weapon slot")
	}
}

// TestPlayerEquipArmor verifies equipping armor
func TestPlayerEquipArmor(t *testing.T) {
	player := NewPlayer(0, 0)
	armor := NewEquipment("Leather Armor", SlotArmor, 0, 2, 10)

	player.Equip(armor)

	if player.EquippedArmor != armor {
		t.Error("Equip() should set EquippedArmor for armor slot")
	}
}

// TestPlayerEquipCharm verifies equipping a charm
func TestPlayerEquipCharm(t *testing.T) {
	player := NewPlayer(0, 0)
	charm := NewEquipment("Hunter's Charm", SlotCharm, 1, 1, 0)

	player.Equip(charm)

	if player.EquippedCharm != charm {
		t.Error("Equip() should set EquippedCharm for charm slot")
	}
}

// TestPlayerEquipReplacesExisting verifies equipping replaces old equipment
func TestPlayerEquipReplacesExisting(t *testing.T) {
	player := NewPlayer(0, 0)
	weapon1 := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)
	weapon2 := NewEquipment("Wyvern Blade", SlotWeapon, 6, 0, 0)

	player.Equip(weapon1)
	oldWeapon := player.Equip(weapon2)

	if player.EquippedWeapon != weapon2 {
		t.Error("Equip() should replace existing weapon")
	}
	if oldWeapon != weapon1 {
		t.Error("Equip() should return the replaced equipment")
	}
}

// TestPlayerUnequip verifies removing equipment
func TestPlayerUnequip(t *testing.T) {
	player := NewPlayer(0, 0)
	weapon := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)
	player.Equip(weapon)

	removed := player.Unequip(SlotWeapon)

	if player.EquippedWeapon != nil {
		t.Error("Unequip() should clear the equipment slot")
	}
	if removed != weapon {
		t.Error("Unequip() should return the removed equipment")
	}
}

// TestPlayerUnequipEmpty verifies unequipping empty slot
func TestPlayerUnequipEmpty(t *testing.T) {
	player := NewPlayer(0, 0)

	removed := player.Unequip(SlotWeapon)

	if removed != nil {
		t.Error("Unequip() on empty slot should return nil")
	}
}

// TestPlayerEffectiveAttackNoEquipment verifies base ATK without equipment
func TestPlayerEffectiveAttackNoEquipment(t *testing.T) {
	player := NewPlayer(0, 0)

	if player.EffectiveAttack() != DefaultPlayerAttack {
		t.Errorf("EffectiveAttack() without equipment = %d, want %d",
			player.EffectiveAttack(), DefaultPlayerAttack)
	}
}

// TestPlayerEffectiveAttackWithWeapon verifies ATK with weapon equipped
func TestPlayerEffectiveAttackWithWeapon(t *testing.T) {
	player := NewPlayer(0, 0)
	weapon := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)
	player.Equip(weapon)

	expected := DefaultPlayerAttack + 3
	if player.EffectiveAttack() != expected {
		t.Errorf("EffectiveAttack() with weapon = %d, want %d",
			player.EffectiveAttack(), expected)
	}
}

// TestPlayerEffectiveAttackWithCharm verifies ATK includes charm bonus
func TestPlayerEffectiveAttackWithCharm(t *testing.T) {
	player := NewPlayer(0, 0)
	charm := NewEquipment("Attack Charm", SlotCharm, 2, 0, 0)
	player.Equip(charm)

	expected := DefaultPlayerAttack + 2
	if player.EffectiveAttack() != expected {
		t.Errorf("EffectiveAttack() with charm = %d, want %d",
			player.EffectiveAttack(), expected)
	}
}

// TestPlayerEffectiveAttackCombined verifies ATK sums all equipment
func TestPlayerEffectiveAttackCombined(t *testing.T) {
	player := NewPlayer(0, 0)
	weapon := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)
	charm := NewEquipment("Attack Charm", SlotCharm, 2, 0, 0)
	player.Equip(weapon)
	player.Equip(charm)

	expected := DefaultPlayerAttack + 3 + 2
	if player.EffectiveAttack() != expected {
		t.Errorf("EffectiveAttack() combined = %d, want %d",
			player.EffectiveAttack(), expected)
	}
}

// TestPlayerEffectiveDefenseNoEquipment verifies base DEF without equipment
func TestPlayerEffectiveDefenseNoEquipment(t *testing.T) {
	player := NewPlayer(0, 0)

	if player.EffectiveDefense() != DefaultPlayerDefense {
		t.Errorf("EffectiveDefense() without equipment = %d, want %d",
			player.EffectiveDefense(), DefaultPlayerDefense)
	}
}

// TestPlayerEffectiveDefenseWithArmor verifies DEF with armor equipped
func TestPlayerEffectiveDefenseWithArmor(t *testing.T) {
	player := NewPlayer(0, 0)
	armor := NewEquipment("Leather Armor", SlotArmor, 0, 2, 0)
	player.Equip(armor)

	expected := DefaultPlayerDefense + 2
	if player.EffectiveDefense() != expected {
		t.Errorf("EffectiveDefense() with armor = %d, want %d",
			player.EffectiveDefense(), expected)
	}
}

// TestPlayerEffectiveDefenseCombined verifies DEF sums all equipment
func TestPlayerEffectiveDefenseCombined(t *testing.T) {
	player := NewPlayer(0, 0)
	armor := NewEquipment("Leather Armor", SlotArmor, 0, 2, 0)
	charm := NewEquipment("Defense Charm", SlotCharm, 0, 1, 0)
	player.Equip(armor)
	player.Equip(charm)

	expected := DefaultPlayerDefense + 2 + 1
	if player.EffectiveDefense() != expected {
		t.Errorf("EffectiveDefense() combined = %d, want %d",
			player.EffectiveDefense(), expected)
	}
}

// TestPlayerEffectiveMaxHPNoEquipment verifies base MaxHP without equipment
func TestPlayerEffectiveMaxHPNoEquipment(t *testing.T) {
	player := NewPlayer(0, 0)

	if player.EffectiveMaxHP() != DefaultPlayerHP {
		t.Errorf("EffectiveMaxHP() without equipment = %d, want %d",
			player.EffectiveMaxHP(), DefaultPlayerHP)
	}
}

// TestPlayerEffectiveMaxHPWithArmor verifies MaxHP with armor equipped
func TestPlayerEffectiveMaxHPWithArmor(t *testing.T) {
	player := NewPlayer(0, 0)
	armor := NewEquipment("Leather Armor", SlotArmor, 0, 0, 10)
	player.Equip(armor)

	expected := DefaultPlayerHP + 10
	if player.EffectiveMaxHP() != expected {
		t.Errorf("EffectiveMaxHP() with armor = %d, want %d",
			player.EffectiveMaxHP(), expected)
	}
}

// TestPlayerHealRespectsEffectiveMaxHP verifies heal respects equipment bonus
func TestPlayerHealRespectsEffectiveMaxHP(t *testing.T) {
	player := NewPlayer(0, 0)
	armor := NewEquipment("Leather Armor", SlotArmor, 0, 0, 20)
	player.Equip(armor)
	player.TakeDamage(50)

	// Heal should respect effective max HP (100 + 20 = 120)
	player.Heal(100)

	expectedMaxHP := DefaultPlayerHP + 20
	if player.HP != expectedMaxHP {
		t.Errorf("After Heal with armor bonus: HP = %d, want %d", player.HP, expectedMaxHP)
	}
}

func TestPlayerHasEquipmentStash(t *testing.T) {
	player := NewPlayer(0, 0)

	if player.EquipmentStash == nil {
		t.Error("NewPlayer should initialize EquipmentStash")
	}

	if player.EquipmentStash.Count() != 0 {
		t.Error("New player's EquipmentStash should be empty")
	}
}

func TestPlayerEquipFromStash(t *testing.T) {
	player := NewPlayer(0, 0)

	sword := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)
	player.EquipmentStash.Add(sword)

	player.EquipFromStash(sword)

	if player.EquippedWeapon != sword {
		t.Error("EquipFromStash should equip the item")
	}

	if player.EquipmentStash.Contains(sword) {
		t.Error("Equipped item should be removed from stash")
	}
}

func TestPlayerEquipFromStashReturnsOldToStash(t *testing.T) {
	player := NewPlayer(0, 0)

	oldSword := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)
	newSword := NewEquipment("Wyvern Blade", SlotWeapon, 6, 0, 0)

	player.Equip(oldSword)
	player.EquipmentStash.Add(newSword)

	player.EquipFromStash(newSword)

	if player.EquippedWeapon != newSword {
		t.Error("New weapon should be equipped")
	}

	if !player.EquipmentStash.Contains(oldSword) {
		t.Error("Old weapon should be returned to stash")
	}

	if player.EquipmentStash.Contains(newSword) {
		t.Error("New weapon should not be in stash after equipping")
	}
}

func TestPlayerUnequipToStash(t *testing.T) {
	player := NewPlayer(0, 0)

	sword := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)
	player.Equip(sword)

	player.UnequipToStash(SlotWeapon)

	if player.EquippedWeapon != nil {
		t.Error("Weapon slot should be empty after unequip")
	}

	if !player.EquipmentStash.Contains(sword) {
		t.Error("Unequipped weapon should be in stash")
	}
}

func TestPlayerUnequipToStashEmptySlot(t *testing.T) {
	player := NewPlayer(0, 0)

	// Should not panic on empty slot
	player.UnequipToStash(SlotWeapon)

	if player.EquipmentStash.Count() != 0 {
		t.Error("Stash should remain empty when unequipping empty slot")
	}
}

func TestPlayerEquipFromStashNotInStash(t *testing.T) {
	player := NewPlayer(0, 0)

	sword := NewEquipment("Iron Sword", SlotWeapon, 3, 0, 0)

	// Equipment not in stash - should not equip
	result := player.EquipFromStash(sword)

	if result {
		t.Error("EquipFromStash should return false for item not in stash")
	}

	if player.EquippedWeapon != nil {
		t.Error("Should not equip item that wasn't in stash")
	}
}

func TestPlayerGetAllEquippedItems(t *testing.T) {
	player := NewPlayer(0, 0)

	if len(player.GetAllEquipped()) != 0 {
		t.Error("New player should have no equipped items")
	}

	weapon := NewEquipment("Sword", SlotWeapon, 3, 0, 0)
	armor := NewEquipment("Armor", SlotArmor, 0, 2, 10)

	player.Equip(weapon)
	player.Equip(armor)

	equipped := player.GetAllEquipped()
	if len(equipped) != 2 {
		t.Errorf("Should have 2 equipped items, got %d", len(equipped))
	}
}
