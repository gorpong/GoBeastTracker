package score

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewLeaderboard(t *testing.T) {
	lb := NewLeaderboard()

	if lb == nil {
		t.Fatal("NewLeaderboard returned nil")
	}

	if len(lb.Entries) != 0 {
		t.Errorf("New leaderboard should be empty, got %d entries", len(lb.Entries))
	}
}

func TestLeaderboardAdd(t *testing.T) {
	lb := NewLeaderboard()

	lb.Add("AAA", 100, 1)

	if len(lb.Entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(lb.Entries))
	}

	if lb.Entries[0].Initials != "AAA" {
		t.Errorf("Initials = %q, want %q", lb.Entries[0].Initials, "AAA")
	}
	if lb.Entries[0].Score != 100 {
		t.Errorf("Score = %d, want %d", lb.Entries[0].Score, 100)
	}
	if lb.Entries[0].Hunt != 1 {
		t.Errorf("Hunt = %d, want %d", lb.Entries[0].Hunt, 1)
	}
}

func TestLeaderboardSortsDescending(t *testing.T) {
	lb := NewLeaderboard()

	lb.Add("LOW", 50, 1)
	lb.Add("HIGH", 200, 2)
	lb.Add("MID", 100, 1)

	if lb.Entries[0].Score != 200 {
		t.Errorf("First entry score = %d, want 200", lb.Entries[0].Score)
	}
	if lb.Entries[1].Score != 100 {
		t.Errorf("Second entry score = %d, want 100", lb.Entries[1].Score)
	}
	if lb.Entries[2].Score != 50 {
		t.Errorf("Third entry score = %d, want 50", lb.Entries[2].Score)
	}
}

func TestLeaderboardMaxEntries(t *testing.T) {
	lb := NewLeaderboard()

	for i := 0; i < 15; i++ {
		lb.Add("TST", i*10, 1)
	}

	if len(lb.Entries) != MaxLeaderboardEntries {
		t.Errorf("Entries = %d, want %d", len(lb.Entries), MaxLeaderboardEntries)
	}

	if lb.Entries[len(lb.Entries)-1].Score != 50 {
		t.Errorf("Lowest score = %d, want 50 (scores 0-40 should be dropped)", 
			lb.Entries[len(lb.Entries)-1].Score)
	}
}

func TestLeaderboardIsHighScore(t *testing.T) {
	lb := NewLeaderboard()

	if !lb.IsHighScore(10) {
		t.Error("Empty leaderboard should accept any score")
	}

	for i := 0; i < MaxLeaderboardEntries; i++ {
		lb.Add("TST", (i+1)*100, 1)
	}

	if lb.IsHighScore(50) {
		t.Error("Score below minimum should not be high score")
	}

	if !lb.IsHighScore(500) {
		t.Error("Score above minimum should be high score")
	}
}

func TestLeaderboardSaveLoad(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_scores.json")

	lb := NewLeaderboard()
	lb.SetFilepath(testFile)

	lb.Add("AAA", 300, 3)
	lb.Add("BBB", 200, 2)
	lb.Add("CCC", 100, 1)

	if err := lb.Save(); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	lb2 := NewLeaderboard()
	lb2.SetFilepath(testFile)

	if err := lb2.Load(); err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(lb2.Entries) != 3 {
		t.Fatalf("Loaded entries = %d, want 3", len(lb2.Entries))
	}

	if lb2.Entries[0].Initials != "AAA" || lb2.Entries[0].Score != 300 {
		t.Errorf("First entry mismatch: got %+v", lb2.Entries[0])
	}
}

func TestLeaderboardLoadNonexistent(t *testing.T) {
	lb := NewLeaderboard()
	lb.SetFilepath("/nonexistent/path/scores.json")

	err := lb.Load()
	if err != nil {
		t.Errorf("Loading nonexistent file should not error, got: %v", err)
	}

	if len(lb.Entries) != 0 {
		t.Error("Leaderboard should be empty after loading nonexistent file")
	}
}

func TestLeaderboardSaveCreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "subdir", "nested", "scores.json")

	lb := NewLeaderboard()
	lb.SetFilepath(testFile)
	lb.Add("TST", 100, 1)

	if err := lb.Save(); err != nil {
		t.Fatalf("Save should create directories: %v", err)
	}

	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("Score file should exist after save")
	}
}

func TestPointConstants(t *testing.T) {
	if PointsPerMonster != 10 {
		t.Errorf("PointsPerMonster = %d, want 10", PointsPerMonster)
	}
	if PointsPerBoss != 50 {
		t.Errorf("PointsPerBoss = %d, want 50", PointsPerBoss)
	}
}
