package video

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// CreateOptions contains video creation options
type CreateOptions struct {
	Transparent bool   // Enable alpha channel for webm/mov
	Codec       string // Video codec for MP4: h264, h265 (default: h264)
}

// Create generates a video file from PNG frames using ffmpeg
func Create(framesDir, outputPath string, fps int) error {
	return CreateWithOptions(framesDir, outputPath, fps, nil)
}

// CreateWithOptions generates a video file with additional options
func CreateWithOptions(framesDir, outputPath string, fps int, opts *CreateOptions) error {
	if opts == nil {
		opts = &CreateOptions{}
	}

	ext := strings.ToLower(filepath.Ext(outputPath))

	// Validate format before checking for ffmpeg
	switch ext {
	case ".gif", ".webp", ".mp4", ".webm", ".mov":
		// Valid format
	default:
		return fmt.Errorf("unsupported output format: %s (supported: mp4, gif, webp, webm, mov)", ext)
	}

	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH: %w", err)
	}

	inputPattern := filepath.Join(framesDir, "frame_%06d.png")

	var args []string

	switch ext {
	case ".gif":
		palettePath := filepath.Join(framesDir, "palette.png")
		paletteCmd := exec.Command("ffmpeg", "-y",
			"-framerate", fmt.Sprintf("%d", fps),
			"-i", inputPattern,
			"-vf", "palettegen=stats_mode=diff",
			palettePath,
		)
		if output, err := paletteCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to generate palette: %w\n%s", err, output)
		}

		args = []string{
			"-y",
			"-framerate", fmt.Sprintf("%d", fps),
			"-i", inputPattern,
			"-i", palettePath,
			"-lavfi", "paletteuse=dither=bayer:bayer_scale=5",
			outputPath,
		}

	case ".webp":
		args = []string{
			"-y",
			"-framerate", fmt.Sprintf("%d", fps),
			"-i", inputPattern,
			"-vcodec", "libwebp",
			"-lossless", "1",
			"-loop", "0",
			"-preset", "default",
			"-an", "-vsync", "0",
			outputPath,
		}

	case ".mp4":
		codecLib := "libx264"
		crf := "23"
		if opts.Codec == "h265" {
			codecLib = "libx265"
			crf = "28" // H.265 uses different CRF scale
		}
		args = []string{
			"-y",
			"-framerate", fmt.Sprintf("%d", fps),
			"-i", inputPattern,
			"-c:v", codecLib,
			"-pix_fmt", "yuv420p",
			"-preset", "medium",
			"-crf", crf,
			outputPath,
		}

	case ".webm":
		if opts.Transparent {
			// VP9 with alpha channel
			args = []string{
				"-y",
				"-framerate", fmt.Sprintf("%d", fps),
				"-i", inputPattern,
				"-c:v", "libvpx-vp9",
				"-pix_fmt", "yuva420p",
				"-auto-alt-ref", "0",
				"-crf", "30",
				"-b:v", "0",
				outputPath,
			}
		} else {
			// VP9 without alpha
			args = []string{
				"-y",
				"-framerate", fmt.Sprintf("%d", fps),
				"-i", inputPattern,
				"-c:v", "libvpx-vp9",
				"-pix_fmt", "yuv420p",
				"-crf", "30",
				"-b:v", "0",
				outputPath,
			}
		}

	case ".mov":
		if opts.Transparent {
			// ProRes 4444 with alpha channel
			args = []string{
				"-y",
				"-framerate", fmt.Sprintf("%d", fps),
				"-i", inputPattern,
				"-c:v", "prores_ks",
				"-profile:v", "4444",
				"-pix_fmt", "yuva444p10le",
				"-alpha_bits", "16",
				outputPath,
			}
		} else {
			// ProRes 422 HQ without alpha
			args = []string{
				"-y",
				"-framerate", fmt.Sprintf("%d", fps),
				"-i", inputPattern,
				"-c:v", "prores_ks",
				"-profile:v", "3",
				"-pix_fmt", "yuv422p10le",
				outputPath,
			}
		}
	}

	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg failed: %w\n%s", err, output)
	}

	return nil
}
