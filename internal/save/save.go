package save

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"beasttracker/internal/entity"
)

const (
	MaxSaveSlots  = 10
	MaxNameLength = 30
	SaveFileExt   = ".json"
)

var (
	ErrMaxSlotsReached = errors.New("maximum save slots reached")
	ErrSaveNotFound    = errors.New("save not found")
	ErrInvalidName     = errors.New("invalid save name")
)

// SaveData contains all persistent game state between hunts
type SaveData struct {
	Name             string                      `json:"name"`
	HuntNumber       int                         `json:"hunt_number"`
	Score            int                         `json:"score"`
	SavedAt          time.Time                   `json:"saved_at"`
	EquippedWeapon   *entity.Equipment           `json:"equipped_weapon,omitempty"`
	EquippedArmor    *entity.Equipment           `json:"equipped_armor,omitempty"`
	EquippedCharm    *entity.Equipment           `json:"equipped_charm,omitempty"`
	StashedEquipment []*entity.Equipment         `json:"stashed_equipment,omitempty"`
	Materials        map[entity.MaterialType]int `json:"materials,omitempty"`
}

// NewSaveData creates a new save data instance
func NewSaveData(name string, huntNumber, score int) *SaveData {
	return &SaveData{
		Name:             name,
		HuntNumber:       huntNumber,
		Score:            score,
		SavedAt:          time.Now(),
		StashedEquipment: make([]*entity.Equipment, 0),
		Materials:        make(map[entity.MaterialType]int),
	}
}

// SaveManager handles save/load operations
type SaveManager struct {
	saveDir string
}

// NewSaveManager creates a new save manager with the specified directory
func NewSaveManager(saveDir string) *SaveManager {
	return &SaveManager{
		saveDir: saveDir,
	}
}

// Save writes save data to disk
func (sm *SaveManager) Save(data *SaveData) error {
	if !IsValidSaveName(data.Name) {
		return ErrInvalidName
	}

	existing, _ := sm.List()
	isOverwrite := false
	for _, save := range existing {
		if save.Name == data.Name {
			isOverwrite = true
			break
		}
	}

	if !isOverwrite && len(existing) >= MaxSaveSlots {
		return ErrMaxSlotsReached
	}

	if err := os.MkdirAll(sm.saveDir, 0755); err != nil {
		return err
	}

	data.SavedAt = time.Now()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	filename := sm.getSaveFilename(data.Name)
	return os.WriteFile(filename, jsonData, 0644)
}

// Load reads save data from disk
func (sm *SaveManager) Load(name string) (*SaveData, error) {
	filename := sm.getSaveFilename(name)

	jsonData, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrSaveNotFound
		}
		return nil, err
	}

	var data SaveData
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// Delete removes a save from disk
func (sm *SaveManager) Delete(name string) error {
	filename := sm.getSaveFilename(name)
	err := os.Remove(filename)
	if os.IsNotExist(err) {
		return ErrSaveNotFound
	}
	return err
}

// List returns all available saves
func (sm *SaveManager) List() ([]*SaveData, error) {
	saves := make([]*SaveData, 0)

	entries, err := os.ReadDir(sm.saveDir)
	if err != nil {
		if os.IsNotExist(err) {
			return saves, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), SaveFileExt) {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), SaveFileExt)
		data, err := sm.Load(name)
		if err == nil {
			saves = append(saves, data)
		}
	}

	return saves, nil
}

// SlotCount returns the number of used save slots
func (sm *SaveManager) SlotCount() int {
	saves, _ := sm.List()
	return len(saves)
}

// HasRoom returns true if there's room for another save
func (sm *SaveManager) HasRoom() bool {
	return sm.SlotCount() < MaxSaveSlots
}

// Exists returns true if a save with the given name exists
func (sm *SaveManager) Exists(name string) bool {
	filename := sm.getSaveFilename(name)
	_, err := os.Stat(filename)
	return err == nil
}

func (sm *SaveManager) getSaveFilename(name string) string {
	safeName := sanitizeFilename(name)
	return filepath.Join(sm.saveDir, safeName+SaveFileExt)
}

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return replacer.Replace(name)
}

// IsValidSaveName validates a save name
func IsValidSaveName(name string) bool {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) == 0 {
		return false
	}
	if len(trimmed) > MaxNameLength {
		return false
	}
	return true
}
