package entity

// ItemType represents different types of items
type ItemType int

const (
	ItemHerbs  ItemType = iota // Basic healing item
	ItemPotion                 // Strong healing item
)

// Healing values for each item type
const (
	HerbsHealing  = 25
	PotionHealing = 60
)

// String returns the display name of the item type
func (it ItemType) String() string {
	switch it {
	case ItemHerbs:
		return "Herbs"
	case ItemPotion:
		return "Potion"
	default:
		return "Unknown"
	}
}

// Glyph returns the display character for this item type
func (it ItemType) Glyph() rune {
	switch it {
	case ItemHerbs:
		return '"'
	case ItemPotion:
		return '!'
	default:
		return '?'
	}
}

// HealingValue returns the HP restored by this item type
func (it ItemType) HealingValue() int {
	switch it {
	case ItemHerbs:
		return HerbsHealing
	case ItemPotion:
		return PotionHealing
	default:
		return 0
	}
}

// Item represents a pickup item in the dungeon
type Item struct {
	Type ItemType
	X    int
	Y    int
}

// NewItem creates a new item at the specified position
func NewItem(itemType ItemType, x, y int) *Item {
	return &Item{
		Type: itemType,
		X:    x,
		Y:    y,
	}
}

// Position returns the item's current coordinates
func (i *Item) Position() (int, int) {
	return i.X, i.Y
}

// Glyph returns the display character for this item
func (i *Item) Glyph() rune {
	return i.Type.Glyph()
}

// HealingValue returns the HP restored by using this item
func (i *Item) HealingValue() int {
	return i.Type.HealingValue()
}

// Name returns the display name of this item
func (i *Item) Name() string {
	return i.Type.String()
}
