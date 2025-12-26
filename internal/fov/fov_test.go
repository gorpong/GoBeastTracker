package fov

import (
	"testing"
)

// mockMap implements the Map interface for testing
type mockMap struct {
	width       int
	height      int
	transparent map[string]bool
}

func newMockMap(width, height int) *mockMap {
	return &mockMap{
		width:       width,
		height:      height,
		transparent: make(map[string]bool),
	}
}

func (m *mockMap) key(x, y int) string {
	return string(rune(x)) + "," + string(rune(y))
}

func (m *mockMap) GetWidth() int  { return m.width }
func (m *mockMap) GetHeight() int { return m.height }

func (m *mockMap) IsTransparent(x, y int) bool {
	if x < 0 || x >= m.width || y < 0 || y >= m.height {
		return false
	}
	key := m.key(x, y)
	transparent, exists := m.transparent[key]
	if !exists {
		return true // Default to transparent (floor)
	}
	return transparent
}

func (m *mockMap) SetTransparent(x, y int, transparent bool) {
	m.transparent[m.key(x, y)] = transparent
}

// TestNewFOVMap verifies FOV map creation
func TestNewFOVMap(t *testing.T) {
	fovMap := NewFOVMap(20, 15)

	if fovMap.Width != 20 {
		t.Errorf("FOVMap Width = %d, want 20", fovMap.Width)
	}
	if fovMap.Height != 15 {
		t.Errorf("FOVMap Height = %d, want 15", fovMap.Height)
	}
	if fovMap.Visible == nil {
		t.Error("Visible map should not be nil")
	}
	if fovMap.Explored == nil {
		t.Error("Explored map should not be nil")
	}
}

// TestFOVOriginVisible verifies that the player's position is always visible
func TestFOVOriginVisible(t *testing.T) {
	gameMap := newMockMap(20, 20)
	fovMap := NewFOVMap(20, 20)

	Compute(fovMap, gameMap, 10, 10, 5)

	if !fovMap.IsVisible(10, 10) {
		t.Error("Origin (player position) should always be visible")
	}
}

// TestFOVOriginExplored verifies that visible tiles become explored
func TestFOVOriginExplored(t *testing.T) {
	gameMap := newMockMap(20, 20)
	fovMap := NewFOVMap(20, 20)

	Compute(fovMap, gameMap, 10, 10, 5)

	if !fovMap.IsExplored(10, 10) {
		t.Error("Origin should be marked as explored after FOV computation")
	}
}

// TestFOVAdjacentTilesVisible verifies adjacent tiles are visible in open space
func TestFOVAdjacentTilesVisible(t *testing.T) {
	gameMap := newMockMap(20, 20)
	fovMap := NewFOVMap(20, 20)

	Compute(fovMap, gameMap, 10, 10, 5)

	// All adjacent tiles should be visible
	adjacents := [][2]int{
		{10, 9},  // Up
		{10, 11}, // Down
		{9, 10},  // Left
		{11, 10}, // Right
	}

	for _, pos := range adjacents {
		if !fovMap.IsVisible(pos[0], pos[1]) {
			t.Errorf("Adjacent tile (%d, %d) should be visible", pos[0], pos[1])
		}
	}
}

// TestFOVRadiusLimit verifies tiles beyond radius are not visible
func TestFOVRadiusLimit(t *testing.T) {
	gameMap := newMockMap(30, 30)
	fovMap := NewFOVMap(30, 30)
	radius := 5

	Compute(fovMap, gameMap, 15, 15, radius)

	// Tile well beyond radius should not be visible
	if fovMap.IsVisible(15, 15+radius+2) {
		t.Errorf("Tile beyond radius should not be visible")
	}
	if fovMap.IsVisible(15+radius+2, 15) {
		t.Errorf("Tile beyond radius should not be visible")
	}
}

// TestFOVWallBlocks verifies walls block line of sight
func TestFOVWallBlocks(t *testing.T) {
	gameMap := newMockMap(20, 20)
	fovMap := NewFOVMap(20, 20)

	// Place a wall directly to the right of player at (11, 10)
	gameMap.SetTransparent(11, 10, false)

	Compute(fovMap, gameMap, 10, 10, 8)

	// The wall itself should be visible
	if !fovMap.IsVisible(11, 10) {
		t.Error("Wall should be visible")
	}

	// Tiles directly behind the wall should NOT be visible
	if fovMap.IsVisible(12, 10) {
		t.Error("Tile behind wall should not be visible")
	}
}

// TestFOVClearsBetweenComputes verifies visible is reset between computations
func TestFOVClearsBetweenComputes(t *testing.T) {
	gameMap := newMockMap(20, 20)
	fovMap := NewFOVMap(20, 20)

	// First compute at (5, 5)
	Compute(fovMap, gameMap, 5, 5, 3)
	if !fovMap.IsVisible(5, 5) {
		t.Error("First origin should be visible")
	}

	// Second compute at (15, 15) - should clear previous visibility
	Compute(fovMap, gameMap, 15, 15, 3)

	// Old position should no longer be visible (but still explored)
	if fovMap.IsVisible(5, 5) {
		t.Error("Previous position should not be visible after new compute")
	}
	if !fovMap.IsExplored(5, 5) {
		t.Error("Previous position should still be explored")
	}
}

// TestFOVExploredPersists verifies explored tiles remain explored
func TestFOVExploredPersists(t *testing.T) {
	gameMap := newMockMap(20, 20)
	fovMap := NewFOVMap(20, 20)

	// First compute - explore some area
	Compute(fovMap, gameMap, 5, 5, 3)

	// Second compute at different location
	Compute(fovMap, gameMap, 15, 15, 3)

	// Previous area should still be explored
	if !fovMap.IsExplored(5, 5) {
		t.Error("Previously visible tiles should remain explored")
	}
}

// TestFOVOutOfBounds verifies out of bounds queries return false
func TestFOVOutOfBounds(t *testing.T) {
	fovMap := NewFOVMap(20, 20)

	if fovMap.IsVisible(-1, 10) {
		t.Error("Out of bounds should not be visible")
	}
	if fovMap.IsVisible(10, -1) {
		t.Error("Out of bounds should not be visible")
	}
	if fovMap.IsVisible(20, 10) {
		t.Error("Out of bounds should not be visible")
	}
	if fovMap.IsVisible(10, 20) {
		t.Error("Out of bounds should not be visible")
	}

	if fovMap.IsExplored(-1, 10) {
		t.Error("Out of bounds should not be explored")
	}
}

// TestFOVSymmetry verifies FOV is roughly symmetric in open space
func TestFOVSymmetry(t *testing.T) {
	gameMap := newMockMap(30, 30)
	fovMap := NewFOVMap(30, 30)

	Compute(fovMap, gameMap, 15, 15, 5)

	// In open space, the FOV should be symmetric
	// Check cardinal directions at same distance
	distance := 3
	if fovMap.IsVisible(15+distance, 15) != fovMap.IsVisible(15-distance, 15) {
		t.Error("FOV should be symmetric on X axis in open space")
	}
	if fovMap.IsVisible(15, 15+distance) != fovMap.IsVisible(15, 15-distance) {
		t.Error("FOV should be symmetric on Y axis in open space")
	}
}
