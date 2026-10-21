// Package app implements the documentation use cases.
package app

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"

	inport "github.com/smedijan/helm-values-docs/internal/ports/in"
	outport "github.com/smedijan/helm-values-docs/internal/ports/out"
)

// Service coordinates domain operations through outbound ports.
type Service struct {
	schemas   outport.SchemaLoader
	values    outport.ValuesLoader
	renderer  outport.Renderer
	templates outport.TemplateComposer
	store     outport.DocumentStore
}

// New constructs the use-case implementation from its required adapters.
func New(schemas outport.SchemaLoader, values outport.ValuesLoader, renderer outport.Renderer, templates outport.TemplateComposer, store outport.DocumentStore) *Service {
	return &Service{schemas: schemas, values: values, renderer: renderer, templates: templates, store: store}
}

// Render implements the render use case and atomically persists the result.
func (s *Service) Render(ctx context.Context, request inport.Request) error {
	data, err := s.render(ctx, request)
	if err != nil {
		return err
	}
	return s.store.WriteAtomic(request.OutputPath, data)
}

// Check implements the check use case without mutating the output document.
func (s *Service) Check(ctx context.Context, request inport.Request) (bool, error) {
	want, err := s.render(ctx, request)
	if err != nil {
		return false, err
	}
	have, err := s.store.Read(request.OutputPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, fmt.Errorf("%w: %w", inport.ErrOutputMissing, err)
		}
		return false, err
	}
	return bytes.Equal(want, have), nil
}

func (s *Service) render(ctx context.Context, request inport.Request) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	model, err := s.schemas.Load(request.SchemaPath)
	if err != nil {
		return nil, err
	}
	values, err := s.values.Load(request.ValuesPaths)
	if err != nil {
		return nil, err
	}
	data, err := s.renderer.Render(model, values)
	if err != nil {
		return nil, err
	}
	if request.TemplatePath == "" {
		return data, nil
	}
	data, err = s.templates.Compose(data, request.TemplatePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}
