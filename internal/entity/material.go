package entity

// MaterialType represents different types of crafting materials
type MaterialType int

const (
	// Common materials (dropped by regular monsters)
	MaterialScales MaterialType = iota
	MaterialClaws
	MaterialFangs
	MaterialHide

	// Rare materials (dropped by boss monsters)
	MaterialWyvernScale
	MaterialOgreHide
	MaterialTrollClaw
	MaterialCyclopsEye
	MaterialMinotaurHorn
)

// String returns the display name of the material type
func (mt MaterialType) String() string {
	switch mt {
	case MaterialScales:
		return "Scales"
	case MaterialClaws:
		return "Claws"
	case MaterialFangs:
		return "Fangs"
	case MaterialHide:
		return "Hide"
	case MaterialWyvernScale:
		return "Wyvern Scale"
	case MaterialOgreHide:
		return "Ogre Hide"
	case MaterialTrollClaw:
		return "Troll Claw"
	case MaterialCyclopsEye:
		return "Cyclops Eye"
	case MaterialMinotaurHorn:
		return "Minotaur Horn"
	default:
		return "Unknown"
	}
}

// Glyph returns the display character for this material type
// Common materials use lowercase, rare materials use uppercase (boss initial)
func (mt MaterialType) Glyph() rune {
	switch mt {
	case MaterialScales:
		return 's'
	case MaterialClaws:
		return 'c'
	case MaterialFangs:
		return 'f'
	case MaterialHide:
		return 'h'
	case MaterialWyvernScale:
		return 'W'
	case MaterialOgreHide:
		return 'O'
	case MaterialTrollClaw:
		return 'T'
	case MaterialCyclopsEye:
		return 'C'
	case MaterialMinotaurHorn:
		return 'M'
	default:
		return '?'
	}
}

// IsRare returns true if this is a rare (boss-dropped) material
func (mt MaterialType) IsRare() bool {
	switch mt {
	case MaterialWyvernScale, MaterialOgreHide, MaterialTrollClaw,
		MaterialCyclopsEye, MaterialMinotaurHorn:
		return true
	default:
		return false
	}
}

// GetCommonMaterialTypes returns a slice of all common material types
func GetCommonMaterialTypes() []MaterialType {
	return []MaterialType{
		MaterialScales,
		MaterialClaws,
		MaterialFangs,
		MaterialHide,
	}
}

// GetRareMaterialTypes returns a slice of all rare material types
func GetRareMaterialTypes() []MaterialType {
	return []MaterialType{
		MaterialWyvernScale,
		MaterialOgreHide,
		MaterialTrollClaw,
		MaterialCyclopsEye,
		MaterialMinotaurHorn,
	}
}

// Material represents a crafting material dropped by monsters
type Material struct {
	Type MaterialType
	X    int
	Y    int
}

// NewMaterial creates a new material at the specified position
func NewMaterial(materialType MaterialType, x, y int) *Material {
	return &Material{
		Type: materialType,
		X:    x,
		Y:    y,
	}
}

// Position returns the material's current coordinates
func (m *Material) Position() (int, int) {
	return m.X, m.Y
}

// Glyph returns the display character for this material
func (m *Material) Glyph() rune {
	return m.Type.Glyph()
}

// Name returns the display name of this material
func (m *Material) Name() string {
	return m.Type.String()
}

// MaterialPouch stores collected materials (unlimited capacity)
type MaterialPouch struct {
	materials map[MaterialType]int
}

// NewMaterialPouch creates a new empty material pouch
func NewMaterialPouch() *MaterialPouch {
	return &MaterialPouch{
		materials: make(map[MaterialType]int),
	}
}

// Add adds the specified quantity of a material to the pouch
func (mp *MaterialPouch) Add(materialType MaterialType, quantity int) {
	mp.materials[materialType] += quantity
}

// Remove attempts to remove the specified quantity of a material.
// Returns true if successful, false if insufficient quantity.
// Does not modify pouch if removal would fail.
func (mp *MaterialPouch) Remove(materialType MaterialType, quantity int) bool {
	current := mp.materials[materialType]
	if current < quantity {
		return false
	}

	mp.materials[materialType] = current - quantity

	// Clean up zero entries
	if mp.materials[materialType] == 0 {
		delete(mp.materials, materialType)
	}

	return true
}

// Count returns the quantity of a specific material type in the pouch
func (mp *MaterialPouch) Count(materialType MaterialType) int {
	return mp.materials[materialType]
}

// Has returns true if the pouch contains at least the specified quantity
func (mp *MaterialPouch) Has(materialType MaterialType, quantity int) bool {
	return mp.materials[materialType] >= quantity
}

// TotalCount returns the total number of materials in the pouch
func (mp *MaterialPouch) TotalCount() int {
	total := 0
	for _, count := range mp.materials {
		total += count
	}
	return total
}

// AllMaterials returns a slice of all material types currently in the pouch
func (mp *MaterialPouch) AllMaterials() []MaterialType {
	result := make([]MaterialType, 0, len(mp.materials))
	for materialType := range mp.materials {
		result = append(result, materialType)
	}
	return result
}

// Clear removes all materials from the pouch
func (mp *MaterialPouch) Clear() {
	mp.materials = make(map[MaterialType]int)
}
