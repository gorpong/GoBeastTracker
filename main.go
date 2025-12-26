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
	messageStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	bossStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorPurple).Bold(true)

	// Main game loop
	for gameState.Running {
		screen.Clear()

		screenW, screenH := screen.Size()

		// Check for game over or victory states
		if gameState.GameState == game.StateGameOver {
			drawGameOver(screen, screenW, screenH)
			screen.Show()
			waitForKeyPress(screen)
			break
		}
		if gameState.GameState == game.StateVictory {
			drawVictory(screen, screenW, screenH, gameState)
			screen.Show()
			waitForKeyPress(screen)
			break
		}

		px, py := gameState.Player.Position()

		// Calculate camera offset (center player on screen)
		// Reserve 2 rows for HUD at top, 2 rows for messages/instructions at bottom
		viewHeight := screenH - 4
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

		// Draw dungeon tiles (offset by 2 for HUD rows)
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
					screen.SetCell(screenX, screenY+2, tile.Glyph(), style)
				}
			}
		}

		// Draw monsters (only if visible, relative to camera, offset by 2 for HUD)
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
			monsterScreenY := my - cameraY + 2
			if monsterScreenX >= 0 && monsterScreenX < screenW && monsterScreenY >= 2 && monsterScreenY < screenH-2 {
				// Use special style for boss
				style := monsterStyle
				if monster.IsBoss {
					style = bossStyle
				}
				screen.SetCell(monsterScreenX, monsterScreenY, monster.Glyph, style)
			}
		}

		// Draw player (relative to camera, offset by 2 for HUD)
		playerScreenX := px - cameraX
		playerScreenY := py - cameraY + 2
		if playerScreenX >= 0 && playerScreenX < screenW && playerScreenY >= 2 && playerScreenY < screenH-2 {
			screen.SetCell(playerScreenX, playerScreenY, gameState.Player.Glyph, playerStyle)
		}

		// Draw HUD at top (2 rows)
		// Row 0: Title and HP
		title := "BeastTracker"
		screen.DrawString(0, 0, title, hudStyle)

		hpInfo := fmt.Sprintf("HP: %d/%d", gameState.Player.HP, gameState.Player.MaxHP)
		screen.DrawString(len(title)+2, 0, hpInfo, getHPStyle(gameState.Player.HP, gameState.Player.MaxHP))

		// Show boss info if exists
		boss := gameState.GetBoss()
		if boss != nil && !boss.Dead {
			bossInfo := fmt.Sprintf("Target: %s HP:%d/%d", boss.Name, boss.HP, boss.MaxHP)
			screen.DrawString(screenW-len(bossInfo)-1, 0, bossInfo, bossStyle)
		}

		// Row 1: ATK/DEF and position
		statsInfo := fmt.Sprintf("ATK:%d DEF:%d", gameState.Player.Attack, gameState.Player.Defense)
		screen.DrawString(0, 1, statsInfo, hudStyle)

		posInfo := fmt.Sprintf("Pos:(%d,%d)", px, py)
		screen.DrawString(screenW-len(posInfo)-1, 1, posInfo, hudStyle)

		// Draw messages above instructions (second to last row)
		messageRow := screenH - 2
		if len(gameState.Messages) > 0 {
			lastMessage := gameState.Messages[len(gameState.Messages)-1]
			if len(lastMessage) > screenW {
				lastMessage = lastMessage[:screenW]
			}
			screen.DrawString(0, messageRow, lastMessage, messageStyle)
		}

		// Draw instructions at bottom
		instructions := "Move/Attack: arrows/hjkl/wasd | Quit: q/ESC"
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

// getHPStyle returns a style based on current HP percentage
func getHPStyle(current, max int) tcell.Style {
	percent := float64(current) / float64(max)
	if percent > 0.6 {
		return tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)
	} else if percent > 0.3 {
		return tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow)
	}
	return tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
}

// drawGameOver draws the game over screen
func drawGameOver(screen *ui.Screen, width, height int) {
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed).Bold(true)
	titleStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)

	lines := []string{
		"",
		"  ####    ###   ##   ## #####      ###  ##   ## ##### ##### ",
		" ##      ## ##  ### ### ##        ## ## ##   ## ##    ##  ##",
		" ## ### ##   ## ## # ## ####      ## ## ##   ## ####  ##### ",
		" ##  ## ####### ##   ## ##        ## ##  ## ##  ##    ##  ##",
		"  ####  ##   ## ##   ## #####      ###    ###   ##### ##  ##",
		"",
		"You have been slain!",
		"",
		"Press any key to exit...",
	}

	startY := height/2 - len(lines)/2
	for i, line := range lines {
		x := (width - len(line)) / 2
		if x < 0 {
			x = 0
		}
		if i < 6 {
			screen.DrawString(x, startY+i, line, style)
		} else {
			screen.DrawString(x, startY+i, line, titleStyle)
		}
	}
}

// drawVictory draws the victory screen
func drawVictory(screen *ui.Screen, width, height int, gameState *game.Game) {
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen).Bold(true)
	titleStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow)

	lines := []string{
		"",
		" ##   ## ####  #### #####  ###  #####  ##   ## ",
		" ##   ##  ##  ##      ##  ## ## ##  ##  ## ##  ",
		"  ## ##   ##  ##      ##  ## ## #####    ###   ",
		"   ###    ##  ##      ##  ## ## ##  ##   ##    ",
		"    #    ####  ####   ##   ###  ##  ##   ##    ",
		"",
		"You have slain the beast!",
		"",
		fmt.Sprintf("Final HP: %d/%d", gameState.Player.HP, gameState.Player.MaxHP),
		"",
		"Press any key to exit...",
	}

	startY := height/2 - len(lines)/2
	for i, line := range lines {
		x := (width - len(line)) / 2
		if x < 0 {
			x = 0
		}
		if i < 6 {
			screen.DrawString(x, startY+i, line, style)
		} else {
			screen.DrawString(x, startY+i, line, titleStyle)
		}
	}
}

// waitForKeyPress waits for any key press
func waitForKeyPress(screen *ui.Screen) {
	for {
		ev := screen.PollEvent()
		switch ev.(type) {
		case *tcell.EventKey:
			return
		}
	}
}
