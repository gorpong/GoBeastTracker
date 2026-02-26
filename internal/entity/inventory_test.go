package entity

import (
	"testing"
)

// TestNewInventory verifies inventory creation with correct capacity
func TestNewInventory(t *testing.T) {
	inventory := NewInventory(6)

	if inventory.Capacity() != 6 {
		t.Errorf("NewInventory(6).Capacity() = %d, want 6", inventory.Capacity())
	}
	if inventory.Count() != 0 {
		t.Errorf("NewInventory(6).Count() = %d, want 0", inventory.Count())
	}
	if inventory.IsFull() {
		t.Error("New inventory should not be full")
	}
}

// TestInventoryAddItem verifies adding items to inventory
func TestInventoryAddItem(t *testing.T) {
	inventory := NewInventory(6)
	herbItem := NewItem(ItemHerbs, 0, 0)

	added := inventory.Add(herbItem)

	if !added {
		t.Error("Add() should return true when inventory has space")
	}
	if inventory.Count() != 1 {
		t.Errorf("After Add: Count() = %d, want 1", inventory.Count())
	}
}

// TestInventoryAddMultipleItems verifies adding multiple items
func TestInventoryAddMultipleItems(t *testing.T) {
	inventory := NewInventory(6)

	for i := 0; i < 4; i++ {
		item := NewItem(ItemHerbs, 0, 0)
		inventory.Add(item)
	}

	if inventory.Count() != 4 {
		t.Errorf("After adding 4 items: Count() = %d, want 4", inventory.Count())
	}
}

// TestInventoryFull verifies inventory respects capacity
func TestInventoryFull(t *testing.T) {
	inventory := NewInventory(3)

	// Fill inventory
	for i := 0; i < 3; i++ {
		item := NewItem(ItemPotion, 0, 0)
		added := inventory.Add(item)
		if !added {
			t.Errorf("Add() should succeed for item %d", i)
		}
	}

	if !inventory.IsFull() {
		t.Error("Inventory should be full after adding 3 items to capacity 3")
	}

	// Try to add one more
	extraItem := NewItem(ItemHerbs, 0, 0)
	added := inventory.Add(extraItem)

	if added {
		t.Error("Add() should return false when inventory is full")
	}
	if inventory.Count() != 3 {
		t.Errorf("Count should still be 3, got %d", inventory.Count())
	}
}

// TestInventoryGetItem verifies retrieving items by slot
func TestInventoryGetItem(t *testing.T) {
	inventory := NewInventory(6)
	herbItem := NewItem(ItemHerbs, 0, 0)
	potionItem := NewItem(ItemPotion, 0, 0)

	inventory.Add(herbItem)
	inventory.Add(potionItem)

	// Slots are 1-indexed for user display
	item1 := inventory.GetSlot(1)
	item2 := inventory.GetSlot(2)
	item3 := inventory.GetSlot(3)

	if item1 != herbItem {
		t.Error("GetSlot(1) should return first added item")
	}
	if item2 != potionItem {
		t.Error("GetSlot(2) should return second added item")
	}
	if item3 != nil {
		t.Error("GetSlot(3) should return nil for empty slot")
	}
}

// TestInventoryGetSlotOutOfBounds verifies bounds checking
func TestInventoryGetSlotOutOfBounds(t *testing.T) {
	inventory := NewInventory(6)
	inventory.Add(NewItem(ItemHerbs, 0, 0))

	if inventory.GetSlot(0) != nil {
		t.Error("GetSlot(0) should return nil (slots are 1-indexed)")
	}
	if inventory.GetSlot(7) != nil {
		t.Error("GetSlot(7) should return nil (beyond capacity)")
	}
	if inventory.GetSlot(-1) != nil {
		t.Error("GetSlot(-1) should return nil")
	}
}

// TestInventoryRemoveItem verifies removing items by slot
func TestInventoryRemoveItem(t *testing.T) {
	inventory := NewInventory(6)
	herbItem := NewItem(ItemHerbs, 0, 0)
	potionItem := NewItem(ItemPotion, 0, 0)

	inventory.Add(herbItem)
	inventory.Add(potionItem)

	removed := inventory.Remove(1)

	if removed != herbItem {
		t.Error("Remove(1) should return the first item")
	}
	if inventory.Count() != 1 {
		t.Errorf("After Remove: Count() = %d, want 1", inventory.Count())
	}

	// Remaining item should now be in slot 1
	if inventory.GetSlot(1) != potionItem {
		t.Error("After removing slot 1, potion should shift to slot 1")
	}
}

// TestInventoryRemoveEmptySlot verifies removing from empty slot
func TestInventoryRemoveEmptySlot(t *testing.T) {
	inventory := NewInventory(6)
	inventory.Add(NewItem(ItemHerbs, 0, 0))

	removed := inventory.Remove(2)

	if removed != nil {
		t.Error("Remove(2) should return nil for empty slot")
	}
	if inventory.Count() != 1 {
		t.Errorf("Count should still be 1, got %d", inventory.Count())
	}
}

// TestInventoryRemoveOutOfBounds verifies bounds checking on remove
func TestInventoryRemoveOutOfBounds(t *testing.T) {
	inventory := NewInventory(6)
	inventory.Add(NewItem(ItemHerbs, 0, 0))

	if inventory.Remove(0) != nil {
		t.Error("Remove(0) should return nil")
	}
	if inventory.Remove(7) != nil {
		t.Error("Remove(7) should return nil")
	}
	if inventory.Remove(-1) != nil {
		t.Error("Remove(-1) should return nil")
	}
}

// TestInventoryItems verifies getting all items
func TestInventoryItems(t *testing.T) {
	inventory := NewInventory(6)
	herb1 := NewItem(ItemHerbs, 0, 0)
	herb2 := NewItem(ItemHerbs, 0, 0)
	potion := NewItem(ItemPotion, 0, 0)

	inventory.Add(herb1)
	inventory.Add(herb2)
	inventory.Add(potion)

	items := inventory.Items()

	if len(items) != 3 {
		t.Errorf("Items() returned %d items, want 3", len(items))
	}
	if items[0] != herb1 || items[1] != herb2 || items[2] != potion {
		t.Error("Items() returned items in wrong order")
	}
}

// TestInventoryIncreaseCapacity verifies capacity can be increased
func TestInventoryIncreaseCapacity(t *testing.T) {
	inventory := NewInventory(6)

	// Fill it up
	for i := 0; i < 6; i++ {
		inventory.Add(NewItem(ItemHerbs, 0, 0))
	}

	if !inventory.IsFull() {
		t.Error("Inventory should be full")
	}

	// Increase capacity
	inventory.SetCapacity(8)

	if inventory.Capacity() != 8 {
		t.Errorf("After SetCapacity(8): Capacity() = %d, want 8", inventory.Capacity())
	}
	if inventory.IsFull() {
		t.Error("Inventory should no longer be full after capacity increase")
	}

	// Should be able to add more
	added := inventory.Add(NewItem(ItemPotion, 0, 0))
	if !added {
		t.Error("Should be able to add item after capacity increase")
	}
}

// TestInventorySetCapacityPreservesItems verifies items are kept when capacity changes
func TestInventorySetCapacityPreservesItems(t *testing.T) {
	inventory := NewInventory(6)
	herb := NewItem(ItemHerbs, 0, 0)
	potion := NewItem(ItemPotion, 0, 0)

	inventory.Add(herb)
	inventory.Add(potion)

	inventory.SetCapacity(8)

	if inventory.Count() != 2 {
		t.Errorf("After SetCapacity: Count() = %d, want 2", inventory.Count())
	}
	if inventory.GetSlot(1) != herb {
		t.Error("First item should be preserved")
	}
	if inventory.GetSlot(2) != potion {
		t.Error("Second item should be preserved")
	}
}
