package renderer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/spagu/asciinema2video/internal/cast"
	"github.com/spagu/asciinema2video/internal/terminal"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
)

// Renderer renders terminal frames to images
type Renderer struct {
	opts       *Options
	charWidth  int
	charHeight int
	font       *truetype.Font
}

// New creates a new renderer with default options
func New(termWidth, termHeight, fontSize int) (*Renderer, error) {
	opts := DefaultOptions()
	opts.TermWidth = termWidth
	opts.TermHeight = termHeight
	opts.FontSize = fontSize
	return NewFromOptions(opts)
}

// NewWithTheme creates a new renderer with specified theme
func NewWithTheme(termWidth, termHeight, fontSize int, theme *terminal.Theme) (*Renderer, error) {
	opts := DefaultOptions()
	opts.TermWidth = termWidth
	opts.TermHeight = termHeight
	opts.FontSize = fontSize
	opts.Theme = theme
	return NewFromOptions(opts)
}

// NewWithOptions creates a new renderer with basic options (backward compatible)
func NewWithOptions(termWidth, termHeight, fontSize int, theme *terminal.Theme, fontPath string) (*Renderer, error) {
	opts := DefaultOptions()
	opts.TermWidth = termWidth
	opts.TermHeight = termHeight
	opts.FontSize = fontSize
	opts.Theme = theme
	opts.FontPath = fontPath
	return NewFromOptions(opts)
}

// NewFromOptions creates a new renderer from Options struct
func NewFromOptions(opts *Options) (*Renderer, error) {
	var f *truetype.Font
	var err error

	if opts.FontPath != "" {
		fontData, err := os.ReadFile(opts.FontPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read font file: %w", err)
		}
		f, err = truetype.Parse(fontData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse font file: %w", err)
		}
	} else {
		f, err = truetype.Parse(gomono.TTF)
		if err != nil {
			return nil, fmt.Errorf("failed to parse font: %w", err)
		}
	}

	face := truetype.NewFace(f, &truetype.Options{
		Size: float64(opts.FontSize),
		DPI:  72,
	})
	charWidth := font.MeasureString(face, "M").Ceil()
	charHeight := int(float64(opts.FontSize) * 1.5)

	if opts.Theme == nil {
		opts.Theme, _ = terminal.GetTheme("default")
	}

	return &Renderer{
		opts:       opts,
		charWidth:  charWidth,
		charHeight: charHeight,
		font:       f,
	}, nil
}

// RenderFrames generates PNG frames from a recording
func (r *Renderer) RenderFrames(recording *cast.Recording, outputDir string, fps int) ([]string, error) {
	term := terminal.NewWithTheme(r.opts.TermWidth, r.opts.TermHeight, r.opts.Theme)

	frameInterval := 1.0 / float64(fps)
	var framePaths []string
	frameNum := 0
	nextFrameTime := 0.0

	eventIdx := 0
	maxTime := 0.0
	if len(recording.Events) > 0 {
		maxTime = recording.Events[len(recording.Events)-1].Time
	}

	for nextFrameTime <= maxTime {
		for eventIdx < len(recording.Events) && recording.Events[eventIdx].Time <= nextFrameTime {
			term.Write(recording.Events[eventIdx].Data)
			eventIdx++
		}

		framePath := filepath.Join(outputDir, fmt.Sprintf("frame_%06d.png", frameNum))
		if err := r.renderFrame(term, framePath); err != nil {
			return nil, fmt.Errorf("failed to render frame %d: %w", frameNum, err)
		}

		framePaths = append(framePaths, framePath)
		frameNum++
		nextFrameTime += frameInterval
	}

	return framePaths, nil
}

func (r *Renderer) renderFrame(term *terminal.Terminal, path string) error {
	padding := r.opts.Padding
	borderWidth := 0
	if r.opts.BorderEnabled {
		borderWidth = r.opts.BorderWidth
	}

	terminalWidth := r.opts.TermWidth*r.charWidth + padding*2
	terminalHeight := r.opts.TermHeight*r.charHeight + padding*2

	imgWidth := terminalWidth + borderWidth*2
	imgHeight := terminalHeight + borderWidth*2

	// Ensure dimensions are divisible by 2 (required for h264/h265)
	if imgWidth%2 != 0 {
		imgWidth++
	}
	if imgHeight%2 != 0 {
		imgHeight++
	}

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	// Fill outer background (transparent or colored)
	outerBg := r.opts.OuterBackground
	if r.opts.Transparent {
		outerBg = color.RGBA{0, 0, 0, 0}
	}
	for y := 0; y < imgHeight; y++ {
		for x := 0; x < imgWidth; x++ {
			img.Set(x, y, outerBg)
		}
	}

	// Draw border if enabled
	if r.opts.BorderEnabled {
		r.drawRoundedRect(img, 0, 0, imgWidth, imgHeight, r.opts.BorderRadius, r.opts.BorderColor)
	}

	// Draw terminal background with rounded corners
	termX := borderWidth
	termY := borderWidth
	if r.opts.BorderRadius > 0 {
		r.drawRoundedRect(img, termX, termY, terminalWidth, terminalHeight, r.opts.BorderRadius-borderWidth, r.opts.Theme.Background)
	} else {
		for y := termY; y < termY+terminalHeight; y++ {
			for x := termX; x < termX+terminalWidth; x++ {
				img.Set(x, y, r.opts.Theme.Background)
			}
		}
	}

	// Setup freetype context
	ctx := freetype.NewContext()
	ctx.SetDPI(72)
	ctx.SetFont(r.font)
	ctx.SetFontSize(float64(r.opts.FontSize))
	ctx.SetClip(img.Bounds())
	ctx.SetDst(img)

	// Render each cell
	contentX := borderWidth + padding
	contentY := borderWidth + padding

	for row := 0; row < r.opts.TermHeight; row++ {
		for col := 0; col < r.opts.TermWidth; col++ {
			cell := term.Screen[row][col]

			cellX := contentX + col*r.charWidth
			cellY := contentY + row*r.charHeight

			// Draw cell background
			for y := cellY; y < cellY+r.charHeight; y++ {
				for x := cellX; x < cellX+r.charWidth; x++ {
					if r.isInsideRoundedRect(x-borderWidth, y-borderWidth, terminalWidth, terminalHeight, r.opts.BorderRadius-borderWidth) {
						img.Set(x, y, cell.BG)
					}
				}
			}

			// Draw character
			if cell.Char != ' ' && cell.Char != 0 {
				ctx.SetSrc(image.NewUniform(cell.FG))
				pt := freetype.Pt(cellX, cellY+r.charHeight-4)
				_, _ = ctx.DrawString(string(cell.Char), pt)
			}
		}
	}

	// Save PNG
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	return png.Encode(file, img)
}

// drawRoundedRect draws a filled rounded rectangle
func (r *Renderer) drawRoundedRect(img *image.RGBA, x, y, w, h, radius int, c color.RGBA) {
	if radius <= 0 {
		for py := y; py < y+h; py++ {
			for px := x; px < x+w; px++ {
				img.Set(px, py, c)
			}
		}
		return
	}

	for py := y; py < y+h; py++ {
		for px := x; px < x+w; px++ {
			if r.isInsideRoundedRect(px-x, py-y, w, h, radius) {
				img.Set(px, py, c)
			}
		}
	}
}

// isInsideRoundedRect checks if a point is inside a rounded rectangle
func (r *Renderer) isInsideRoundedRect(x, y, w, h, radius int) bool {
	if radius <= 0 {
		return x >= 0 && x < w && y >= 0 && y < h
	}

	// Check corners
	// Top-left
	if x < radius && y < radius {
		dx := radius - x
		dy := radius - y
		return dx*dx+dy*dy <= radius*radius
	}
	// Top-right
	if x >= w-radius && y < radius {
		dx := x - (w - radius - 1)
		dy := radius - y
		return dx*dx+dy*dy <= radius*radius
	}
	// Bottom-left
	if x < radius && y >= h-radius {
		dx := radius - x
		dy := y - (h - radius - 1)
		return dx*dx+dy*dy <= radius*radius
	}
	// Bottom-right
	if x >= w-radius && y >= h-radius {
		dx := x - (w - radius - 1)
		dy := y - (h - radius - 1)
		return dx*dx+dy*dy <= radius*radius
	}

	return x >= 0 && x < w && y >= 0 && y < h
}
