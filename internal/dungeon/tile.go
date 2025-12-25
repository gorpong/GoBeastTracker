package dungeon

// TileType represents different types of terrain
type TileType int

const (
	TileFloor TileType = iota
	TileWall
	TileDoor
)

// Walkable returns true if entities can walk on this tile type
func (t TileType) Walkable() bool {
	switch t {
	case TileFloor, TileDoor:
		return true
	default:
		return false
	}
}

// Transparent returns true if this tile type allows vision through it
func (t TileType) Transparent() bool {
	switch t {
	case TileFloor, TileDoor:
		return true
	default:
		return false
	}
}

// Glyph returns the display character for this tile type
func (t TileType) Glyph() rune {
	switch t {
	case TileFloor:
		return '.'
	case TileWall:
		return '#'
	case TileDoor:
		return '+'
	default:
		return '?'
	}
}

// Tile represents a single cell in the dungeon
type Tile struct {
	Type     TileType
	Explored bool // Has the player seen this tile before?
	Visible  bool // Is the tile currently in the player's FOV?
}

// NewTile creates a new tile of the specified type
func NewTile(t TileType) *Tile {
	return &Tile{
		Type:     t,
		Explored: false,
		Visible:  false,
	}
}

// Glyph returns the display character for this tile
func (t *Tile) Glyph() rune {
	return t.Type.Glyph()
}

// Walkable returns true if entities can walk on this tile
func (t *Tile) Walkable() bool {
	return t.Type.Walkable()
}

// Transparent returns true if this tile allows vision through it
func (t *Tile) Transparent() bool {
	return t.Type.Transparent()
}
