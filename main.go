package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"

	"beasttracker/internal/dungeon"
	"beasttracker/internal/game"
	"beasttracker/internal/ui"
)

const (
	// Dungeon dimensions (larger than screen for scrolling)
	dungeonWidth  = 100
	dungeonHeight = 40
)

func main() {
	// Initialize tcell screen
	tcellScreen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create screen: %v\n", err)
		os.Exit(1)
	}

	if err := tcellScreen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize screen: %v\n", err)
		os.Exit(1)
	}

	// Wrap in our Screen abstraction
	screen, err := ui.NewScreen(tcellScreen)
	if err != nil {
		tcellScreen.Fini()
		fmt.Fprintf(os.Stderr, "Failed to create UI screen: %v\n", err)
		os.Exit(1)
	}
	defer screen.Fini()

	// Create game with current time as seed for variety
	seed := time.Now().UnixNano()
	gameState := game.NewGame(dungeonWidth, dungeonHeight, seed)

	// Define styles for visible tiles
	floorStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorDarkGray)
	wallStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGray)
	// Define styles for explored but not visible tiles (dimmed)
	exploredFloorStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorDarkSlateGray)
	exploredWallStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorDarkSlateGray)
	// Other styles
	playerStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow).Bold(true)
	hudStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)

	// Main game loop
	for gameState.Running {
		screen.Clear()

		screenW, screenH := screen.Size()
		px, py := gameState.Player.Position()

		// Calculate camera offset (center player on screen)
		// Reserve 1 row for HUD at top, 1 row for instructions at bottom
		viewHeight := screenH - 2
		cameraX := px - screenW/2
		cameraY := py - viewHeight/2

		// Clamp camera to dungeon bounds
		if cameraX < 0 {
			cameraX = 0
		}
		if cameraY < 0 {
			cameraY = 0
		}
		if cameraX > dungeonWidth-screenW {
			cameraX = dungeonWidth - screenW
		}
		if cameraY > dungeonHeight-viewHeight {
			cameraY = dungeonHeight - viewHeight
		}

		// Draw dungeon tiles (offset by 1 for HUD row)
		for screenY := 0; screenY < viewHeight; screenY++ {
			for screenX := 0; screenX < screenW; screenX++ {
				worldX := screenX + cameraX
				worldY := screenY + cameraY

				// Only draw tiles that have been explored
				if !gameState.IsExplored(worldX, worldY) {
					continue
				}

				tile := gameState.Dungeon.GetTile(worldX, worldY)
				if tile != nil {
					var style tcell.Style
					isVisible := gameState.IsVisible(worldX, worldY)

					switch tile.Type {
					case dungeon.TileFloor:
						if isVisible {
							style = floorStyle
						} else {
							style = exploredFloorStyle
						}
					case dungeon.TileWall:
						if isVisible {
							style = wallStyle
						} else {
							style = exploredWallStyle
						}
					default:
						if isVisible {
							style = floorStyle
						} else {
							style = exploredFloorStyle
						}
					}
					screen.SetCell(screenX, screenY+1, tile.Glyph(), style)
				}
			}
		}

		// Draw monsters (only if visible, relative to camera, offset by 1 for HUD)
		monsterStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
		for _, monster := range gameState.Monsters {
			if monster.Dead {
				continue
			}
			mx, my := monster.Position()
			// Only draw monsters that are visible to the player
			if !gameState.IsVisible(mx, my) {
				continue
			}
			monsterScreenX := mx - cameraX
			monsterScreenY := my - cameraY + 1
			if monsterScreenX >= 0 && monsterScreenX < screenW && monsterScreenY >= 1 && monsterScreenY < screenH-1 {
				screen.SetCell(monsterScreenX, monsterScreenY, monster.Glyph, monsterStyle)
			}
		}

		// Draw player (relative to camera, offset by 1 for HUD)
		playerScreenX := px - cameraX
		playerScreenY := py - cameraY + 1
		if playerScreenX >= 0 && playerScreenX < screenW && playerScreenY >= 1 && playerScreenY < screenH-1 {
			screen.SetCell(playerScreenX, playerScreenY, gameState.Player.Glyph, playerStyle)
		}

		// Draw HUD at top
		title := "BeastTracker - Phase 5"
		screen.DrawString(0, 0, title, hudStyle)

		// Draw position info
		posInfo := fmt.Sprintf("Pos: (%d, %d) Monsters: %d", px, py, len(gameState.Monsters))
		screen.DrawString(screenW-len(posInfo)-1, 0, posInfo, hudStyle)

		// Draw instructions at bottom
		instructions := "Move: arrows/hjkl/wasd | Quit: q/ESC"
		screen.DrawString(0, screenH-1, instructions, floorStyle)

		screen.Show()

		// Handle input
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			action := ui.ParseAction(ev.Key(), ev.Rune())
			dir := ui.ParseDirection(ev.Key(), ev.Rune())
			gameState.HandleInput(action, dir)
		}
	}
}
