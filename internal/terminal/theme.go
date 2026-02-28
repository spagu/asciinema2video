package terminal

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"
)

// Theme defines terminal colors
type Theme struct {
	Name       string       `json:"name"`
	Foreground color.RGBA   `json:"-"`
	Background color.RGBA   `json:"-"`
	Colors     []color.RGBA `json:"-"`

	// JSON fields for parsing
	ForegroundHex string   `json:"foreground"`
	BackgroundHex string   `json:"background"`
	ColorsHex     []string `json:"colors"`
}

// Predefined themes
var themes = map[string]*Theme{
	"default": {
		Name:       "default",
		Foreground: color.RGBA{229, 229, 229, 255},
		Background: color.RGBA{30, 30, 30, 255},
		Colors:     defaultColors,
	},
	"monokai": {
		Name:       "monokai",
		Foreground: color.RGBA{248, 248, 242, 255},
		Background: color.RGBA{39, 40, 34, 255},
		Colors: []color.RGBA{
			{39, 40, 34, 255},    // Black
			{249, 38, 114, 255},  // Red
			{166, 226, 46, 255},  // Green
			{244, 191, 117, 255}, // Yellow
			{102, 217, 239, 255}, // Blue
			{174, 129, 255, 255}, // Magenta
			{161, 239, 228, 255}, // Cyan
			{248, 248, 242, 255}, // White
			{117, 113, 94, 255},  // Bright Black
			{249, 38, 114, 255},  // Bright Red
			{166, 226, 46, 255},  // Bright Green
			{244, 191, 117, 255}, // Bright Yellow
			{102, 217, 239, 255}, // Bright Blue
			{174, 129, 255, 255}, // Bright Magenta
			{161, 239, 228, 255}, // Bright Cyan
			{248, 248, 242, 255}, // Bright White
		},
	},
	"dracula": {
		Name:       "dracula",
		Foreground: color.RGBA{248, 248, 242, 255},
		Background: color.RGBA{40, 42, 54, 255},
		Colors: []color.RGBA{
			{40, 42, 54, 255},    // Black
			{255, 85, 85, 255},   // Red
			{80, 250, 123, 255},  // Green
			{241, 250, 140, 255}, // Yellow
			{189, 147, 249, 255}, // Blue
			{255, 121, 198, 255}, // Magenta
			{139, 233, 253, 255}, // Cyan
			{248, 248, 242, 255}, // White
			{98, 114, 164, 255},  // Bright Black
			{255, 110, 103, 255}, // Bright Red
			{90, 247, 142, 255},  // Bright Green
			{244, 249, 157, 255}, // Bright Yellow
			{202, 169, 250, 255}, // Bright Blue
			{255, 146, 208, 255}, // Bright Magenta
			{154, 237, 254, 255}, // Bright Cyan
			{255, 255, 255, 255}, // Bright White
		},
	},
	"solarized-dark": {
		Name:       "solarized-dark",
		Foreground: color.RGBA{131, 148, 150, 255},
		Background: color.RGBA{0, 43, 54, 255},
		Colors: []color.RGBA{
			{7, 54, 66, 255},     // Black
			{220, 50, 47, 255},   // Red
			{133, 153, 0, 255},   // Green
			{181, 137, 0, 255},   // Yellow
			{38, 139, 210, 255},  // Blue
			{211, 54, 130, 255},  // Magenta
			{42, 161, 152, 255},  // Cyan
			{238, 232, 213, 255}, // White
			{0, 43, 54, 255},     // Bright Black
			{203, 75, 22, 255},   // Bright Red
			{88, 110, 117, 255},  // Bright Green
			{101, 123, 131, 255}, // Bright Yellow
			{131, 148, 150, 255}, // Bright Blue
			{108, 113, 196, 255}, // Bright Magenta
			{147, 161, 161, 255}, // Bright Cyan
			{253, 246, 227, 255}, // Bright White
		},
	},
	"nord": {
		Name:       "nord",
		Foreground: color.RGBA{216, 222, 233, 255},
		Background: color.RGBA{46, 52, 64, 255},
		Colors: []color.RGBA{
			{59, 66, 82, 255},    // Black
			{191, 97, 106, 255},  // Red
			{163, 190, 140, 255}, // Green
			{235, 203, 139, 255}, // Yellow
			{129, 161, 193, 255}, // Blue
			{180, 142, 173, 255}, // Magenta
			{136, 192, 208, 255}, // Cyan
			{229, 233, 240, 255}, // White
			{76, 86, 106, 255},   // Bright Black
			{191, 97, 106, 255},  // Bright Red
			{163, 190, 140, 255}, // Bright Green
			{235, 203, 139, 255}, // Bright Yellow
			{129, 161, 193, 255}, // Bright Blue
			{180, 142, 173, 255}, // Bright Magenta
			{143, 188, 187, 255}, // Bright Cyan
			{236, 239, 244, 255}, // Bright White
		},
	},
	"gruvbox": {
		Name:       "gruvbox",
		Foreground: color.RGBA{235, 219, 178, 255},
		Background: color.RGBA{40, 40, 40, 255},
		Colors: []color.RGBA{
			{40, 40, 40, 255},    // Black
			{204, 36, 29, 255},   // Red
			{152, 151, 26, 255},  // Green
			{215, 153, 33, 255},  // Yellow
			{69, 133, 136, 255},  // Blue
			{177, 98, 134, 255},  // Magenta
			{104, 157, 106, 255}, // Cyan
			{168, 153, 132, 255}, // White
			{146, 131, 116, 255}, // Bright Black
			{251, 73, 52, 255},   // Bright Red
			{184, 187, 38, 255},  // Bright Green
			{250, 189, 47, 255},  // Bright Yellow
			{131, 165, 152, 255}, // Bright Blue
			{211, 134, 155, 255}, // Bright Magenta
			{142, 192, 124, 255}, // Bright Cyan
			{235, 219, 178, 255}, // Bright White
		},
	},
}

// GetTheme returns a theme by name
func GetTheme(name string) (*Theme, error) {
	if theme, ok := themes[name]; ok {
		return theme, nil
	}
	return nil, fmt.Errorf("unknown theme: %s", name)
}

// LoadThemeFromFile loads a custom theme from JSON file
func LoadThemeFromFile(path string) (*Theme, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read theme file: %w", err)
	}

	var theme Theme
	if err := json.Unmarshal(data, &theme); err != nil {
		return nil, fmt.Errorf("failed to parse theme: %w", err)
	}

	// Parse hex colors
	theme.Foreground = parseHexColor(theme.ForegroundHex)
	theme.Background = parseHexColor(theme.BackgroundHex)

	theme.Colors = make([]color.RGBA, len(theme.ColorsHex))
	for i, hex := range theme.ColorsHex {
		theme.Colors[i] = parseHexColor(hex)
	}

	// Ensure we have 16 colors
	for len(theme.Colors) < 16 {
		theme.Colors = append(theme.Colors, defaultColors[len(theme.Colors)])
	}

	return &theme, nil
}

// ListThemes returns available theme names
func ListThemes() []string {
	names := make([]string, 0, len(themes))
	for name := range themes {
		names = append(names, name)
	}
	return names
}

// Get256Color returns a color from the 256 color palette
func Get256Color(n int) color.RGBA {
	if n < 16 {
		return defaultColors[n]
	}
	if n < 232 {
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

func parseHexColor(hex string) color.RGBA {
	if len(hex) == 0 {
		return color.RGBA{0, 0, 0, 255}
	}
	if hex[0] == '#' {
		hex = hex[1:]
	}
	if len(hex) != 6 {
		return color.RGBA{0, 0, 0, 255}
	}

	var r, g, b uint8
	_, _ = fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	return color.RGBA{r, g, b, 255}
}
