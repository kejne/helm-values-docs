// Package in defines the application's inbound use-case interface.
package in

import (
	"context"
	"errors"
)

// ErrOutputMissing marks a check request whose generated document does not exist.
var ErrOutputMissing = errors.New("output document is missing")

// Request identifies the files used to produce a document.
type Request struct {
	SchemaPath   string
	ValuesPaths  []string
	OutputPath   string
	TemplatePath string
}

// Documentation generates and verifies chart values documentation.
type Documentation interface {
	Render(context.Context, Request) error
	Check(context.Context, Request) (bool, error)
}
