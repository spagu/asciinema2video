package terminal

import (
	"image/color"
	"testing"
)

func TestNewTerminal(t *testing.T) {
	term := New(80, 24)

	if term.Width != 80 {
		t.Errorf("expected width 80, got %d", term.Width)
	}
	if term.Height != 24 {
		t.Errorf("expected height 24, got %d", term.Height)
	}
	if len(term.Screen) != 24 {
		t.Errorf("expected 24 rows, got %d", len(term.Screen))
	}
	if len(term.Screen[0]) != 80 {
		t.Errorf("expected 80 cols, got %d", len(term.Screen[0]))
	}
}

func TestWriteSimpleText(t *testing.T) {
	term := New(80, 24)
	term.Write("Hello")

	expected := "Hello"
	for i, ch := range expected {
		if term.Screen[0][i].Char != ch {
			t.Errorf("expected '%c' at position %d, got '%c'", ch, i, term.Screen[0][i].Char)
		}
	}
}

func TestWriteNewline(t *testing.T) {
	term := New(80, 24)
	term.Write("Line1\r\nLine2") // \r\n resets X position and moves down

	if term.Screen[0][0].Char != 'L' {
		t.Error("expected 'L' at row 0")
	}
	if term.Screen[1][0].Char != 'L' {
		t.Error("expected 'L' at row 1")
	}
}

func TestWriteCarriageReturn(t *testing.T) {
	term := New(80, 24)
	term.Write("Hello\rWorld")

	expected := "World"
	for i, ch := range expected {
		if term.Screen[0][i].Char != ch {
			t.Errorf("expected '%c' at position %d, got '%c'", ch, i, term.Screen[0][i].Char)
		}
	}
}

func TestWriteBackspace(t *testing.T) {
	term := New(80, 24)
	term.Write("Hello\b\bXX")

	expected := "HelXX"
	for i, ch := range expected {
		if term.Screen[0][i].Char != ch {
			t.Errorf("expected '%c' at position %d, got '%c'", ch, i, term.Screen[0][i].Char)
		}
	}
}

func TestWriteTab(t *testing.T) {
	term := New(80, 24)
	term.Write("A\tB")

	if term.Screen[0][0].Char != 'A' {
		t.Error("expected 'A' at position 0")
	}
	if term.Screen[0][8].Char != 'B' {
		t.Error("expected 'B' at position 8 (next tab stop)")
	}
}

func TestSGRColors(t *testing.T) {
	term := New(80, 24)
	term.Write("\x1b[31mRed\x1b[0m")

	// Check that first character is red
	if term.Screen[0][0].FG != defaultColors[1] {
		t.Errorf("expected red foreground, got %v", term.Screen[0][0].FG)
	}
}

func TestSGRBold(t *testing.T) {
	term := New(80, 24)
	term.Write("\x1b[1mBold\x1b[0m")

	if !term.Screen[0][0].Bold {
		t.Error("expected bold attribute")
	}
}

func TestSGRReset(t *testing.T) {
	term := New(80, 24)
	term.Write("\x1b[31;1mStyled\x1b[0mNormal")

	// After reset, should be default
	normalIdx := 6
	theme, _ := GetTheme("default")
	if term.Screen[0][normalIdx].FG != theme.Foreground {
		t.Error("expected default foreground after reset")
	}
	if term.Screen[0][normalIdx].Bold {
		t.Error("expected no bold after reset")
	}
}

func TestCursorMovement(t *testing.T) {
	term := New(80, 24)
	term.Write("\x1b[5;10HX")

	if term.Screen[4][9].Char != 'X' {
		t.Error("expected 'X' at row 5, col 10 (0-indexed: 4,9)")
	}
}

func TestEraseDisplay(t *testing.T) {
	term := New(80, 24)
	term.Write("AAAAA")
	term.Write("\x1b[2J")

	for i := 0; i < 5; i++ {
		if term.Screen[0][i].Char != ' ' {
			t.Errorf("expected space at position %d after erase", i)
		}
	}
}

func TestEraseLine(t *testing.T) {
	term := New(80, 24)
	term.Write("AAAAA\x1b[3D\x1b[K")

	// First 2 chars should remain, rest erased
	if term.Screen[0][0].Char != 'A' {
		t.Error("expected 'A' at position 0")
	}
	if term.Screen[0][1].Char != 'A' {
		t.Error("expected 'A' at position 1")
	}
	if term.Screen[0][2].Char != ' ' {
		t.Error("expected space at position 2 after erase")
	}
}

func TestScrollUp(t *testing.T) {
	term := New(80, 3)
	term.Write("Line1\r\nLine2\r\nLine3\r\nLine4")

	// After scroll, Line1 should be gone, Line2 is now at row 0
	if term.Screen[0][0].Char != 'L' {
		t.Error("expected 'L' at row 0 after scroll")
	}
}

func TestGet256Color(t *testing.T) {
	// Test basic colors (0-15)
	c := Get256Color(1)
	if c != defaultColors[1] {
		t.Errorf("expected red, got %v", c)
	}

	// Test grayscale (232-255)
	gray := Get256Color(240)
	if gray.R != gray.G || gray.G != gray.B {
		t.Error("expected grayscale color")
	}
}

func TestInverseVideo(t *testing.T) {
	term := New(80, 24)
	term.Write("\x1b[7mInverse\x1b[27m")

	// In inverse mode, FG and BG should be swapped
	cell := term.Screen[0][0]
	theme, _ := GetTheme("default")
	if cell.FG == theme.Foreground {
		t.Error("expected inverted foreground")
	}
}

func TestDefaultColors(t *testing.T) {
	if len(defaultColors) != 16 {
		t.Errorf("expected 16 default colors, got %d", len(defaultColors))
	}

	// Check black is actually black
	black := defaultColors[0]
	if black.R != 0 || black.G != 0 || black.B != 0 {
		t.Error("expected black to be (0,0,0)")
	}

	// Check white is actually white-ish
	white := defaultColors[15]
	if white.R != 255 || white.G != 255 || white.B != 255 {
		t.Error("expected bright white to be (255,255,255)")
	}
}

func TestColorRGBA(t *testing.T) {
	c := color.RGBA{255, 0, 0, 255}
	r, g, b, a := c.RGBA()
	if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 || a>>8 != 255 {
		t.Error("RGBA conversion failed")
	}
}
