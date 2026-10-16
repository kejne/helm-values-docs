// Package out defines interfaces required by the application from adapters.
package out

import "github.com/kejne/helm-values-docs/internal/domain"

// SchemaLoader reads and decodes a schema document.
type SchemaLoader interface {
	Load(string) (*domain.Schema, error)
}

// ValuesLoader reads and merges ordered values documents.
type ValuesLoader interface {
	Load([]string) (map[string]any, error)
}

// Renderer turns the domain model and effective values into Markdown.
type Renderer interface {
	Render(*domain.Schema, map[string]any) ([]byte, error)
}

// TemplateComposer optionally wraps rendered Markdown in a template.
type TemplateComposer interface {
	Compose([]byte, string) ([]byte, error)
}

// DocumentStore persists and reads generated documents.
type DocumentStore interface {
	WriteAtomic(string, []byte) error
	Read(string) ([]byte, error)
}
