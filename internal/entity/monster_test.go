package entity

import (
	"testing"

	"beasttracker/internal/ui"
)

// TestNewMonster verifies monster creation with correct initial values
func TestNewMonster(t *testing.T) {
	m := NewMonster("Goblin", 'g', 10, 10, 20, 3)

	if m.Name != "Goblin" {
		t.Errorf("Monster Name = %q, want 'Goblin'", m.Name)
	}
	if m.Glyph != 'g' {
		t.Errorf("Monster Glyph = %q, want 'g'", m.Glyph)
	}
	if m.X != 10 {
		t.Errorf("Monster X = %d, want 10", m.X)
	}
	if m.Y != 10 {
		t.Errorf("Monster Y = %d, want 10", m.Y)
	}
	if m.HP != 20 {
		t.Errorf("Monster HP = %d, want 20", m.HP)
	}
	if m.MaxHP != 20 {
		t.Errorf("Monster MaxHP = %d, want 20", m.MaxHP)
	}
	if m.Attack != 3 {
		t.Errorf("Monster Attack = %d, want 3", m.Attack)
	}
	if m.AI != AIWander {
		t.Errorf("Monster AI = %v, want AIWander", m.AI)
	}
	if m.Dead {
		t.Error("New monster should not be dead")
	}
}

// TestMonsterPosition verifies Position() returns correct coordinates
func TestMonsterPosition(t *testing.T) {
	m := NewMonster("Rat", 'r', 5, 7, 10, 2)
	x, y := m.Position()

	if x != 5 || y != 7 {
		t.Errorf("Position() = (%d, %d), want (5, 7)", x, y)
	}
}

// TestMonsterSetPosition verifies SetPosition correctly updates coordinates
func TestMonsterSetPosition(t *testing.T) {
	m := NewMonster("Spider", 's', 0, 0, 15, 2)
	m.SetPosition(20, 30)

	if m.X != 20 || m.Y != 30 {
		t.Errorf("After SetPosition(20, 30): (%d, %d), want (20, 30)", m.X, m.Y)
	}
}

// TestMonsterMove verifies monster moves correctly
func TestMonsterMove(t *testing.T) {
	m := NewMonster("Bat", 'b', 10, 10, 8, 1)

	m.Move(ui.DirRight)
	if m.X != 11 || m.Y != 10 {
		t.Errorf("After move right: (%d, %d), want (11, 10)", m.X, m.Y)
	}

	m.Move(ui.DirDown)
	if m.X != 11 || m.Y != 11 {
		t.Errorf("After move down: (%d, %d), want (11, 11)", m.X, m.Y)
	}
}

// TestMonsterTakeDamage verifies damage reduces HP
func TestMonsterTakeDamage(t *testing.T) {
	m := NewMonster("Orc", 'o', 0, 0, 30, 5)

	m.TakeDamage(10)
	if m.HP != 20 {
		t.Errorf("After taking 10 damage: HP = %d, want 20", m.HP)
	}
	if m.Dead {
		t.Error("Monster should not be dead at 20 HP")
	}
}

// TestMonsterDeath verifies monster dies when HP reaches 0
func TestMonsterDeath(t *testing.T) {
	m := NewMonster("Slime", 'S', 0, 0, 15, 2)

	m.TakeDamage(15)
	if m.HP != 0 {
		t.Errorf("After taking fatal damage: HP = %d, want 0", m.HP)
	}
	if !m.Dead {
		t.Error("Monster should be dead at 0 HP")
	}
}

// TestMonsterOverkillDamage verifies HP doesn't go negative
func TestMonsterOverkillDamage(t *testing.T) {
	m := NewMonster("Zombie", 'Z', 0, 0, 10, 3)

	m.TakeDamage(50)
	if m.HP < 0 {
		t.Errorf("HP should not be negative: HP = %d", m.HP)
	}
	if m.HP != 0 {
		t.Errorf("HP should be 0 after overkill: HP = %d", m.HP)
	}
	if !m.Dead {
		t.Error("Monster should be dead")
	}
}

// TestMonsterIsAlive verifies IsAlive() returns correct status
func TestMonsterIsAlive(t *testing.T) {
	m := NewMonster("Wolf", 'w', 0, 0, 25, 4)

	if !m.IsAlive() {
		t.Error("New monster should be alive")
	}

	m.TakeDamage(25)
	if m.IsAlive() {
		t.Error("Monster should not be alive after fatal damage")
	}
}

// TestMonsterAIType verifies AI type setting
func TestMonsterAIType(t *testing.T) {
	m := NewMonster("Dragon", 'D', 0, 0, 100, 15)
	m.AI = AIChase

	if m.AI != AIChase {
		t.Errorf("AI = %v, want AIChase", m.AI)
	}
}

// TestAITypeString verifies AI type string representation
func TestAITypeString(t *testing.T) {
	tests := []struct {
		ai   AIType
		want string
	}{
		{AIWander, "Wander"},
		{AIChase, "Chase"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.ai.String(); got != tt.want {
				t.Errorf("AIType.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestNewBossMonster verifies boss monster creation
func TestNewBossMonster(t *testing.T) {
	boss := NewBossMonster("Wyvern", 'W', 10, 10, 100, 15)

	if boss.Name != "Wyvern" {
		t.Errorf("Boss Name = %q, want 'Wyvern'", boss.Name)
	}
	if !boss.IsBoss {
		t.Error("Boss IsBoss should be true")
	}
	if boss.HP != 100 {
		t.Errorf("Boss HP = %d, want 100", boss.HP)
	}
}

// TestBossMonsterHigherStats verifies boss has higher stats than regular
func TestBossMonsterHigherStats(t *testing.T) {
	regular := NewMonster("Goblin", 'g', 0, 0, 15, 3)
	boss := NewBossMonster("Wyvern", 'W', 0, 0, 100, 15)

	// Boss should have significantly more HP
	if boss.HP <= regular.HP*2 {
		t.Errorf("Boss HP (%d) should be much higher than regular (%d)", boss.HP, regular.HP)
	}

	// Boss attack should be higher
	if boss.Attack <= regular.Attack {
		t.Errorf("Boss Attack (%d) should be higher than regular (%d)", boss.Attack, regular.Attack)
	}
}

// TestDropTableCreation verifies drop table creation
func TestDropTableCreation(t *testing.T) {
	dropTable := NewDropTable(
		[]MaterialType{MaterialScales},
		[]MaterialType{MaterialClaws, MaterialFangs},
	)

	if len(dropTable.Guaranteed) != 1 {
		t.Errorf("DropTable Guaranteed length = %d, want 1", len(dropTable.Guaranteed))
	}
	if len(dropTable.Possible) != 2 {
		t.Errorf("DropTable Possible length = %d, want 2", len(dropTable.Possible))
	}
}

// TestDropTableGenerateDropsGuaranteed verifies guaranteed drops always occur
func TestDropTableGenerateDropsGuaranteed(t *testing.T) {
	dropTable := NewDropTable(
		[]MaterialType{MaterialScales, MaterialClaws},
		[]MaterialType{},
	)

	// Run multiple times to verify guaranteed drops
	for i := 0; i < 10; i++ {
		drops := dropTable.GenerateDrops()

		if len(drops) != 2 {
			t.Errorf("Iteration %d: drops length = %d, want 2 (guaranteed)", i, len(drops))
		}

		hasScales := false
		hasClaws := false
		for _, drop := range drops {
			if drop == MaterialScales {
				hasScales = true
			}
			if drop == MaterialClaws {
				hasClaws = true
			}
		}

		if !hasScales {
			t.Errorf("Iteration %d: missing guaranteed Scales", i)
		}
		if !hasClaws {
			t.Errorf("Iteration %d: missing guaranteed Claws", i)
		}
	}
}

// TestDropTableGenerateDropsPossible verifies possible drops are random
func TestDropTableGenerateDropsPossible(t *testing.T) {
	dropTable := NewDropTable(
		[]MaterialType{},
		[]MaterialType{MaterialFangs},
	)

	// Run many times - should get some drops and some empty
	gotDrop := false
	gotEmpty := false

	for i := 0; i < 50; i++ {
		drops := dropTable.GenerateDrops()
		if len(drops) > 0 {
			gotDrop = true
		} else {
			gotEmpty = true
		}

		if gotDrop && gotEmpty {
			break
		}
	}

	if !gotDrop {
		t.Error("Possible drops should sometimes produce materials")
	}
	if !gotEmpty {
		t.Error("Possible drops should sometimes be empty (50% chance)")
	}
}

// TestMonsterDropTable verifies monsters have drop tables
func TestMonsterDropTable(t *testing.T) {
	monster := NewMonster("Goblin", 'g', 10, 10, 15, 3)

	// Regular monster should have nil drop table by default
	if monster.DropTable != nil {
		t.Error("New monster should have nil DropTable by default")
	}

	// Set a drop table
	dropTable := NewDropTable(
		[]MaterialType{},
		[]MaterialType{MaterialScales, MaterialClaws},
	)
	monster.DropTable = dropTable

	if monster.DropTable == nil {
		t.Error("Monster DropTable should be set after assignment")
	}
}

// TestBossMonsterDropTable verifies boss drop tables include rare materials
func TestBossMonsterDropTable(t *testing.T) {
	boss := NewBossMonster("Wyvern", 'W', 10, 10, 100, 15)

	// Set boss drop table with rare material
	dropTable := NewDropTable(
		[]MaterialType{MaterialWyvernScale},
		[]MaterialType{MaterialScales, MaterialClaws},
	)
	boss.DropTable = dropTable

	// Verify guaranteed rare drop
	hasRare := false
	for _, mat := range boss.DropTable.Guaranteed {
		if mat.IsRare() {
			hasRare = true
			break
		}
	}

	if !hasRare {
		t.Error("Boss drop table should include a guaranteed rare material")
	}
}

// TestGetRegularMonsterDropTable verifies drop table assignment for regular monsters
func TestGetRegularMonsterDropTable(t *testing.T) {
	dropTable := GetRegularMonsterDropTable()

	if dropTable == nil {
		t.Fatal("GetRegularMonsterDropTable() returned nil")
	}

	// Regular monsters should not have guaranteed drops
	if len(dropTable.Guaranteed) != 0 {
		t.Errorf("Regular monster should have no guaranteed drops, got %d", len(dropTable.Guaranteed))
	}

	// Should have possible drops
	if len(dropTable.Possible) == 0 {
		t.Error("Regular monster should have possible drops")
	}

	// Should only have common materials
	for _, mat := range dropTable.Possible {
		if mat.IsRare() {
			t.Errorf("Regular monster drop table should not include rare material: %s", mat.String())
		}
	}
}

// TestGetBossDropTable verifies drop table for specific boss types
func TestGetBossDropTable(t *testing.T) {
	tests := []struct {
		bossName     string
		expectedRare MaterialType
	}{
		{"Wyvern", MaterialWyvernScale},
		{"Ogre", MaterialOgreHide},
		{"Troll", MaterialTrollClaw},
		{"Cyclops", MaterialCyclopsEye},
		{"Minotaur", MaterialMinotaurHorn},
	}

	for _, tt := range tests {
		t.Run(tt.bossName, func(t *testing.T) {
			dropTable := GetBossDropTable(tt.bossName)

			if dropTable == nil {
				t.Fatalf("GetBossDropTable(%q) returned nil", tt.bossName)
			}

			// Should have guaranteed rare material
			hasExpectedRare := false
			for _, mat := range dropTable.Guaranteed {
				if mat == tt.expectedRare {
					hasExpectedRare = true
					break
				}
			}

			if !hasExpectedRare {
				t.Errorf("Boss %s should guarantee %s drop", tt.bossName, tt.expectedRare.String())
			}
		})
	}
}

func TestMonsterChaseAI(t *testing.T) {
	boss := NewBossMonster("TestBoss", 'B', 10, 10, 100, 10)

	if boss.AI != AIChase {
		t.Errorf("Boss AI = %v, want AIChase", boss.AI)
	}
}

func TestMonsterGetChaseDirection(t *testing.T) {
	monster := NewMonster("Test", 't', 5, 5, 10, 2)

	tests := []struct {
		targetX, targetY int
		wantDir          ui.Direction
	}{
		{8, 5, ui.DirRight},
		{2, 5, ui.DirLeft},
		{5, 8, ui.DirDown},
		{5, 2, ui.DirUp},
		{8, 8, ui.DirRight}, // Diagonal equal, horizontal wins
		{8, 6, ui.DirRight},
		{6, 8, ui.DirDown},
		{5, 5, ui.DirNone},
	}

	for _, tc := range tests {
		dir := monster.GetChaseDirection(tc.targetX, tc.targetY)
		if dir != tc.wantDir {
			t.Errorf("GetChaseDirection(%d,%d) from (5,5) = %v, want %v",
				tc.targetX, tc.targetY, dir, tc.wantDir)
		}
	}
}

func TestMonsterGetFleeDirection(t *testing.T) {
	monster := NewMonster("Test", 't', 5, 5, 10, 2)

	tests := []struct {
		targetX, targetY int
		wantDir          ui.Direction
	}{
		{8, 5, ui.DirLeft},  // Target right, flee left
		{2, 5, ui.DirRight}, // Target left, flee right
		{5, 8, ui.DirUp},    // Target below, flee up
		{5, 2, ui.DirDown},  // Target above, flee down
		{5, 5, ui.DirNone},
	}

	for _, tc := range tests {
		dir := monster.GetFleeDirection(tc.targetX, tc.targetY)
		if dir != tc.wantDir {
			t.Errorf("GetFleeDirection(%d,%d) from (5,5) = %v, want %v",
				tc.targetX, tc.targetY, dir, tc.wantDir)
		}
	}
}

func TestMonsterIsLowHP(t *testing.T) {
	monster := NewMonster("Test", 't', 0, 0, 100, 5)

	if monster.IsLowHP() {
		t.Error("Full HP monster should not be low HP")
	}

	monster.HP = 29
	if !monster.IsLowHP() {
		t.Error("Monster at 29% HP should be low HP")
	}

	monster.HP = 30
	if monster.IsLowHP() {
		t.Error("Monster at exactly 30% HP should not be low HP")
	}
}

func TestBossIsEnraged(t *testing.T) {
	boss := NewBossMonster("Test", 'T', 0, 0, 100, 10)

	if boss.IsEnraged() {
		t.Error("Full HP boss should not be enraged")
	}

	boss.HP = 49
	if !boss.IsEnraged() {
		t.Error("Boss at 49% HP should be enraged")
	}

	boss.HP = 50
	if boss.IsEnraged() {
		t.Error("Boss at exactly 50% HP should not be enraged")
	}
}

func TestBossEffectiveAttack(t *testing.T) {
	boss := NewBossMonsterWithBehavior("Ogre", 'O', 0, 0, 100, 10, BossAggressive)

	if boss.GetEffectiveAttack() != 10 {
		t.Errorf("Non-enraged boss attack = %d, want 10", boss.GetEffectiveAttack())
	}

	boss.HP = 40
	if boss.GetEffectiveAttack() != 20 {
		t.Errorf("Enraged aggressive boss attack = %d, want 20", boss.GetEffectiveAttack())
	}
}

func TestBossEffectiveAttackNonAggressive(t *testing.T) {
	boss := NewBossMonsterWithBehavior("Minotaur", 'M', 0, 0, 100, 10, BossNormal)

	boss.HP = 40
	if boss.GetEffectiveAttack() != 10 {
		t.Errorf("Enraged non-aggressive boss attack = %d, want 10 (no bonus)", boss.GetEffectiveAttack())
	}
}

func TestBossTeleportCooldown(t *testing.T) {
	boss := NewBossMonsterWithBehavior("Wyvern", 'W', 0, 0, 100, 10, BossTeleport)

	if !boss.CanTeleport() {
		t.Error("Boss with teleport behavior should be able to teleport initially")
	}

	boss.TeleportCooldown = 3
	if boss.CanTeleport() {
		t.Error("Boss on cooldown should not be able to teleport")
	}

	boss.TickCooldowns()
	if boss.TeleportCooldown != 2 {
		t.Errorf("Cooldown after tick = %d, want 2", boss.TeleportCooldown)
	}
}

func TestBossSummonCooldown(t *testing.T) {
	boss := NewBossMonsterWithBehavior("Cyclops", 'C', 0, 0, 100, 10, BossSummoner)

	if !boss.CanSummon() {
		t.Error("Boss with summoner behavior should be able to summon initially")
	}

	boss.SummonCooldown = 5
	if boss.CanSummon() {
		t.Error("Boss on cooldown should not be able to summon")
	}
}

func TestNewMonsterWithAI(t *testing.T) {
	aggressive := NewMonsterWithAI("Wolf", 'w', 5, 5, 20, 5, AIAggressive)

	if aggressive.AI != AIAggressive {
		t.Errorf("Monster AI = %v, want AIAggressive", aggressive.AI)
	}
}

func TestGetBossBehaviorForType(t *testing.T) {
	tests := []struct {
		bossName string
		want     BossBehavior
	}{
		{"Wyvern", BossTeleport},
		{"Ogre", BossAggressive},
		{"Troll", BossAggressive},
		{"Cyclops", BossSummoner},
		{"Minotaur", BossNormal},
		{"Unknown", BossNormal},
	}

	for _, tc := range tests {
		got := GetBossBehaviorForType(tc.bossName)
		if got != tc.want {
			t.Errorf("GetBossBehaviorForType(%q) = %v, want %v", tc.bossName, got, tc.want)
		}
	}
}

func TestAITypeStrings(t *testing.T) {
	tests := []struct {
		ai   AIType
		want string
	}{
		{AIWander, "Wander"},
		{AIChase, "Chase"},
		{AIAggressive, "Aggressive"},
		{AIDefensive, "Defensive"},
		{AIFleeing, "Fleeing"},
	}

	for _, tc := range tests {
		if got := tc.ai.String(); got != tc.want {
			t.Errorf("AIType(%d).String() = %q, want %q", tc.ai, got, tc.want)
		}
	}
}

func TestBossBehaviorStrings(t *testing.T) {
	tests := []struct {
		behavior BossBehavior
		want     string
	}{
		{BossNormal, "Normal"},
		{BossAggressive, "Aggressive"},
		{BossTeleport, "Teleport"},
		{BossSummoner, "Summoner"},
	}

	for _, tc := range tests {
		if got := tc.behavior.String(); got != tc.want {
			t.Errorf("BossBehavior(%d).String() = %q, want %q", tc.behavior, got, tc.want)
		}
	}
}

func TestMonsterDefense(t *testing.T) {
	regular := NewMonster("Goblin", 'g', 0, 0, 15, 3)
	boss := NewBossMonster("Wyvern", 'W', 0, 0, 100, 10)

	if regular.Defense != 0 {
		t.Errorf("Regular monster defense = %d, want 0", regular.Defense)
	}

	if boss.Defense != 2 {
		t.Errorf("Boss defense = %d, want 2", boss.Defense)
	}
}
