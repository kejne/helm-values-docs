// Package domain contains the renderer's technology-independent schema model.
package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

// Schema is the supported subset of a JSON Schema used by the renderer.
type Schema struct {
	Types                []string
	Description          string
	Title                string
	Default              any
	HasDefault           bool
	Enum                 []any
	Examples             []any
	Properties           map[string]*Schema
	Required             map[string]bool
	Items                *Schema
	AdditionalProperties *bool
	AdditionalSchema     *Schema
	Reference            string
	Constraints          map[string]any
	Unsupported          []string
}

// PropertyNames returns property names in canonical order.
func (s *Schema) PropertyNames() []string {
	result := make([]string, 0, len(s.Properties))
	for name := range s.Properties {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

// JSON returns a compact deterministic representation of a value. Values that
// cannot be encoded as JSON are represented by their stable Go type name.
func JSON(value any) string {
	if value == nil {
		return "null"
	}
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%T", value)
	}
	return string(bytes.TrimSpace(data))
}
