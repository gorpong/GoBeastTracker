package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"

	"beasttracker/internal/game"
	"beasttracker/internal/ui"
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

	// Get screen dimensions and create game
	width, height := screen.Size()
	g := game.NewGame(width, height)

	// Define styles
	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	playerStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow).Bold(true)
	hudStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)

	// Main game loop
	for g.Running {
		screen.Clear()

		// Draw floor (dots for empty space)
		w, h := screen.Size()
		for y := 1; y < h-1; y++ {
			for x := 0; x < w; x++ {
				screen.SetCell(x, y, '.', defStyle)
			}
		}

		// Draw player
		px, py := g.Player.Position()
		screen.SetCell(px, py, g.Player.Glyph, playerStyle)

		// Draw HUD at top
		title := "BeastTracker - Phase 1"
		screen.DrawString(0, 0, title, hudStyle)

		// Draw position info
		posInfo := fmt.Sprintf("Pos: (%d, %d)", px, py)
		screen.DrawString(w-len(posInfo)-1, 0, posInfo, hudStyle)

		// Draw instructions at bottom
		instructions := "Move: arrows/hjkl/wasd | Quit: q/ESC"
		screen.DrawString(0, h-1, instructions, defStyle)

		screen.Show()

		// Handle input
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
			// Update game dimensions on resize
			newW, newH := screen.Size()
			g.Width = newW
			g.Height = newH
		case *tcell.EventKey:
			action := ui.ParseAction(ev.Key(), ev.Rune())
			dir := ui.ParseDirection(ev.Key(), ev.Rune())
			g.HandleInput(action, dir)
		}
	}
}
