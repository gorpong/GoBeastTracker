package fov

// Map is the interface that the game map must implement for FOV calculation
type Map interface {
	GetWidth() int
	GetHeight() int
	IsTransparent(x, y int) bool
}

// FOVMap tracks which tiles are visible and explored
type FOVMap struct {
	Width    int
	Height   int
	Visible  [][]bool
	Explored [][]bool
}

// NewFOVMap creates a new FOV map with the specified dimensions
func NewFOVMap(width, height int) *FOVMap {
	visible := make([][]bool, width)
	explored := make([][]bool, width)
	for x := 0; x < width; x++ {
		visible[x] = make([]bool, height)
		explored[x] = make([]bool, height)
	}

	return &FOVMap{
		Width:    width,
		Height:   height,
		Visible:  visible,
		Explored: explored,
	}
}

// IsVisible returns true if the tile at (x, y) is currently visible
func (f *FOVMap) IsVisible(x, y int) bool {
	if !f.inBounds(x, y) {
		return false
	}
	return f.Visible[x][y]
}

// IsExplored returns true if the tile at (x, y) has been seen before
func (f *FOVMap) IsExplored(x, y int) bool {
	if !f.inBounds(x, y) {
		return false
	}
	return f.Explored[x][y]
}

func (f *FOVMap) inBounds(x, y int) bool {
	return x >= 0 && x < f.Width && y >= 0 && y < f.Height
}

func (f *FOVMap) setVisible(x, y int, visible bool) {
	if f.inBounds(x, y) {
		f.Visible[x][y] = visible
		if visible {
			f.Explored[x][y] = true
		}
	}
}

func (f *FOVMap) clearVisible() {
	for x := 0; x < f.Width; x++ {
		for y := 0; y < f.Height; y++ {
			f.Visible[x][y] = false
		}
	}
}

// Compute calculates the field of view from (originX, originY) with given radius
// using a recursive shadowcasting algorithm
func Compute(fovMap *FOVMap, gameMap Map, originX, originY, radius int) {
	fovMap.clearVisible()

	// Origin is always visible
	fovMap.setVisible(originX, originY, true)

	// Cast light in all 8 octants
	for octant := 0; octant < 8; octant++ {
		castLight(fovMap, gameMap, originX, originY, radius, 1, 1.0, 0.0, octant)
	}
}

// Multipliers for transforming coordinates into each octant
// Each octant has 4 multipliers: xx, xy, yx, yy
var octantMultipliers = [8][4]int{
	{1, 0, 0, 1},   // Octant 0: E-NE
	{0, 1, 1, 0},   // Octant 1: N-NE
	{0, -1, 1, 0},  // Octant 2: N-NW
	{-1, 0, 0, 1},  // Octant 3: W-NW
	{-1, 0, 0, -1}, // Octant 4: W-SW
	{0, -1, -1, 0}, // Octant 5: S-SW
	{0, 1, -1, 0},  // Octant 6: S-SE
	{1, 0, 0, -1},  // Octant 7: E-SE
}

// castLight recursively casts light in a single octant
func castLight(fovMap *FOVMap, gameMap Map, originX, originY, radius int,
	row int, startSlope, endSlope float64, octant int) {

	if startSlope < endSlope {
		return
	}

	mult := octantMultipliers[octant]
	xx, xy, yx, yy := mult[0], mult[1], mult[2], mult[3]

	radiusSquared := radius * radius
	newStart := startSlope

	for j := row; j <= radius; j++ {
		blocked := false

		for dx := -j; dx <= 0; dx++ {
			dy := -j

			// Transform coordinates based on octant
			mapX := originX + dx*xx + dy*xy
			mapY := originY + dx*yx + dy*yy

			// Calculate slopes for this cell
			leftSlope := (float64(dx) - 0.5) / (float64(dy) + 0.5)
			rightSlope := (float64(dx) + 0.5) / (float64(dy) - 0.5)

			if startSlope < rightSlope {
				continue
			}
			if endSlope > leftSlope {
				break
			}

			// Check if within radius (using squared distance for efficiency)
			distSquared := dx*dx + dy*dy
			if distSquared <= radiusSquared {
				fovMap.setVisible(mapX, mapY, true)
			}

			// Check for blocking
			if blocked {
				if !gameMap.IsTransparent(mapX, mapY) {
					newStart = rightSlope
					continue
				} else {
					blocked = false
					startSlope = newStart
				}
			} else {
				if !gameMap.IsTransparent(mapX, mapY) && j < radius {
					blocked = true
					castLight(fovMap, gameMap, originX, originY, radius,
						j+1, startSlope, leftSlope, octant)
					newStart = rightSlope
				}
			}
		}

		if blocked {
			break
		}
	}
}
