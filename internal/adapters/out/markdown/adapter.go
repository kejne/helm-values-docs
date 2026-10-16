package markdown

import "github.com/kejne/helm-values-docs/internal/domain"

// Adapter exposes the Markdown renderer through the outbound port.
type Adapter struct{}

func (Adapter) Render(schema *domain.Schema, values map[string]any) ([]byte, error) {
	return Render(schema, values)
}
