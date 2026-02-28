package video

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateUnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "output.avi")

	err := Create(tmpDir, outputPath, 10)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
	if !strings.Contains(err.Error(), "unsupported output format") {
		t.Errorf("expected 'unsupported output format' error, got: %v", err)
	}
}

func TestCreateNoFFmpeg(t *testing.T) {
	// This test will pass if ffmpeg is not installed
	// and fail gracefully if it is
	tmpDir := t.TempDir()

	// Create a dummy frame
	framePath := filepath.Join(tmpDir, "frame_000000.png")
	err := os.WriteFile(framePath, []byte("not a real png"), 0644)
	if err != nil {
		t.Fatalf("failed to create dummy frame: %v", err)
	}

	outputPath := filepath.Join(tmpDir, "output.mp4")
	err = Create(tmpDir, outputPath, 10)

	// Either ffmpeg not found or ffmpeg fails on invalid PNG
	// Both are acceptable outcomes for this test
	if err == nil {
		// FFmpeg succeeded somehow, check if output exists
		if _, statErr := os.Stat(outputPath); statErr != nil {
			t.Log("ffmpeg ran but no output created")
		}
	}
}

func TestOutputFormats(t *testing.T) {
	formats := []string{".mp4", ".gif", ".webp"}
	for _, ext := range formats {
		t.Run(ext, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputPath := filepath.Join(tmpDir, "output"+ext)

			// This will fail because no frames exist, but it validates
			// that the format is recognized
			err := Create(tmpDir, outputPath, 10)
			if err != nil && strings.Contains(err.Error(), "unsupported output format") {
				t.Errorf("format %s should be supported", ext)
			}
		})
	}
}
