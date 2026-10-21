// Package filesystem adapts local files to document storage and template ports.
package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

// Store reads and atomically writes local documents.
type Store struct{}

func (Store) Read(path string) ([]byte, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- path is explicitly provided by the caller.
	if err != nil {
		return nil, fmt.Errorf("read output %q: %w", path, err)
	}
	return data, nil
}

func (Store) WriteAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	file, err := os.CreateTemp(dir, ".hvd-*")
	if err != nil {
		return fmt.Errorf("create temporary output for %q: %w", path, err)
	}
	temporary := file.Name()
	defer func() { _ = os.Remove(temporary) }()
	if _, err := file.Write(data); err != nil {
		_ = file.Close()
		return fmt.Errorf("write output %q: %w", path, err)
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return fmt.Errorf("sync output %q: %w", path, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close output %q: %w", path, err)
	}
	if err := os.Rename(temporary, path); err != nil {
		return fmt.Errorf("replace output %q: %w", path, err)
	}
	return nil
}
