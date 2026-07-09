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

// Store defines the interface for reading and writing compilation database records.
type Store interface {
	Add(records []Record) error
	Records() ([]Record, error)
}

// FileStore handles reading and writing compilation database records in a
// thread-safe manner using file locking.
type FileStore struct {
	file string
}

// Add merges new compilation records into the database file.
//
// To prevent database corruption and guarantee reliability during concurrent compiler
// execution, updates are serialized using an advisory lock, and the write is performed
// atomically via a temporary swap file to ensure the database is never left in a partially-written state.
func (s *FileStore) Add(records []Record) error {
	//nolint:gosec // compilation database directories must be user-accessible (0755)
	if err := os.MkdirAll(filepath.Dir(s.file), 0o755); err != nil {
		return fmt.Errorf("creating database directory: %w", err)
	}

	mu := flock.New(s.file + ".lock")
	if err := mu.Lock(); err != nil {
		return fmt.Errorf("acquiring lock: %w", err)
	}
	//nolint:errcheck // lock release is best-effort on defer
	defer func() { _ = mu.Unlock() }()

	stored := []Record{}
	if data, err := os.ReadFile(s.file); err == nil {
		//nolint:errcheck // corrupt database is overwritten
		_ = json.Unmarshal(data, &stored)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("reading database file: %w", err)
	}

	stored = mergeRecords(stored, records)

	swpFile := s.file + ".swp"
	//nolint:gosec // database file path is designated by user configuration
	f, err := os.Create(swpFile)
	if err != nil {
		return fmt.Errorf("creating swap file: %w", err)
	}

	if err := json.NewEncoder(f).Encode(stored); err != nil {
		//nolint:errcheck // cleanup is best-effort on error path
		_ = os.Remove(swpFile)
		//nolint:errcheck // file close is best-effort on error path
		_ = f.Close()
		return fmt.Errorf("encoding database JSON: %w", err)
	}

	if err := f.Close(); err != nil {
		//nolint:errcheck // cleanup is best-effort on error path
		_ = os.Remove(swpFile)
		return fmt.Errorf("closing swap file: %w", err)
	}

	if err := os.Rename(swpFile, s.file); err != nil {
		//nolint:errcheck // cleanup is best-effort on error path
		_ = os.Remove(swpFile)
		return fmt.Errorf("committing database swap: %w", err)
	}
	return nil
}

// Records returns all compilation records stored in the database file.
func (s *FileStore) Records() ([]Record, error) {
	//nolint:gosec // compilation database directories must be user-accessible (0755)
	if err := os.MkdirAll(filepath.Dir(s.file), 0o755); err != nil {
		return nil, fmt.Errorf("creating database directory: %w", err)
	}

	mu := flock.New(s.file + ".lock")
	if err := mu.RLock(); err != nil {
		return nil, fmt.Errorf("acquiring read lock: %w", err)
	}
	//nolint:errcheck // lock release is best-effort on defer
	defer func() { _ = mu.Unlock() }()

	records := []Record{}
	if data, err := os.ReadFile(s.file); err == nil {
		if err := json.Unmarshal(data, &records); err != nil {
			return nil, fmt.Errorf("unmarshalling database JSON: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("reading database file: %w", err)
	}
	return records, nil
}

var _ Store = (*FileStore)(nil)

// NewFileStore creates a new FileStore instance managing the specified database file.
func NewFileStore(file string) *FileStore {
	return &FileStore{file: file}
}
