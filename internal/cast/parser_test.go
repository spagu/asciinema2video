package cast

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseValidCast(t *testing.T) {
	content := `{"version": 2, "width": 80, "height": 24, "timestamp": 1234567890}
[0.5, "o", "Hello"]
[1.0, "o", " World"]
`
	tmpFile := createTempFile(t, content)

	recording, err := Parse(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if recording.Header.Version != 2 {
		t.Errorf("expected version 2, got %d", recording.Header.Version)
	}
	if recording.Header.Width != 80 {
		t.Errorf("expected width 80, got %d", recording.Header.Width)
	}
	if recording.Header.Height != 24 {
		t.Errorf("expected height 24, got %d", recording.Header.Height)
	}
	if len(recording.Events) != 2 {
		t.Errorf("expected 2 events, got %d", len(recording.Events))
	}
}

func TestParseInvalidVersion(t *testing.T) {
	content := `{"version": 1, "width": 80, "height": 24}
[0.5, "o", "Hello"]
`
	tmpFile := createTempFile(t, content)

	_, err := Parse(tmpFile)
	if err == nil {
		t.Fatal("expected error for unsupported version")
	}
}

func TestParseEmptyFile(t *testing.T) {
	tmpFile := createTempFile(t, "")

	_, err := Parse(tmpFile)
	if err == nil {
		t.Fatal("expected error for empty file")
	}
}

func TestParseNonExistentFile(t *testing.T) {
	_, err := Parse("/nonexistent/file.cast")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestParseEventTypes(t *testing.T) {
	content := `{"version": 2, "width": 80, "height": 24}
[0.5, "o", "output"]
[1.0, "i", "input"]
[1.5, "o", "more output"]
`
	tmpFile := createTempFile(t, content)

	recording, err := Parse(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Only output events should be captured
	if len(recording.Events) != 2 {
		t.Errorf("expected 2 output events, got %d", len(recording.Events))
	}
}

func createTempFile(t *testing.T, content string) string {
	t.Helper()
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.cast")
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	return tmpFile
}
