package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"

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

	// Set default style
	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	titleStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow).Bold(true)

	// Main loop
	for {
		screen.Clear()

		width, height := screen.Size()

		// Draw title centered
		title := "BeastTracker"
		titleX := (width - len(title)) / 2
		titleY := height / 3
		screen.DrawString(titleX, titleY, title, titleStyle)

		// Draw instructions
		instructions := "Press 'q' or ESC to quit"
		instrX := (width - len(instructions)) / 2
		screen.DrawString(instrX, titleY+2, instructions, defStyle)

		// Draw version
		version := "v0.1.0 - Phase 0"
		verX := (width - len(version)) / 2
		screen.DrawString(verX, titleY+4, version, defStyle)

		screen.Show()

		// Handle input
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' {
				return
			}
		}
	}
}
