package craft

import (
	"strings"
	"testing"

	"beasttracker/internal/entity"
)

// TestNewRecipe verifies recipe creation
func TestNewRecipe(t *testing.T) {
	ingredients := map[entity.MaterialType]int{
		entity.MaterialScales: 3,
		entity.MaterialClaws:  2,
	}
	result := entity.NewEquipment("Iron Sword", entity.SlotWeapon, 3, 0, 0)

	recipe := NewRecipe("Iron Sword", result, ingredients)

	if recipe.Name != "Iron Sword" {
		t.Errorf("Recipe Name = %q, want \"Iron Sword\"", recipe.Name)
	}
	if recipe.Result != result {
		t.Error("Recipe Result does not match provided equipment")
	}
	if len(recipe.Ingredients) != 2 {
		t.Errorf("Recipe Ingredients count = %d, want 2", len(recipe.Ingredients))
	}
}

// TestRecipeCanCraftWithSufficientMaterials verifies crafting check passes
func TestRecipeCanCraftWithSufficientMaterials(t *testing.T) {
	ingredients := map[entity.MaterialType]int{
		entity.MaterialScales: 3,
		entity.MaterialClaws:  2,
	}
	result := entity.NewEquipment("Iron Sword", entity.SlotWeapon, 3, 0, 0)
	recipe := NewRecipe("Iron Sword", result, ingredients)

	pouch := entity.NewMaterialPouch()
	pouch.Add(entity.MaterialScales, 5)
	pouch.Add(entity.MaterialClaws, 3)

	if !recipe.CanCraft(pouch) {
		t.Error("CanCraft() should return true with sufficient materials")
	}
}

// TestRecipeCanCraftWithExactMaterials verifies crafting check passes with exact amounts
func TestRecipeCanCraftWithExactMaterials(t *testing.T) {
	ingredients := map[entity.MaterialType]int{
		entity.MaterialScales: 3,
		entity.MaterialClaws:  2,
	}
	result := entity.NewEquipment("Iron Sword", entity.SlotWeapon, 3, 0, 0)
	recipe := NewRecipe("Iron Sword", result, ingredients)

	pouch := entity.NewMaterialPouch()
	pouch.Add(entity.MaterialScales, 3)
	pouch.Add(entity.MaterialClaws, 2)

	if !recipe.CanCraft(pouch) {
		t.Error("CanCraft() should return true with exact materials")
	}
}

// TestRecipeCanCraftWithInsufficientMaterials verifies crafting check fails
func TestRecipeCanCraftWithInsufficientMaterials(t *testing.T) {
	ingredients := map[entity.MaterialType]int{
		entity.MaterialScales: 3,
		entity.MaterialClaws:  2,
	}
	result := entity.NewEquipment("Iron Sword", entity.SlotWeapon, 3, 0, 0)
	recipe := NewRecipe("Iron Sword", result, ingredients)

	pouch := entity.NewMaterialPouch()
	pouch.Add(entity.MaterialScales, 2) // Not enough
	pouch.Add(entity.MaterialClaws, 2)

	if recipe.CanCraft(pouch) {
		t.Error("CanCraft() should return false with insufficient materials")
	}
}

// TestRecipeCanCraftWithMissingMaterial verifies crafting check fails when missing a material
func TestRecipeCanCraftWithMissingMaterial(t *testing.T) {
	ingredients := map[entity.MaterialType]int{
		entity.MaterialScales: 3,
		entity.MaterialClaws:  2,
	}
	result := entity.NewEquipment("Iron Sword", entity.SlotWeapon, 3, 0, 0)
	recipe := NewRecipe("Iron Sword", result, ingredients)

	pouch := entity.NewMaterialPouch()
	pouch.Add(entity.MaterialScales, 5)
	// Missing claws entirely

	if recipe.CanCraft(pouch) {
		t.Error("CanCraft() should return false with missing material type")
	}
}

// TestRecipeCraft verifies successful crafting consumes materials
func TestRecipeCraft(t *testing.T) {
	ingredients := map[entity.MaterialType]int{
		entity.MaterialScales: 3,
		entity.MaterialClaws:  2,
	}
	result := entity.NewEquipment("Iron Sword", entity.SlotWeapon, 3, 0, 0)
	recipe := NewRecipe("Iron Sword", result, ingredients)

	pouch := entity.NewMaterialPouch()
	pouch.Add(entity.MaterialScales, 5)
	pouch.Add(entity.MaterialClaws, 3)

	craftedItem, ok := recipe.Craft(pouch)

	if !ok {
		t.Error("Craft() should return true when crafting succeeds")
	}
	if craftedItem != result {
		t.Error("Craft() should return the recipe's result equipment")
	}

	// Verify materials were consumed
	if pouch.Count(entity.MaterialScales) != 2 {
		t.Errorf("After crafting: Scales = %d, want 2 (5-3)", pouch.Count(entity.MaterialScales))
	}
	if pouch.Count(entity.MaterialClaws) != 1 {
		t.Errorf("After crafting: Claws = %d, want 1 (3-2)", pouch.Count(entity.MaterialClaws))
	}
}

// TestRecipeCraftFailsWithInsufficientMaterials verifies crafting fails gracefully
func TestRecipeCraftFailsWithInsufficientMaterials(t *testing.T) {
	ingredients := map[entity.MaterialType]int{
		entity.MaterialScales: 3,
	}
	result := entity.NewEquipment("Iron Sword", entity.SlotWeapon, 3, 0, 0)
	recipe := NewRecipe("Iron Sword", result, ingredients)

	pouch := entity.NewMaterialPouch()
	pouch.Add(entity.MaterialScales, 2) // Not enough

	craftedItem, ok := recipe.Craft(pouch)

	if ok {
		t.Error("Craft() should return false when insufficient materials")
	}
	if craftedItem != nil {
		t.Error("Craft() should return nil when crafting fails")
	}

	// Materials should be unchanged
	if pouch.Count(entity.MaterialScales) != 2 {
		t.Errorf("After failed craft: Scales = %d, want 2 (unchanged)", pouch.Count(entity.MaterialScales))
	}
}

// TestRecipeIngredientsString verifies ingredient list formatting
func TestRecipeIngredientsString(t *testing.T) {
	ingredients := map[entity.MaterialType]int{
		entity.MaterialScales: 3,
		entity.MaterialClaws:  2,
	}
	result := entity.NewEquipment("Iron Sword", entity.SlotWeapon, 3, 0, 0)
	recipe := NewRecipe("Iron Sword", result, ingredients)

	desc := recipe.IngredientsString()

	// Should contain both ingredients (order may vary due to map)
	if !strings.Contains(desc, "3 Scales") {
		t.Errorf("IngredientsString() = %q, should contain \"3 Scales\"", desc)
	}
	if !strings.Contains(desc, "2 Claws") {
		t.Errorf("IngredientsString() = %q, should contain \"2 Claws\"", desc)
	}
}

// TestGetAllRecipes verifies recipe list is populated
func TestGetAllRecipes(t *testing.T) {
	recipes := GetAllRecipes()

	if len(recipes) == 0 {
		t.Error("GetAllRecipes() should return at least one recipe")
	}

	// Verify we have both basic and boss recipes
	hasBasic := false
	hasBoss := false

	for _, recipe := range recipes {
		for matType := range recipe.Ingredients {
			if matType.IsRare() {
				hasBoss = true
			} else {
				hasBasic = true
			}
		}
	}

	if !hasBasic {
		t.Error("GetAllRecipes() should include recipes with common materials")
	}
	if !hasBoss {
		t.Error("GetAllRecipes() should include recipes with boss materials")
	}
}

// TestGetBasicRecipes verifies only basic recipes returned
func TestGetBasicRecipes(t *testing.T) {
	recipes := GetBasicRecipes()

	if len(recipes) == 0 {
		t.Error("GetBasicRecipes() should return at least one recipe")
	}

	for _, recipe := range recipes {
		for matType := range recipe.Ingredients {
			if matType.IsRare() {
				t.Errorf("GetBasicRecipes() should not include rare materials, found %s in %s",
					matType.String(), recipe.Name)
			}
		}
	}
}

// TestGetBossRecipes verifies boss recipes require rare materials
func TestGetBossRecipes(t *testing.T) {
	recipes := GetBossRecipes()

	if len(recipes) == 0 {
		t.Error("GetBossRecipes() should return at least one recipe")
	}

	for _, recipe := range recipes {
		hasRare := false
		for matType := range recipe.Ingredients {
			if matType.IsRare() {
				hasRare = true
				break
			}
		}
		if !hasRare {
			t.Errorf("GetBossRecipes() recipe %s should require at least one rare material",
				recipe.Name)
		}
	}
}

// TestRecipeByName verifies finding recipe by name
func TestRecipeByName(t *testing.T) {
	recipe := GetRecipeByName("Iron Sword")

	if recipe == nil {
		t.Fatal("GetRecipeByName(\"Iron Sword\") returned nil")
	}
	if recipe.Name != "Iron Sword" {
		t.Errorf("Recipe Name = %q, want \"Iron Sword\"", recipe.Name)
	}
}

// TestRecipeByNameNotFound verifies nil returned for unknown recipe
func TestRecipeByNameNotFound(t *testing.T) {
	recipe := GetRecipeByName("Nonexistent Weapon")

	if recipe != nil {
		t.Error("GetRecipeByName() should return nil for unknown recipe")
	}
}
