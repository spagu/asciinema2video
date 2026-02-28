# asciinema2video

[![Go Version](https://img.shields.io/badge/Go-1.26-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-BSD--3--Clause-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/spagu/asciinema2video)](https://github.com/spagu/asciinema2video/releases)

Simple CLI tool to convert [asciinema](https://asciinema.org/) recordings (`.cast` files) to video formats (MP4, GIF, WebP, WebM, MOV).

## Demo

<video src="examples/themes/demo-solarized.webm" autoplay loop muted playsinline></video>

*Solarized Dark theme with rounded corners*

## Features

- Multiple output formats: MP4, GIF, WebP, WebM, MOV (ProRes 4444)
- Transparent background support (WebM, MOV)
- Customizable color themes (6 built-in + custom JSON)
- Rounded corners with configurable radius
- Border support with custom color
- Custom TTF fonts
- Configurable terminal size, FPS, padding
- Video codec selection (H.264/H.265)

## Requirements

- **FFmpeg** - required for video encoding

### Installing FFmpeg

**Ubuntu/Debian:**
```bash
sudo apt update && sudo apt install ffmpeg
```

**macOS (Homebrew):**
```bash
brew install ffmpeg
```

**Windows (Chocolatey):**
```bash
choco install ffmpeg
```

**Windows (Scoop):**
```bash
scoop install ffmpeg
```

**Arch Linux:**
```bash
sudo pacman -S ffmpeg
```

**Fedora:**
```bash
sudo dnf install ffmpeg
```

**FreeBSD:**
```bash
pkg install ffmpeg
```

## Installation

### Quick install (curl)

```bash
curl -sSL https://raw.githubusercontent.com/spagu/asciinema2video/main/install.sh | sh
```

### From releases

Download the latest release for your platform from [GitHub Releases](https://github.com/spagu/asciinema2video/releases).

### Homebrew (macOS/Linux)

```bash
brew install spagu/tap/asciinema2video
```

### Snap (Ubuntu/Linux)

```bash
snap install asciinema2video
```

### DEB package (Ubuntu/Debian)

```bash
wget https://github.com/spagu/asciinema2video/releases/download/v1.0.0/asciinema2video_1.0.0_amd64.deb
sudo dpkg -i asciinema2video_1.0.0_amd64.deb
```

### RPM package (Fedora/RHEL/CentOS)

```bash
wget https://github.com/spagu/asciinema2video/releases/download/v1.0.0/asciinema2video_1.0.0_amd64.rpm
sudo rpm -i asciinema2video_1.0.0_amd64.rpm
```

### FreeBSD

```bash
# Download from releases
wget https://github.com/spagu/asciinema2video/releases/download/v1.0.0/asciinema2video_1.0.0_freebsd_amd64.tar.gz
tar xzf asciinema2video_1.0.0_freebsd_amd64.tar.gz
sudo mv asciinema2video /usr/local/bin/
```

### Go install

```bash
go install github.com/spagu/asciinema2video/cmd/asciinema2video@latest
```

### Build from source

```bash
git clone https://github.com/spagu/asciinema2video.git
cd asciinema2video
make build
```

## Usage

```bash
# Basic conversion to MP4
asciinema2video -i demo.cast -o demo.mp4

# Convert to GIF
asciinema2video -i demo.cast -o demo.gif

# Convert to WebP
asciinema2video -i demo.cast -o demo.webp

# Convert to WebM (with optional transparency)
asciinema2video -i demo.cast -o demo.webm
asciinema2video -i demo.cast -o demo.webm --transparent

# Convert to MOV (ProRes 4444 with alpha)
asciinema2video -i demo.cast -o demo.mov --transparent

# With custom theme
asciinema2video -i demo.cast -o demo.mp4 --theme dracula

# With rounded corners and border
asciinema2video -i demo.cast -o demo.mp4 --border --border-radius 12

# Transparent background with rounded corners (for overlays)
asciinema2video -i demo.cast -o demo.webm --border --border-radius 20 --transparent

# Custom settings
asciinema2video -i demo.cast -o demo.mp4 --fps 15 --font-size 16 --padding 20

# Standard terminal size (80x24)
asciinema2video -i demo.cast -o demo.mp4 --cols 80 --rows 24

# List available themes
asciinema2video --list-themes
```

## Flags

### Input/Output

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--input` | `-i` | | Input .cast file (required) |
| `--output` | `-o` | `{input}.mp4` | Output file (mp4, gif, webp, webm, mov) |

### Rendering

| Flag | Default | Description |
|------|---------|-------------|
| `--fps` | `10` | Frames per second |
| `--font-size` | `14` | Font size in pixels |
| `--font` | | Path to custom TTF font file |
| `--padding` | `16` | Padding around terminal content |
| `--cols` | `0` | Override terminal width (0 = from cast) |
| `--rows` | `0` | Override terminal height (0 = from cast) |

### Theme

| Flag | Default | Description |
|------|---------|-------------|
| `--theme` | `default` | Color theme name |
| `--theme-file` | | Path to custom theme JSON file |
| `--list-themes` | | List available themes |

### Border

| Flag | Default | Description |
|------|---------|-------------|
| `--border` | `false` | Enable border around terminal |
| `--border-width` | `2` | Border width in pixels |
| `--border-color` | `#646464` | Border color (hex) |
| `--border-radius` | `0` | Border radius for rounded corners |

### Background

| Flag | Default | Description |
|------|---------|-------------|
| `--outer-bg` | `#000000` | Outer background color (hex) |
| `--transparent` | `false` | Enable transparent background (webm, mov) |

### Video

| Flag | Default | Description |
|------|---------|-------------|
| `--codec` | `h264` | Video codec for MP4 (h264, h265) |

### Other

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--quiet` | `-q` | `false` | Suppress output messages |

## Themes

Built-in themes: `default`, `monokai`, `dracula`, `solarized-dark`, `nord`, `gruvbox`

See [Theme Gallery](examples/themes/) for all theme demos or [THEMES.md](docs/THEMES.md) for custom theme documentation.

Use with: `--theme dracula` or `--theme-file my-theme.json`

## Output Formats

| Format | Extension | Transparency | Notes |
|--------|-----------|--------------|-------|
| MP4 | `.mp4` | No | H.264, most compatible |
| GIF | `.gif` | No | Animated, large files |
| WebP | `.webp` | No | Animated, smaller than GIF |
| WebM | `.webm` | Yes | VP9, good for web |
| MOV | `.mov` | Yes | ProRes 4444, for video editing |

## How it works

1. Parses asciinema v2 format `.cast` file
2. Emulates VT100 terminal with ANSI color support
3. Renders terminal frames to PNG images
4. Uses FFmpeg to encode frames into video

## Supported features

- ANSI colors (16 basic + 256 extended)
- Bold text
- Cursor movement
- Screen/line clearing
- UTF-8 characters
- Custom color themes
- Rounded corners with transparency

## Dependencies

Minimal external dependencies:
- `github.com/spf13/cobra` - CLI framework
- `github.com/golang/freetype` - Font rendering
- `golang.org/x/image` - Go Mono font

## Development

```bash
# Build
make build

# Run tests
make test

# Format & lint
make check-all

# Build for all platforms
make build-all

# Create release
make release
```

## License

BSD-3-Clause
