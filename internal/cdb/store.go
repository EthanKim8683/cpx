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

// Store handles reading and writing compilation database records in a
// thread-safe manner using file locking.
type Store struct {
	file string
}

// Add merges new compilation records into the database file, performing
// atomic updates using a swap file and file locks.
func (s *Store) Add(records []Record) error {
	if err := os.MkdirAll(filepath.Dir(s.file), 0755); err != nil {
		return fmt.Errorf("creating database directory: %w", err)
	}

	mu := flock.New(s.file + ".lock")
	if err := mu.Lock(); err != nil {
		return fmt.Errorf("acquiring lock: %w", err)
	}
	defer mu.Unlock()

	stored := []Record{}
	if data, err := os.ReadFile(s.file); err == nil {
		if err := json.Unmarshal(data, &stored); err != nil {
			return fmt.Errorf("parsing database JSON: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("reading database file: %w", err)
	}

	stored = mergeRecords(stored, records)

	swpFile := s.file + ".swp"
	f, err := os.Create(swpFile)
	if err != nil {
		return fmt.Errorf("creating swap file: %w", err)
	}

	if err := json.NewEncoder(f).Encode(stored); err != nil {
		_ = os.Remove(swpFile)
		f.Close()
		return fmt.Errorf("encoding database JSON: %w", err)
	}

	if err := f.Close(); err != nil {
		_ = os.Remove(swpFile)
		return fmt.Errorf("closing swap file: %w", err)
	}

	if err := os.Rename(swpFile, s.file); err != nil {
		_ = os.Remove(swpFile)
		return fmt.Errorf("committing database swap: %w", err)
	}
	return nil
}

// NewStore creates a new Store instance managing the specified database file.
func NewStore(file string) *Store {
	return &Store{file: file}
}
