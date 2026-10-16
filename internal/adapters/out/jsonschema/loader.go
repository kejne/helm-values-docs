package jsonschema

import (
	"fmt"
	"os"

	"github.com/kejne/helm-values-docs/internal/domain"
)

// Loader reads schema files and decodes them into the domain model.
type Loader struct{}

func (Loader) Load(path string) (*domain.Schema, error) {
	data, err := os.ReadFile(path) // #nosec G304 -- path is provided explicitly by the caller.
	if err != nil {
		return nil, fmt.Errorf("read schema %q: %w", path, err)
	}
	model, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("schema %q: %w", path, err)
	}
	return model, nil
}
