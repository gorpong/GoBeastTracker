package entity

import (
	"testing"
)

// TestItemTypeProperties verifies item types have correct properties
func TestItemTypeProperties(t *testing.T) {
	tests := []struct {
		name        string
		itemType    ItemType
		wantGlyph   rune
		wantHealing int
		wantName    string
	}{
		{
			name:        "Herbs heal 25 HP",
			itemType:    ItemHerbs,
			wantGlyph:   '"',
			wantHealing: 25,
			wantName:    "Herbs",
		},
		{
			name:        "Potion heals 60 HP",
			itemType:    ItemPotion,
			wantGlyph:   '!',
			wantHealing: 60,
			wantName:    "Potion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.itemType.Glyph(); got != tt.wantGlyph {
				t.Errorf("ItemType.Glyph() = %q, want %q", got, tt.wantGlyph)
			}
			if got := tt.itemType.HealingValue(); got != tt.wantHealing {
				t.Errorf("ItemType.HealingValue() = %d, want %d", got, tt.wantHealing)
			}
			if got := tt.itemType.String(); got != tt.wantName {
				t.Errorf("ItemType.String() = %q, want %q", got, tt.wantName)
			}
		})
	}
}

// TestNewItem verifies item creation with correct initial values
func TestNewItem(t *testing.T) {
	item := NewItem(ItemPotion, 10, 15)

	if item.Type != ItemPotion {
		t.Errorf("NewItem Type = %v, want ItemPotion", item.Type)
	}
	if item.X != 10 {
		t.Errorf("NewItem X = %d, want 10", item.X)
	}
	if item.Y != 15 {
		t.Errorf("NewItem Y = %d, want 15", item.Y)
	}
}

// TestItemPosition verifies Position returns correct coordinates
func TestItemPosition(t *testing.T) {
	item := NewItem(ItemHerbs, 5, 8)

	x, y := item.Position()
	if x != 5 || y != 8 {
		t.Errorf("Position() = (%d, %d), want (5, 8)", x, y)
	}
}

// TestItemGlyph verifies item returns correct glyph from its type
func TestItemGlyph(t *testing.T) {
	herbItem := NewItem(ItemHerbs, 0, 0)
	potionItem := NewItem(ItemPotion, 0, 0)

	if herbItem.Glyph() != '"' {
		t.Errorf("Herbs Glyph() = %q, want '\"'", herbItem.Glyph())
	}
	if potionItem.Glyph() != '!' {
		t.Errorf("Potion Glyph() = %q, want '!'", potionItem.Glyph())
	}
}

// TestItemHealingValue verifies item returns correct healing from its type
func TestItemHealingValue(t *testing.T) {
	herbItem := NewItem(ItemHerbs, 0, 0)
	potionItem := NewItem(ItemPotion, 0, 0)

	if herbItem.HealingValue() != 25 {
		t.Errorf("Herbs HealingValue() = %d, want 25", herbItem.HealingValue())
	}
	if potionItem.HealingValue() != 60 {
		t.Errorf("Potion HealingValue() = %d, want 60", potionItem.HealingValue())
	}
}

// TestItemName verifies item returns correct name from its type
func TestItemName(t *testing.T) {
	herbItem := NewItem(ItemHerbs, 0, 0)
	potionItem := NewItem(ItemPotion, 0, 0)

	if herbItem.Name() != "Herbs" {
		t.Errorf("Herbs Name() = %q, want \"Herbs\"", herbItem.Name())
	}
	if potionItem.Name() != "Potion" {
		t.Errorf("Potion Name() = %q, want \"Potion\"", potionItem.Name())
	}
}
