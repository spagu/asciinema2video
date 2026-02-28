package renderer

import (
	"image/color"

	"github.com/spagu/asciinema2video/internal/terminal"
)

// Options contains all rendering options
type Options struct {
	TermWidth  int
	TermHeight int
	FontSize   int
	FontPath   string
	Theme      *terminal.Theme

	// Border options
	BorderEnabled bool
	BorderWidth   int
	BorderColor   color.RGBA
	BorderRadius  int // 0 = no rounding

	// Background options (outside terminal, for rounded corners)
	OuterBackground color.RGBA
	Transparent     bool // Use transparent background (for webm, mov)

	// Padding
	Padding int
}

// DefaultOptions returns options with sensible defaults
func DefaultOptions() *Options {
	theme, _ := terminal.GetTheme("default")
	return &Options{
		TermWidth:       80,
		TermHeight:      24,
		FontSize:        14,
		Theme:           theme,
		BorderEnabled:   false,
		BorderWidth:     2,
		BorderColor:     color.RGBA{100, 100, 100, 255},
		BorderRadius:    0,
		OuterBackground: color.RGBA{0, 0, 0, 0}, // Transparent
		Transparent:     false,
		Padding:         16,
	}
}

// WithBorder enables border with specified width and color
func (o *Options) WithBorder(width int, c color.RGBA) *Options {
	o.BorderEnabled = true
	o.BorderWidth = width
	o.BorderColor = c
	return o
}

// WithRoundedCorners sets border radius
func (o *Options) WithRoundedCorners(radius int) *Options {
	o.BorderRadius = radius
	return o
}

// WithOuterBackground sets background color outside terminal (visible with rounded corners)
func (o *Options) WithOuterBackground(c color.RGBA) *Options {
	o.OuterBackground = c
	return o
}

// WithTransparency enables transparent background
func (o *Options) WithTransparency(enabled bool) *Options {
	o.Transparent = enabled
	if enabled {
		o.OuterBackground = color.RGBA{0, 0, 0, 0}
	}
	return o
}
