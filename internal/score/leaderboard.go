package score

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

const (
	MaxLeaderboardEntries = 10
	DefaultScoreFile      = "assets/data/scores.json"

	PointsPerMonster = 10
	PointsPerBoss    = 50
)

type Entry struct {
	Initials string `json:"initials"`
	Score    int    `json:"score"`
	Hunt     int    `json:"hunt"`
}

type Leaderboard struct {
	Entries  []Entry `json:"entries"`
	filepath string
}

func NewLeaderboard() *Leaderboard {
	return &Leaderboard{
		Entries:  make([]Entry, 0, MaxLeaderboardEntries),
		filepath: DefaultScoreFile,
	}
}

func (lb *Leaderboard) SetFilepath(path string) {
	lb.filepath = path
}

func (lb *Leaderboard) Load() error {
	data, err := os.ReadFile(lb.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, lb)
}

func (lb *Leaderboard) Save() error {
	dir := filepath.Dir(lb.filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(lb, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(lb.filepath, data, 0644)
}

func (lb *Leaderboard) IsHighScore(score int) bool {
	if len(lb.Entries) < MaxLeaderboardEntries {
		return true
	}
	return score > lb.Entries[len(lb.Entries)-1].Score
}

func (lb *Leaderboard) Add(initials string, score, hunt int) {
	entry := Entry{
		Initials: initials,
		Score:    score,
		Hunt:     hunt,
	}

	lb.Entries = append(lb.Entries, entry)

	sort.Slice(lb.Entries, func(i, j int) bool {
		return lb.Entries[i].Score > lb.Entries[j].Score
	})

	if len(lb.Entries) > MaxLeaderboardEntries {
		lb.Entries = lb.Entries[:MaxLeaderboardEntries]
	}
}

func (lb *Leaderboard) GetEntries() []Entry {
	return lb.Entries
}

func (lb *Leaderboard) Clear() {
	lb.Entries = make([]Entry, 0, MaxLeaderboardEntries)
}
