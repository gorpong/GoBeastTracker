package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"

	"beasttracker/internal/dungeon"
	"beasttracker/internal/entity"
	"beasttracker/internal/game"
	"beasttracker/internal/score"
	"beasttracker/internal/ui"
)

const (
	dungeonWidth  = 100
	dungeonHeight = 40
)

type AppState int

const (
	AppStateSplash AppState = iota
	AppStatePlaying
	AppStateVictory
	AppStateGameOver
	AppStateInitials
)

type App struct {
	screen       *ui.Screen
	tcellScreen  tcell.Screen
	gameState    *game.Game
	leaderboard  *score.Leaderboard
	appState     AppState
	initials     string
	totalScore   int
	currentHunt  int
	styles       styles
}

type styles struct {
	floor           tcell.Style
	wall            tcell.Style
	exploredFloor   tcell.Style
	exploredWall    tcell.Style
	player          tcell.Style
	hud             tcell.Style
	message         tcell.Style
	boss            tcell.Style
	monster         tcell.Style
	item            tcell.Style
	material        tcell.Style
	rareMaterial    tcell.Style
	inventory       tcell.Style
	menu            tcell.Style
	menuHeader      tcell.Style
	craftable       tcell.Style
	uncraftable     tcell.Style
	selected        tcell.Style
	title           tcell.Style
	subtitle        tcell.Style
}

func newStyles() styles {
	return styles{
		floor:         tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorDarkGray),
		wall:          tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGray),
		exploredFloor: tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorDarkSlateGray),
		exploredWall:  tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorDarkSlateGray),
		player:        tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow).Bold(true),
		hud:           tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen),
		message:       tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite),
		boss:          tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorPurple).Bold(true),
		monster:       tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed),
		item:          tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorAqua),
		material:      tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorOrange),
		rareMaterial:  tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGold).Bold(true),
		inventory:     tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorSilver),
		menu:          tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorWhite),
		menuHeader:    tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorYellow).Bold(true),
		craftable:     tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorGreen),
		uncraftable:   tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorDarkGray),
		selected:      tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite).Bold(true),
		title:         tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow).Bold(true),
		subtitle:      tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite),
	}
}

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

	leaderboard := score.NewLeaderboard()
	leaderboard.Load()

	app := &App{
		screen:      screen,
		tcellScreen: tcellScreen,
		leaderboard: leaderboard,
		appState:    AppStateSplash,
		styles:      newStyles(),
		currentHunt: 1,
		totalScore:  0,
	}

	app.run()
}

func (a *App) run() {
	for {
		a.screen.Clear()
		screenW, screenH := a.screen.Size()

		switch a.appState {
		case AppStateSplash:
			a.drawSplash(screenW, screenH)
		case AppStatePlaying:
			a.drawGame(screenW, screenH)
		case AppStateVictory:
			a.drawVictory(screenW, screenH)
		case AppStateGameOver:
			a.drawGameOver(screenW, screenH)
		case AppStateInitials:
			a.drawInitialsEntry(screenW, screenH)
		}

		a.screen.Show()

		ev := a.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			a.screen.Sync()
		case *tcell.EventKey:
			if a.handleKeyEvent(ev) {
				return
			}
		}
	}
}

func (a *App) handleKeyEvent(ev *tcell.EventKey) bool {
	switch a.appState {
	case AppStateSplash:
		return a.handleSplashKey(ev)
	case AppStatePlaying:
		return a.handlePlayingKey(ev)
	case AppStateVictory:
		return a.handleVictoryKey(ev)
	case AppStateGameOver:
		return a.handleGameOverKey(ev)
	case AppStateInitials:
		return a.handleInitialsKey(ev)
	}
	return false
}

func (a *App) handleSplashKey(ev *tcell.EventKey) bool {
	if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' || ev.Rune() == 'Q' {
		return true
	}

	if ev.Key() == tcell.KeyEnter || ev.Rune() == ' ' {
		a.startNewGame()
	}

	return false
}

func (a *App) handlePlayingKey(ev *tcell.EventKey) bool {
	if a.gameState == nil {
		return false
	}

	switch a.gameState.InputMode {
	case game.InputModeDropping:
		a.handleDropModeKey(ev)
	case game.InputModeDropMenu:
		a.handleDropMenuKey(ev)
	case game.InputModeInventory:
		a.handleInventoryKey(ev)
	case game.InputModeCrafting:
		a.handleCraftingKey(ev)
	default:
		a.handleNormalModeKey(ev)
	}

	if a.gameState.GameState == game.StateVictory {
		a.totalScore += a.gameState.Score
		a.appState = AppStateVictory
	} else if a.gameState.GameState == game.StateGameOver {
		a.totalScore += a.gameState.Score
		a.appState = AppStateGameOver
	}

	if !a.gameState.Running {
		return true
	}

	return false
}

func (a *App) handleVictoryKey(ev *tcell.EventKey) bool {
	if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' || ev.Rune() == 'Q' {
		a.checkHighScore()
		return false
	}

	if ev.Rune() == 'n' || ev.Rune() == 'N' {
		a.startNextHunt()
		return false
	}

	return false
}

func (a *App) handleGameOverKey(ev *tcell.EventKey) bool {
	a.checkHighScore()
	return false
}

func (a *App) handleInitialsKey(ev *tcell.EventKey) bool {
	if ev.Key() == tcell.KeyEnter && len(a.initials) > 0 {
		a.leaderboard.Add(a.initials, a.totalScore, a.currentHunt)
		a.leaderboard.Save()
		a.resetToSplash()
		return false
	}

	if ev.Key() == tcell.KeyEscape {
		a.resetToSplash()
		return false
	}

	if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
		if len(a.initials) > 0 {
			a.initials = a.initials[:len(a.initials)-1]
		}
		return false
	}

	if ui.IsLetter(ev.Rune()) && len(a.initials) < 3 {
		a.initials += string(ui.ToUpper(ev.Rune()))
	}

	return false
}

func (a *App) checkHighScore() {
	if a.leaderboard.IsHighScore(a.totalScore) {
		a.initials = ""
		a.appState = AppStateInitials
	} else {
		a.resetToSplash()
	}
}

func (a *App) resetToSplash() {
	a.appState = AppStateSplash
	a.gameState = nil
	a.currentHunt = 1
	a.totalScore = 0
}

func (a *App) startNewGame() {
	seed := time.Now().UnixNano()
	a.gameState = game.NewGame(dungeonWidth, dungeonHeight, seed)
	a.currentHunt = 1
	a.totalScore = 0
	a.appState = AppStatePlaying
}

func (a *App) startNextHunt() {
	seed := time.Now().UnixNano()
	a.currentHunt++
	a.gameState = game.NewGameWithHunt(dungeonWidth, dungeonHeight, seed, a.currentHunt, a.gameState.Player)
	a.appState = AppStatePlaying
}

func (a *App) handleNormalModeKey(ev *tcell.EventKey) {
	action := ui.ParseAction(ev.Key(), ev.Rune())
	dir := ui.ParseDirection(ev.Key(), ev.Rune())

	if action == ui.ActionUseItem {
		slot, ok := ui.ParseSlotNumber(ev.Rune())
		if ok {
			a.gameState.UseItemInSlot(slot)
			return
		}
	}

	a.gameState.HandleInput(action, dir)
}

func (a *App) handleDropModeKey(ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyEscape {
		a.gameState.InputMode = game.InputModeNormal
		return
	}

	a.gameState.HandleDropModeInput(ev.Rune())
}

func (a *App) handleDropMenuKey(ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyEscape {
		a.gameState.InputMode = game.InputModeNormal
		return
	}

	slot, ok := ui.ParseSlotNumber(ev.Rune())
	if ok {
		a.gameState.HandleDropModeInput(ev.Rune())
		return
	}
	_ = slot
}

func (a *App) handleInventoryKey(ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyEscape || ev.Rune() == 'i' || ev.Rune() == 'I' {
		a.gameState.InputMode = game.InputModeNormal
		return
	}

	slot, ok := ui.ParseSlotNumber(ev.Rune())
	if ok {
		a.gameState.UseItemInSlot(slot)
		return
	}
}

func (a *App) handleCraftingKey(ev *tcell.EventKey) {
	recipes := a.gameState.GetAllRecipes()

	switch ev.Key() {
	case tcell.KeyEscape:
		a.gameState.InputMode = game.InputModeNormal
	case tcell.KeyUp:
		if a.gameState.CraftingCursor > 0 {
			a.gameState.CraftingCursor--
		}
	case tcell.KeyDown:
		if a.gameState.CraftingCursor < len(recipes)-1 {
			a.gameState.CraftingCursor++
		}
	case tcell.KeyEnter:
		selectedRecipe := recipes[a.gameState.CraftingCursor]
		a.gameState.CraftRecipe(selectedRecipe.Name)
	}

	if ev.Key() == tcell.KeyRune {
		switch ev.Rune() {
		case 'k', 'K':
			if a.gameState.CraftingCursor > 0 {
				a.gameState.CraftingCursor--
			}
		case 'j', 'J':
			if a.gameState.CraftingCursor < len(recipes)-1 {
				a.gameState.CraftingCursor++
			}
		case 'c', 'C':
			a.gameState.InputMode = game.InputModeNormal
		}
	}
}

func (a *App) drawSplash(screenW, screenH int) {
	title := []string{
		"  ____                 _  _____               _             ",
		" | __ )  ___  __ _ ___| ||_   _| __ __ _  ___| | _____ _ __ ",
		" |  _ \\ / _ \\/ _` / __| __|| || '__/ _` |/ __| |/ / _ \\ '__|",
		" | |_) |  __/ (_| \\__ \\ |_ | || | | (_| | (__|   <  __/ |   ",
		" |____/ \\___|\\__,_|___/\\__||_||_|  \\__,_|\\___|_|\\_\\___|_|   ",
	}

	startY := 3
	for i, line := range title {
		x := (screenW - len(line)) / 2
		if x < 0 {
			x = 0
		}
		a.screen.DrawString(x, startY+i, line, a.styles.title)
	}

	subtitle := "A Monster Hunter-Inspired Roguelike"
	a.screen.DrawString((screenW-len(subtitle))/2, startY+7, subtitle, a.styles.subtitle)

	instructions := []string{
		"",
		"Hunt down powerful boss monsters in procedurally generated dungeons!",
		"",
		"Controls:",
		"  Move: Arrow keys / HJKL / WASD",
		"  Use Item: 1-9",
		"  Drop Item: X + number",
		"  Craft: C",
		"  Quit: Q / ESC",
		"",
		"Press ENTER or SPACE to start a new hunt!",
	}

	instY := startY + 10
	for i, line := range instructions {
		x := (screenW - len(line)) / 2
		if x < 0 {
			x = 0
		}
		a.screen.DrawString(x, instY+i, line, a.styles.message)
	}

	a.drawLeaderboard(screenW, screenH, instY+len(instructions)+2)
}

func (a *App) drawLeaderboard(screenW, screenH, startY int) {
	entries := a.leaderboard.GetEntries()

	header := "=== HIGH SCORES ==="
	a.screen.DrawString((screenW-len(header))/2, startY, header, a.styles.menuHeader)

	if len(entries) == 0 {
		noScores := "No high scores yet!"
		a.screen.DrawString((screenW-len(noScores))/2, startY+2, noScores, a.styles.inventory)
		return
	}

	for i, entry := range entries {
		if startY+2+i >= screenH-1 {
			break
		}
		line := fmt.Sprintf("%2d. %s  %6d  (Hunt %d)", i+1, entry.Initials, entry.Score, entry.Hunt)
		a.screen.DrawString((screenW-len(line))/2, startY+2+i, line, a.styles.message)
	}
}

func (a *App) drawGame(screenW, screenH int) {
	gs := a.gameState

	px, py := gs.Player.Position()

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

			if !gs.IsExplored(worldX, worldY) {
				continue
			}

			tile := gs.Dungeon.GetTile(worldX, worldY)
			if tile != nil {
				var style tcell.Style
				isVisible := gs.IsVisible(worldX, worldY)

				switch tile.Type {
				case dungeon.TileFloor:
					if isVisible {
						style = a.styles.floor
					} else {
						style = a.styles.exploredFloor
					}
				case dungeon.TileWall:
					if isVisible {
						style = a.styles.wall
					} else {
						style = a.styles.exploredWall
					}
				default:
					if isVisible {
						style = a.styles.floor
					} else {
						style = a.styles.exploredFloor
					}
				}
				a.screen.SetCell(screenX, screenY+3, tile.Glyph(), style)
			}
		}
	}

	for _, material := range gs.Materials {
		matX, matY := material.Position()
		if !gs.IsVisible(matX, matY) {
			continue
		}
		matScreenX := matX - cameraX
		matScreenY := matY - cameraY + 3
		if matScreenX >= 0 && matScreenX < screenW && matScreenY >= 3 && matScreenY < screenH-2 {
			style := a.styles.material
			if material.Type.IsRare() {
				style = a.styles.rareMaterial
			}
			a.screen.SetCell(matScreenX, matScreenY, material.Glyph(), style)
		}
	}

	for _, item := range gs.Items {
		ix, iy := item.Position()
		if !gs.IsVisible(ix, iy) {
			continue
		}
		itemScreenX := ix - cameraX
		itemScreenY := iy - cameraY + 3
		if itemScreenX >= 0 && itemScreenX < screenW && itemScreenY >= 3 && itemScreenY < screenH-2 {
			a.screen.SetCell(itemScreenX, itemScreenY, item.Glyph(), a.styles.item)
		}
	}

	for _, monster := range gs.Monsters {
		if monster.Dead {
			continue
		}
		mx, my := monster.Position()
		if !gs.IsVisible(mx, my) {
			continue
		}
		monsterScreenX := mx - cameraX
		monsterScreenY := my - cameraY + 3
		if monsterScreenX >= 0 && monsterScreenX < screenW && monsterScreenY >= 3 && monsterScreenY < screenH-2 {
			style := a.styles.monster
			if monster.IsBoss {
				style = a.styles.boss
			}
			a.screen.SetCell(monsterScreenX, monsterScreenY, monster.Glyph, style)
		}
	}

	playerScreenX := px - cameraX
	playerScreenY := py - cameraY + 3
	if playerScreenX >= 0 && playerScreenX < screenW && playerScreenY >= 3 && playerScreenY < screenH-2 {
		a.screen.SetCell(playerScreenX, playerScreenY, gs.Player.Glyph, a.styles.player)
	}

	title := fmt.Sprintf("BeastTracker - Hunt %d", gs.HuntNumber)
	a.screen.DrawString(0, 0, title, a.styles.hud)

	hpInfo := fmt.Sprintf("HP:%d/%d", gs.Player.HP, gs.Player.EffectiveMaxHP())
	a.screen.DrawString(len(title)+2, 0, hpInfo, getHPStyle(gs.Player.HP, gs.Player.EffectiveMaxHP()))

	scoreInfo := fmt.Sprintf("Score:%d", a.totalScore+gs.Score)
	a.screen.DrawString(len(title)+2+len(hpInfo)+2, 0, scoreInfo, a.styles.hud)

	boss := gs.GetBoss()
	if boss != nil && !boss.Dead {
		bossInfo := fmt.Sprintf("Target: %s HP:%d/%d", boss.Name, boss.HP, boss.MaxHP)
		a.screen.DrawString(screenW-len(bossInfo)-1, 0, bossInfo, a.styles.boss)
	}

	statsInfo := fmt.Sprintf("ATK:%d DEF:%d", gs.Player.EffectiveAttack(), gs.Player.EffectiveDefense())
	a.screen.DrawString(0, 1, statsInfo, a.styles.hud)

	visibleMonsters := gs.GetVisibleMonsters()
	monsterInfoX := len(statsInfo) + 2
	for _, monster := range visibleMonsters {
		info := fmt.Sprintf("%s:%d/%d ", monster.Name, monster.HP, monster.MaxHP)
		if monsterInfoX+len(info) < screenW-15 {
			a.screen.DrawString(monsterInfoX, 1, info, a.styles.monster)
			monsterInfoX += len(info)
		}
	}

	posInfo := fmt.Sprintf("Pos:(%d,%d)", px, py)
	a.screen.DrawString(screenW-len(posInfo)-1, 1, posInfo, a.styles.hud)

	a.drawInventoryBar(gs.Player.Inventory, 0, 2, screenW/2)
	a.drawMaterialPouch(gs.Player.MaterialPouch, screenW/2, 2, screenW/2)

	messageRow := screenH - 2
	if len(gs.Messages) > 0 {
		lastMessage := gs.Messages[len(gs.Messages)-1]
		if len(lastMessage) > screenW {
			lastMessage = lastMessage[:screenW]
		}
		a.screen.DrawString(0, messageRow, lastMessage, a.styles.message)
	}

	instructions := getInstructionsForMode(gs.InputMode)
	a.screen.DrawString(0, screenH-1, instructions, a.styles.floor)

	if gs.InputMode == game.InputModeDropMenu {
		a.drawDropMenu(gs.Player.Inventory, screenW, screenH)
	}

	if gs.InputMode == game.InputModeCrafting {
		a.drawCraftingMenu(gs, screenW, screenH)
	}
}

func (a *App) drawInventoryBar(inv *entity.Inventory, x, y, maxWidth int) {
	label := "Inv:"
	a.screen.DrawString(x, y, label, a.styles.inventory)

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
			style = a.styles.item
		} else {
			slotStr = fmt.Sprintf("[%d:-]", slot)
			style = a.styles.inventory
		}

		a.screen.DrawString(slotX, y, slotStr, style)
		slotX += len(slotStr) + 1
	}
}

func (a *App) drawMaterialPouch(pouch *entity.MaterialPouch, x, y, maxWidth int) {
	label := "Mat:"
	a.screen.DrawString(x, y, label, a.styles.inventory)

	materials := pouch.AllMaterials()
	if len(materials) == 0 {
		a.screen.DrawString(x+len(label)+1, y, "(empty)", a.styles.inventory)
		return
	}

	matX := x + len(label) + 1
	for _, matType := range materials {
		count := pouch.Count(matType)
		info := fmt.Sprintf("%c:%d ", matType.Glyph(), count)
		if matX+len(info) >= x+maxWidth {
			break
		}
		a.screen.DrawString(matX, y, info, a.styles.material)
		matX += len(info)
	}
}

func (a *App) drawDropMenu(inv *entity.Inventory, screenW, screenH int) {
	menuWidth := 30
	menuHeight := inv.Capacity() + 4
	menuX := (screenW - menuWidth) / 2
	menuY := (screenH - menuHeight) / 2

	for dy := 0; dy < menuHeight; dy++ {
		for dx := 0; dx < menuWidth; dx++ {
			a.screen.SetCell(menuX+dx, menuY+dy, ' ', a.styles.menu)
		}
	}

	header := "== Drop Item =="
	headerX := menuX + (menuWidth-len(header))/2
	a.screen.DrawString(headerX, menuY+1, header, a.styles.menuHeader)

	for slot := 1; slot <= inv.Capacity(); slot++ {
		item := inv.GetSlot(slot)
		var line string
		var style tcell.Style

		if item != nil {
			line = fmt.Sprintf(" %d. %s", slot, item.Name())
			style = a.styles.item
		} else {
			line = fmt.Sprintf(" %d. (empty)", slot)
			style = a.styles.menu
		}

		a.screen.DrawString(menuX+2, menuY+2+slot, line, style)
	}

	instrLine := "Press 1-9 or ESC to cancel"
	instrX := menuX + (menuWidth-len(instrLine))/2
	a.screen.DrawString(instrX, menuY+menuHeight-1, instrLine, a.styles.menu)
}

func (a *App) drawCraftingMenu(gs *game.Game, screenW, screenH int) {
	recipes := gs.GetAllRecipes()

	menuWidth := 50
	menuHeight := len(recipes) + 8
	if menuHeight > screenH-4 {
		menuHeight = screenH - 4
	}
	menuX := (screenW - menuWidth) / 2
	menuY := (screenH - menuHeight) / 2

	for dy := 0; dy < menuHeight; dy++ {
		for dx := 0; dx < menuWidth; dx++ {
			a.screen.SetCell(menuX+dx, menuY+dy, ' ', a.styles.menu)
		}
	}

	header := "== Crafting =="
	headerX := menuX + (menuWidth-len(header))/2
	a.screen.DrawString(headerX, menuY+1, header, a.styles.menuHeader)

	visibleRecipes := menuHeight - 6
	startIdx := 0
	if gs.CraftingCursor >= visibleRecipes {
		startIdx = gs.CraftingCursor - visibleRecipes + 1
	}

	for i := 0; i < visibleRecipes && startIdx+i < len(recipes); i++ {
		recipeIdx := startIdx + i
		recipe := recipes[recipeIdx]
		canCraft := recipe.CanCraft(gs.Player.MaterialPouch)

		var style tcell.Style
		if recipeIdx == gs.CraftingCursor {
			style = a.styles.selected
		} else if canCraft {
			style = a.styles.craftable
		} else {
			style = a.styles.uncraftable
		}

		prefix := "  "
		if recipeIdx == gs.CraftingCursor {
			prefix = "> "
		}

		line := fmt.Sprintf("%s%s (%s)", prefix, recipe.Name, recipe.Result.StatsString())
		if len(line) > menuWidth-4 {
			line = line[:menuWidth-4]
		}
		a.screen.DrawString(menuX+2, menuY+3+i, line, style)
	}

	selectedRecipe := recipes[gs.CraftingCursor]
	ingredientLine := "Needs: " + selectedRecipe.IngredientsString()
	if len(ingredientLine) > menuWidth-4 {
		ingredientLine = ingredientLine[:menuWidth-4]
	}
	a.screen.DrawString(menuX+2, menuY+menuHeight-3, ingredientLine, a.styles.menu)

	instrLine := "↑↓:Select Enter:Craft ESC:Close"
	instrX := menuX + (menuWidth-len(instrLine))/2
	a.screen.DrawString(instrX, menuY+menuHeight-1, instrLine, a.styles.menu)
}

func (a *App) drawVictory(screenW, screenH int) {
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
		fmt.Sprintf("Hunt %d Complete!", a.gameState.HuntNumber),
		fmt.Sprintf("Hunt Score: %d", a.gameState.Score),
		fmt.Sprintf("Total Score: %d", a.totalScore+a.gameState.Score),
		fmt.Sprintf("Final HP: %d/%d", a.gameState.Player.HP, a.gameState.Player.EffectiveMaxHP()),
		"",
		"Press N for next hunt, Q to end run",
	}

	startY := screenH/2 - len(lines)/2
	for i, line := range lines {
		x := (screenW - len(line)) / 2
		if x < 0 {
			x = 0
		}
		if i < 6 {
			a.screen.DrawString(x, startY+i, line, tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true))
		} else {
			a.screen.DrawString(x, startY+i, line, a.styles.title)
		}
	}
}

func (a *App) drawGameOver(screenW, screenH int) {
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
		fmt.Sprintf("Hunts Completed: %d", a.currentHunt-1),
		fmt.Sprintf("Final Score: %d", a.totalScore),
		"",
		"Press any key to continue...",
	}

	startY := screenH/2 - len(lines)/2
	for i, line := range lines {
		x := (screenW - len(line)) / 2
		if x < 0 {
			x = 0
		}
		if i < 6 {
			a.screen.DrawString(x, startY+i, line, tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true))
		} else {
			a.screen.DrawString(x, startY+i, line, a.styles.subtitle)
		}
	}
}

func (a *App) drawInitialsEntry(screenW, screenH int) {
	lines := []string{
		"",
		"=== NEW HIGH SCORE! ===",
		"",
		fmt.Sprintf("Score: %d", a.totalScore),
		fmt.Sprintf("Hunt: %d", a.currentHunt),
		"",
		"Enter your initials:",
		"",
		fmt.Sprintf("[ %s ]", padInitials(a.initials)),
		"",
		"Press ENTER to confirm, ESC to skip",
	}

	startY := screenH/2 - len(lines)/2
	for i, line := range lines {
		x := (screenW - len(line)) / 2
		if x < 0 {
			x = 0
		}
		if i == 1 {
			a.screen.DrawString(x, startY+i, line, a.styles.title)
		} else if i == 8 {
			a.screen.DrawString(x, startY+i, line, a.styles.selected)
		} else {
			a.screen.DrawString(x, startY+i, line, a.styles.subtitle)
		}
	}
}

func padInitials(initials string) string {
	for len(initials) < 3 {
		initials += "_"
	}
	return initials
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

func getHPStyle(current, max int) tcell.Style {
	percent := float64(current) / float64(max)
	if percent > 0.6 {
		return tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)
	} else if percent > 0.3 {
		return tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow)
	}
	return tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
}
