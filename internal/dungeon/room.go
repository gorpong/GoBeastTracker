package dungeon

// Room represents a rectangular room in the dungeon
type Room struct {
	X      int // Top-left X coordinate
	Y      int // Top-left Y coordinate
	Width  int
	Height int
}

// NewRoom creates a new room with the specified position and dimensions
func NewRoom(x, y, width, height int) *Room {
	return &Room{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// Center returns the center coordinates of the room
func (r *Room) Center() (int, int) {
	cx := r.X + r.Width/2
	cy := r.Y + r.Height/2
	return cx, cy
}

// Contains returns true if the point (x, y) is inside the room
func (r *Room) Contains(x, y int) bool {
	return x >= r.X && x < r.X+r.Width &&
		y >= r.Y && y < r.Y+r.Height
}

// Intersects returns true if this room overlaps with another room
func (r *Room) Intersects(other *Room) bool {
	return r.X < other.X+other.Width &&
		r.X+r.Width > other.X &&
		r.Y < other.Y+other.Height &&
		r.Y+r.Height > other.Y
}

// IntersectsWithPadding returns true if this room overlaps with another room
// when both rooms are expanded by the padding amount
func (r *Room) IntersectsWithPadding(other *Room, padding int) bool {
	return r.X-padding < other.X+other.Width+padding &&
		r.X+r.Width+padding > other.X-padding &&
		r.Y-padding < other.Y+other.Height+padding &&
		r.Y+r.Height+padding > other.Y-padding
}
