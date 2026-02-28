package terminal

import (
	"image/color"
	"regexp"
	"strconv"
	"strings"
)

// Cell represents a single character cell
type Cell struct {
	Char rune
	FG   color.RGBA
	BG   color.RGBA
	Bold bool
}

// Terminal emulates a VT100-compatible terminal
type Terminal struct {
	Width   int
	Height  int
	Screen  [][]Cell
	cursorX int
	cursorY int
	fg      color.RGBA
	bg      color.RGBA
	bold    bool
	inverse bool
	theme   *Theme
}

// Default 256 color palette (basic 16 colors)
var defaultColors = []color.RGBA{
	{0, 0, 0, 255},       // 0: Black
	{205, 49, 49, 255},   // 1: Red
	{13, 188, 121, 255},  // 2: Green
	{229, 229, 16, 255},  // 3: Yellow
	{36, 114, 200, 255},  // 4: Blue
	{188, 63, 188, 255},  // 5: Magenta
	{17, 168, 205, 255},  // 6: Cyan
	{229, 229, 229, 255}, // 7: White
	{102, 102, 102, 255}, // 8: Bright Black
	{241, 76, 76, 255},   // 9: Bright Red
	{35, 209, 139, 255},  // 10: Bright Green
	{245, 245, 67, 255},  // 11: Bright Yellow
	{59, 142, 234, 255},  // 12: Bright Blue
	{214, 112, 214, 255}, // 13: Bright Magenta
	{41, 184, 219, 255},  // 14: Bright Cyan
	{255, 255, 255, 255}, // 15: Bright White
}

var escapeRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]|\x1b\][^\x07]*\x07|\x1b\[\?[0-9;]*[hlsr]|\x1b[=>]|\x1b\][^\x1b]*\x1b\\`)

// New creates a new terminal with given dimensions and default theme
func New(width, height int) *Terminal {
	theme, _ := GetTheme("default")
	return NewWithTheme(width, height, theme)
}

// NewWithTheme creates a new terminal with given dimensions and theme
func NewWithTheme(width, height int, theme *Theme) *Terminal {
	if theme == nil {
		theme, _ = GetTheme("default")
	}
	t := &Terminal{
		Width:  width,
		Height: height,
		fg:     theme.Foreground,
		bg:     theme.Background,
		theme:  theme,
	}
	t.Screen = make([][]Cell, height)
	for i := range t.Screen {
		t.Screen[i] = make([]Cell, width)
		for j := range t.Screen[i] {
			t.Screen[i][j] = Cell{Char: ' ', FG: theme.Foreground, BG: theme.Background}
		}
	}
	return t
}

// Write processes output data
func (t *Terminal) Write(data string) {
	i := 0
	for i < len(data) {
		if data[i] == '\x1b' {
			// Find escape sequence
			match := escapeRegex.FindString(data[i:])
			if match != "" {
				t.processEscape(match)
				i += len(match)
				continue
			}
		}

		r := rune(data[i])
		switch r {
		case '\r':
			t.cursorX = 0
		case '\n':
			t.cursorY++
			if t.cursorY >= t.Height {
				t.scrollUp()
				t.cursorY = t.Height - 1
			}
		case '\b':
			if t.cursorX > 0 {
				t.cursorX--
			}
		case '\t':
			t.cursorX = ((t.cursorX / 8) + 1) * 8
			if t.cursorX >= t.Width {
				t.cursorX = t.Width - 1
			}
		case '\x07': // Bell - ignore
		default:
			if r >= 32 {
				t.putChar(r)
			}
		}
		i++
	}
}

func (t *Terminal) putChar(r rune) {
	if t.cursorX >= t.Width {
		t.cursorX = 0
		t.cursorY++
		if t.cursorY >= t.Height {
			t.scrollUp()
			t.cursorY = t.Height - 1
		}
	}

	fg, bg := t.fg, t.bg
	if t.inverse {
		fg, bg = bg, fg
	}

	t.Screen[t.cursorY][t.cursorX] = Cell{
		Char: r,
		FG:   fg,
		BG:   bg,
		Bold: t.bold,
	}
	t.cursorX++
}

func (t *Terminal) scrollUp() {
	copy(t.Screen, t.Screen[1:])
	t.Screen[t.Height-1] = make([]Cell, t.Width)
	for j := range t.Screen[t.Height-1] {
		t.Screen[t.Height-1][j] = Cell{Char: ' ', FG: t.theme.Foreground, BG: t.theme.Background}
	}
}

func (t *Terminal) processEscape(seq string) {
	if len(seq) < 2 {
		return
	}

	// CSI sequences
	if strings.HasPrefix(seq, "\x1b[") {
		t.processCSI(seq[2:])
	}
}

func (t *Terminal) processCSI(seq string) {
	if len(seq) == 0 {
		return
	}

	// Get command (last character)
	cmd := seq[len(seq)-1]
	params := seq[:len(seq)-1]

	// Handle ? prefix
	if strings.HasPrefix(params, "?") {
		return // Ignore private mode sequences
	}

	switch cmd {
	case 'm': // SGR - Select Graphic Rendition
		t.processSGR(params)
	case 'H', 'f': // Cursor position
		t.processCursorPos(params)
	case 'A': // Cursor up
		n := parseNum(params, 1)
		t.cursorY -= n
		if t.cursorY < 0 {
			t.cursorY = 0
		}
	case 'B': // Cursor down
		n := parseNum(params, 1)
		t.cursorY += n
		if t.cursorY >= t.Height {
			t.cursorY = t.Height - 1
		}
	case 'C': // Cursor forward
		n := parseNum(params, 1)
		t.cursorX += n
		if t.cursorX >= t.Width {
			t.cursorX = t.Width - 1
		}
	case 'D': // Cursor back
		n := parseNum(params, 1)
		t.cursorX -= n
		if t.cursorX < 0 {
			t.cursorX = 0
		}
	case 'J': // Erase in display
		n := parseNum(params, 0)
		t.eraseDisplay(n)
	case 'K': // Erase in line
		n := parseNum(params, 0)
		t.eraseLine(n)
	}
}

func (t *Terminal) processSGR(params string) {
	if params == "" {
		t.resetSGR()
		return
	}

	parts := strings.Split(params, ";")
	i := 0
	for i < len(parts) {
		n, _ := strconv.Atoi(parts[i])
		switch {
		case n == 0:
			t.resetSGR()
		case n == 1:
			t.bold = true
		case n == 7:
			t.inverse = true
		case n == 22:
			t.bold = false
		case n == 27:
			t.inverse = false
		case n >= 30 && n <= 37:
			t.fg = t.theme.Colors[n-30]
		case n == 38:
			// Extended foreground color
			if i+2 < len(parts) {
				mode, _ := strconv.Atoi(parts[i+1])
				if mode == 5 && i+2 < len(parts) {
					colorNum, _ := strconv.Atoi(parts[i+2])
					t.fg = t.get256Color(colorNum)
					i += 2
				}
			}
		case n == 39:
			t.fg = t.theme.Foreground
		case n >= 40 && n <= 47:
			t.bg = t.theme.Colors[n-40]
		case n == 48:
			// Extended background color
			if i+2 < len(parts) {
				mode, _ := strconv.Atoi(parts[i+1])
				if mode == 5 && i+2 < len(parts) {
					colorNum, _ := strconv.Atoi(parts[i+2])
					t.bg = t.get256Color(colorNum)
					i += 2
				}
			}
		case n == 49:
			t.bg = t.theme.Background
		case n >= 90 && n <= 97:
			t.fg = t.theme.Colors[n-90+8]
		case n >= 100 && n <= 107:
			t.bg = t.theme.Colors[n-100+8]
		}
		i++
	}
}

func (t *Terminal) resetSGR() {
	t.fg = t.theme.Foreground
	t.bg = t.theme.Background
	t.bold = false
	t.inverse = false
}

func (t *Terminal) processCursorPos(params string) {
	parts := strings.Split(params, ";")
	row, col := 1, 1
	if len(parts) >= 1 && parts[0] != "" {
		row, _ = strconv.Atoi(parts[0])
	}
	if len(parts) >= 2 && parts[1] != "" {
		col, _ = strconv.Atoi(parts[1])
	}
	t.cursorY = row - 1
	t.cursorX = col - 1
	if t.cursorY < 0 {
		t.cursorY = 0
	}
	if t.cursorY >= t.Height {
		t.cursorY = t.Height - 1
	}
	if t.cursorX < 0 {
		t.cursorX = 0
	}
	if t.cursorX >= t.Width {
		t.cursorX = t.Width - 1
	}
}

func (t *Terminal) eraseDisplay(mode int) {
	switch mode {
	case 0: // Erase from cursor to end
		for x := t.cursorX; x < t.Width; x++ {
			t.Screen[t.cursorY][x] = Cell{Char: ' ', FG: t.theme.Foreground, BG: t.theme.Background}
		}
		for y := t.cursorY + 1; y < t.Height; y++ {
			for x := 0; x < t.Width; x++ {
				t.Screen[y][x] = Cell{Char: ' ', FG: t.theme.Foreground, BG: t.theme.Background}
			}
		}
	case 1: // Erase from start to cursor
		for y := 0; y < t.cursorY; y++ {
			for x := 0; x < t.Width; x++ {
				t.Screen[y][x] = Cell{Char: ' ', FG: t.theme.Foreground, BG: t.theme.Background}
			}
		}
		for x := 0; x <= t.cursorX; x++ {
			t.Screen[t.cursorY][x] = Cell{Char: ' ', FG: t.theme.Foreground, BG: t.theme.Background}
		}
	case 2: // Erase entire display
		for y := 0; y < t.Height; y++ {
			for x := 0; x < t.Width; x++ {
				t.Screen[y][x] = Cell{Char: ' ', FG: t.theme.Foreground, BG: t.theme.Background}
			}
		}
	}
}

func (t *Terminal) eraseLine(mode int) {
	switch mode {
	case 0: // Erase from cursor to end of line
		for x := t.cursorX; x < t.Width; x++ {
			t.Screen[t.cursorY][x] = Cell{Char: ' ', FG: t.theme.Foreground, BG: t.theme.Background}
		}
	case 1: // Erase from start of line to cursor
		for x := 0; x <= t.cursorX; x++ {
			t.Screen[t.cursorY][x] = Cell{Char: ' ', FG: t.theme.Foreground, BG: t.theme.Background}
		}
	case 2: // Erase entire line
		for x := 0; x < t.Width; x++ {
			t.Screen[t.cursorY][x] = Cell{Char: ' ', FG: t.theme.Foreground, BG: t.theme.Background}
		}
	}
}

func parseNum(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

func (t *Terminal) get256Color(n int) color.RGBA {
	if n < 16 {
		return t.theme.Colors[n]
	}
	if n < 232 {
		// 216 color cube
		n -= 16
		r := uint8((n / 36) * 51)
		g := uint8(((n / 6) % 6) * 51)
		b := uint8((n % 6) * 51)
		return color.RGBA{r, g, b, 255}
	}
	// Grayscale: n is 232-255, result is 8-238
	grayVal := (n-232)*10 + 8
	if grayVal > 255 {
		grayVal = 255
	}
	return color.RGBA{uint8(grayVal), uint8(grayVal), uint8(grayVal), 255}
}
