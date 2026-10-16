package filesystem_test

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kejne/helm-values-docs/internal/adapters/out/filesystem"
)

func TestStorePersistsAndReadsGeneratedDocument(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "README.md")
	store := filesystem.Store{}

	for _, document := range [][]byte{[]byte("first version\n"), []byte("updated version\n")} {
		if err := store.WriteAtomic(path, document); err != nil {
			t.Fatalf("WriteAtomic() error = %v", err)
		}
		got, err := store.Read(path)
		if err != nil {
			t.Fatalf("Read() error = %v", err)
		}
		if !bytes.Equal(got, document) {
			t.Fatalf("Read() = %q, want %q", got, document)
		}
	}
}

func TestStoreReportsMissingDocument(t *testing.T) {
	_, err := (filesystem.Store{}).Read(filepath.Join(t.TempDir(), "missing.md"))
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Read() error = %v, want an error matching fs.ErrNotExist", err)
	}
}

func TestStoreDoesNotReplaceDirectoryWithDocument(t *testing.T) {
	dir := t.TempDir()
	output := filepath.Join(dir, "README.md")
	if err := os.Mkdir(output, 0o700); err != nil {
		t.Fatal(err)
	}

	if err := (filesystem.Store{}).WriteAtomic(output, []byte("generated\n")); err == nil {
		t.Fatal("WriteAtomic() replaced a directory")
	}
	info, err := os.Stat(output)
	if err != nil || !info.IsDir() {
		t.Fatalf("destination after failed write: info=%v error=%v, want existing directory", info, err)
	}
}

func TestTemplateComposerWrapsDocumentAndReportsInvalidTemplates(t *testing.T) {
	dir := t.TempDir()
	templatePath := filepath.Join(dir, "README.tmpl.md")
	composer := filesystem.TemplateAdapter{}

	if err := os.WriteFile(templatePath, []byte("# Values\n{{.API}}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := composer.Compose([]byte("| setting | true |\n"), templatePath)
	if err != nil {
		t.Fatalf("Compose() error = %v", err)
	}
	if want := "# Values\n| setting | true |\n\n"; string(got) != want {
		t.Fatalf("Compose() = %q, want %q", got, want)
	}

	if err := os.WriteFile(templatePath, []byte("{{.Missing}}"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err = composer.Compose([]byte("document"), templatePath)
	if err == nil || !strings.Contains(err.Error(), templatePath) {
		t.Fatalf("Compose() error = %v, want an execution error naming the template", err)
	}
}
