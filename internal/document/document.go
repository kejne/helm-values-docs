// Package document provides the shared input and output pipeline.
package document

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/kejne/helm-values-docs/internal/input"
	"github.com/kejne/helm-values-docs/internal/markdown"
	"github.com/kejne/helm-values-docs/internal/schema"
)

// TemplateData is the stable data contract exposed to composition templates.
type TemplateData struct {
	API string
}

// RenderFiles reads the schema and ordered values files and renders Markdown.
func RenderFiles(schemaPath string, valuesPaths []string) ([]byte, error) {
	return RenderFilesWithTemplate(schemaPath, valuesPaths, "")
}

// RenderFilesWithTemplate renders Markdown and optionally composes it with a template.
func RenderFilesWithTemplate(schemaPath string, valuesPaths []string, templatePath string) ([]byte, error) {
	schemaData, err := os.ReadFile(schemaPath) // #nosec G304 -- paths are explicitly supplied by the caller.
	if err != nil {
		return nil, fmt.Errorf("read schema %q: %w", schemaPath, err)
	}
	model, err := schema.Parse(schemaData)
	if err != nil {
		return nil, fmt.Errorf("schema %q: %w", schemaPath, err)
	}
	values, err := input.LoadValues(valuesPaths)
	if err != nil {
		return nil, err
	}
	api, err := markdown.Render(model, values)
	if err != nil {
		return nil, err
	}
	if templatePath == "" {
		return api, nil
	}
	return compose(api, templatePath)
}

func compose(api []byte, templatePath string) ([]byte, error) {
	templateData, err := os.ReadFile(templatePath) // #nosec G304 -- paths are explicitly supplied by the caller.
	if err != nil {
		return nil, fmt.Errorf("read template %q: %w", templatePath, err)
	}
	tmpl, err := template.New(filepath.Base(templatePath)).Option("missingkey=error").Parse(string(templateData))
	if err != nil {
		return nil, fmt.Errorf("parse template %q: %w", templatePath, err)
	}
	var composed bytes.Buffer
	if err := tmpl.Execute(&composed, TemplateData{API: string(api)}); err != nil {
		return nil, fmt.Errorf("execute template %q: %w", templatePath, err)
	}
	return composed.Bytes(), nil
}

// CheckFiles reports whether the existing output exactly matches freshly rendered data.
// It never writes or changes the output file.
func CheckFiles(schemaPath string, valuesPaths []string, outputPath string) (bool, error) {
	return CheckFilesWithTemplate(schemaPath, valuesPaths, outputPath, "")
}

// CheckFilesWithTemplate compares standalone or composed rendered bytes to output.
func CheckFilesWithTemplate(schemaPath string, valuesPaths []string, outputPath, templatePath string) (bool, error) {
	want, err := RenderFilesWithTemplate(schemaPath, valuesPaths, templatePath)
	if err != nil {
		return false, err
	}
	return CheckRendered(want, outputPath)
}

// CheckRendered compares already-rendered bytes with an output file.
func CheckRendered(want []byte, outputPath string) (bool, error) {
	have, err := os.ReadFile(outputPath) // #nosec G304 -- paths are explicitly supplied by the caller.
	if err != nil {
		return false, fmt.Errorf("read output %q: %w", outputPath, err)
	}
	return bytes.Equal(want, have), nil
}

// WriteAtomic writes data by replacing path only after the complete document is ready.
func WriteAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	file, err := os.CreateTemp(dir, ".helm-values-docs-*")
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
