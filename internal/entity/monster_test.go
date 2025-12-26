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
