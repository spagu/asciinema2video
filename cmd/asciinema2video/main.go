package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spagu/asciinema2video/internal/cast"
	"github.com/spagu/asciinema2video/internal/renderer"
	"github.com/spagu/asciinema2video/internal/terminal"
	"github.com/spagu/asciinema2video/internal/video"
	"github.com/spf13/cobra"
)

var (
	version    = "1.0.0"
	inputFile  string
	outputFile string
	fontSize   int
	fps        int
	themeName  string
	themeFile  string
	listThemes bool
	cols       int
	rows       int
	fontFile   string

	// Border options
	borderEnabled bool
	borderWidth   int
	borderColor   string
	borderRadius  int

	// Background options
	outerBackground string
	transparent     bool
	padding         int

	// Output options
	quiet bool
	codec string
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "asciinema2video",
		Short:   "Convert asciinema .cast files to video",
		Long:    "A simple tool to convert asciinema recordings (.cast) to MP4, GIF, WebP, WebM, or MOV video files.",
		Version: version,
		Example: `  asciinema2video -i demo.cast -o demo.mp4
  asciinema2video -i demo.cast -o demo.gif
  asciinema2video -i demo.cast -o demo.webm --transparent
  asciinema2video -i demo.cast -o demo.mov --transparent
  asciinema2video -i demo.cast -o demo.mp4 --theme dracula
  asciinema2video -i demo.cast -o demo.mp4 --codec h265
  asciinema2video -i demo.cast -o demo.mp4 --border --border-radius 12
  asciinema2video -i demo.cast -o demo.webm --border --border-radius 20 --transparent
  asciinema2video -i demo.cast -o demo.mp4 --quiet
  asciinema2video --list-themes`,
		RunE: run,
	}

	// Input/Output
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "input .cast file (required)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output file (mp4, gif, webp, webm, mov)")

	// Rendering
	rootCmd.Flags().IntVar(&fontSize, "font-size", 14, "font size in pixels")
	rootCmd.Flags().IntVar(&fps, "fps", 10, "frames per second")
	rootCmd.Flags().StringVar(&fontFile, "font", "", "path to custom TTF font file (default: Go Mono)")
	rootCmd.Flags().IntVar(&padding, "padding", 16, "padding around terminal content in pixels")

	// Terminal size
	rootCmd.Flags().IntVar(&cols, "cols", 0, "override terminal width (0 = use from cast file)")
	rootCmd.Flags().IntVar(&rows, "rows", 0, "override terminal height (0 = use from cast file)")

	// Theme
	rootCmd.Flags().StringVar(&themeName, "theme", "default", "color theme name")
	rootCmd.Flags().StringVar(&themeFile, "theme-file", "", "path to custom theme JSON file")
	rootCmd.Flags().BoolVar(&listThemes, "list-themes", false, "list available themes")

	// Border
	rootCmd.Flags().BoolVar(&borderEnabled, "border", false, "enable border around terminal")
	rootCmd.Flags().IntVar(&borderWidth, "border-width", 2, "border width in pixels")
	rootCmd.Flags().StringVar(&borderColor, "border-color", "#646464", "border color (hex)")
	rootCmd.Flags().IntVar(&borderRadius, "border-radius", 0, "border radius for rounded corners")

	// Background
	rootCmd.Flags().StringVar(&outerBackground, "outer-bg", "#000000", "outer background color (hex, or 'transparent')")
	rootCmd.Flags().BoolVar(&transparent, "transparent", false, "enable transparent background (for webm, mov)")

	// Video codec
	rootCmd.Flags().StringVar(&codec, "codec", "h264", "video codec for MP4 (h264, h265)")

	// Output control
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "suppress output messages")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func logf(format string, args ...interface{}) {
	if !quiet {
		fmt.Printf(format, args...)
	}
}

func logln(args ...interface{}) {
	if !quiet {
		fmt.Println(args...)
	}
}

func run(cmd *cobra.Command, args []string) error {
	startTime := time.Now()

	if listThemes {
		themes := terminal.ListThemes()
		sort.Strings(themes)
		fmt.Println("Available themes:")
		for _, t := range themes {
			fmt.Printf("  - %s\n", t)
		}
		return nil
	}

	if inputFile == "" {
		return fmt.Errorf("input file is required")
	}

	if outputFile == "" {
		base := strings.TrimSuffix(filepath.Base(inputFile), ".cast")
		outputFile = base + ".mp4"
	}

	// Validate codec
	codec = strings.ToLower(codec)
	if codec != "h264" && codec != "h265" {
		return fmt.Errorf("invalid codec: %s (supported: h264, h265)", codec)
	}

	// Load theme
	var theme *terminal.Theme
	var err error

	if themeFile != "" {
		theme, err = terminal.LoadThemeFromFile(themeFile)
		if err != nil {
			return fmt.Errorf("error loading theme file: %w", err)
		}
		logf("Using custom theme from: %s\n", themeFile)
	} else {
		theme, err = terminal.GetTheme(themeName)
		if err != nil {
			return fmt.Errorf("error loading theme: %w", err)
		}
		logf("Using theme: %s\n", themeName)
	}

	recording, err := cast.Parse(inputFile)
	if err != nil {
		return fmt.Errorf("error parsing cast file: %w", err)
	}

	// Use custom dimensions or from cast file
	termWidth := recording.Header.Width
	termHeight := recording.Header.Height
	if cols > 0 {
		termWidth = cols
	}
	if rows > 0 {
		termHeight = rows
	}

	logf("Recording: %dx%d, %d events\n", termWidth, termHeight, len(recording.Events))

	// Build renderer options
	opts := renderer.DefaultOptions()
	opts.TermWidth = termWidth
	opts.TermHeight = termHeight
	opts.FontSize = fontSize
	opts.FontPath = fontFile
	opts.Theme = theme
	opts.Padding = padding

	// Border options
	opts.BorderEnabled = borderEnabled
	opts.BorderWidth = borderWidth
	opts.BorderColor = parseHexColor(borderColor)
	opts.BorderRadius = borderRadius

	// Background options
	opts.Transparent = transparent
	if transparent {
		opts.OuterBackground = color.RGBA{0, 0, 0, 0}
	} else if outerBackground == "transparent" {
		opts.OuterBackground = color.RGBA{0, 0, 0, 0}
		opts.Transparent = true
	} else {
		opts.OuterBackground = parseHexColor(outerBackground)
	}

	r, err := renderer.NewFromOptions(opts)
	if err != nil {
		return fmt.Errorf("error creating renderer: %w", err)
	}

	tempDir, err := os.MkdirTemp("", "asciinema2video-*")
	if err != nil {
		return fmt.Errorf("error creating temp dir: %w", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	logln("Rendering frames...")
	framePaths, err := r.RenderFrames(recording, tempDir, fps)
	if err != nil {
		return fmt.Errorf("error rendering frames: %w", err)
	}

	logf("Generated %d frames\n", len(framePaths))

	logf("Creating %s...\n", outputFile)

	videoOpts := &video.CreateOptions{
		Transparent: opts.Transparent,
		Codec:       codec,
	}
	err = video.CreateWithOptions(tempDir, outputFile, fps, videoOpts)
	if err != nil {
		log.Fatalf("Error creating video: %v", err)
	}

	// Get output file info
	elapsed := time.Since(startTime)
	fileInfo, err := os.Stat(outputFile)
	if err != nil {
		return fmt.Errorf("error getting output file info: %w", err)
	}

	logf("Done! Output: %s (%s) in %s\n", outputFile, formatFileSize(fileInfo.Size()), formatDuration(elapsed))
	return nil
}

func parseHexColor(hex string) color.RGBA {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return color.RGBA{0, 0, 0, 255}
	}

	r, _ := strconv.ParseUint(hex[0:2], 16, 8)
	g, _ := strconv.ParseUint(hex[2:4], 16, 8)
	b, _ := strconv.ParseUint(hex[4:6], 16, 8)

	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}

func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
}
