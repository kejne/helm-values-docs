package filesystem

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// TemplateAdapter composes a document with an optional Markdown template.
type TemplateAdapter struct{}

type templateData struct {
	API string
}

func (TemplateAdapter) Compose(api []byte, path string) ([]byte, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- path is explicitly provided by the caller.
	if err != nil {
		return nil, fmt.Errorf("read template %q: %w", path, err)
	}
	tmpl, err := template.New(filepath.Base(path)).Option("missingkey=error").Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse template %q: %w", path, err)
	}
	var composed bytes.Buffer
	if err := tmpl.Execute(&composed, templateData{API: string(api)}); err != nil {
		return nil, fmt.Errorf("execute template %q: %w", path, err)
	}
	return composed.Bytes(), nil
}
