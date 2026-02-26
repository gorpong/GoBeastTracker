package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"

	"beasttracker/internal/dungeon"
	"beasttracker/internal/entity"
	"beasttracker/internal/game"
	"beasttracker/internal/ui"
)

const (
	dungeonWidth  = 100
	dungeonHeight = 40
)

func main() {
	tcellScreen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create screen: %v\n", err)
		os.Exit(1)
	}

	if err := tcellScreen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize screen: %v\n", err)
		os.Exit(1)
	}

	screen, err := ui.NewScreen(tcellScreen)
	if err != nil {
		tcellScreen.Fini()
		fmt.Fprintf(os.Stderr, "Failed to create UI screen: %v\n", err)
		os.Exit(1)
	}
	defer screen.Fini()

	seed := time.Now().UnixNano()
	gameState := game.NewGame(dungeonWidth, dungeonHeight, seed)

	// Define styles
	floorStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorDarkGray)
	wallStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGray)
	exploredFloorStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorDarkSlateGray)
	exploredWallStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorDarkSlateGray)
	playerStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow).Bold(true)
	hudStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)
	messageStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	bossStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorPurple).Bold(true)
	monsterStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
	itemStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorAqua)
	inventoryStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorSilver)
	menuStyle := tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite)
	menuHeaderStyle := tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorYellow).Bold(true)

	// Main game loop
	for gameState.Running {
		screen.Clear()

		screenW, screenH := screen.Size()

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

		// Reserve 3 rows for HUD at top, 2 rows for messages/instructions at bottom
		viewHeight := screenH - 5
		cameraX := px - screenW/2
		cameraY := py - viewHeight/2

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

		// Draw dungeon tiles (offset by 3 for HUD rows)
		for screenY := 0; screenY < viewHeight; screenY++ {
			for screenX := 0; screenX < screenW; screenX++ {
				worldX := screenX + cameraX
				worldY := screenY + cameraY

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
					screen.SetCell(screenX, screenY+3, tile.Glyph(), style)
				}
			}
		}

		// Draw items (only if visible)
		for _, item := range gameState.Items {
			ix, iy := item.Position()
			if !gameState.IsVisible(ix, iy) {
				continue
			}
			itemScreenX := ix - cameraX
			itemScreenY := iy - cameraY + 3
			if itemScreenX >= 0 && itemScreenX < screenW && itemScreenY >= 3 && itemScreenY < screenH-2 {
				screen.SetCell(itemScreenX, itemScreenY, item.Glyph(), itemStyle)
			}
		}

		// Draw monsters (only if visible)
		for _, monster := range gameState.Monsters {
			if monster.Dead {
				continue
			}
			mx, my := monster.Position()
			if !gameState.IsVisible(mx, my) {
				continue
			}
			monsterScreenX := mx - cameraX
			monsterScreenY := my - cameraY + 3
			if monsterScreenX >= 0 && monsterScreenX < screenW && monsterScreenY >= 3 && monsterScreenY < screenH-2 {
				style := monsterStyle
				if monster.IsBoss {
					style = bossStyle
				}
				screen.SetCell(monsterScreenX, monsterScreenY, monster.Glyph, style)
			}
		}

		// Draw player
		playerScreenX := px - cameraX
		playerScreenY := py - cameraY + 3
		if playerScreenX >= 0 && playerScreenX < screenW && playerScreenY >= 3 && playerScreenY < screenH-2 {
			screen.SetCell(playerScreenX, playerScreenY, gameState.Player.Glyph, playerStyle)
		}

		// Draw HUD (3 rows at top)
		// Row 0: Title and HP
		title := "BeastTracker"
		screen.DrawString(0, 0, title, hudStyle)

		hpInfo := fmt.Sprintf("HP: %d/%d", gameState.Player.HP, gameState.Player.MaxHP)
		screen.DrawString(len(title)+2, 0, hpInfo, getHPStyle(gameState.Player.HP, gameState.Player.MaxHP))

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

		// Row 2: Inventory display
		drawInventoryBar(screen, gameState.Player.Inventory, 0, 2, screenW, inventoryStyle, itemStyle)

		// Draw messages (second to last row)
		messageRow := screenH - 2
		if len(gameState.Messages) > 0 {
			lastMessage := gameState.Messages[len(gameState.Messages)-1]
			if len(lastMessage) > screenW {
				lastMessage = lastMessage[:screenW]
			}
			screen.DrawString(0, messageRow, lastMessage, messageStyle)
		}

		// Draw instructions (last row) - context sensitive
		instructions := getInstructionsForMode(gameState.InputMode)
		screen.DrawString(0, screenH-1, instructions, floorStyle)

		// Draw drop menu overlay if in drop menu mode
		if gameState.InputMode == game.InputModeDropMenu {
			drawDropMenu(screen, gameState.Player.Inventory, screenW, screenH, menuStyle, menuHeaderStyle, itemStyle)
		}

		screen.Show()

		// Handle input
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			handleKeyEvent(gameState, ev)
		}
	}
}

// drawInventoryBar renders the inventory as a horizontal bar
func drawInventoryBar(screen *ui.Screen, inv *entity.Inventory, x, y, maxWidth int, defaultStyle, itemStyle tcell.Style) {
	label := "Inv: "
	screen.DrawString(x, y, label, defaultStyle)

	slotX := x + len(label)
	for slot := 1; slot <= inv.Capacity(); slot++ {
		if slotX >= maxWidth-4 {
			break
		}

		item := inv.GetSlot(slot)
		var slotStr string
		var style tcell.Style

		if item != nil {
			slotStr = fmt.Sprintf("[%d:%c]", slot, item.Glyph())
			style = itemStyle
		} else {
			slotStr = fmt.Sprintf("[%d:-]", slot)
			style = defaultStyle
		}

		screen.DrawString(slotX, y, slotStr, style)
		slotX += len(slotStr) + 1
	}
}

// drawDropMenu renders the drop menu overlay
func drawDropMenu(screen *ui.Screen, inv *entity.Inventory, screenW, screenH int, menuStyle, headerStyle, itemStyle tcell.Style) {
	// Menu dimensions
	menuWidth := 30
	menuHeight := inv.Capacity() + 4
	menuX := (screenW - menuWidth) / 2
	menuY := (screenH - menuHeight) / 2

	// Draw menu background
	for dy := 0; dy < menuHeight; dy++ {
		for dx := 0; dx < menuWidth; dx++ {
			screen.SetCell(menuX+dx, menuY+dy, ' ', menuStyle)
		}
	}

	// Draw header
	header := "== Drop Item =="
	headerX := menuX + (menuWidth-len(header))/2
	screen.DrawString(headerX, menuY+1, header, headerStyle)

	// Draw inventory items
	for slot := 1; slot <= inv.Capacity(); slot++ {
		item := inv.GetSlot(slot)
		var line string
		var style tcell.Style

		if item != nil {
			line = fmt.Sprintf(" %d. %s", slot, item.Name())
			style = itemStyle
		} else {
			line = fmt.Sprintf(" %d. (empty)", slot)
			style = menuStyle
		}

		screen.DrawString(menuX+2, menuY+2+slot, line, style)
	}

	// Draw instructions
	instrLine := "Press 1-9 or ESC to cancel"
	instrX := menuX + (menuWidth-len(instrLine))/2
	screen.DrawString(instrX, menuY+menuHeight-1, instrLine, menuStyle)
}

// getInstructionsForMode returns context-sensitive instructions
func getInstructionsForMode(mode game.InputMode) string {
	switch mode {
	case game.InputModeDropping:
		return "Drop mode: Press 1-9 to drop, x/i for menu, other key to cancel"
	case game.InputModeDropMenu:
		return "Drop menu: Press 1-9 to drop, ESC to cancel"
	case game.InputModeInventory:
		return "Inventory: Press 1-9 to use, ESC to close"
	default:
		return "Move: arrows/hjkl/wasd | Items: 1-9 use, i inv, x drop | Quit: q/ESC"
	}
}

// handleKeyEvent processes keyboard input based on current game mode
func handleKeyEvent(gameState *game.Game, ev *tcell.EventKey) {
	switch gameState.InputMode {
	case game.InputModeDropping:
		handleDropModeKey(gameState, ev)
	case game.InputModeDropMenu:
		handleDropMenuKey(gameState, ev)
	case game.InputModeInventory:
		handleInventoryKey(gameState, ev)
	default:
		handleNormalModeKey(gameState, ev)
	}
}

// handleNormalModeKey processes input in normal gameplay mode
func handleNormalModeKey(gameState *game.Game, ev *tcell.EventKey) {
	action := ui.ParseAction(ev.Key(), ev.Rune())
	dir := ui.ParseDirection(ev.Key(), ev.Rune())

	// Handle item use with number keys
	if action == ui.ActionUseItem {
		slot, ok := ui.ParseSlotNumber(ev.Rune())
		if ok {
			gameState.UseItemInSlot(slot)
			return
		}
	}

	gameState.HandleInput(action, dir)
}

// handleDropModeKey processes input while in drop mode (waiting for slot)
func handleDropModeKey(gameState *game.Game, ev *tcell.EventKey) {
	r := ev.Rune()

	// ESC cancels drop mode
	if ev.Key() == tcell.KeyEscape {
		gameState.InputMode = game.InputModeNormal
		return
	}

	gameState.HandleDropModeInput(r)
}

// handleDropMenuKey processes input while drop menu is displayed
func handleDropMenuKey(gameState *game.Game, ev *tcell.EventKey) {
	// ESC closes menu
	if ev.Key() == tcell.KeyEscape {
		gameState.InputMode = game.InputModeNormal
		return
	}

	// Number keys drop item
	slot, ok := ui.ParseSlotNumber(ev.Rune())
	if ok {
		gameState.HandleDropModeInput(ev.Rune())
		// HandleDropModeInput already sets mode to Normal after drop
		return
	}

	// Enter does nothing without selection (could add cursor later)
	_ = slot
}

// handleInventoryKey processes input while inventory is displayed
func handleInventoryKey(gameState *game.Game, ev *tcell.EventKey) {
	// ESC or i closes inventory
	if ev.Key() == tcell.KeyEscape || ev.Rune() == 'i' || ev.Rune() == 'I' {
		gameState.InputMode = game.InputModeNormal
		return
	}

	// Number keys use item
	slot, ok := ui.ParseSlotNumber(ev.Rune())
	if ok {
		gameState.UseItemInSlot(slot)
		return
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
