# Styling Options

asciinema2video provides extensive styling options to customize the appearance of your terminal recordings.

## Font Settings

### Font Size

Control the size of terminal text:

```bash
# Default: 14 pixels
asciinema2video -i demo.cast -o output.mp4 --font-size 14

# Larger text for presentations
asciinema2video -i demo.cast -o output.mp4 --font-size 20

# Smaller text for compact output
asciinema2video -i demo.cast -o output.mp4 --font-size 12
```

### Custom Font

Use any TrueType font (.ttf):

```bash
# Use system font
asciinema2video -i demo.cast -o output.mp4 --font /usr/share/fonts/truetype/firacode/FiraCode-Regular.ttf

# Use downloaded font
asciinema2video -i demo.cast -o output.mp4 --font ./JetBrainsMono-Regular.ttf
```

**Recommended monospace fonts:**
- Go Mono (default)
- JetBrains Mono
- Fira Code
- Source Code Pro
- Cascadia Code
- Hack

## Terminal Size

### Override Dimensions

Override the terminal dimensions from the cast file:

```bash
# Standard terminal size (80x24)
asciinema2video -i demo.cast -o output.mp4 --cols 80 --rows 24

# Wide terminal
asciinema2video -i demo.cast -o output.mp4 --cols 120 --rows 30

# Custom dimensions
asciinema2video -i demo.cast -o output.mp4 --cols 100 --rows 40
```

## Padding

Control the space around terminal content:

```bash
# Default: 16 pixels
asciinema2video -i demo.cast -o output.mp4 --padding 16

# More padding for a spacious look
asciinema2video -i demo.cast -o output.mp4 --padding 32

# Minimal padding
asciinema2video -i demo.cast -o output.mp4 --padding 8
```

## Border

### Enable Border

```bash
asciinema2video -i demo.cast -o output.mp4 --border
```

### Border Width

```bash
# Thin border (1px)
asciinema2video -i demo.cast -o output.mp4 --border --border-width 1

# Default border (2px)
asciinema2video -i demo.cast -o output.mp4 --border --border-width 2

# Thick border (4px)
asciinema2video -i demo.cast -o output.mp4 --border --border-width 4
```

### Border Color

```bash
# Gray border (default)
asciinema2video -i demo.cast -o output.mp4 --border --border-color "#646464"

# Accent color border
asciinema2video -i demo.cast -o output.mp4 --border --border-color "#bd93f9"

# White border
asciinema2video -i demo.cast -o output.mp4 --border --border-color "#ffffff"
```

## Rounded Corners

### Border Radius

```bash
# No rounding (default)
asciinema2video -i demo.cast -o output.mp4 --border-radius 0

# Subtle rounding
asciinema2video -i demo.cast -o output.mp4 --border --border-radius 8

# Medium rounding
asciinema2video -i demo.cast -o output.mp4 --border --border-radius 12

# Large rounding
asciinema2video -i demo.cast -o output.mp4 --border --border-radius 20
```

**Note:** Rounded corners work best with border enabled.

## Background

### Outer Background

The area outside the terminal (visible with rounded corners):

```bash
# Black (default)
asciinema2video -i demo.cast -o output.mp4 --border --border-radius 12 --outer-bg "#000000"

# Match your website/presentation
asciinema2video -i demo.cast -o output.mp4 --border --border-radius 12 --outer-bg "#1a1a2e"

# Transparent
asciinema2video -i demo.cast -o output.webm --border --border-radius 12 --outer-bg transparent
```

### Transparent Background

Enable alpha channel for compositing:

```bash
# WebM with transparency
asciinema2video -i demo.cast -o output.webm --transparent

# MOV with ProRes 4444 alpha
asciinema2video -i demo.cast -o output.mov --transparent

# Combined with rounded corners
asciinema2video -i demo.cast -o output.webm --transparent --border --border-radius 20
```

**Supported formats for transparency:**
- WebM (VP9 with alpha)
- MOV (ProRes 4444)

## FPS (Frames Per Second)

Control playback smoothness:

```bash
# Default: 10 FPS
asciinema2video -i demo.cast -o output.mp4 --fps 10

# Smoother playback
asciinema2video -i demo.cast -o output.mp4 --fps 15

# High quality
asciinema2video -i demo.cast -o output.mp4 --fps 30

# Cinema quality
asciinema2video -i demo.cast -o output.mp4 --fps 60
```

**Trade-offs:**
- Higher FPS = smoother playback, larger file size
- Lower FPS = smaller file size, may look choppy

## Complete Examples

### Documentation Screenshot Style

```bash
asciinema2video -i demo.cast -o demo.mp4 \
  --theme default \
  --font-size 14 \
  --cols 80 \
  --rows 24 \
  --padding 16 \
  --border \
  --border-radius 8
```

### Presentation Style

```bash
asciinema2video -i demo.cast -o demo.mp4 \
  --theme dracula \
  --font-size 20 \
  --fps 15 \
  --padding 32 \
  --border \
  --border-width 3 \
  --border-color "#bd93f9" \
  --border-radius 16
```

### Web Overlay Style

```bash
asciinema2video -i demo.cast -o demo.webm \
  --theme nord \
  --transparent \
  --border \
  --border-radius 20 \
  --padding 24
```

### Professional Video Editing

```bash
asciinema2video -i demo.cast -o demo.mov \
  --transparent \
  --fps 30 \
  --font-size 18 \
  --border \
  --border-radius 12
```
