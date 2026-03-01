package entity

import (
	"math/rand"
	"time"

	"beasttracker/internal/ui"
)

type AIType int

const (
	AIWander AIType = iota
	AIChase
	AIAggressive
	AIDefensive
	AIFleeing
)

func (a AIType) String() string {
	switch a {
	case AIWander:
		return "Wander"
	case AIChase:
		return "Chase"
	case AIAggressive:
		return "Aggressive"
	case AIDefensive:
		return "Defensive"
	case AIFleeing:
		return "Fleeing"
	default:
		return "Unknown"
	}
}

type BossBehavior int

const (
	BossNormal BossBehavior = iota
	BossAggressive
	BossTeleport
	BossSummoner
)

func (b BossBehavior) String() string {
	switch b {
	case BossNormal:
		return "Normal"
	case BossAggressive:
		return "Aggressive"
	case BossTeleport:
		return "Teleport"
	case BossSummoner:
		return "Summoner"
	default:
		return "Unknown"
	}
}

type Monster struct {
	Name             string
	Glyph            rune
	X                int
	Y                int
	HP               int
	MaxHP            int
	Attack           int
	Defense          int
	AI               AIType
	Dead             bool
	IsBoss           bool
	BossBehavior     BossBehavior
	DropTable        *DropTable
	TeleportCooldown int
	SummonCooldown   int
}

func NewMonster(name string, glyph rune, x, y, hp, attack int) *Monster {
	return &Monster{
		Name:    name,
		Glyph:   glyph,
		X:       x,
		Y:       y,
		HP:      hp,
		MaxHP:   hp,
		Attack:  attack,
		Defense: 0,
		AI:      AIWander,
		Dead:    false,
		IsBoss:  false,
	}
}

func NewMonsterWithAI(name string, glyph rune, x, y, hp, attack int, ai AIType) *Monster {
	m := NewMonster(name, glyph, x, y, hp, attack)
	m.AI = ai
	return m
}

func NewBossMonster(name string, glyph rune, x, y, hp, attack int) *Monster {
	return &Monster{
		Name:         name,
		Glyph:        glyph,
		X:            x,
		Y:            y,
		HP:           hp,
		MaxHP:        hp,
		Attack:       attack,
		Defense:      2,
		AI:           AIChase,
		Dead:         false,
		IsBoss:       true,
		BossBehavior: BossNormal,
	}
}

func NewBossMonsterWithBehavior(name string, glyph rune, x, y, hp, attack int, behavior BossBehavior) *Monster {
	m := NewBossMonster(name, glyph, x, y, hp, attack)
	m.BossBehavior = behavior
	return m
}

func (m *Monster) Position() (int, int) {
	return m.X, m.Y
}

func (m *Monster) SetPosition(x, y int) {
	m.X = x
	m.Y = y
}

func (m *Monster) Move(dir ui.Direction) {
	dx, dy := dir.Delta()
	m.X += dx
	m.Y += dy
}

func (m *Monster) TakeDamage(damage int) {
	m.HP -= damage
	if m.HP <= 0 {
		m.HP = 0
		m.Dead = true
	}
}

func (m *Monster) IsAlive() bool {
	return !m.Dead
}

func (m *Monster) GetChaseDirection(targetX, targetY int) ui.Direction {
	dx := targetX - m.X
	dy := targetY - m.Y

	if dx == 0 && dy == 0 {
		return ui.DirNone
	}

	if abs(dx) >= abs(dy) && dx != 0 {
		if dx > 0 {
			return ui.DirRight
		}
		return ui.DirLeft
	}

	if dy > 0 {
		return ui.DirDown
	}
	if dy < 0 {
		return ui.DirUp
	}

	return ui.DirNone
}

func (m *Monster) GetFleeDirection(targetX, targetY int) ui.Direction {
	dx := targetX - m.X
	dy := targetY - m.Y

	if dx == 0 && dy == 0 {
		return ui.DirNone
	}

	if abs(dx) >= abs(dy) && dx != 0 {
		if dx > 0 {
			return ui.DirLeft
		}
		return ui.DirRight
	}

	if dy > 0 {
		return ui.DirUp
	}
	if dy < 0 {
		return ui.DirDown
	}

	return ui.DirNone
}

func (m *Monster) IsLowHP() bool {
	return float64(m.HP)/float64(m.MaxHP) < 0.3
}

func (m *Monster) IsEnraged() bool {
	return float64(m.HP)/float64(m.MaxHP) < 0.5
}

func (m *Monster) CanTeleport() bool {
	return m.BossBehavior == BossTeleport && m.TeleportCooldown == 0
}

func (m *Monster) CanSummon() bool {
	return m.BossBehavior == BossSummoner && m.SummonCooldown == 0
}

func (m *Monster) TickCooldowns() {
	if m.TeleportCooldown > 0 {
		m.TeleportCooldown--
	}
	if m.SummonCooldown > 0 {
		m.SummonCooldown--
	}
}

func (m *Monster) GetEffectiveAttack() int {
	if m.IsBoss && m.BossBehavior == BossAggressive && m.IsEnraged() {
		return m.Attack * 2
	}
	return m.Attack
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// DropTable defines what materials a monster can drop on death
type DropTable struct {
	// For regular monsters: chance to drop anything at all
	DropChance float64
	// For bosses: guaranteed rare material (zero value means none)
	GuaranteedRare MaterialType
	// Whether there's a guaranteed drop
	HasGuaranteed bool
	// Possible common materials to drop
	PossibleDrops []MaterialType
	// For bosses: chance to get bonus common drops
	BonusChance float64
	// Maximum bonus drops
	MaxBonusDrops int
}

// NewDropTable creates a drop table for regular monsters (legacy compatibility)
func NewDropTable(guaranteed, possible []MaterialType) *DropTable {
	dt := &DropTable{
		DropChance:    0.30,
		PossibleDrops: possible,
		HasGuaranteed: len(guaranteed) > 0,
		BonusChance:   0.0,
		MaxBonusDrops: 0,
	}
	if len(guaranteed) > 0 {
		dt.GuaranteedRare = guaranteed[0]
	}
	return dt
}

// GenerateDrops generates the actual drops based on the drop table
func (dt *DropTable) GenerateDrops() []MaterialType {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	drops := make([]MaterialType, 0)

	// Add guaranteed rare drop for bosses
	if dt.HasGuaranteed {
		drops = append(drops, dt.GuaranteedRare)

		// Boss bonus drops
		if dt.BonusChance > 0 && rng.Float64() < dt.BonusChance {
			numBonus := rng.Intn(dt.MaxBonusDrops) + 1
			for i := 0; i < numBonus; i++ {
				if len(dt.PossibleDrops) > 0 {
					randomMat := dt.PossibleDrops[rng.Intn(len(dt.PossibleDrops))]
					drops = append(drops, randomMat)
				}
			}
		}
		return drops
	}

	// Regular monster: DropChance to drop anything, then 1 random item
	if rng.Float64() < dt.DropChance {
		if len(dt.PossibleDrops) > 0 {
			randomMat := dt.PossibleDrops[rng.Intn(len(dt.PossibleDrops))]
			drops = append(drops, randomMat)
		}
	}

	return drops
}

// Guaranteed returns the guaranteed drops (for test compatibility)
func (dt *DropTable) Guaranteed() []MaterialType {
	if dt.HasGuaranteed {
		return []MaterialType{dt.GuaranteedRare}
	}
	return []MaterialType{}
}

// Possible returns the possible drops (for test compatibility)
func (dt *DropTable) Possible() []MaterialType {
	return dt.PossibleDrops
}

// GetRegularMonsterDropTable returns the drop table for regular monsters
// 30% chance to drop anything, then 1 random common material
func GetRegularMonsterDropTable() *DropTable {
	return &DropTable{
		DropChance:    0.30,
		HasGuaranteed: false,
		PossibleDrops: GetCommonMaterialTypes(),
		BonusChance:   0.0,
		MaxBonusDrops: 0,
	}
}

// GetBossDropTable returns the drop table for a specific boss type
// 1 guaranteed rare + 50% chance for 1-2 commons
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
		rareMaterial = MaterialWyvernScale
	}

	return &DropTable{
		DropChance:     0.0, // Not used for bosses
		GuaranteedRare: rareMaterial,
		HasGuaranteed:  true,
		PossibleDrops:  GetCommonMaterialTypes(),
		BonusChance:    0.50,
		MaxBonusDrops:  2,
	}
}

// GetBossBehaviorForType returns the behavior associated with each boss type
func GetBossBehaviorForType(bossName string) BossBehavior {
	switch bossName {
	case "Wyvern":
		return BossTeleport
	case "Ogre":
		return BossAggressive
	case "Troll":
		return BossAggressive
	case "Cyclops":
		return BossSummoner
	case "Minotaur":
		return BossNormal
	default:
		return BossNormal
	}
}
