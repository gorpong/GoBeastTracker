package game

import (
	"beasttracker/internal/entity"
	"beasttracker/internal/ui"
)

// Game holds all game state
type Game struct {
	Width   int
	Height  int
	Player  *entity.Player
	Running bool
}

// NewGame creates a new game with the specified dimensions
func NewGame(width, height int) *Game {
	// Spawn player in center of the area
	playerX := width / 2
	playerY := height / 2

	return &Game{
		Width:   width,
		Height:  height,
		Player:  entity.NewPlayer(playerX, playerY),
		Running: true,
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
// Movement is blocked if it would go out of bounds.
func (g *Game) tryMovePlayer(dir ui.Direction) {
	dx, dy := dir.Delta()
	newX := g.Player.X + dx
	newY := g.Player.Y + dy

	// Check bounds
	if newX < 0 || newX >= g.Width {
		return
	}
	if newY < 0 || newY >= g.Height {
		return
	}

	g.Player.SetPosition(newX, newY)
}
