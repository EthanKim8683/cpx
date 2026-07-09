package cdb

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
)

// Record represents a single compilation entry in the compilation database.
type Record struct {
	File    string
	Dir     string
	Shim    string
	Command Command
}

func mergeRecords(a, b []Record) []Record {
	m := make(map[string]Record, len(a)+len(b))
	for _, record := range a {
		m[record.File] = record
	}
	for _, record := range b {
		m[record.File] = record
	}
	merged := make([]Record, 0, len(m))
	for _, record := range m {
		merged = append(merged, record)
	}
	return merged
}

// Recorder defines the interface for recording compilation records to a database.
type Recorder interface {
	Record(records []Record) error
}

// FileRecorder handles reading and writing compilation database records in a
// thread-safe manner using file locking.
type FileRecorder struct {
	file string
}

// Record merges new compilation records into the database file.
//
// To prevent database corruption and guarantee reliability during concurrent compiler
// execution, updates are serialized using an advisory lock, and the write is performed
// atomically via a temporary swap file to ensure the database is never left in a partially-written state.
func (r *FileRecorder) Record(records []Record) error {
	if err := os.MkdirAll(filepath.Dir(r.file), 0755); err != nil { //nolint:gosec // compilation database directories must be user-accessible (0755)
		return fmt.Errorf("creating database directory: %w", err)
	}

	mu := flock.New(r.file + ".lock")
	if err := mu.Lock(); err != nil {
		return fmt.Errorf("acquiring lock: %w", err)
	}
	defer func() { _ = mu.Unlock() }() //nolint:errcheck // lock release is best-effort on defer

	recorded := []Record{}
	if data, err := os.ReadFile(r.file); err == nil {
		_ = json.Unmarshal(data, &recorded) //nolint:errcheck // corrupt database is overwritten
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("reading database file: %w", err)
	}

	recorded = mergeRecords(recorded, records)

	swpFile := r.file + ".swp"
	f, err := os.Create(swpFile) //nolint:gosec // database file path is designated by user configuration
	if err != nil {
		return fmt.Errorf("creating swap file: %w", err)
	}

	if err := json.NewEncoder(f).Encode(recorded); err != nil {
		_ = os.Remove(swpFile) //nolint:errcheck // cleanup is best-effort on error path
		_ = f.Close()          //nolint:errcheck // file close is best-effort on error path
		return fmt.Errorf("encoding database JSON: %w", err)
	}

	if err := f.Close(); err != nil {
		_ = os.Remove(swpFile) //nolint:errcheck // cleanup is best-effort on error path
		return fmt.Errorf("closing swap file: %w", err)
	}

	if err := os.Rename(swpFile, r.file); err != nil {
		_ = os.Remove(swpFile) //nolint:errcheck // cleanup is best-effort on error path
		return fmt.Errorf("committing database swap: %w", err)
	}
	return nil
}

// NewFileRecorder creates a new FileRecorder instance managing the specified database file.
func NewFileRecorder(file string) *FileRecorder {
	return &FileRecorder{file: file}
}
