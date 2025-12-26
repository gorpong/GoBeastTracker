package dungeon

import (
	"fmt"
	"math/rand"
)

const (
	minRoomSize = 5
	maxRoomSize = 12
	maxRooms    = 15
	roomPadding = 1 // Minimum space between rooms
)

// Dungeon represents the game map
type Dungeon struct {
	Width  int
	Height int
	Tiles  [][]*Tile
	Rooms  []*Room
}

// NewDungeon creates an empty dungeon filled with walls
func NewDungeon(width, height int) *Dungeon {
	tiles := make([][]*Tile, width)
	for x := 0; x < width; x++ {
		tiles[x] = make([]*Tile, height)
		for y := 0; y < height; y++ {
			tiles[x][y] = NewTile(TileWall)
		}
	}

	return &Dungeon{
		Width:  width,
		Height: height,
		Tiles:  tiles,
		Rooms:  make([]*Room, 0),
	}
}

// GetTile returns the tile at (x, y), or nil if out of bounds
func (d *Dungeon) GetTile(x, y int) *Tile {
	if !d.InBounds(x, y) {
		return nil
	}
	return d.Tiles[x][y]
}

// InBounds returns true if (x, y) is within the dungeon
func (d *Dungeon) InBounds(x, y int) bool {
	return x >= 0 && x < d.Width && y >= 0 && y < d.Height
}

// IsWalkable returns true if the tile at (x, y) can be walked on
func (d *Dungeon) IsWalkable(x, y int) bool {
	tile := d.GetTile(x, y)
	if tile == nil {
		return false
	}
	return tile.Walkable()
}

// GenerateDungeon creates a new dungeon with rooms and corridors
func GenerateDungeon(width, height int, seed int64) *Dungeon {
	rng := rand.New(rand.NewSource(seed))
	generatedDungeon := NewDungeon(width, height)

	// Generate rooms
	for i := 0; i < maxRooms; i++ {
		roomWidth := rng.Intn(maxRoomSize-minRoomSize+1) + minRoomSize
		roomHeight := rng.Intn(maxRoomSize-minRoomSize+1) + minRoomSize

		// Leave border around the dungeon
		x := rng.Intn(width-roomWidth-2) + 1
		y := rng.Intn(height-roomHeight-2) + 1

		newRoom := NewRoom(x, y, roomWidth, roomHeight)

		// Check for overlap with existing rooms
		overlaps := false
		for _, other := range generatedDungeon.Rooms {
			if newRoom.IntersectsWithPadding(other, roomPadding) {
				overlaps = true
				break
			}
		}

		if !overlaps {
			generatedDungeon.carveRoom(newRoom)

			// Connect to previous room with corridor
			if len(generatedDungeon.Rooms) > 0 {
				prevRoom := generatedDungeon.Rooms[len(generatedDungeon.Rooms)-1]
				generatedDungeon.carveCorridor(prevRoom, newRoom, rng)
			}

			generatedDungeon.Rooms = append(generatedDungeon.Rooms, newRoom)
		}
	}

	return generatedDungeon
}

// carveRoom carves out a room (fills with floor tiles)
func (d *Dungeon) carveRoom(room *Room) {
	for x := room.X; x < room.X+room.Width; x++ {
		for y := room.Y; y < room.Y+room.Height; y++ {
			if d.InBounds(x, y) {
				d.Tiles[x][y].Type = TileFloor
			}
		}
	}
}

// carveCorridor carves a corridor between two rooms
func (d *Dungeon) carveCorridor(room1, room2 *Room, rng *rand.Rand) {
	x1, y1 := room1.Center()
	x2, y2 := room2.Center()

	// Randomly choose whether to go horizontal-first or vertical-first
	if rng.Intn(2) == 0 {
		d.carveHorizontalTunnel(x1, x2, y1)
		d.carveVerticalTunnel(y1, y2, x2)
	} else {
		d.carveVerticalTunnel(y1, y2, x1)
		d.carveHorizontalTunnel(x1, x2, y2)
	}
}

// carveHorizontalTunnel carves a horizontal line of floor tiles
func (d *Dungeon) carveHorizontalTunnel(x1, x2, y int) {
	minX, maxX := x1, x2
	if x1 > x2 {
		minX, maxX = x2, x1
	}

	for x := minX; x <= maxX; x++ {
		if d.InBounds(x, y) {
			d.Tiles[x][y].Type = TileFloor
		}
	}
}

// carveVerticalTunnel carves a vertical line of floor tiles
func (d *Dungeon) carveVerticalTunnel(y1, y2, x int) {
	minY, maxY := y1, y2
	if y1 > y2 {
		minY, maxY = y2, y1
	}

	for y := minY; y <= maxY; y++ {
		if d.InBounds(x, y) {
			d.Tiles[x][y].Type = TileFloor
		}
	}
}

// floodFill performs a flood fill from (x, y) and marks visited positions
func (d *Dungeon) floodFill(x, y int, visited map[string]bool) {
	key := fmt.Sprintf("%d,%d", x, y)
	if visited[key] {
		return
	}
	if !d.InBounds(x, y) {
		return
	}
	if !d.IsWalkable(x, y) {
		return
	}

	visited[key] = true

	// Recursively fill adjacent tiles
	d.floodFill(x+1, y, visited)
	d.floodFill(x-1, y, visited)
	d.floodFill(x, y+1, visited)
	d.floodFill(x, y-1, visited)
}
