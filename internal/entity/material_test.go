package entity

import (
	"testing"
)

// TestMaterialTypeProperties verifies material types have correct properties
func TestMaterialTypeProperties(t *testing.T) {
	tests := []struct {
		name         string
		materialType MaterialType
		wantName     string
		wantGlyph    rune
		wantRare     bool
	}{
		{
			name:         "Scales are common",
			materialType: MaterialScales,
			wantName:     "Scales",
			wantGlyph:    's',
			wantRare:     false,
		},
		{
			name:         "Claws are common",
			materialType: MaterialClaws,
			wantName:     "Claws",
			wantGlyph:    'c',
			wantRare:     false,
		},
		{
			name:         "Fangs are common",
			materialType: MaterialFangs,
			wantName:     "Fangs",
			wantGlyph:    'f',
			wantRare:     false,
		},
		{
			name:         "Hide is common",
			materialType: MaterialHide,
			wantName:     "Hide",
			wantGlyph:    'h',
			wantRare:     false,
		},
		{
			name:         "Wyvern Scale is rare",
			materialType: MaterialWyvernScale,
			wantName:     "Wyvern Scale",
			wantGlyph:    'W',
			wantRare:     true,
		},
		{
			name:         "Ogre Hide is rare",
			materialType: MaterialOgreHide,
			wantName:     "Ogre Hide",
			wantGlyph:    'O',
			wantRare:     true,
		},
		{
			name:         "Troll Claw is rare",
			materialType: MaterialTrollClaw,
			wantName:     "Troll Claw",
			wantGlyph:    'T',
			wantRare:     true,
		},
		{
			name:         "Cyclops Eye is rare",
			materialType: MaterialCyclopsEye,
			wantName:     "Cyclops Eye",
			wantGlyph:    'C',
			wantRare:     true,
		},
		{
			name:         "Minotaur Horn is rare",
			materialType: MaterialMinotaurHorn,
			wantName:     "Minotaur Horn",
			wantGlyph:    'M',
			wantRare:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.materialType.String(); got != tt.wantName {
				t.Errorf("MaterialType.String() = %q, want %q", got, tt.wantName)
			}
			if got := tt.materialType.Glyph(); got != tt.wantGlyph {
				t.Errorf("MaterialType.Glyph() = %q, want %q", got, tt.wantGlyph)
			}
			if got := tt.materialType.IsRare(); got != tt.wantRare {
				t.Errorf("MaterialType.IsRare() = %v, want %v", got, tt.wantRare)
			}
		})
	}
}

// TestNewMaterial verifies material creation with correct initial values
func TestNewMaterial(t *testing.T) {
	material := NewMaterial(MaterialScales, 10, 15)

	if material.Type != MaterialScales {
		t.Errorf("NewMaterial Type = %v, want MaterialScales", material.Type)
	}
	if material.X != 10 {
		t.Errorf("NewMaterial X = %d, want 10", material.X)
	}
	if material.Y != 15 {
		t.Errorf("NewMaterial Y = %d, want 15", material.Y)
	}
}

// TestMaterialPosition verifies Position returns correct coordinates
func TestMaterialPosition(t *testing.T) {
	material := NewMaterial(MaterialClaws, 5, 8)

	x, y := material.Position()
	if x != 5 || y != 8 {
		t.Errorf("Position() = (%d, %d), want (5, 8)", x, y)
	}
}

// TestMaterialGlyph verifies material returns correct glyph from its type
func TestMaterialGlyph(t *testing.T) {
	scaleMaterial := NewMaterial(MaterialScales, 0, 0)
	wyvernMaterial := NewMaterial(MaterialWyvernScale, 0, 0)

	if scaleMaterial.Glyph() != 's' {
		t.Errorf("Scales Glyph() = %q, want 's'", scaleMaterial.Glyph())
	}
	if wyvernMaterial.Glyph() != 'W' {
		t.Errorf("WyvernScale Glyph() = %q, want 'W'", wyvernMaterial.Glyph())
	}
}

// TestMaterialName verifies material returns correct name from its type
func TestMaterialName(t *testing.T) {
	hideMaterial := NewMaterial(MaterialHide, 0, 0)
	ogreMaterial := NewMaterial(MaterialOgreHide, 0, 0)

	if hideMaterial.Name() != "Hide" {
		t.Errorf("Hide Name() = %q, want \"Hide\"", hideMaterial.Name())
	}
	if ogreMaterial.Name() != "Ogre Hide" {
		t.Errorf("OgreHide Name() = %q, want \"Ogre Hide\"", ogreMaterial.Name())
	}
}

// TestNewMaterialPouch verifies pouch creation
func TestNewMaterialPouch(t *testing.T) {
	pouch := NewMaterialPouch()

	if pouch == nil {
		t.Fatal("NewMaterialPouch() returned nil")
	}
	if pouch.TotalCount() != 0 {
		t.Errorf("New pouch TotalCount() = %d, want 0", pouch.TotalCount())
	}
	if len(pouch.AllMaterials()) != 0 {
		t.Errorf("New pouch AllMaterials() length = %d, want 0", len(pouch.AllMaterials()))
	}
}

// TestMaterialPouchAdd verifies adding materials to pouch
func TestMaterialPouchAdd(t *testing.T) {
	pouch := NewMaterialPouch()

	pouch.Add(MaterialScales, 1)

	if pouch.Count(MaterialScales) != 1 {
		t.Errorf("After Add(Scales, 1): Count(Scales) = %d, want 1", pouch.Count(MaterialScales))
	}
	if pouch.TotalCount() != 1 {
		t.Errorf("After Add(Scales, 1): TotalCount() = %d, want 1", pouch.TotalCount())
	}
}

// TestMaterialPouchAddMultiple verifies adding multiple of same material
func TestMaterialPouchAddMultiple(t *testing.T) {
	pouch := NewMaterialPouch()

	pouch.Add(MaterialClaws, 3)
	pouch.Add(MaterialClaws, 2)

	if pouch.Count(MaterialClaws) != 5 {
		t.Errorf("After adding 3+2 Claws: Count(Claws) = %d, want 5", pouch.Count(MaterialClaws))
	}
}

// TestMaterialPouchAddDifferentTypes verifies adding different material types
func TestMaterialPouchAddDifferentTypes(t *testing.T) {
	pouch := NewMaterialPouch()

	pouch.Add(MaterialScales, 3)
	pouch.Add(MaterialFangs, 2)
	pouch.Add(MaterialHide, 4)

	if pouch.Count(MaterialScales) != 3 {
		t.Errorf("Count(Scales) = %d, want 3", pouch.Count(MaterialScales))
	}
	if pouch.Count(MaterialFangs) != 2 {
		t.Errorf("Count(Fangs) = %d, want 2", pouch.Count(MaterialFangs))
	}
	if pouch.Count(MaterialHide) != 4 {
		t.Errorf("Count(Hide) = %d, want 4", pouch.Count(MaterialHide))
	}
	if pouch.TotalCount() != 9 {
		t.Errorf("TotalCount() = %d, want 9", pouch.TotalCount())
	}
}

// TestMaterialPouchRemove verifies removing materials from pouch
func TestMaterialPouchRemove(t *testing.T) {
	pouch := NewMaterialPouch()
	pouch.Add(MaterialScales, 5)

	removed := pouch.Remove(MaterialScales, 3)

	if !removed {
		t.Error("Remove() should return true when sufficient materials exist")
	}
	if pouch.Count(MaterialScales) != 2 {
		t.Errorf("After Remove(Scales, 3): Count(Scales) = %d, want 2", pouch.Count(MaterialScales))
	}
}

// TestMaterialPouchRemoveAll verifies removing all of a material type
func TestMaterialPouchRemoveAll(t *testing.T) {
	pouch := NewMaterialPouch()
	pouch.Add(MaterialFangs, 3)

	removed := pouch.Remove(MaterialFangs, 3)

	if !removed {
		t.Error("Remove() should return true when removing exact amount")
	}
	if pouch.Count(MaterialFangs) != 0 {
		t.Errorf("After removing all: Count(Fangs) = %d, want 0", pouch.Count(MaterialFangs))
	}
}

// TestMaterialPouchRemoveInsufficient verifies removal fails with insufficient materials
func TestMaterialPouchRemoveInsufficient(t *testing.T) {
	pouch := NewMaterialPouch()
	pouch.Add(MaterialHide, 2)

	removed := pouch.Remove(MaterialHide, 5)

	if removed {
		t.Error("Remove() should return false when insufficient materials")
	}
	// Materials should be unchanged
	if pouch.Count(MaterialHide) != 2 {
		t.Errorf("After failed remove: Count(Hide) = %d, want 2 (unchanged)", pouch.Count(MaterialHide))
	}
}

// TestMaterialPouchRemoveNonexistent verifies removal fails for materials not in pouch
func TestMaterialPouchRemoveNonexistent(t *testing.T) {
	pouch := NewMaterialPouch()

	removed := pouch.Remove(MaterialClaws, 1)

	if removed {
		t.Error("Remove() should return false for materials not in pouch")
	}
}

// TestMaterialPouchHas verifies Has() checks for sufficient quantity
func TestMaterialPouchHas(t *testing.T) {
	pouch := NewMaterialPouch()
	pouch.Add(MaterialScales, 5)

	if !pouch.Has(MaterialScales, 3) {
		t.Error("Has(Scales, 3) should return true when pouch has 5")
	}
	if !pouch.Has(MaterialScales, 5) {
		t.Error("Has(Scales, 5) should return true when pouch has exactly 5")
	}
	if pouch.Has(MaterialScales, 6) {
		t.Error("Has(Scales, 6) should return false when pouch has only 5")
	}
	if pouch.Has(MaterialClaws, 1) {
		t.Error("Has(Claws, 1) should return false when pouch has no claws")
	}
}

// TestMaterialPouchAllMaterials verifies listing all materials
func TestMaterialPouchAllMaterials(t *testing.T) {
	pouch := NewMaterialPouch()
	pouch.Add(MaterialScales, 3)
	pouch.Add(MaterialFangs, 1)
	pouch.Add(MaterialWyvernScale, 1)

	materials := pouch.AllMaterials()

	if len(materials) != 3 {
		t.Errorf("AllMaterials() length = %d, want 3", len(materials))
	}

	// Verify all expected types are present
	found := make(map[MaterialType]bool)
	for _, mt := range materials {
		found[mt] = true
	}

	if !found[MaterialScales] {
		t.Error("AllMaterials() should include Scales")
	}
	if !found[MaterialFangs] {
		t.Error("AllMaterials() should include Fangs")
	}
	if !found[MaterialWyvernScale] {
		t.Error("AllMaterials() should include WyvernScale")
	}
}

// TestMaterialPouchClear verifies clearing the pouch
func TestMaterialPouchClear(t *testing.T) {
	pouch := NewMaterialPouch()
	pouch.Add(MaterialScales, 5)
	pouch.Add(MaterialClaws, 3)
	pouch.Add(MaterialWyvernScale, 1)

	pouch.Clear()

	if pouch.TotalCount() != 0 {
		t.Errorf("After Clear(): TotalCount() = %d, want 0", pouch.TotalCount())
	}
	if len(pouch.AllMaterials()) != 0 {
		t.Errorf("After Clear(): AllMaterials() length = %d, want 0", len(pouch.AllMaterials()))
	}
}

// TestCommonMaterialTypes verifies GetCommonMaterialTypes returns only common materials
func TestCommonMaterialTypes(t *testing.T) {
	commonTypes := GetCommonMaterialTypes()

	if len(commonTypes) != 4 {
		t.Errorf("GetCommonMaterialTypes() length = %d, want 4", len(commonTypes))
	}

	for _, mt := range commonTypes {
		if mt.IsRare() {
			t.Errorf("GetCommonMaterialTypes() included rare material: %s", mt.String())
		}
	}
}

// TestRareMaterialTypes verifies GetRareMaterialTypes returns only rare materials
func TestRareMaterialTypes(t *testing.T) {
	rareTypes := GetRareMaterialTypes()

	if len(rareTypes) != 5 {
		t.Errorf("GetRareMaterialTypes() length = %d, want 5", len(rareTypes))
	}

	for _, mt := range rareTypes {
		if !mt.IsRare() {
			t.Errorf("GetRareMaterialTypes() included common material: %s", mt.String())
		}
	}
}
