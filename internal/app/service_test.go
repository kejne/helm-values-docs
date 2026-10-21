package app_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/smedijan/helm-values-docs/internal/adapters/out/filesystem"
	"github.com/smedijan/helm-values-docs/internal/adapters/out/jsonschema"
	"github.com/smedijan/helm-values-docs/internal/adapters/out/markdown"
	"github.com/smedijan/helm-values-docs/internal/adapters/out/yamlvalues"
	"github.com/smedijan/helm-values-docs/internal/app"
	inport "github.com/smedijan/helm-values-docs/internal/ports/in"
)

func TestDocumentationLifecycle(t *testing.T) {
	dir := t.TempDir()
	schemaPath := filepath.Join(dir, "values.schema.json")
	valuesPath := filepath.Join(dir, "values.yaml")
	outputPath := filepath.Join(dir, "README.md")
	writeFile(t, schemaPath, `{"type":"object","required":["enabled"],"properties":{"enabled":{"type":"boolean","description":"Enable the feature"}}}`)
	writeFile(t, valuesPath, "enabled: true\n")

	service := app.New(jsonschema.Loader{}, yamlvalues.Adapter{}, markdown.Adapter{}, filesystem.TemplateAdapter{}, filesystem.Store{})
	request := inport.Request{
		SchemaPath:  schemaPath,
		ValuesPaths: []string{valuesPath},
		OutputPath:  outputPath,
	}
	if err := service.Render(context.Background(), request); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	document, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(document), `| Description | Enable the feature |`) {
		t.Fatalf("rendered documentation does not describe the configured value:\n%s", document)
	}

	current, err := service.Check(context.Background(), request)
	if err != nil || !current {
		t.Fatalf("Check() = (%v, %v), want (true, nil)", current, err)
	}

	writeFile(t, outputPath, "stale documentation\n")
	current, err = service.Check(context.Background(), request)
	if err != nil || current {
		t.Fatalf("Check() for stale output = (%v, %v), want (false, nil)", current, err)
	}

	if err := os.Remove(outputPath); err != nil {
		t.Fatal(err)
	}
	_, err = service.Check(context.Background(), request)
	if !errors.Is(err, inport.ErrOutputMissing) {
		t.Fatalf("Check() for missing output error = %v, want ErrOutputMissing", err)
	}
}

func TestRenderAcceptsAnEmptyTemplate(t *testing.T) {
	dir := t.TempDir()
	schemaPath := filepath.Join(dir, "values.schema.json")
	valuesPath := filepath.Join(dir, "values.yaml")
	templatePath := filepath.Join(dir, "empty.tmpl")
	outputPath := filepath.Join(dir, "README.md")
	writeFile(t, schemaPath, `{"type":"object"}`)
	writeFile(t, valuesPath, "{}\n")
	writeFile(t, templatePath, "")

	service := app.New(jsonschema.Loader{}, yamlvalues.Adapter{}, markdown.Adapter{}, filesystem.TemplateAdapter{}, filesystem.Store{})
	err := service.Render(context.Background(), inport.Request{
		SchemaPath:   schemaPath,
		ValuesPaths:  []string{valuesPath},
		OutputPath:   outputPath,
		TemplatePath: templatePath,
	})
	if err != nil {
		t.Fatalf("Render() with an empty template error = %v", err)
	}
	if got, err := os.ReadFile(outputPath); err != nil || len(got) != 0 {
		t.Fatalf("empty template output = %q, error = %v; want an empty document", got, err)
	}
}

func TestRenderDoesNotReplaceOutputWhenTemplateFails(t *testing.T) {
	dir := t.TempDir()
	schemaPath := filepath.Join(dir, "values.schema.json")
	valuesPath := filepath.Join(dir, "values.yaml")
	templatePath := filepath.Join(dir, "invalid.tmpl")
	outputPath := filepath.Join(dir, "README.md")
	writeFile(t, schemaPath, `{"type":"object"}`)
	writeFile(t, valuesPath, "{}\n")
	writeFile(t, outputPath, "previous documentation\n")
	writeFile(t, templatePath, "{{.Missing}}")

	service := app.New(jsonschema.Loader{}, yamlvalues.Adapter{}, markdown.Adapter{}, filesystem.TemplateAdapter{}, filesystem.Store{})
	err := service.Render(context.Background(), inport.Request{
		SchemaPath:   schemaPath,
		ValuesPaths:  []string{valuesPath},
		OutputPath:   outputPath,
		TemplatePath: templatePath,
	})
	if err == nil || !strings.Contains(err.Error(), templatePath) {
		t.Fatalf("Render() error = %v, want a template execution error naming %q", err, templatePath)
	}
	if got, err := os.ReadFile(outputPath); err != nil || string(got) != "previous documentation\n" {
		t.Fatalf("output after failed render = %q, error = %v", got, err)
	}
}

func TestRenderHonorsCancellationWithoutCreatingOutput(t *testing.T) {
	dir := t.TempDir()
	schemaPath := filepath.Join(dir, "values.schema.json")
	valuesPath := filepath.Join(dir, "values.yaml")
	outputPath := filepath.Join(dir, "README.md")
	writeFile(t, schemaPath, `{"type":"object"}`)
	writeFile(t, valuesPath, "{}\n")

	service := app.New(jsonschema.Loader{}, yamlvalues.Adapter{}, markdown.Adapter{}, filesystem.TemplateAdapter{}, filesystem.Store{})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := service.Render(ctx, inport.Request{
		SchemaPath:  schemaPath,
		ValuesPaths: []string{valuesPath},
		OutputPath:  outputPath,
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Render() error = %v, want context.Canceled", err)
	}
	if _, err := os.Stat(outputPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("output exists after canceled render; stat error = %v", err)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
