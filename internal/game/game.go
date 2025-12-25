package game

import (
	"beasttracker/internal/dungeon"
	"beasttracker/internal/entity"
	"beasttracker/internal/ui"
)

// Game holds all game state
type Game struct {
	Width   int
	Height  int
	Player  *entity.Player
	Dungeon *dungeon.Dungeon
	Running bool
	Seed    int64
}

// NewGame creates a new game with the specified dimensions and RNG seed
func NewGame(width, height int, seed int64) *Game {
	// Generate dungeon
	d := dungeon.GenerateDungeon(width, height, seed)

	// Spawn player in the center of the first room
	var playerX, playerY int
	if len(d.Rooms) > 0 {
		playerX, playerY = d.Rooms[0].Center()
	} else {
		// Fallback to center if no rooms (shouldn't happen)
		playerX = width / 2
		playerY = height / 2
	}

	return &Game{
		Width:   width,
		Height:  height,
		Player:  entity.NewPlayer(playerX, playerY),
		Dungeon: d,
		Running: true,
		Seed:    seed,
	}
}

// HandleInput processes player input and updates game state
func (g *Game) HandleInput(action ui.Action, dir ui.Direction) {
	switch action {
	case ui.ActionQuit:
		g.Running = false
	case ui.ActionMove:
		g.tryMovePlayer(dir)
	}
}

// tryMovePlayer attempts to move the player in the given direction.
// Movement is blocked if it would hit a wall or go out of bounds.
func (g *Game) tryMovePlayer(dir ui.Direction) {
	dx, dy := dir.Delta()
	newX := g.Player.X + dx
	newY := g.Player.Y + dy

	// Check if target position is walkable (includes bounds check)
	if !g.Dungeon.IsWalkable(newX, newY) {
		return
	}

	g.Player.SetPosition(newX, newY)
}
