package dungeon

import (
	"testing"
)

// TestNewRoom verifies room creation with correct bounds
func TestNewRoom(t *testing.T) {
	room := NewRoom(5, 10, 8, 6)

	if room.X != 5 {
		t.Errorf("Room X = %d, want 5", room.X)
	}
	if room.Y != 10 {
		t.Errorf("Room Y = %d, want 10", room.Y)
	}
	if room.Width != 8 {
		t.Errorf("Room Width = %d, want 8", room.Width)
	}
	if room.Height != 6 {
		t.Errorf("Room Height = %d, want 6", room.Height)
	}
}

// TestRoomCenter verifies center calculation
func TestRoomCenter(t *testing.T) {
	tests := []struct {
		name   string
		x, y   int
		w, h   int
		wantCX int
		wantCY int
	}{
		{"Even dimensions", 0, 0, 10, 10, 5, 5},
		{"Odd dimensions", 0, 0, 7, 5, 3, 2},
		{"Offset room", 10, 20, 8, 6, 14, 23},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := NewRoom(tt.x, tt.y, tt.w, tt.h)
			cx, cy := room.Center()
			if cx != tt.wantCX || cy != tt.wantCY {
				t.Errorf("Center() = (%d, %d), want (%d, %d)", cx, cy, tt.wantCX, tt.wantCY)
			}
		})
	}
}

// TestRoomContains verifies point-in-room detection
func TestRoomContains(t *testing.T) {
	room := NewRoom(10, 10, 5, 5) // Room from (10,10) to (14,14)

	tests := []struct {
		name string
		x, y int
		want bool
	}{
		{"Inside center", 12, 12, true},
		{"Top-left corner", 10, 10, true},
		{"Bottom-right corner", 14, 14, true},
		{"Just outside left", 9, 12, false},
		{"Just outside right", 15, 12, false},
		{"Just outside top", 12, 9, false},
		{"Just outside bottom", 12, 15, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := room.Contains(tt.x, tt.y); got != tt.want {
				t.Errorf("Contains(%d, %d) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

// TestRoomIntersects verifies room overlap detection
func TestRoomIntersects(t *testing.T) {
	room1 := NewRoom(10, 10, 5, 5) // (10,10) to (14,14)

	tests := []struct {
		name string
		x, y int
		w, h int
		want bool
	}{
		{"Overlapping", 12, 12, 5, 5, true},
		{"Adjacent right (no overlap)", 15, 10, 5, 5, false},
		{"Adjacent bottom (no overlap)", 10, 15, 5, 5, false},
		{"Far away", 50, 50, 5, 5, false},
		{"Completely inside", 11, 11, 2, 2, true},
		{"Completely surrounds", 5, 5, 20, 20, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room2 := NewRoom(tt.x, tt.y, tt.w, tt.h)
			if got := room1.Intersects(room2); got != tt.want {
				t.Errorf("Intersects() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRoomIntersectsWithPadding verifies overlap detection with spacing
func TestRoomIntersectsWithPadding(t *testing.T) {
	room1 := NewRoom(10, 10, 5, 5) // (10,10) to (14,14)

	tests := []struct {
		name    string
		x, y    int
		w, h    int
		padding int
		want    bool
	}{
		{"Adjacent with padding 1", 15, 10, 5, 5, 1, true},  // Would touch with padding
		{"Adjacent with padding 0", 15, 10, 5, 5, 0, false}, // No overlap without padding
		{"Close with padding 2", 17, 10, 5, 5, 2, true},     // Overlap with padding (room1 extends to 16, room2 starts at 15)
		{"Far with padding 2", 20, 10, 5, 5, 2, false},      // Too far even with padding
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room2 := NewRoom(tt.x, tt.y, tt.w, tt.h)
			if got := room1.IntersectsWithPadding(room2, tt.padding); got != tt.want {
				t.Errorf("IntersectsWithPadding(%d) = %v, want %v", tt.padding, got, tt.want)
			}
		})
	}
}
