# Color Themes

asciinema2video supports customizable color themes for terminal rendering.

## Theme Gallery

| Theme | Description | Demo |
|-------|-------------|------|
| `default` | Dark theme with standard terminal colors | [demo-default.webm](../examples/themes/demo-default.webm) |
| `monokai` | Popular code editor color scheme | [demo-monokai.webm](../examples/themes/demo-monokai.webm) |
| `dracula` | Dark purple-based color scheme | [demo-dracula.webm](../examples/themes/demo-dracula.webm) |
| `solarized-dark` | Ethan Schoonover's Solarized palette | [demo-solarized.webm](../examples/themes/demo-solarized.webm) |
| `nord` | Arctic, north-bluish color palette | [demo-nord.webm](../examples/themes/demo-nord.webm) |
| `gruvbox` | Retro groove color scheme | [demo-gruvbox.webm](../examples/themes/demo-gruvbox.webm) |

## Built-in Themes

| Theme | Description |
|-------|-------------|
| `default` | Dark theme with standard terminal colors |
| `monokai` | Popular code editor color scheme |
| `dracula` | Dark purple-based color scheme |
| `solarized-dark` | Ethan Schoonover's Solarized palette |
| `nord` | Arctic, north-bluish color palette |
| `gruvbox` | Retro groove color scheme |

### Usage

```bash
# Use built-in theme
asciinema2video -i demo.cast -o output.mp4 --theme dracula

# List available themes
asciinema2video --list-themes
```

## Custom Themes

Create a JSON file with your custom color scheme.

### Theme File Structure

```json
{
  "name": "my-custom-theme",
  "foreground": "#e0e0e0",
  "background": "#1e1e1e",
  "colors": [
    "#000000",
    "#ff0000",
    "#00ff00",
    "#ffff00",
    "#0000ff",
    "#ff00ff",
    "#00ffff",
    "#ffffff",
    "#808080",
    "#ff8080",
    "#80ff80",
    "#ffff80",
    "#8080ff",
    "#ff80ff",
    "#80ffff",
    "#ffffff"
  ]
}
```

### Color Palette

The `colors` array contains 16 ANSI colors:

| Index | Name | ANSI Code | Usage |
|-------|------|-----------|-------|
| 0 | Black | 30/40 | Background, borders |
| 1 | Red | 31/41 | Errors, deletions |
| 2 | Green | 32/42 | Success, additions |
| 3 | Yellow | 33/43 | Warnings, highlights |
| 4 | Blue | 34/44 | Info, links |
| 5 | Magenta | 35/45 | Special, keywords |
| 6 | Cyan | 36/46 | Secondary info |
| 7 | White | 37/47 | Default text |
| 8 | Bright Black | 90/100 | Comments, muted text |
| 9 | Bright Red | 91/101 | Bold errors |
| 10 | Bright Green | 92/102 | Bold success |
| 11 | Bright Yellow | 93/103 | Bold warnings |
| 12 | Bright Blue | 94/104 | Bold info |
| 13 | Bright Magenta | 95/105 | Bold special |
| 14 | Bright Cyan | 96/106 | Bold secondary |
| 15 | Bright White | 97/107 | Bold text |

### Usage

```bash
asciinema2video -i demo.cast -o output.mp4 --theme-file my-theme.json
```

## Color Format

All colors use hexadecimal format:

- `#RRGGBB` - Standard 6-character hex (e.g., `#ff5555`)
- `#RGB` is NOT supported, use full 6 characters

## Theme Examples

### One Dark

```json
{
  "name": "one-dark",
  "foreground": "#abb2bf",
  "background": "#282c34",
  "colors": [
    "#282c34", "#e06c75", "#98c379", "#e5c07b",
    "#61afef", "#c678dd", "#56b6c2", "#abb2bf",
    "#5c6370", "#e06c75", "#98c379", "#e5c07b",
    "#61afef", "#c678dd", "#56b6c2", "#ffffff"
  ]
}
```

### Tokyo Night

```json
{
  "name": "tokyo-night",
  "foreground": "#a9b1d6",
  "background": "#1a1b26",
  "colors": [
    "#15161e", "#f7768e", "#9ece6a", "#e0af68",
    "#7aa2f7", "#bb9af7", "#7dcfff", "#a9b1d6",
    "#414868", "#f7768e", "#9ece6a", "#e0af68",
    "#7aa2f7", "#bb9af7", "#7dcfff", "#c0caf5"
  ]
}
```

### Catppuccin Mocha

```json
{
  "name": "catppuccin-mocha",
  "foreground": "#cdd6f4",
  "background": "#1e1e2e",
  "colors": [
    "#45475a", "#f38ba8", "#a6e3a1", "#f9e2af",
    "#89b4fa", "#f5c2e7", "#94e2d5", "#bac2de",
    "#585b70", "#f38ba8", "#a6e3a1", "#f9e2af",
    "#89b4fa", "#f5c2e7", "#94e2d5", "#a6adc8"
  ]
}
```

## Tips

1. **Contrast**: Ensure good contrast between foreground and background
2. **Consistency**: Keep bright colors visually related to their normal variants
3. **Testing**: Test your theme with various terminal output to verify readability
4. **WCAG**: Consider WCAG 2.2 contrast requirements for accessibility
