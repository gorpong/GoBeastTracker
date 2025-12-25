package ui

import (
	"github.com/gdamore/tcell/v2"
)

// Screen wraps tcell.Screen to provide a simplified interface for the game.
// This abstraction allows us to mock the screen for testing.
type Screen struct {
	tcell tcell.Screen
}

// NewScreen creates a new Screen wrapper around a tcell.Screen.
// Pass nil to create a real terminal screen, or pass a tcell.SimulationScreen for testing.
func NewScreen(s tcell.Screen) (*Screen, error) {
	return &Screen{tcell: s}, nil
}

// Size returns the current terminal dimensions (width, height).
func (s *Screen) Size() (int, int) {
	return s.tcell.Size()
}

// Clear clears the screen buffer.
func (s *Screen) Clear() {
	s.tcell.Clear()
}

// SetCell sets a single cell at (x, y) with the given rune and style.
func (s *Screen) SetCell(x, y int, ch rune, style tcell.Style) {
	s.tcell.SetContent(x, y, ch, nil, style)
}

// DrawString draws a string starting at (x, y) with the given style.
func (s *Screen) DrawString(x, y int, str string, style tcell.Style) {
	for i, ch := range str {
		s.tcell.SetContent(x+i, y, ch, nil, style)
	}
}

// Show synchronizes the internal buffer to the terminal.
func (s *Screen) Show() {
	s.tcell.Show()
}

// PollEvent waits for and returns the next event.
func (s *Screen) PollEvent() tcell.Event {
	return s.tcell.PollEvent()
}

// Fini finalizes the screen and restores terminal state.
func (s *Screen) Fini() {
	s.tcell.Fini()
}

// Sync synchronizes the terminal, useful after resize.
func (s *Screen) Sync() {
	s.tcell.Sync()
}
