package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gdamore/tcell/v2"

	"beasttracker/internal/dungeon"
	"beasttracker/internal/entity"
	"beasttracker/internal/game"
	"beasttracker/internal/save"
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
	AppStateSaveGame
	AppStateLoadGame
)

type App struct {
	screen         *ui.Screen
	tcellScreen    tcell.Screen
	gameState      *game.Game
	leaderboard    *score.Leaderboard
	saveManager    *save.SaveManager
	appState       AppState
	previousState  AppState
	initials       string
	saveName       string
	totalScore     int
	currentHunt    int
	checkpoint     *save.SaveData
	styles         styles
	saveList       []*save.SaveData
	saveCursor     int
	equipSlotView  entity.EquipmentSlot
	equipCursor    int
	showVictoryMsg bool
}

type styles struct {
	floor         tcell.Style
	wall          tcell.Style
	exploredFloor tcell.Style
	exploredWall  tcell.Style
	player        tcell.Style
	hud           tcell.Style
	message       tcell.Style
	boss          tcell.Style
	monster       tcell.Style
	item          tcell.Style
	material      tcell.Style
	rareMaterial  tcell.Style
	inventory     tcell.Style
	menu          tcell.Style
	menuHeader    tcell.Style
	craftable     tcell.Style
	uncraftable   tcell.Style
	selected      tcell.Style
	title         tcell.Style
	subtitle      tcell.Style
	victory       tcell.Style
	danger        tcell.Style
	equipped      tcell.Style
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
		victory:       tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen).Bold(true),
		danger:        tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed).Bold(true),
		equipped:      tcell.StyleDefault.Background(tcell.ColorDarkBlue).Foreground(tcell.ColorLime).Bold(true),
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

	savePath := getSavePath()
	saveManager := save.NewSaveManager(savePath)

	leaderboard := score.NewLeaderboard()
	leaderboard.SetFilepath(getScorePath())
	leaderboard.Load()

	app := &App{
		screen:        screen,
		tcellScreen:   tcellScreen,
		leaderboard:   leaderboard,
		saveManager:   saveManager,
		appState:      AppStateSplash,
		styles:        newStyles(),
		currentHunt:   1,
		totalScore:    0,
		equipSlotView: entity.SlotWeapon,
	}

	app.run()
}

func getSavePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		return filepath.Join("assets", "data", "saves")
	}
	return filepath.Join(homeDir, ".beasttracker", "saves")
}

func getScorePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join("assets", "data", "scores.json")
	}
	return filepath.Join(homeDir, ".beasttracker", "scores.json")
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
			if a.gameState.InputMode == game.InputModeEquipment {
				a.drawEquipmentScreen(screenW, screenH)
			}
		case AppStateVictory:
			a.drawGame(screenW, screenH)
			if a.gameState.InputMode == game.InputModeEquipment {
				a.drawEquipmentScreen(screenW, screenH)
			} else if a.showVictoryMsg && a.gameState.InputMode == game.InputModeNormal {
				a.drawVictoryOverlay(screenW, screenH)
			}
		case AppStateGameOver:
			a.drawGameOverScreen(screenW, screenH)
		case AppStateInitials:
			a.drawInitialsEntry(screenW, screenH)
		case AppStateSaveGame:
			a.drawSaveGameScreen(screenW, screenH)
		case AppStateLoadGame:
			a.drawLoadGameScreen(screenW, screenH)
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
	case AppStateSaveGame:
		return a.handleSaveGameKey(ev)
	case AppStateLoadGame:
		return a.handleLoadGameKey(ev)
	}
	return false
}

func (a *App) handleSplashKey(ev *tcell.EventKey) bool {
	if ev.Key() == tcell.KeyEscape {
		return true
	}

	switch ev.Rune() {
	case 'q', 'Q':
		return true
	case 'n', 'N':
		a.startNewGame()
	case 'l', 'L':
		a.openLoadScreen()
	}

	if ev.Key() == tcell.KeyEnter {
		a.startNewGame()
	}

	return false
}

func (a *App) handlePlayingKey(ev *tcell.EventKey) bool {
	if a.gameState == nil {
		return false
	}

	// Clear victory overlay when entering any sub-mode
	if a.gameState.InputMode != game.InputModeNormal {
		a.showVictoryMsg = false
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
	case game.InputModeEquipment:
		a.handleEquipmentModeKey(ev)
	default:
		a.handleNormalModeKey(ev)
	}

	if a.gameState.GameState == game.StateVictory && a.appState != AppStateVictory {
		a.totalScore += a.gameState.Score
		a.gameState.Score = 0
		a.showVictoryMsg = true
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
	// If in a sub-mode, handle that instead
	if a.gameState.InputMode != game.InputModeNormal {
		a.showVictoryMsg = false

		switch a.gameState.InputMode {
		case game.InputModeDropping:
			a.handleDropModeKey(ev)
		case game.InputModeDropMenu:
			a.handleDropMenuKey(ev)
		case game.InputModeInventory:
			a.handleInventoryKey(ev)
		case game.InputModeCrafting:
			a.handleCraftingKey(ev)
		case game.InputModeEquipment:
			a.handleEquipmentModeKey(ev)
		}
		return false
	}

	action := ui.ParseAction(ev.Key(), ev.Rune())

	switch action {
	case ui.ActionQuit:
		a.checkHighScore()
		return false
	case ui.ActionNextHunt:
		a.openSaveScreen()
		return false
	case ui.ActionCraft:
		a.showVictoryMsg = false
		a.gameState.InputMode = game.InputModeCrafting
		a.gameState.CraftingCursor = 0
		return false
	case ui.ActionEquipment:
		a.showVictoryMsg = false
		a.previousState = a.appState
		a.gameState.InputMode = game.InputModeEquipment
		a.equipCursor = 0
		a.equipSlotView = entity.SlotWeapon
		return false
	case ui.ActionInventory:
		a.showVictoryMsg = false
		a.gameState.InputMode = game.InputModeInventory
		return false
	case ui.ActionDropMode:
		a.showVictoryMsg = false
		a.gameState.InputMode = game.InputModeDropping
		return false
	case ui.ActionMove:
		a.showVictoryMsg = false
		dir := ui.ParseDirection(ev.Key(), ev.Rune())
		a.gameState.HandleInput(action, dir)
		return false
	case ui.ActionUseItem:
		a.showVictoryMsg = false
		slot, ok := ui.ParseSlotNumber(ev.Rune())
		if ok {
			a.gameState.UseItemInSlot(slot)
		}
		return false
	}

	// Any other key dismisses the modal but stays in victory/explore state
	if ev.Key() != tcell.KeyRune || (ev.Rune() != 'q' && ev.Rune() != 'Q') {
		a.showVictoryMsg = false
	}

	return false
}

func (a *App) handleGameOverKey(ev *tcell.EventKey) bool {
	action := ui.ParseAction(ev.Key(), ev.Rune())

	switch action {
	case ui.ActionRestart:
		if a.checkpoint != nil {
			a.restartFromCheckpoint()
		}
		return false
	case ui.ActionQuit:
		a.checkHighScore()
		return false
	}

	if ev.Key() == tcell.KeyEnter {
		a.checkHighScore()
	}

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

func (a *App) handleSaveGameKey(ev *tcell.EventKey) bool {
	if ev.Key() == tcell.KeyEscape {
		a.appState = AppStateVictory
		return false
	}

	if ev.Key() == tcell.KeyEnter && len(a.saveName) > 0 {
		a.saveAndContinue()
		return false
	}

	if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
		if len(a.saveName) > 0 {
			a.saveName = a.saveName[:len(a.saveName)-1]
		}
		return false
	}

	if ev.Key() == tcell.KeyRune && len(a.saveName) < save.MaxNameLength {
		r := ev.Rune()
		if ui.IsLetter(r) || (r >= '0' && r <= '9') || r == ' ' || r == '-' || r == '_' {
			a.saveName += string(r)
		}
	}

	return false
}

func (a *App) handleLoadGameKey(ev *tcell.EventKey) bool {
	if ev.Key() == tcell.KeyEscape {
		a.appState = AppStateSplash
		return false
	}

	if len(a.saveList) == 0 {
		if ev.Key() == tcell.KeyEnter {
			a.appState = AppStateSplash
		}
		return false
	}

	switch ev.Key() {
	case tcell.KeyUp:
		if a.saveCursor > 0 {
			a.saveCursor--
		}
	case tcell.KeyDown:
		if a.saveCursor < len(a.saveList)-1 {
			a.saveCursor++
		}
	case tcell.KeyEnter:
		a.loadSelectedSave()
	case tcell.KeyDelete:
		a.deleteSelectedSave()
	}

	if ev.Rune() == 'k' || ev.Rune() == 'K' {
		if a.saveCursor > 0 {
			a.saveCursor--
		}
	}
	if ev.Rune() == 'j' || ev.Rune() == 'J' {
		if a.saveCursor < len(a.saveList)-1 {
			a.saveCursor++
		}
	}

	return false
}

func (a *App) handleEquipmentKey(ev *tcell.EventKey) bool {
	if ev.Key() == tcell.KeyEscape || ev.Rune() == 'e' || ev.Rune() == 'E' {
		a.gameState.InputMode = game.InputModeNormal
		// Return to correct app state - don't change appState, just close the menu
		return false
	}

	equipList := a.gameState.GetEquipmentList(a.equipSlotView)

	switch ev.Key() {
	case tcell.KeyUp:
		if a.equipCursor > 0 {
			a.equipCursor--
		}
	case tcell.KeyDown:
		if a.equipCursor < len(equipList)-1 {
			a.equipCursor++
		}
	case tcell.KeyLeft:
		a.equipSlotView = (a.equipSlotView + 2) % 3
		a.equipCursor = 0
	case tcell.KeyRight:
		a.equipSlotView = (a.equipSlotView + 1) % 3
		a.equipCursor = 0
	case tcell.KeyEnter:
		if len(equipList) > 0 && a.equipCursor < len(equipList) {
			selectedEquip := equipList[a.equipCursor]
			if !a.gameState.IsEquipped(selectedEquip) {
				a.gameState.Player.EquipFromStash(selectedEquip)
				a.gameState.AddMessage(fmt.Sprintf("Equipped %s.", selectedEquip.Name))
			}
		}
	}

	if ev.Key() == tcell.KeyRune {
		switch ev.Rune() {
		case 'k', 'K':
			if a.equipCursor > 0 {
				a.equipCursor--
			}
		case 'j', 'J':
			if a.equipCursor < len(equipList)-1 {
				a.equipCursor++
			}
		case 'h', 'H':
			a.equipSlotView = (a.equipSlotView + 2) % 3
			a.equipCursor = 0
		case 'l', 'L':
			a.equipSlotView = (a.equipSlotView + 1) % 3
			a.equipCursor = 0
		case 'u', 'U':
			a.gameState.Player.UnequipToStash(a.equipSlotView)
			a.gameState.AddMessage(fmt.Sprintf("Unequipped %s slot.", a.equipSlotView.String()))
		}
	}

	return false
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

	if action == ui.ActionEquipment {
		a.showVictoryMsg = false
		a.gameState.InputMode = game.InputModeEquipment
		a.equipCursor = 0
		a.equipSlotView = entity.SlotWeapon
		return
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

	_, ok := ui.ParseSlotNumber(ev.Rune())
	if ok {
		a.gameState.HandleDropModeInput(ev.Rune())
	}
}

func (a *App) handleInventoryKey(ev *tcell.EventKey) {
	if ev.Key() == tcell.KeyEscape || ev.Rune() == 'i' || ev.Rune() == 'I' {
		a.gameState.InputMode = game.InputModeNormal
		return
	}

	slot, ok := ui.ParseSlotNumber(ev.Rune())
	if ok {
		a.gameState.UseItemInSlot(slot)
	}
}

func (a *App) handleCraftingKey(ev *tcell.EventKey) {
	recipes := a.gameState.GetAllRecipes()

	switch ev.Key() {
	case tcell.KeyEscape:
		a.gameState.InputMode = game.InputModeNormal
		// Don't quit - just return to normal/victory explore mode
		return
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

func (a *App) handleEquipmentModeKey(ev *tcell.EventKey) {
	a.handleEquipmentKey(ev)
}

func (a *App) checkHighScore() {
	if a.leaderboard.IsHighScore(a.totalScore) && a.totalScore > 0 {
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
	a.checkpoint = nil
	a.showVictoryMsg = false
}

func (a *App) startNewGame() {
	seed := time.Now().UnixNano()
	a.gameState = game.NewGame(dungeonWidth, dungeonHeight, seed)
	a.currentHunt = 1
	a.totalScore = 0
	a.checkpoint = a.gameState.CreateCheckpoint()
	a.appState = AppStatePlaying
	a.showVictoryMsg = false
}

func (a *App) startNextHunt() {
	seed := time.Now().UnixNano()
	a.currentHunt++
	a.gameState = game.NewGameWithHunt(dungeonWidth, dungeonHeight, seed, a.currentHunt, a.gameState.Player)
	a.checkpoint = a.gameState.CreateCheckpoint()
	a.appState = AppStatePlaying
	a.showVictoryMsg = false
}

func (a *App) openSaveScreen() {
	a.saveName = ""
	a.appState = AppStateSaveGame
}

func (a *App) openLoadScreen() {
	a.saveList, _ = a.saveManager.List()
	a.saveCursor = 0
	a.appState = AppStateLoadGame
}

func (a *App) saveAndContinue() {
	checkpoint := a.gameState.CreateCheckpoint()
	checkpoint.Name = a.saveName
	checkpoint.HuntNumber = a.currentHunt + 1
	checkpoint.Score = a.totalScore

	err := a.saveManager.Save(checkpoint)
	if err != nil {
		a.gameState.AddMessage(fmt.Sprintf("Save failed: %v", err))
		a.appState = AppStateVictory
		return
	}

	a.gameState.AddMessage(fmt.Sprintf("Game saved as '%s'", a.saveName))
	a.startNextHunt()
}

func (a *App) loadSelectedSave() {
	if len(a.saveList) == 0 || a.saveCursor >= len(a.saveList) {
		return
	}

	selectedSave := a.saveList[a.saveCursor]
	seed := time.Now().UnixNano()

	a.gameState = game.NewGameFromCheckpoint(dungeonWidth, dungeonHeight, seed, selectedSave)
	a.currentHunt = selectedSave.HuntNumber
	a.totalScore = selectedSave.Score
	a.checkpoint = selectedSave
	a.appState = AppStatePlaying
	a.showVictoryMsg = false
}

func (a *App) deleteSelectedSave() {
	if len(a.saveList) == 0 || a.saveCursor >= len(a.saveList) {
		return
	}

	selectedSave := a.saveList[a.saveCursor]
	a.saveManager.Delete(selectedSave.Name)
	a.saveList, _ = a.saveManager.List()

	if a.saveCursor >= len(a.saveList) && a.saveCursor > 0 {
		a.saveCursor--
	}
}

func (a *App) restartFromCheckpoint() {
	if a.checkpoint == nil {
		return
	}

	seed := time.Now().UnixNano()
	a.gameState = game.NewGameFromCheckpoint(dungeonWidth, dungeonHeight, seed, a.checkpoint)
	a.totalScore = a.checkpoint.Score
	a.appState = AppStatePlaying
	a.showVictoryMsg = false
}

// Drawing functions

func (a *App) drawSplash(screenW, screenH int) {
	title := []string{
		"  ____                 _  _____               _             ",
		" | __ )  ___  __ _ ___| ||_   _| __ __ _  ___| | _____ _ __ ",
		" |  _ \\ / _ \\/ _` / __| __|| || '__/ _` |/ __| |/ / _ \\ '__|",
		" | |_) |  __/ (_| \\__ \\ |_ | || | | (_| | (__|   <  __/ |   ",
		" |____/ \\___|\\__,_|___/\\__||_||_|  \\__,_|\\___|_|\\_\\___|_|   ",
	}

	startY := 2
	for i, line := range title {
		x := (screenW - len(line)) / 2
		if x < 0 {
			x = 0
		}
		a.screen.DrawString(x, startY+i, line, a.styles.title)
	}

	subtitle := "A Monster Hunter-Inspired Roguelike"
	a.screen.DrawString((screenW-len(subtitle))/2, startY+7, subtitle, a.styles.subtitle)

	menuItems := []string{
		"",
		"[N] New Game",
		"[L] Load Game",
		"[Q] Quit",
		"",
		"Controls:",
		"  Move: Arrow keys / HJKL / WASD",
		"  Use Item: 1-9",
		"  Drop Item: X",
		"  Craft: C",
		"  Equipment: E",
	}

	menuY := startY + 9
	for i, line := range menuItems {
		x := (screenW - len(line)) / 2
		if x < 0 {
			x = 0
		}
		a.screen.DrawString(x, menuY+i, line, a.styles.message)
	}

	a.drawLeaderboard(screenW, screenH, menuY+len(menuItems)+1)
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

	a.drawHUD(screenW, screenH)
}

func (a *App) drawHUD(screenW, screenH int) {
	gs := a.gameState

	title := fmt.Sprintf("Hunt %d", a.currentHunt)
	a.screen.DrawString(0, 0, title, a.styles.hud)

	hpInfo := fmt.Sprintf("HP:%d/%d", gs.Player.HP, gs.Player.EffectiveMaxHP())
	a.screen.DrawString(len(title)+2, 0, hpInfo, getHPStyle(gs.Player.HP, gs.Player.EffectiveMaxHP()))

	scoreInfo := fmt.Sprintf("Score:%d", a.totalScore+gs.Score)
	a.screen.DrawString(len(title)+2+len(hpInfo)+2, 0, scoreInfo, a.styles.hud)

	boss := gs.GetBoss()
	if boss != nil {
		bossInfo := fmt.Sprintf("Target: %s HP:%d/%d", boss.Name, boss.HP, boss.MaxHP)
		a.screen.DrawString(screenW-len(bossInfo)-1, 0, bossInfo, a.styles.boss)
	} else if gs.GameState == game.StateVictory {
		victoryInfo := "BOSS DEFEATED!"
		a.screen.DrawString(screenW-len(victoryInfo)-1, 0, victoryInfo, a.styles.victory)
	}

	statsInfo := fmt.Sprintf("ATK:%d DEF:%d", gs.Player.EffectiveAttack(), gs.Player.EffectiveDefense())
	a.screen.DrawString(0, 1, statsInfo, a.styles.hud)

	visibleMonsters := gs.GetVisibleMonsters()
	monsterInfoX := len(statsInfo) + 2
	for _, monster := range visibleMonsters {
		info := fmt.Sprintf("%c:%d ", monster.Glyph, monster.HP)
		if monsterInfoX+len(info) < screenW-15 {
			a.screen.DrawString(monsterInfoX, 1, info, a.styles.monster)
			monsterInfoX += len(info)
		}
	}

	px, py := gs.Player.Position()
	posInfo := fmt.Sprintf("(%d,%d)", px, py)
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

	instructions := a.getInstructionsForMode()
	a.screen.DrawString(0, screenH-1, instructions, a.styles.floor)

	if gs.InputMode == game.InputModeDropMenu {
		a.drawDropMenu(gs.Player.Inventory, screenW, screenH)
	}

	if gs.InputMode == game.InputModeCrafting {
		a.drawCraftingMenu(gs, screenW, screenH)
	}
}

func (a *App) getInstructionsForMode() string {
	if a.gameState == nil {
		return ""
	}

	if a.appState == AppStateVictory && a.showVictoryMsg {
		return "VICTORY! N:Next Hunt | Q:End Run | Any other key:Continue exploring"
	}

	switch a.gameState.InputMode {
	case game.InputModeDropping:
		return "Drop: 1-9 quick drop | x/i menu | other:cancel"
	case game.InputModeDropMenu:
		return "Drop menu: 1-9 to drop | ESC:cancel"
	case game.InputModeInventory:
		return "Inventory: 1-9 to use | ESC:close"
	case game.InputModeCrafting:
		return "Craft: ↑↓ select | Enter:craft | ESC:close"
	case game.InputModeEquipment:
		return "Equipment: ←→ slot | ↑↓ select | Enter:equip | U:unequip | ESC:close"
	default:
		return "Move:↑↓←→/hjkl | 1-9:use | x:drop | c:craft | e:equip | q:quit"
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
		style := a.styles.material
		if matType.IsRare() {
			style = a.styles.rareMaterial
		}
		a.screen.DrawString(matX, y, info, style)
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

func (a *App) drawVictoryOverlay(screenW, screenH int) {
	if !a.showVictoryMsg {
		return
	}

	menuWidth := 40
	menuHeight := 10
	menuX := (screenW - menuWidth) / 2
	menuY := (screenH - menuHeight) / 2

	for dy := 0; dy < menuHeight; dy++ {
		for dx := 0; dx < menuWidth; dx++ {
			a.screen.SetCell(menuX+dx, menuY+dy, ' ', a.styles.menu)
		}
	}

	header := "*** VICTORY! ***"
	a.screen.DrawString(menuX+(menuWidth-len(header))/2, menuY+1, header, a.styles.victory)

	msg := "You have slain the beast!"
	a.screen.DrawString(menuX+(menuWidth-len(msg))/2, menuY+3, msg, a.styles.subtitle)

	huntInfo := fmt.Sprintf("Hunt %d Complete!", a.currentHunt)
	a.screen.DrawString(menuX+(menuWidth-len(huntInfo))/2, menuY+5, huntInfo, a.styles.message)

	scoreInfo := fmt.Sprintf("Total Score: %d", a.totalScore)
	a.screen.DrawString(menuX+(menuWidth-len(scoreInfo))/2, menuY+6, scoreInfo, a.styles.message)

	instructions := "N:Next Hunt | Q:Quit | Other:Explore"
	a.screen.DrawString(menuX+(menuWidth-len(instructions))/2, menuY+8, instructions, a.styles.menu)
}

func (a *App) drawGameOverScreen(screenW, screenH int) {
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
	}

	if a.currentHunt > 1 {
		lines = append(lines, fmt.Sprintf("Hunts Completed: %d", a.currentHunt-1))
	} else {
		lines = append(lines, "Hunts Completed: 0")
	}

	lines = append(lines, fmt.Sprintf("Final Score: %d", a.totalScore))
	lines = append(lines, "")

	if a.checkpoint != nil {
		lines = append(lines, "[R] Restart from checkpoint")
	}
	lines = append(lines, "[Enter/Q] Continue")

	startY := screenH/2 - len(lines)/2
	for i, line := range lines {
		x := (screenW - len(line)) / 2
		if x < 0 {
			x = 0
		}
		if i >= 1 && i <= 5 {
			a.screen.DrawString(x, startY+i, line, a.styles.danger)
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

func (a *App) drawSaveGameScreen(screenW, screenH int) {
	menuWidth := 50
	menuHeight := 12
	menuX := (screenW - menuWidth) / 2
	menuY := (screenH - menuHeight) / 2

	for dy := 0; dy < menuHeight; dy++ {
		for dx := 0; dx < menuWidth; dx++ {
			a.screen.SetCell(menuX+dx, menuY+dy, ' ', a.styles.menu)
		}
	}

	header := "== Save Game =="
	a.screen.DrawString(menuX+(menuWidth-len(header))/2, menuY+1, header, a.styles.menuHeader)

	info := fmt.Sprintf("Hunt %d | Score: %d", a.currentHunt+1, a.totalScore)
	a.screen.DrawString(menuX+(menuWidth-len(info))/2, menuY+3, info, a.styles.message)

	prompt := "Enter save name:"
	a.screen.DrawString(menuX+4, menuY+5, prompt, a.styles.subtitle)

	nameDisplay := fmt.Sprintf("[%s_]", a.saveName)
	if len(a.saveName) >= save.MaxNameLength {
		nameDisplay = fmt.Sprintf("[%s]", a.saveName)
	}
	a.screen.DrawString(menuX+4, menuY+7, nameDisplay, a.styles.selected)

	if !a.saveManager.HasRoom() && !a.saveManager.Exists(a.saveName) {
		warning := "Save slots full! Use existing name to overwrite."
		a.screen.DrawString(menuX+(menuWidth-len(warning))/2, menuY+9, warning, a.styles.danger)
	}

	instructions := "Enter:Save | ESC:Cancel"
	a.screen.DrawString(menuX+(menuWidth-len(instructions))/2, menuY+menuHeight-1, instructions, a.styles.menu)
}

func (a *App) drawLoadGameScreen(screenW, screenH int) {
	menuWidth := 55
	menuHeight := 16
	menuX := (screenW - menuWidth) / 2
	menuY := (screenH - menuHeight) / 2

	for dy := 0; dy < menuHeight; dy++ {
		for dx := 0; dx < menuWidth; dx++ {
			a.screen.SetCell(menuX+dx, menuY+dy, ' ', a.styles.menu)
		}
	}

	header := "== Load Game =="
	a.screen.DrawString(menuX+(menuWidth-len(header))/2, menuY+1, header, a.styles.menuHeader)

	if len(a.saveList) == 0 {
		noSaves := "No saved games found."
		a.screen.DrawString(menuX+(menuWidth-len(noSaves))/2, menuY+5, noSaves, a.styles.inventory)

		instructions := "Press Enter or ESC to return"
		a.screen.DrawString(menuX+(menuWidth-len(instructions))/2, menuY+menuHeight-1, instructions, a.styles.menu)
		return
	}

	visibleSaves := menuHeight - 6
	startIdx := 0
	if a.saveCursor >= visibleSaves {
		startIdx = a.saveCursor - visibleSaves + 1
	}

	for i := 0; i < visibleSaves && startIdx+i < len(a.saveList); i++ {
		saveIdx := startIdx + i
		saveData := a.saveList[saveIdx]

		var style tcell.Style
		if saveIdx == a.saveCursor {
			style = a.styles.selected
		} else {
			style = a.styles.message
		}

		prefix := "  "
		if saveIdx == a.saveCursor {
			prefix = "> "
		}

		line := fmt.Sprintf("%s%-20s Hunt:%d Score:%d", prefix, saveData.Name, saveData.HuntNumber, saveData.Score)
		if len(line) > menuWidth-4 {
			line = line[:menuWidth-4]
		}
		a.screen.DrawString(menuX+2, menuY+3+i, line, style)
	}

	instructions := "↑↓:Select | Enter:Load | Del:Delete | ESC:Back"
	a.screen.DrawString(menuX+(menuWidth-len(instructions))/2, menuY+menuHeight-1, instructions, a.styles.menu)
}

func (a *App) drawEquipmentScreen(screenW, screenH int) {
	menuWidth := 50
	menuHeight := 18
	menuX := (screenW - menuWidth) / 2
	menuY := (screenH - menuHeight) / 2

	for dy := 0; dy < menuHeight; dy++ {
		for dx := 0; dx < menuWidth; dx++ {
			a.screen.SetCell(menuX+dx, menuY+dy, ' ', a.styles.menu)
		}
	}

	header := "== Equipment =="
	a.screen.DrawString(menuX+(menuWidth-len(header))/2, menuY+1, header, a.styles.menuHeader)

	// Tab headers
	slots := []entity.EquipmentSlot{entity.SlotWeapon, entity.SlotArmor, entity.SlotCharm}
	tabX := menuX + 4
	for _, slot := range slots {
		tabName := fmt.Sprintf("[%s]", slot.String())
		style := a.styles.menu
		if slot == a.equipSlotView {
			style = a.styles.selected
		}
		a.screen.DrawString(tabX, menuY+3, tabName, style)
		tabX += len(tabName) + 2
	}

	// Current stats
	statsLine := fmt.Sprintf("ATK:%d DEF:%d MaxHP:%d",
		a.gameState.Player.EffectiveAttack(),
		a.gameState.Player.EffectiveDefense(),
		a.gameState.Player.EffectiveMaxHP())
	a.screen.DrawString(menuX+4, menuY+5, statsLine, a.styles.hud)

	// Equipment list
	equipList := a.gameState.GetEquipmentList(a.equipSlotView)

	if len(equipList) == 0 {
		noEquip := "(No equipment available)"
		a.screen.DrawString(menuX+4, menuY+7, noEquip, a.styles.inventory)
	} else {
		visibleItems := menuHeight - 10
		startIdx := 0
		if a.equipCursor >= visibleItems {
			startIdx = a.equipCursor - visibleItems + 1
		}

		for i := 0; i < visibleItems && startIdx+i < len(equipList); i++ {
			itemIdx := startIdx + i
			equip := equipList[itemIdx]

			var style tcell.Style
			isEquipped := a.gameState.IsEquipped(equip)

			if itemIdx == a.equipCursor {
				style = a.styles.selected
			} else if isEquipped {
				style = a.styles.equipped
			} else {
				style = a.styles.message
			}

			prefix := "  "
			if itemIdx == a.equipCursor {
				prefix = "> "
			}

			suffix := ""
			if isEquipped {
				suffix = " [E]"
			}

			line := fmt.Sprintf("%s%s (%s)%s", prefix, equip.Name, equip.StatsString(), suffix)
			if len(line) > menuWidth-4 {
				line = line[:menuWidth-4]
			}
			a.screen.DrawString(menuX+2, menuY+7+i, line, style)
		}
	}

	instructions := "←→:Slot | ↑↓:Select | Enter:Equip | U:Unequip | ESC:Close"
	a.screen.DrawString(menuX+(menuWidth-len(instructions))/2, menuY+menuHeight-1, instructions, a.styles.menu)
}

func padInitials(initials string) string {
	for len(initials) < 3 {
		initials += "_"
	}
	return initials
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
