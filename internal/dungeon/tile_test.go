package dungeon

import (
	"testing"
)

// TestTileTypeProperties verifies tile types have correct walkability and transparency
func TestTileTypeProperties(t *testing.T) {
	tests := []struct {
		name        string
		tileType    TileType
		walkable    bool
		transparent bool
		glyph       rune
	}{
		{"Floor", TileFloor, true, true, '.'},
		{"Wall", TileWall, false, false, '#'},
		{"Door", TileDoor, true, true, '+'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.tileType.Walkable() != tt.walkable {
				t.Errorf("%s.Walkable() = %v, want %v", tt.name, tt.tileType.Walkable(), tt.walkable)
			}
			if tt.tileType.Transparent() != tt.walkable {
				t.Errorf("%s.Transparent() = %v, want %v", tt.name, tt.tileType.Transparent(), tt.transparent)
			}
			if tt.tileType.Glyph() != tt.glyph {
				t.Errorf("%s.Glyph() = %q, want %q", tt.name, tt.tileType.Glyph(), tt.glyph)
			}
		})
	}
}

// TestNewTile verifies tile creation with default state
func TestNewTile(t *testing.T) {
	tile := NewTile(TileFloor)

	if tile.Type != TileFloor {
		t.Errorf("NewTile type = %v, want TileFloor", tile.Type)
	}
	if tile.Explored {
		t.Error("New tile should not be explored")
	}
	if tile.Visible {
		t.Error("New tile should not be visible")
	}
}

// TestTileGlyph verifies tile returns correct glyph for its type
func TestTileGlyph(t *testing.T) {
	floor := NewTile(TileFloor)
	wall := NewTile(TileWall)
	door := NewTile(TileDoor)

	if floor.Glyph() != '.' {
		t.Errorf("Floor glyph = %q, want '.'", floor.Glyph())
	}
	if wall.Glyph() != '#' {
		t.Errorf("Wall glyph = %q, want '#'", wall.Glyph())
	}
	if door.Glyph() != '+' {
		t.Errorf("Door glyph = %q, want '+'", door.Glyph())
	}
}

// TestTileWalkable verifies walkability check
func TestTileWalkable(t *testing.T) {
	floor := NewTile(TileFloor)
	wall := NewTile(TileWall)
	door := NewTile(TileDoor)

	if !floor.Walkable() {
		t.Error("Floor should be walkable")
	}
	if wall.Walkable() {
		t.Error("Wall should not be walkable")
	}
	if !door.Walkable() {
		t.Error("Door should be walkable")
	}
}

// TestTileTransparent verifies transparency for FOV
func TestTileTransparent(t *testing.T) {
	floor := NewTile(TileFloor)
	wall := NewTile(TileWall)
	door := NewTile(TileDoor)

	if !floor.Transparent() {
		t.Error("Floor should be transparent")
	}
	if wall.Transparent() {
		t.Error("Wall should not be transparent")
	}
	if !door.Transparent() {
		t.Error("Door should be transparent")
	}
}
