package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRenderCommand(t *testing.T) {
	dir := t.TempDir()
	schema := filepath.Join(dir, "values.schema.json")
	values := filepath.Join(dir, "values.yaml")
	output := filepath.Join(dir, "README.md")
	if err := os.WriteFile(schema, []byte(`{"type":"object","properties":{"enabled":{"type":"boolean"}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(values, []byte("enabled: true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var stderr bytes.Buffer
	if code := Run([]string{"render", "--schema", schema, "--values", values, "--output", output}, &bytes.Buffer{}, &stderr); code != 0 {
		t.Fatalf("Run returned %d: %s", code, stderr.String())
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `| ["enabled"] | boolean | no | true |`) {
		t.Fatalf("output does not contain rendered value:\n%s", data)
	}
}

func TestRenderRequiresInputs(t *testing.T) {
	var stderr bytes.Buffer
	if code := Run([]string{"render"}, &bytes.Buffer{}, &stderr); code != 2 {
		t.Fatalf("Run returned %d, want 2", code)
	}
}

func TestCheckCommand(t *testing.T) {
	dir := t.TempDir()
	schema := filepath.Join(dir, "values.schema.json")
	values := filepath.Join(dir, "values.yaml")
	output := filepath.Join(dir, "README.md")
	if err := os.WriteFile(schema, []byte(`{"type":"object","properties":{"enabled":{"type":"boolean"}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(values, []byte("enabled: true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var stderr bytes.Buffer
	args := []string{"render", "--schema", schema, "--values", values, "--output", output}
	if code := Run(args, &bytes.Buffer{}, &stderr); code != 0 {
		t.Fatalf("render returned %d: %s", code, stderr.String())
	}
	before, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	stderr.Reset()
	checkArgs := []string{"check", "--schema", schema, "--values", values, "--output", output}
	if code := Run(checkArgs, &bytes.Buffer{}, &stderr); code != 0 {
		t.Fatalf("current check returned %d: %s", code, stderr.String())
	}
	if after, err := os.ReadFile(output); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(after, before) {
		t.Fatal("check changed output")
	}
	if err := os.WriteFile(output, []byte("stale\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	stderr.Reset()
	if code := Run(checkArgs, &bytes.Buffer{}, &stderr); code != 1 {
		t.Fatalf("stale check returned %d, want 1: %s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "stale") {
		t.Fatalf("stale diagnostic = %q", stderr.String())
	}
	if err := os.Remove(output); err != nil {
		t.Fatal(err)
	}
	stderr.Reset()
	if code := Run(checkArgs, &bytes.Buffer{}, &stderr); code != 1 {
		t.Fatalf("missing check returned %d, want 1: %s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "missing") {
		t.Fatalf("missing diagnostic = %q", stderr.String())
	}
}

func TestCheckDoesNotMutateOutputMetadata(t *testing.T) {
	dir := t.TempDir()
	schema := filepath.Join(dir, "values.schema.json")
	values := filepath.Join(dir, "values.yaml")
	output := filepath.Join(dir, "README.md")
	if err := os.WriteFile(schema, []byte(`{"type":"object","properties":{"enabled":{"type":"boolean"}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(values, []byte("enabled: true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var stderr bytes.Buffer
	args := []string{"render", "--schema", schema, "--values", values, "--output", output}
	if code := Run(args, &bytes.Buffer{}, &stderr); code != 0 {
		t.Fatalf("render returned %d: %s", code, stderr.String())
	}
	const mode = 0o640
	if err := os.Chmod(output, mode); err != nil {
		t.Fatal(err)
	}
	stamp := time.Unix(123, 456)
	if err := os.Chtimes(output, stamp, stamp); err != nil {
		t.Fatal(err)
	}
	before, err := os.Stat(output)
	if err != nil {
		t.Fatal(err)
	}
	stderr.Reset()
	checkArgs := []string{"check", "--schema", schema, "--values", values, "--output", output}
	if code := Run(checkArgs, &bytes.Buffer{}, &stderr); code != 0 {
		t.Fatalf("check returned %d: %s", code, stderr.String())
	}
	after, err := os.Stat(output)
	if err != nil {
		t.Fatal(err)
	}
	if before.Mode() != after.Mode() || !before.ModTime().Equal(after.ModTime()) {
		t.Fatalf("check changed output metadata: before mode=%#o time=%s, after mode=%#o time=%s", before.Mode(), before.ModTime(), after.Mode(), after.ModTime())
	}
}

func TestCheckMissingSchemaReturnsInputError(t *testing.T) {
	dir := t.TempDir()
	schema := filepath.Join(dir, "missing-schema.json")
	values := filepath.Join(dir, "values.yaml")
	output := filepath.Join(dir, "README.md")
	if err := os.WriteFile(values, []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(output, []byte("existing\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var stderr bytes.Buffer
	code := Run([]string{"check", "--schema", schema, "--values", values, "--output", output}, &bytes.Buffer{}, &stderr)
	if code != 2 {
		t.Fatalf("check returned %d, want 2: %s", code, stderr.String())
	}
	if strings.Contains(stderr.String(), "is missing") {
		t.Fatalf("missing schema was reported as missing output: %s", stderr.String())
	}
}

func TestCheckInvalidInputReturnsUsageError(t *testing.T) {
	dir := t.TempDir()
	schema := filepath.Join(dir, "values.schema.json")
	values := filepath.Join(dir, "values.yaml")
	output := filepath.Join(dir, "README.md")
	if err := os.WriteFile(schema, []byte("not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(values, []byte("enabled: true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var stderr bytes.Buffer
	code := Run([]string{"check", "--schema", schema, "--values", values, "--output", output}, &bytes.Buffer{}, &stderr)
	if code != 2 {
		t.Fatalf("check returned %d, want 2: %s", code, stderr.String())
	}
}

func TestTemplateCompositionGoldenAndCheck(t *testing.T) {
	dir := t.TempDir()
	schema := filepath.Join(dir, "values.schema.json")
	values := filepath.Join(dir, "values.yaml")
	templatePath := filepath.Join(dir, "README.tmpl.md")
	output := filepath.Join(dir, "README.md")
	if err := os.WriteFile(schema, []byte(`{"type":"object","properties":{"enabled":{"type":"boolean"}}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(values, []byte("enabled: true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(templatePath, []byte("# Chart values\n\n{{.API}}\n## End\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var stderr bytes.Buffer
	args := []string{"render", "--schema", schema, "--values", values, "--output", output, "--template", templatePath}
	if code := Run(args, &bytes.Buffer{}, &stderr); code != 0 {
		t.Fatalf("render returned %d: %s", code, stderr.String())
	}
	got, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	assertApproval(t, "composed.md", got)
	stderr.Reset()
	checkArgs := []string{"check", "--schema", schema, "--values", values, "--output", output, "--template", templatePath}
	if code := Run(checkArgs, &bytes.Buffer{}, &stderr); code != 0 {
		t.Fatalf("check returned %d: %s", code, stderr.String())
	}
}

func assertApproval(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", name)
	if os.Getenv("UPDATE_APPROVALS") == "1" {
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatalf("read approval directory: %v", err)
	}
	for _, entry := range entries {
		if entry.Name() != name {
			t.Fatalf("unexpected approval artifact %s", filepath.Join(filepath.Dir(path), entry.Name()))
		}
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read approval %s: %v", path, err)
	}
	if string(got) != string(want) {
		t.Fatalf("approval %s differs:\\n got:\\n%s\\nwant:\\n%s", path, got, want)
	}
}

func TestTemplateErrorsIncludePath(t *testing.T) {
	dir := t.TempDir()
	schema := filepath.Join(dir, "values.schema.json")
	values := filepath.Join(dir, "values.yaml")
	templatePath := filepath.Join(dir, "README.tmpl.md")
	output := filepath.Join(dir, "README.md")
	if err := os.WriteFile(schema, []byte(`{"type":"object"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(values, []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(templatePath, []byte("{{.Missing}}"), 0o600); err != nil {
		t.Fatal(err)
	}
	var stderr bytes.Buffer
	code := Run([]string{"render", "--schema", schema, "--values", values, "--output", output, "--template", templatePath}, &bytes.Buffer{}, &stderr)
	if code != 2 {
		t.Fatalf("render returned %d, want 2: %s", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), templatePath) {
		t.Fatalf("template diagnostic %q does not name path %q", stderr.String(), templatePath)
	}
}
