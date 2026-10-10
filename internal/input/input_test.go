package input

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadValuesMergesInOrder(t *testing.T) {
	dir := t.TempDir()
	base := filepath.Join(dir, "base.yaml")
	override := filepath.Join(dir, "override.yaml")
	if err := os.WriteFile(base, []byte("image:\n  repository: app\n  tag: old\nreplicas: 1\nitems: [a, b]\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(override, []byte("image:\n  tag: new\nreplicas: 3\nitems: [c]\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := LoadValues([]string{base, override})
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{
		"image":    map[string]any{"repository": "app", "tag": "new"},
		"replicas": 3,
		"items":    []any{"c"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("merged values = %#v, want %#v", got, want)
	}
}

func TestLoadValuesRejectsNonMappingRoot(t *testing.T) {
	path := filepath.Join(t.TempDir(), "values.yaml")
	if err := os.WriteFile(path, []byte("- value\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := LoadValues([]string{path})
	if err == nil {
		t.Fatal("LoadValues accepted a sequence root")
	}
}
