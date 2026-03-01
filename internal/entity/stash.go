package entity

// EquipmentStash stores unequipped gear (unlimited capacity)
type EquipmentStash struct {
	items []*Equipment
}

// NewEquipmentStash creates a new empty equipment stash
func NewEquipmentStash() *EquipmentStash {
	return &EquipmentStash{
		items: make([]*Equipment, 0),
	}
}

// Add adds equipment to the stash
func (es *EquipmentStash) Add(equipment *Equipment) {
	if equipment != nil {
		es.items = append(es.items, equipment)
	}
}

// Remove removes equipment from the stash, returns true if found and removed
func (es *EquipmentStash) Remove(equipment *Equipment) bool {
	for i, item := range es.items {
		if item == equipment {
			es.items = append(es.items[:i], es.items[i+1:]...)
			return true
		}
	}
	return false
}

// Count returns the number of items in the stash
func (es *EquipmentStash) Count() int {
	return len(es.items)
}

// GetAll returns all equipment in the stash
func (es *EquipmentStash) GetAll() []*Equipment {
	result := make([]*Equipment, len(es.items))
	copy(result, es.items)
	return result
}

// GetBySlot returns all equipment for a specific slot type
func (es *EquipmentStash) GetBySlot(slot EquipmentSlot) []*Equipment {
	result := make([]*Equipment, 0)
	for _, item := range es.items {
		if item.Slot == slot {
			result = append(result, item)
		}
	}
	return result
}

// Contains returns true if the equipment is in the stash
func (es *EquipmentStash) Contains(equipment *Equipment) bool {
	for _, item := range es.items {
		if item == equipment {
			return true
		}
	}
	return false
}

// FindByName returns equipment with the given name, or nil if not found
func (es *EquipmentStash) FindByName(name string) *Equipment {
	for _, item := range es.items {
		if item.Name == name {
			return item
		}
	}
	return nil
}

// Clear removes all equipment from the stash
func (es *EquipmentStash) Clear() {
	es.items = make([]*Equipment, 0)
}
