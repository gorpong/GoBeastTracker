package entity

// DefaultInventoryCapacity is the starting inventory size
const DefaultInventoryCapacity = 6

// Inventory manages a collection of items with limited capacity
type Inventory struct {
	items    []*Item
	capacity int
}

// NewInventory creates a new inventory with the specified capacity
func NewInventory(capacity int) *Inventory {
	return &Inventory{
		items:    make([]*Item, 0, capacity),
		capacity: capacity,
	}
}

// Capacity returns the maximum number of items the inventory can hold
func (inv *Inventory) Capacity() int {
	return inv.capacity
}

// Count returns the current number of items in the inventory
func (inv *Inventory) Count() int {
	return len(inv.items)
}

// IsFull returns true if the inventory cannot hold more items
func (inv *Inventory) IsFull() bool {
	return len(inv.items) >= inv.capacity
}

// Add attempts to add an item to the inventory.
// Returns true if successful, false if inventory is full.
func (inv *Inventory) Add(item *Item) bool {
	if inv.IsFull() {
		return false
	}
	inv.items = append(inv.items, item)
	return true
}

// GetSlot returns the item at the specified slot (1-indexed), or nil if empty/invalid.
func (inv *Inventory) GetSlot(slot int) *Item {
	index := slot - 1 // Convert to 0-indexed
	if index < 0 || index >= len(inv.items) {
		return nil
	}
	return inv.items[index]
}

// Remove removes and returns the item at the specified slot (1-indexed).
// Returns nil if the slot is empty or invalid.
// Remaining items shift down to fill the gap.
func (inv *Inventory) Remove(slot int) *Item {
	index := slot - 1 // Convert to 0-indexed
	if index < 0 || index >= len(inv.items) {
		return nil
	}

	item := inv.items[index]
	// Shift remaining items down
	inv.items = append(inv.items[:index], inv.items[index+1:]...)
	return item
}

// Items returns a copy of all items in the inventory
func (inv *Inventory) Items() []*Item {
	result := make([]*Item, len(inv.items))
	copy(result, inv.items)
	return result
}

// SetCapacity changes the inventory capacity.
// Existing items are preserved (capacity should not be reduced below current count).
func (inv *Inventory) SetCapacity(newCapacity int) {
	inv.capacity = newCapacity
}
