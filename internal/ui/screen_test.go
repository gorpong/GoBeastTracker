package ui

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

// TestNewScreen verifies that we can create a new screen wrapper
// with a simulation screen (for testing without a real terminal)
func TestNewScreen(t *testing.T) {
	simScreen := tcell.NewSimulationScreen("")
	screen, err := NewScreen(simScreen)

	if err != nil {
		t.Fatalf("NewScreen() returned error: %v", err)
	}

	if screen == nil {
		t.Fatal("NewScreen() returned nil screen")
	}

	if screen.tcell == nil {
		t.Fatal("Screen.tcell should not be nil")
	}
}

// TestScreenSize verifies that the screen reports correct dimensions
func TestScreenSize(t *testing.T) {
	simScreen := tcell.NewSimulationScreen("")
	if err := simScreen.Init(); err != nil {
		t.Fatalf("Failed to init simulation screen: %v", err)
	}
	defer simScreen.Fini()

	// Simulation screen defaults to 80x25
	simScreen.SetSize(80, 25)

	screen, err := NewScreen(simScreen)
	if err != nil {
		t.Fatalf("NewScreen() returned error: %v", err)
	}

	width, height := screen.Size()

	if width != 80 {
		t.Errorf("Expected width 80, got %d", width)
	}
	if height != 25 {
		t.Errorf("Expected height 25, got %d", height)
	}
}

// TestScreenClear verifies that Clear() doesn't panic and can be called
func TestScreenClear(t *testing.T) {
	simScreen := tcell.NewSimulationScreen("")
	if err := simScreen.Init(); err != nil {
		t.Fatalf("Failed to init simulation screen: %v", err)
	}
	defer simScreen.Fini()

	screen, err := NewScreen(simScreen)
	if err != nil {
		t.Fatalf("NewScreen() returned error: %v", err)
	}

	// Should not panic
	screen.Clear()
}

// TestScreenSetCell verifies that we can set a cell with a rune and style
func TestScreenSetCell(t *testing.T) {
	simScreen := tcell.NewSimulationScreen("")
	if err := simScreen.Init(); err != nil {
		t.Fatalf("Failed to init simulation screen: %v", err)
	}
	defer simScreen.Fini()

	screen, err := NewScreen(simScreen)
	if err != nil {
		t.Fatalf("NewScreen() returned error: %v", err)
	}

	// Set a cell
	screen.SetCell(5, 5, '@', tcell.StyleDefault)
	screen.Show()

	// Verify the cell was set correctly
	primary, _, style, _ := simScreen.GetContent(5, 5)

	if primary != '@' {
		t.Errorf("Expected '@' at (5,5), got '%c'", primary)
	}

	if style != tcell.StyleDefault {
		t.Errorf("Expected default style, got different style")
	}
}

// TestScreenDrawString verifies that we can draw a string at a position
func TestScreenDrawString(t *testing.T) {
	simScreen := tcell.NewSimulationScreen("")
	if err := simScreen.Init(); err != nil {
		t.Fatalf("Failed to init simulation screen: %v", err)
	}
	defer simScreen.Fini()

	screen, err := NewScreen(simScreen)
	if err != nil {
		t.Fatalf("NewScreen() returned error: %v", err)
	}

	testStr := "BeastTracker"
	screen.DrawString(10, 5, testStr, tcell.StyleDefault)
	screen.Show()

	// Verify each character was drawn
	for i, ch := range testStr {
		primary, _, _, _ := simScreen.GetContent(10+i, 5)
		if primary != ch {
			t.Errorf("Expected '%c' at (%d,5), got '%c'", ch, 10+i, primary)
		}
	}
}
