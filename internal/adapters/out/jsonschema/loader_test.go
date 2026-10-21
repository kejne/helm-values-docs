package jsonschema_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/smedijan/helm-values-docs/internal/adapters/out/jsonschema"
)

func TestLoaderReadsAndDecodesSchema(t *testing.T) {
	path := filepath.Join(t.TempDir(), "values.schema.json")
	if err := os.WriteFile(path, []byte(`{"type":"object","properties":{"enabled":{"type":"boolean"}}}`), 0o600); err != nil {
		t.Fatal(err)
	}

	model, err := (jsonschema.Loader{}).Load(path)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if got := model.Properties["enabled"].Types; len(got) != 1 || got[0] != "boolean" {
		t.Fatalf("loaded property types = %v, want [boolean]", got)
	}
}

func TestLoaderPreservesFileAndSchemaErrors(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.json")
	if _, err := (jsonschema.Loader{}).Load(missing); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Load() missing-file error = %v, want fs.ErrNotExist", err)
	}

	path := filepath.Join(t.TempDir(), "invalid.json")
	if err := os.WriteFile(path, []byte("not JSON"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := (jsonschema.Loader{}).Load(path); err == nil || !strings.Contains(err.Error(), path) {
		t.Fatalf("Load() malformed-schema error = %v, want an error naming %q", err, path)
	}
}
