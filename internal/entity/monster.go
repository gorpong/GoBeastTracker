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
	AIAggressive // Always chases, attacks twice
	AIDefensive  // Retreats when low HP
	AIFleeing    // Runs away from player
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
	BossNormal     BossBehavior = iota
	BossAggressive              // Double damage when below 50% HP
	BossTeleport                // Can teleport near player
	BossSummoner                // Spawns minions
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

// GetChaseDirection returns the direction to move toward the target position
func (m *Monster) GetChaseDirection(targetX, targetY int) ui.Direction {
	dx := targetX - m.X
	dy := targetY - m.Y

	if dx == 0 && dy == 0 {
		return ui.DirNone
	}

	// Horizontal wins ties
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

// GetFleeDirection returns the direction to move away from the target position
func (m *Monster) GetFleeDirection(targetX, targetY int) ui.Direction {
	dx := targetX - m.X
	dy := targetY - m.Y

	if dx == 0 && dy == 0 {
		return ui.DirNone
	}

	// Move opposite to chase direction
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

// IsLowHP returns true if monster is below 30% HP
func (m *Monster) IsLowHP() bool {
	return float64(m.HP)/float64(m.MaxHP) < 0.3
}

// IsEnraged returns true if boss is below 50% HP (for aggressive behavior)
func (m *Monster) IsEnraged() bool {
	return float64(m.HP)/float64(m.MaxHP) < 0.5
}

// CanTeleport returns true if teleport is off cooldown
func (m *Monster) CanTeleport() bool {
	return m.BossBehavior == BossTeleport && m.TeleportCooldown == 0
}

// CanSummon returns true if summon is off cooldown
func (m *Monster) CanSummon() bool {
	return m.BossBehavior == BossSummoner && m.SummonCooldown == 0
}

// TickCooldowns reduces cooldown timers by 1
func (m *Monster) TickCooldowns() {
	if m.TeleportCooldown > 0 {
		m.TeleportCooldown--
	}
	if m.SummonCooldown > 0 {
		m.SummonCooldown--
	}
}

// GetEffectiveAttack returns attack value, doubled if enraged aggressive boss
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

type DropTable struct {
	Guaranteed []MaterialType
	Possible   []MaterialType
}

func NewDropTable(guaranteed, possible []MaterialType) *DropTable {
	return &DropTable{
		Guaranteed: guaranteed,
		Possible:   possible,
	}
}

func (dt *DropTable) GenerateDrops() []MaterialType {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	drops := make([]MaterialType, 0)

	drops = append(drops, dt.Guaranteed...)

	for _, matType := range dt.Possible {
		if rng.Intn(2) == 0 {
			drops = append(drops, matType)
		}
	}

	return drops
}

func GetRegularMonsterDropTable() *DropTable {
	return NewDropTable(
		[]MaterialType{},
		GetCommonMaterialTypes(),
	)
}

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

	return NewDropTable(
		[]MaterialType{rareMaterial},
		GetCommonMaterialTypes(),
	)
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
