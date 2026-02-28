# Examples

## Custom Theme

The `theme.json` file demonstrates how to create a custom color theme.

### Theme Structure

```json
{
  "name": "my-theme",
  "foreground": "#f8f8f2",
  "background": "#282a36",
  "colors": [
    "#000000",  // 0: Black
    "#ff0000",  // 1: Red
    "#00ff00",  // 2: Green
    "#ffff00",  // 3: Yellow
    "#0000ff",  // 4: Blue
    "#ff00ff",  // 5: Magenta
    "#00ffff",  // 6: Cyan
    "#ffffff",  // 7: White
    "#808080",  // 8: Bright Black (Gray)
    "#ff8080",  // 9: Bright Red
    "#80ff80",  // 10: Bright Green
    "#ffff80",  // 11: Bright Yellow
    "#8080ff",  // 12: Bright Blue
    "#ff80ff",  // 13: Bright Magenta
    "#80ffff",  // 14: Bright Cyan
    "#ffffff"   // 15: Bright White
  ]
}
```

### Fields

| Field | Description |
|-------|-------------|
| `name` | Theme name (for display purposes) |
| `foreground` | Default text color (hex format) |
| `background` | Terminal background color (hex format) |
| `colors` | Array of 16 ANSI colors (hex format) |

### Color Order

The colors array must contain exactly 16 colors in ANSI order:

- **0-7**: Standard colors (Black, Red, Green, Yellow, Blue, Magenta, Cyan, White)
- **8-15**: Bright/bold variants of the above colors

### Usage

```bash
asciinema2video -i recording.cast -o output.mp4 --theme-file examples/theme.json
```

## Popular Theme Examples

### Dracula

```json
{
  "name": "dracula",
  "foreground": "#f8f8f2",
  "background": "#282a36",
  "colors": [
    "#21222c", "#ff5555", "#50fa7b", "#f1fa8c",
    "#bd93f9", "#ff79c6", "#8be9fd", "#f8f8f2",
    "#6272a4", "#ff6e67", "#5af78e", "#f4f99d",
    "#caa9fa", "#ff92d0", "#9aedfe", "#ffffff"
  ]
}
```

### Nord

```json
{
  "name": "nord",
  "foreground": "#d8dee9",
  "background": "#2e3440",
  "colors": [
    "#3b4252", "#bf616a", "#a3be8c", "#ebcb8b",
    "#81a1c1", "#b48ead", "#88c0d0", "#e5e9f0",
    "#4c566a", "#bf616a", "#a3be8c", "#ebcb8b",
    "#81a1c1", "#b48ead", "#8fbcbb", "#eceff4"
  ]
}
```

### Solarized Dark

```json
{
  "name": "solarized-dark",
  "foreground": "#839496",
  "background": "#002b36",
  "colors": [
    "#073642", "#dc322f", "#859900", "#b58900",
    "#268bd2", "#d33682", "#2aa198", "#eee8d5",
    "#002b36", "#cb4b16", "#586e75", "#657b83",
    "#839496", "#6c71c4", "#93a1a1", "#fdf6e3"
  ]
}
```

## Border and Styling Examples

### Rounded corners with border

```bash
asciinema2video -i recording.cast -o output.mp4 \
  --border \
  --border-width 3 \
  --border-color "#ff79c6" \
  --border-radius 15
```

### Transparent background (for video overlays)

```bash
asciinema2video -i recording.cast -o output.webm \
  --transparent \
  --border \
  --border-radius 20
```

### Custom outer background color

```bash
asciinema2video -i recording.cast -o output.mp4 \
  --border \
  --border-radius 12 \
  --outer-bg "#1a1a2e"
```

### High-quality video for professional editing

```bash
asciinema2video -i recording.cast -o output.mov \
  --transparent \
  --fps 30 \
  --font-size 18
```

## Complete Example

```bash
asciinema2video -i demo.cast -o demo.mp4 \
  --theme dracula \
  --font-size 16 \
  --fps 15 \
  --padding 20 \
  --border \
  --border-width 2 \
  --border-color "#bd93f9" \
  --border-radius 12 \
  --cols 80 \
  --rows 24
```

## Output Format Comparison

| Format | Extension | Alpha | Best For |
|--------|-----------|-------|----------|
| MP4 | `.mp4` | No | General sharing, most compatible |
| GIF | `.gif` | No | Quick previews, documentation |
| WebP | `.webp` | No | Web, smaller than GIF |
| WebM | `.webm` | Yes | Web with transparency |
| MOV | `.mov` | Yes | Professional video editing (ProRes 4444) |
