package craft

import (
	"fmt"
	"sort"
	"strings"

	"beasttracker/internal/entity"
)

// Recipe defines materials needed and resulting equipment
type Recipe struct {
	Name        string
	Result      *entity.Equipment
	Ingredients map[entity.MaterialType]int
}

// NewRecipe creates a new crafting recipe
func NewRecipe(name string, result *entity.Equipment, ingredients map[entity.MaterialType]int) *Recipe {
	return &Recipe{
		Name:        name,
		Result:      result,
		Ingredients: ingredients,
	}
}

// CanCraft returns true if the pouch contains sufficient materials
func (r *Recipe) CanCraft(pouch *entity.MaterialPouch) bool {
	for matType, required := range r.Ingredients {
		if !pouch.Has(matType, required) {
			return false
		}
	}
	return true
}

// Craft attempts to craft the recipe, consuming materials from the pouch.
// Returns the crafted equipment and true if successful, nil and false otherwise.
func (r *Recipe) Craft(pouch *entity.MaterialPouch) (*entity.Equipment, bool) {
	if !r.CanCraft(pouch) {
		return nil, false
	}

	// Consume materials
	for matType, required := range r.Ingredients {
		pouch.Remove(matType, required)
	}

	return r.Result, true
}

// IngredientsString returns a formatted string of required ingredients
func (r *Recipe) IngredientsString() string {
	// Sort by material type for consistent ordering
	types := make([]entity.MaterialType, 0, len(r.Ingredients))
	for matType := range r.Ingredients {
		types = append(types, matType)
	}
	sort.Slice(types, func(i, j int) bool {
		return types[i] < types[j]
	})

	parts := make([]string, 0, len(r.Ingredients))
	for _, matType := range types {
		count := r.Ingredients[matType]
		parts = append(parts, fmt.Sprintf("%d %s", count, matType.String()))
	}

	return strings.Join(parts, ", ")
}

// allRecipes holds all available recipes
var allRecipes []*Recipe

// init initializes the recipe list
func init() {
	allRecipes = createAllRecipes()
}

// createAllRecipes builds the complete recipe list
func createAllRecipes() []*Recipe {
	recipes := make([]*Recipe, 0)

	// Basic Recipes (Common Materials Only)
	recipes = append(recipes, NewRecipe(
		"Iron Sword",
		entity.NewEquipment("Iron Sword", entity.SlotWeapon, 3, 0, 0),
		map[entity.MaterialType]int{
			entity.MaterialScales: 3,
			entity.MaterialClaws:  2,
		},
	))

	recipes = append(recipes, NewRecipe(
		"Bone Knife",
		entity.NewEquipment("Bone Knife", entity.SlotWeapon, 2, 0, 0),
		map[entity.MaterialType]int{
			entity.MaterialFangs:  2,
			entity.MaterialScales: 1,
		},
	))

	recipes = append(recipes, NewRecipe(
		"Leather Armor",
		entity.NewEquipment("Leather Armor", entity.SlotArmor, 0, 2, 10),
		map[entity.MaterialType]int{
			entity.MaterialHide:  4,
			entity.MaterialFangs: 1,
		},
	))

	recipes = append(recipes, NewRecipe(
		"Hide Vest",
		entity.NewEquipment("Hide Vest", entity.SlotArmor, 0, 1, 5),
		map[entity.MaterialType]int{
			entity.MaterialHide: 3,
		},
	))

	recipes = append(recipes, NewRecipe(
		"Hunter's Charm",
		entity.NewEquipment("Hunter's Charm", entity.SlotCharm, 1, 1, 0),
		map[entity.MaterialType]int{
			entity.MaterialFangs: 2,
			entity.MaterialClaws: 2,
		},
	))

	// Boss Recipes (Require Rare Materials)
	recipes = append(recipes, NewRecipe(
		"Wyvern Blade",
		entity.NewEquipment("Wyvern Blade", entity.SlotWeapon, 6, 0, 0),
		map[entity.MaterialType]int{
			entity.MaterialWyvernScale: 1,
			entity.MaterialScales:      3,
		},
	))

	recipes = append(recipes, NewRecipe(
		"Ogre Armor",
		entity.NewEquipment("Ogre Armor", entity.SlotArmor, 0, 5, 25),
		map[entity.MaterialType]int{
			entity.MaterialOgreHide: 1,
			entity.MaterialHide:     4,
		},
	))

	recipes = append(recipes, NewRecipe(
		"Troll Gauntlets",
		entity.NewEquipment("Troll Gauntlets", entity.SlotCharm, 4, 1, 0),
		map[entity.MaterialType]int{
			entity.MaterialTrollClaw: 1,
			entity.MaterialClaws:     2,
		},
	))

	recipes = append(recipes, NewRecipe(
		"Cyclops Monocle",
		entity.NewEquipment("Cyclops Monocle", entity.SlotCharm, 2, 2, 0),
		map[entity.MaterialType]int{
			entity.MaterialCyclopsEye: 1,
			entity.MaterialFangs:      2,
		},
	))

	recipes = append(recipes, NewRecipe(
		"Minotaur Horn Helm",
		entity.NewEquipment("Minotaur Horn Helm", entity.SlotArmor, 0, 3, 15),
		map[entity.MaterialType]int{
			entity.MaterialMinotaurHorn: 1,
			entity.MaterialHide:         3,
		},
	))

	return recipes
}

// GetAllRecipes returns all available recipes
func GetAllRecipes() []*Recipe {
	return allRecipes
}

// GetBasicRecipes returns recipes that only use common materials
func GetBasicRecipes() []*Recipe {
	result := make([]*Recipe, 0)
	for _, recipe := range allRecipes {
		hasRare := false
		for matType := range recipe.Ingredients {
			if matType.IsRare() {
				hasRare = true
				break
			}
		}
		if !hasRare {
			result = append(result, recipe)
		}
	}
	return result
}

// GetBossRecipes returns recipes that require rare (boss) materials
func GetBossRecipes() []*Recipe {
	result := make([]*Recipe, 0)
	for _, recipe := range allRecipes {
		for matType := range recipe.Ingredients {
			if matType.IsRare() {
				result = append(result, recipe)
				break
			}
		}
	}
	return result
}

// GetRecipeByName finds a recipe by its name, returns nil if not found
func GetRecipeByName(name string) *Recipe {
	for _, recipe := range allRecipes {
		if recipe.Name == name {
			return recipe
		}
	}
	return nil
}
