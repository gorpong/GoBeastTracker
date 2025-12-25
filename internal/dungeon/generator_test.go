package dungeon

import (
	"fmt"
	"testing"
)

// TestNewDungeon verifies dungeon creation with correct dimensions
func TestNewDungeon(t *testing.T) {
	d := NewDungeon(100, 40)

	if d.Width != 100 {
		t.Errorf("Dungeon Width = %d, want 100", d.Width)
	}
	if d.Height != 40 {
		t.Errorf("Dungeon Height = %d, want 40", d.Height)
	}
	if d.Tiles == nil {
		t.Fatal("Dungeon Tiles is nil")
	}
	if len(d.Tiles) != 100 {
		t.Errorf("Tiles width = %d, want 100", len(d.Tiles))
	}
	if len(d.Tiles[0]) != 40 {
		t.Errorf("Tiles height = %d, want 40", len(d.Tiles[0]))
	}
}

// TestDungeonInitializedAsWalls verifies all tiles start as walls
func TestDungeonInitializedAsWalls(t *testing.T) {
	d := NewDungeon(20, 20)

	for x := 0; x < d.Width; x++ {
		for y := 0; y < d.Height; y++ {
			if d.Tiles[x][y].Type != TileWall {
				t.Errorf("Tile at (%d,%d) = %v, want TileWall", x, y, d.Tiles[x][y].Type)
			}
		}
	}
}

// TestDungeonGetTile verifies tile access
func TestDungeonGetTile(t *testing.T) {
	d := NewDungeon(20, 20)

	// Valid position
	tile := d.GetTile(10, 10)
	if tile == nil {
		t.Error("GetTile(10, 10) returned nil for valid position")
	}

	// Out of bounds
	if d.GetTile(-1, 0) != nil {
		t.Error("GetTile(-1, 0) should return nil")
	}
	if d.GetTile(0, -1) != nil {
		t.Error("GetTile(0, -1) should return nil")
	}
	if d.GetTile(20, 0) != nil {
		t.Error("GetTile(20, 0) should return nil")
	}
	if d.GetTile(0, 20) != nil {
		t.Error("GetTile(0, 20) should return nil")
	}
}

// TestDungeonIsWalkable verifies walkability check
func TestDungeonIsWalkable(t *testing.T) {
	d := NewDungeon(20, 20)

	// All walls initially - not walkable
	if d.IsWalkable(10, 10) {
		t.Error("Wall tile should not be walkable")
	}

	// Carve a floor
	d.Tiles[10][10].Type = TileFloor
	if !d.IsWalkable(10, 10) {
		t.Error("Floor tile should be walkable")
	}

	// Out of bounds - not walkable
	if d.IsWalkable(-1, 0) {
		t.Error("Out of bounds should not be walkable")
	}
}

// TestGenerateDungeonDeterministic verifies same seed produces same dungeon
func TestGenerateDungeonDeterministic(t *testing.T) {
	seed := int64(12345)

	d1 := GenerateDungeon(100, 40, seed)
	d2 := GenerateDungeon(100, 40, seed)

	// Compare all tiles
	for x := 0; x < d1.Width; x++ {
		for y := 0; y < d1.Height; y++ {
			if d1.Tiles[x][y].Type != d2.Tiles[x][y].Type {
				t.Fatalf("Dungeons differ at (%d,%d): %v vs %v",
					x, y, d1.Tiles[x][y].Type, d2.Tiles[x][y].Type)
			}
		}
	}
}

// TestGenerateDungeonHasRooms verifies dungeon contains rooms
func TestGenerateDungeonHasRooms(t *testing.T) {
	d := GenerateDungeon(100, 40, 12345)

	if len(d.Rooms) == 0 {
		t.Error("Generated dungeon should have at least one room")
	}

	// Should have multiple rooms
	if len(d.Rooms) < 3 {
		t.Errorf("Expected at least 3 rooms, got %d", len(d.Rooms))
	}
}

// TestGenerateDungeonHasFloors verifies dungeon has walkable floor tiles
func TestGenerateDungeonHasFloors(t *testing.T) {
	d := GenerateDungeon(100, 40, 12345)

	floorCount := 0
	for x := 0; x < d.Width; x++ {
		for y := 0; y < d.Height; y++ {
			if d.Tiles[x][y].Type == TileFloor {
				floorCount++
			}
		}
	}

	if floorCount == 0 {
		t.Error("Generated dungeon should have floor tiles")
	}

	// Sanity check: should have significant floor area
	totalTiles := d.Width * d.Height
	floorPercent := float64(floorCount) / float64(totalTiles) * 100
	if floorPercent < 10 {
		t.Errorf("Floor coverage too low: %.1f%% (expected > 10%%)", floorPercent)
	}
}

// TestGenerateDungeonRoomsConnected verifies all rooms are reachable
func TestGenerateDungeonRoomsConnected(t *testing.T) {
	d := GenerateDungeon(100, 40, 12345)

	if len(d.Rooms) < 2 {
		t.Skip("Need at least 2 rooms to test connectivity")
	}

	// Simple flood fill from first room center
	cx, cy := d.Rooms[0].Center()
	visited := make(map[string]bool)
	d.floodFill(cx, cy, visited)

	// Check all room centers are reachable
	for i, room := range d.Rooms {
		rcx, rcy := room.Center()
		key := coordKey(rcx, rcy)
		if !visited[key] {
			t.Errorf("Room %d at (%d,%d) is not reachable from room 0", i, rcx, rcy)
		}
	}
}

// TestDungeonInBounds verifies bounds checking
func TestDungeonInBounds(t *testing.T) {
	d := NewDungeon(100, 40)

	tests := []struct {
		name string
		x, y int
		want bool
	}{
		{"Inside", 50, 20, true},
		{"Top-left", 0, 0, true},
		{"Bottom-right", 99, 39, true},
		{"Negative X", -1, 20, false},
		{"Negative Y", 50, -1, false},
		{"Beyond Width", 100, 20, false},
		{"Beyond Height", 50, 40, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := d.InBounds(tt.x, tt.y); got != tt.want {
				t.Errorf("InBounds(%d, %d) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

// Helper for creating map keys
func coordKey(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}
