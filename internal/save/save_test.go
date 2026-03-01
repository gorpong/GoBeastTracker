package save

import (
	"fmt"
	"testing"

	"beasttracker/internal/entity"
)

func TestNewSaveManager(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSaveManager(tmpDir)

	if sm == nil {
		t.Fatal("NewSaveManager returned nil")
	}
}

func TestSaveDataCreation(t *testing.T) {
	data := NewSaveData("TestSave", 3, 150)

	if data.Name != "TestSave" {
		t.Errorf("Name = %q, want \"TestSave\"", data.Name)
	}
	if data.HuntNumber != 3 {
		t.Errorf("HuntNumber = %d, want 3", data.HuntNumber)
	}
	if data.Score != 150 {
		t.Errorf("Score = %d, want 150", data.Score)
	}
}

func TestSaveDataWithEquipment(t *testing.T) {
	data := NewSaveData("TestSave", 1, 0)

	weapon := entity.NewEquipment("Iron Sword", entity.SlotWeapon, 3, 0, 0)
	data.EquippedWeapon = weapon

	if data.EquippedWeapon == nil {
		t.Error("EquippedWeapon should be set")
	}
	if data.EquippedWeapon.Name != "Iron Sword" {
		t.Errorf("Weapon name = %q, want \"Iron Sword\"", data.EquippedWeapon.Name)
	}
}

func TestSaveManagerSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSaveManager(tmpDir)

	data := NewSaveData("TestSave", 5, 300)
	data.EquippedWeapon = entity.NewEquipment("Wyvern Blade", entity.SlotWeapon, 6, 0, 0)
	data.Materials = map[entity.MaterialType]int{
		entity.MaterialScales: 5,
		entity.MaterialClaws:  3,
	}

	err := sm.Save(data)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := sm.Load("TestSave")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Name != "TestSave" {
		t.Errorf("Loaded name = %q, want \"TestSave\"", loaded.Name)
	}
	if loaded.HuntNumber != 5 {
		t.Errorf("Loaded hunt = %d, want 5", loaded.HuntNumber)
	}
	if loaded.Score != 300 {
		t.Errorf("Loaded score = %d, want 300", loaded.Score)
	}
	if loaded.EquippedWeapon == nil || loaded.EquippedWeapon.Name != "Wyvern Blade" {
		t.Error("Loaded weapon mismatch")
	}
	if loaded.Materials[entity.MaterialScales] != 5 {
		t.Errorf("Loaded scales = %d, want 5", loaded.Materials[entity.MaterialScales])
	}
}

func TestSaveManagerList(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSaveManager(tmpDir)

	sm.Save(NewSaveData("Save1", 1, 100))
	sm.Save(NewSaveData("Save2", 2, 200))
	sm.Save(NewSaveData("Save3", 3, 300))

	saves, err := sm.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(saves) != 3 {
		t.Errorf("Save count = %d, want 3", len(saves))
	}
}

func TestSaveManagerDelete(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSaveManager(tmpDir)

	sm.Save(NewSaveData("ToDelete", 1, 100))

	err := sm.Delete("ToDelete")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	saves, _ := sm.List()
	for _, save := range saves {
		if save.Name == "ToDelete" {
			t.Error("Deleted save should not appear in list")
		}
	}
}

func TestSaveManagerMaxSlots(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSaveManager(tmpDir)

	for i := 1; i <= 10; i++ {
		err := sm.Save(NewSaveData(fmt.Sprintf("Save%d", i), i, i*100))
		if err != nil {
			t.Fatalf("Save %d failed: %v", i, err)
		}
	}

	err := sm.Save(NewSaveData("Save11", 11, 1100))
	if err == nil {
		t.Error("11th save should fail due to max slots")
	}
}

func TestSaveManagerOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSaveManager(tmpDir)

	sm.Save(NewSaveData("MySave", 1, 100))
	sm.Save(NewSaveData("MySave", 5, 500))

	saves, _ := sm.List()
	if len(saves) != 1 {
		t.Errorf("Save count = %d, want 1 (overwrite)", len(saves))
	}

	loaded, _ := sm.Load("MySave")
	if loaded.HuntNumber != 5 {
		t.Errorf("Loaded hunt = %d, want 5 (overwritten value)", loaded.HuntNumber)
	}
}

func TestSaveManagerLoadNonexistent(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSaveManager(tmpDir)

	_, err := sm.Load("DoesNotExist")
	if err == nil {
		t.Error("Loading nonexistent save should return error")
	}
}

func TestSaveManagerSlotCount(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSaveManager(tmpDir)

	if sm.SlotCount() != 0 {
		t.Errorf("Initial slot count = %d, want 0", sm.SlotCount())
	}

	sm.Save(NewSaveData("Save1", 1, 100))
	sm.Save(NewSaveData("Save2", 2, 200))

	if sm.SlotCount() != 2 {
		t.Errorf("Slot count = %d, want 2", sm.SlotCount())
	}
}

func TestSaveManagerHasRoom(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSaveManager(tmpDir)

	if !sm.HasRoom() {
		t.Error("Empty manager should have room")
	}

	for i := 1; i <= 10; i++ {
		sm.Save(NewSaveData(fmt.Sprintf("Save%d", i), i, i*100))
	}

	if sm.HasRoom() {
		t.Error("Full manager should not have room")
	}
}

func TestSaveDataStashEquipment(t *testing.T) {
	data := NewSaveData("TestSave", 1, 0)

	sword := entity.NewEquipment("Iron Sword", entity.SlotWeapon, 3, 0, 0)
	armor := entity.NewEquipment("Leather Armor", entity.SlotArmor, 0, 2, 10)

	data.StashedEquipment = append(data.StashedEquipment, sword)
	data.StashedEquipment = append(data.StashedEquipment, armor)

	if len(data.StashedEquipment) != 2 {
		t.Errorf("Stashed equipment count = %d, want 2", len(data.StashedEquipment))
	}
}

func TestSaveNameValidation(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"ValidName", true},
		{"My Save 1", true},
		{"", false},
		{"   ", false},
		{"ThisNameIsWayTooLongForASaveSlotAndShouldBeRejected", false},
		{"Valid-Name_123", true},
	}

	for _, tc := range tests {
		result := IsValidSaveName(tc.name)
		if result != tc.valid {
			t.Errorf("IsValidSaveName(%q) = %v, want %v", tc.name, result, tc.valid)
		}
	}
}

func TestSaveManagerExists(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSaveManager(tmpDir)

	if sm.Exists("NonExistent") {
		t.Error("Exists should return false for nonexistent save")
	}

	sm.Save(NewSaveData("MySave", 1, 100))

	if !sm.Exists("MySave") {
		t.Error("Exists should return true for existing save")
	}
}
