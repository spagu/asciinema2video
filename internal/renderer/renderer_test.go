package renderer

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/spagu/asciinema2video/internal/cast"
)

func TestNewRenderer(t *testing.T) {
	r, err := New(80, 24, 14)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if r.opts.TermWidth != 80 {
		t.Errorf("expected termWidth 80, got %d", r.opts.TermWidth)
	}
	if r.opts.TermHeight != 24 {
		t.Errorf("expected termHeight 24, got %d", r.opts.TermHeight)
	}
	if r.opts.FontSize != 14 {
		t.Errorf("expected fontSize 14, got %d", r.opts.FontSize)
	}
	if r.charWidth <= 0 {
		t.Error("expected positive charWidth")
	}
	if r.charHeight <= 0 {
		t.Error("expected positive charHeight")
	}
}

func TestRenderFrames(t *testing.T) {
	r, err := New(40, 10, 12)
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	recording := &cast.Recording{
		Header: cast.Header{
			Version: 2,
			Width:   40,
			Height:  10,
		},
		Events: []cast.Event{
			{Time: 0.0, Type: "o", Data: "Hello"},
			{Time: 0.5, Type: "o", Data: " World"},
		},
	}

	tmpDir := t.TempDir()
	fps := 10

	framePaths, err := r.RenderFrames(recording, tmpDir, fps)
	if err != nil {
		t.Fatalf("failed to render frames: %v", err)
	}

	if len(framePaths) == 0 {
		t.Error("expected at least one frame")
	}

	// Check that frames exist
	for _, path := range framePaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("frame file does not exist: %s", path)
		}
	}
}

func TestRenderFramesPNG(t *testing.T) {
	r, err := New(20, 5, 10)
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	recording := &cast.Recording{
		Header: cast.Header{
			Version: 2,
			Width:   20,
			Height:  5,
		},
		Events: []cast.Event{
			{Time: 0.0, Type: "o", Data: "Test"},
		},
	}

	tmpDir := t.TempDir()
	framePaths, err := r.RenderFrames(recording, tmpDir, 1)
	if err != nil {
		t.Fatalf("failed to render: %v", err)
	}

	if len(framePaths) < 1 {
		t.Fatal("expected at least one frame")
	}

	// Verify PNG file
	framePath := framePaths[0]
	data, err := os.ReadFile(framePath)
	if err != nil {
		t.Fatalf("failed to read frame: %v", err)
	}

	// PNG magic bytes
	if len(data) < 8 {
		t.Fatal("frame file too small")
	}
	pngMagic := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	for i, b := range pngMagic {
		if data[i] != b {
			t.Fatal("frame is not a valid PNG file")
		}
	}
}

func TestRenderEmptyRecording(t *testing.T) {
	r, err := New(80, 24, 14)
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	recording := &cast.Recording{
		Header: cast.Header{
			Version: 2,
			Width:   80,
			Height:  24,
		},
		Events: []cast.Event{},
	}

	tmpDir := t.TempDir()
	framePaths, err := r.RenderFrames(recording, tmpDir, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should still generate at least the initial frame
	if len(framePaths) != 1 {
		t.Errorf("expected 1 frame for empty recording, got %d", len(framePaths))
	}
}

func TestFrameNaming(t *testing.T) {
	r, err := New(10, 5, 10)
	if err != nil {
		t.Fatalf("failed to create renderer: %v", err)
	}

	recording := &cast.Recording{
		Header: cast.Header{Version: 2, Width: 10, Height: 5},
		Events: []cast.Event{
			{Time: 0.0, Type: "o", Data: "A"},
			{Time: 0.2, Type: "o", Data: "B"},
		},
	}

	tmpDir := t.TempDir()
	framePaths, err := r.RenderFrames(recording, tmpDir, 5)
	if err != nil {
		t.Fatalf("failed to render: %v", err)
	}

	for i, path := range framePaths {
		expectedName := filepath.Join(tmpDir, "frame_"+padNumber(i)+".png")
		if path != expectedName {
			t.Errorf("expected frame name %s, got %s", expectedName, path)
		}
	}
}

func padNumber(n int) string {
	return fmt.Sprintf("%06d", n)
}
