package cast

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// Header represents asciinema v2 header
type Header struct {
	Version   int               `json:"version"`
	Width     int               `json:"width"`
	Height    int               `json:"height"`
	Timestamp int64             `json:"timestamp"`
	Env       map[string]string `json:"env"`
}

// Event represents a single output event
type Event struct {
	Time float64
	Type string
	Data string
}

// Recording represents a complete asciinema recording
type Recording struct {
	Header Header
	Events []Event
}

// Parse reads and parses a .cast file
func Parse(filename string) (*Recording, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	recording := &Recording{}

	// Parse header (first line)
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty file")
	}

	if err := json.Unmarshal(scanner.Bytes(), &recording.Header); err != nil {
		return nil, fmt.Errorf("failed to parse header: %w", err)
	}

	if recording.Header.Version != 2 {
		return nil, fmt.Errorf("unsupported version: %d (only v2 supported)", recording.Header.Version)
	}

	// Parse events
	for scanner.Scan() {
		var raw []json.RawMessage
		if err := json.Unmarshal(scanner.Bytes(), &raw); err != nil {
			continue // Skip malformed lines
		}

		if len(raw) < 3 {
			continue
		}

		var event Event
		if err := json.Unmarshal(raw[0], &event.Time); err != nil {
			continue
		}
		if err := json.Unmarshal(raw[1], &event.Type); err != nil {
			continue
		}
		if err := json.Unmarshal(raw[2], &event.Data); err != nil {
			continue
		}

		// Only process output events
		if event.Type == "o" {
			recording.Events = append(recording.Events, event)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return recording, nil
}
