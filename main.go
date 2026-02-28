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
	materialStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorOrange)
	rareMaterialStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGold).Bold(true)
	inventoryStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorSilver)
	menuStyle := tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite)
	menuHeaderStyle := tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorYellow).Bold(true)
	craftableStyle := tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorGreen)
	uncraftableStyle := tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorDarkGray)
	selectedStyle := tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite).Bold(true)

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

		for _, material := range gameState.Materials {
			matX, matY := material.Position()
			if !gameState.IsVisible(matX, matY) {
				continue
			}
			matScreenX := matX - cameraX
			matScreenY := matY - cameraY + 3
			if matScreenX >= 0 && matScreenX < screenW && matScreenY >= 3 && matScreenY < screenH-2 {
				style := materialStyle
				if material.Type.IsRare() {
					style = rareMaterialStyle
				}
				screen.SetCell(matScreenX, matScreenY, material.Glyph(), style)
			}
		}

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

		playerScreenX := px - cameraX
		playerScreenY := py - cameraY + 3
		if playerScreenX >= 0 && playerScreenX < screenW && playerScreenY >= 3 && playerScreenY < screenH-2 {
			screen.SetCell(playerScreenX, playerScreenY, gameState.Player.Glyph, playerStyle)
		}

		title := "BeastTracker"
		screen.DrawString(0, 0, title, hudStyle)

		hpInfo := fmt.Sprintf("HP:%d/%d", gameState.Player.HP, gameState.Player.EffectiveMaxHP())
		screen.DrawString(len(title)+2, 0, hpInfo, getHPStyle(gameState.Player.HP, gameState.Player.EffectiveMaxHP()))

		boss := gameState.GetBoss()
		if boss != nil && !boss.Dead {
			bossInfo := fmt.Sprintf("Target: %s HP:%d/%d", boss.Name, boss.HP, boss.MaxHP)
			screen.DrawString(screenW-len(bossInfo)-1, 0, bossInfo, bossStyle)
		}

		statsInfo := fmt.Sprintf("ATK:%d DEF:%d", gameState.Player.EffectiveAttack(), gameState.Player.EffectiveDefense())
		screen.DrawString(0, 1, statsInfo, hudStyle)

		visibleMonsters := gameState.GetVisibleMonsters()
		monsterInfoX := len(statsInfo) + 2
		for _, monster := range visibleMonsters {
			info := fmt.Sprintf("%s:%d/%d ", monster.Name, monster.HP, monster.MaxHP)
			if monsterInfoX+len(info) < screenW-15 {
				screen.DrawString(monsterInfoX, 1, info, monsterStyle)
				monsterInfoX += len(info)
			}
		}

		posInfo := fmt.Sprintf("Pos:(%d,%d)", px, py)
		screen.DrawString(screenW-len(posInfo)-1, 1, posInfo, hudStyle)

		drawInventoryBar(screen, gameState.Player.Inventory, 0, 2, screenW/2, inventoryStyle, itemStyle)
		drawMaterialPouch(screen, gameState.Player.MaterialPouch, screenW/2, 2, screenW/2, inventoryStyle, materialStyle)

		messageRow := screenH - 2
		if len(gameState.Messages) > 0 {
			lastMessage := gameState.Messages[len(gameState.Messages)-1]
			if len(lastMessage) > screenW {
				lastMessage = lastMessage[:screenW]
			}
			screen.DrawString(0, messageRow, lastMessage, messageStyle)
		}

		instructions := getInstructionsForMode(gameState.InputMode)
		screen.DrawString(0, screenH-1, instructions, floorStyle)

		if gameState.InputMode == game.InputModeDropMenu {
			drawDropMenu(screen, gameState.Player.Inventory, screenW, screenH, menuStyle, menuHeaderStyle, itemStyle)
		}

		if gameState.InputMode == game.InputModeCrafting {
			drawCraftingMenu(screen, gameState, screenW, screenH, menuStyle, menuHeaderStyle, craftableStyle, uncraftableStyle, selectedStyle)
		}

		screen.Show()

		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			handleKeyEvent(gameState, ev)
		}
	}
}

func drawInventoryBar(screen *ui.Screen, inv *entity.Inventory, x, y, maxWidth int, defaultStyle, itemStyle tcell.Style) {
	label := "Inv:"
	screen.DrawString(x, y, label, defaultStyle)

	slotX := x + len(label)
	for slot := 1; slot <= inv.Capacity(); slot++ {
		if slotX >= x+maxWidth-4 {
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

func drawMaterialPouch(screen *ui.Screen, pouch *entity.MaterialPouch, x, y, maxWidth int, defaultStyle, materialStyle tcell.Style) {
	label := "Mat:"
	screen.DrawString(x, y, label, defaultStyle)

	materials := pouch.AllMaterials()
	if len(materials) == 0 {
		screen.DrawString(x+len(label)+1, y, "(empty)", defaultStyle)
		return
	}

	matX := x + len(label) + 1
	for _, matType := range materials {
		count := pouch.Count(matType)
		info := fmt.Sprintf("%c:%d ", matType.Glyph(), count)
		if matX+len(info) >= x+maxWidth {
			break
		}
		screen.DrawString(matX, y, info, materialStyle)
		matX += len(info)
	}
}

func drawDropMenu(screen *ui.Screen, inv *entity.Inventory, screenW, screenH int, menuStyle, headerStyle, itemStyle tcell.Style) {
	menuWidth := 30
	menuHeight := inv.Capacity() + 4
	menuX := (screenW - menuWidth) / 2
	menuY := (screenH - menuHeight) / 2

	for dy := 0; dy < menuHeight; dy++ {
		for dx := 0; dx < menuWidth; dx++ {
			screen.SetCell(menuX+dx, menuY+dy, ' ', menuStyle)
		}
	}

	header := "== Drop Item =="
	headerX := menuX + (menuWidth-len(header))/2
	screen.DrawString(headerX, menuY+1, header, headerStyle)

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

	instrLine := "Press 1-9 or ESC to cancel"
	instrX := menuX + (menuWidth-len(instrLine))/2
	screen.DrawString(instrX, menuY+menuHeight-1, instrLine, menuStyle)
}

func drawCraftingMenu(screen *ui.Screen, gameState *game.Game, screenW, screenH int, menuStyle, headerStyle, craftableStyle, uncraftableStyle, selectedStyle tcell.Style) {
	recipes := gameState.GetAllRecipes()

	menuWidth := 50
	menuHeight := len(recipes) + 8
	if menuHeight > screenH-4 {
		menuHeight = screenH - 4
	}
	menuX := (screenW - menuWidth) / 2
	menuY := (screenH - menuHeight) / 2

	for dy := 0; dy < menuHeight; dy++ {
		for dx := 0; dx < menuWidth; dx++ {
			screen.SetCell(menuX+dx, menuY+dy, ' ', menuStyle)
		}
	}

	header := "== Crafting =="
	headerX := menuX + (menuWidth-len(header))/2
	screen.DrawString(headerX, menuY+1, header, headerStyle)

	visibleRecipes := menuHeight - 6
	startIdx := 0
	if gameState.CraftingCursor >= visibleRecipes {
		startIdx = gameState.CraftingCursor - visibleRecipes + 1
	}

	for i := 0; i < visibleRecipes && startIdx+i < len(recipes); i++ {
		recipeIdx := startIdx + i
		recipe := recipes[recipeIdx]
		canCraft := recipe.CanCraft(gameState.Player.MaterialPouch)

		var style tcell.Style
		if recipeIdx == gameState.CraftingCursor {
			style = selectedStyle
		} else if canCraft {
			style = craftableStyle
		} else {
			style = uncraftableStyle
		}

		prefix := "  "
		if recipeIdx == gameState.CraftingCursor {
			prefix = "> "
		}

		line := fmt.Sprintf("%s%s (%s)", prefix, recipe.Name, recipe.Result.StatsString())
		if len(line) > menuWidth-4 {
			line = line[:menuWidth-4]
		}
		screen.DrawString(menuX+2, menuY+3+i, line, style)
	}

	selectedRecipe := recipes[gameState.CraftingCursor]
	ingredientLine := "Needs: " + selectedRecipe.IngredientsString()
	if len(ingredientLine) > menuWidth-4 {
		ingredientLine = ingredientLine[:menuWidth-4]
	}
	screen.DrawString(menuX+2, menuY+menuHeight-3, ingredientLine, menuStyle)

	instrLine := "↑↓:Select Enter:Craft ESC:Close"
	instrX := menuX + (menuWidth-len(instrLine))/2
	screen.DrawString(instrX, menuY+menuHeight-1, instrLine, menuStyle)
}

func getInstructionsForMode(mode game.InputMode) string {
	switch mode {
	case game.InputModeDropping:
		return "Drop mode: Press 1-9 to drop, x/i for menu, other key to cancel"
	case game.InputModeDropMenu:
		return "Drop menu: Press 1-9 to drop, ESC to cancel"
	case game.InputModeInventory:
		return "Inventory: Press 1-9 to use, ESC to close"
	case game.InputModeCrafting:
		return "Crafting: ↑↓ select, Enter craft, ESC close"
	default:
		return "Move:arrows/hjkl | 1-9:use | x:drop | c:craft | q:quit"
	}
}

func handleKeyEvent(gameState *game.Game, ev *tcell.EventKey) {
	switch gameState.InputMode {
	case game.InputModeDropping:
		handleDropModeKey(gameState, ev)
	case game.InputModeDropMenu:
		handleDropMenuKey(gameState, ev)
	case game.InputModeInventory:
		handleInventoryKey(gameState, ev)
	case game.InputModeCrafting:
		handleCraftingKey(gameState, ev)
	default:
		handleNormalModeKey(gameState, ev)
	}
}

func handleNormalModeKey(gameState *game.Game, ev *tcell.EventKey) {
	action := ui.ParseAction(ev.Key(), ev.Rune())
	dir := ui.ParseDirection(ev.Key(), ev.Rune())

	if action == ui.ActionUseItem {
		slot, ok := ui.ParseSlotNumber(ev.Rune())
		if ok {
			gameState.UseItemInSlot(slot)
			return
		}
	}

	gameState.HandleInput(action, dir)
}

func handleDropModeKey(gameState *game.Game, ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyEscape {
		gameState.InputMode = game.InputModeNormal
		return
	}

	gameState.HandleDropModeInput(ev.Rune())
}

func handleDropMenuKey(gameState *game.Game, ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyEscape {
		gameState.InputMode = game.InputModeNormal
		return
	}

	slot, ok := ui.ParseSlotNumber(ev.Rune())
	if ok {
		gameState.HandleDropModeInput(ev.Rune())
		return
	}
	_ = slot
}

func handleInventoryKey(gameState *game.Game, ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyEscape || ev.Rune() == 'i' || ev.Rune() == 'I' {
		gameState.InputMode = game.InputModeNormal
		return
	}

	slot, ok := ui.ParseSlotNumber(ev.Rune())
	if ok {
		gameState.UseItemInSlot(slot)
		return
	}
}

func handleCraftingKey(gameState *game.Game, ev *tcell.EventKey) {
	recipes := gameState.GetAllRecipes()

	switch ev.Key() {
	case tcell.KeyEscape:
		gameState.InputMode = game.InputModeNormal
	case tcell.KeyUp:
		if gameState.CraftingCursor > 0 {
			gameState.CraftingCursor--
		}
	case tcell.KeyDown:
		if gameState.CraftingCursor < len(recipes)-1 {
			gameState.CraftingCursor++
		}
	case tcell.KeyEnter:
		selectedRecipe := recipes[gameState.CraftingCursor]
		gameState.CraftRecipe(selectedRecipe.Name)
	}

	if ev.Key() == tcell.KeyRune {
		switch ev.Rune() {
		case 'k', 'K':
			if gameState.CraftingCursor > 0 {
				gameState.CraftingCursor--
			}
		case 'j', 'J':
			if gameState.CraftingCursor < len(recipes)-1 {
				gameState.CraftingCursor++
			}
		case 'c', 'C':
			gameState.InputMode = game.InputModeNormal
		}
	}
}

func getHPStyle(current, max int) tcell.Style {
	percent := float64(current) / float64(max)
	if percent > 0.6 {
		return tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)
	} else if percent > 0.3 {
		return tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow)
	}
	return tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
}

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
		fmt.Sprintf("Final HP: %d/%d", gameState.Player.HP, gameState.Player.EffectiveMaxHP()),
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

func waitForKeyPress(screen *ui.Screen) {
	for {
		ev := screen.PollEvent()
		switch ev.(type) {
		case *tcell.EventKey:
			return
		}
	}
}
