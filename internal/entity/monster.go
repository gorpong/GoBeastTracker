package entity

import (
	"math/rand"
	"time"

	"beasttracker/internal/ui"
)

// AIType represents the behavior pattern of a monster
type AIType int

const (
	AIWander AIType = iota // Wanders randomly
	AIChase                // Chases player when in sight
)

// String returns the string representation of an AIType
func (a AIType) String() string {
	switch a {
	case AIWander:
		return "Wander"
	case AIChase:
		return "Chase"
	default:
		return "Unknown"
	}
}

// Monster represents an enemy in the game
type Monster struct {
	Name      string
	Glyph     rune
	X         int
	Y         int
	HP        int
	MaxHP     int
	Attack    int
	AI        AIType
	Dead      bool
	IsBoss    bool
	DropTable *DropTable // Materials dropped on death
}

// NewMonster creates a new monster with the specified attributes
func NewMonster(name string, glyph rune, x, y, hp, attack int) *Monster {
	return &Monster{
		Name:   name,
		Glyph:  glyph,
		X:      x,
		Y:      y,
		HP:     hp,
		MaxHP:  hp,
		Attack: attack,
		AI:     AIWander,
		Dead:   false,
		IsBoss: false,
	}
}

// NewBossMonster creates a boss monster (target monster for the hunt)
func NewBossMonster(name string, glyph rune, x, y, hp, attack int) *Monster {
	return &Monster{
		Name:   name,
		Glyph:  glyph,
		X:      x,
		Y:      y,
		HP:     hp,
		MaxHP:  hp,
		Attack: attack,
		AI:     AIWander, // Can be upgraded to AIChase in Phase 10
		Dead:   false,
		IsBoss: true,
	}
}

// Position returns the monster's current coordinates
func (m *Monster) Position() (int, int) {
	return m.X, m.Y
}

// SetPosition sets the monster's position to the specified coordinates
func (m *Monster) SetPosition(x, y int) {
	m.X = x
	m.Y = y
}

// Move moves the monster in the specified direction
func (m *Monster) Move(dir ui.Direction) {
	dx, dy := dir.Delta()
	m.X += dx
	m.Y += dy
}

// TakeDamage reduces the monster's HP by the specified amount
func (m *Monster) TakeDamage(damage int) {
	m.HP -= damage
	if m.HP <= 0 {
		m.HP = 0
		m.Dead = true
	}
}

// IsAlive returns true if the monster is still alive
func (m *Monster) IsAlive() bool {
	return !m.Dead
}

// DropTable defines what materials a monster can drop on death
type DropTable struct {
	Guaranteed []MaterialType // Always drops these
	Possible   []MaterialType // 50% chance each
}

// NewDropTable creates a new drop table with guaranteed and possible drops
func NewDropTable(guaranteed, possible []MaterialType) *DropTable {
	return &DropTable{
		Guaranteed: guaranteed,
		Possible:   possible,
	}
}

// GenerateDrops generates the actual drops based on the drop table
func (dt *DropTable) GenerateDrops() []MaterialType {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	drops := make([]MaterialType, 0)

	// Add all guaranteed drops
	drops = append(drops, dt.Guaranteed...)

	// Roll for each possible drop (50% chance)
	for _, matType := range dt.Possible {
		if rng.Intn(2) == 0 {
			drops = append(drops, matType)
		}
	}

	return drops
}

// GetRegularMonsterDropTable returns the drop table for regular monsters
func GetRegularMonsterDropTable() *DropTable {
	return NewDropTable(
		[]MaterialType{},         // No guaranteed drops
		GetCommonMaterialTypes(), // Can drop any common material
	)
}

// GetBossDropTable returns the drop table for a specific boss type
func GetBossDropTable(bossName string) *DropTable {
	var rareMaterial MaterialType

	switch bossName {
	case "Wyvern":
		rareMaterial = MaterialWyvernScale
	case "Ogre":
		rareMaterial = MaterialOgreHide
	case "Troll":
		rareMaterial = MaterialTrollClaw
	case "Cyclops":
		rareMaterial = MaterialCyclopsEye
	case "Minotaur":
		rareMaterial = MaterialMinotaurHorn
	default:
		// Unknown boss, default to Wyvern
		rareMaterial = MaterialWyvernScale
	}

	return NewDropTable(
		[]MaterialType{rareMaterial}, // Guaranteed rare drop
		GetCommonMaterialTypes(),     // Can also drop common materials
	)
}
