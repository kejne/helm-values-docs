// Package jsonschema translates JSON Schema documents into the domain model.
package jsonschema

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/kejne/helm-values-docs/internal/domain"
)

// Parse decodes a JSON Schema document.
func Parse(data []byte) (*domain.Schema, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("decode JSON schema: %w", err)
	}
	return parseObject(raw, "$")
}

func parseObject(raw map[string]json.RawMessage, path string) (*domain.Schema, error) {
	result := &domain.Schema{
		Properties:  make(map[string]*domain.Schema),
		Required:    make(map[string]bool),
		Constraints: make(map[string]any),
	}
	if value, ok := raw["type"]; ok {
		if err := json.Unmarshal(value, &result.Types); err != nil {
			var one string
			if err := json.Unmarshal(value, &one); err != nil {
				return nil, fmt.Errorf("schema %s: type must be a string or array of strings", path)
			}
			result.Types = []string{one}
		}
		if len(result.Types) == 0 {
			return nil, fmt.Errorf("schema %s: type must not be empty", path)
		}
	}
	if value, ok := raw["description"]; ok {
		if err := json.Unmarshal(value, &result.Description); err != nil {
			return nil, fmt.Errorf("schema %s: description must be a string", path)
		}
	}
	if value, ok := raw["default"]; ok {
		if err := json.Unmarshal(value, &result.Default); err != nil {
			return nil, fmt.Errorf("schema %s: invalid default: %w", path, err)
		}
		result.HasDefault = true
	}
	for _, field := range []string{"enum", "examples"} {
		if value, ok := raw[field]; ok {
			var items []any
			if err := json.Unmarshal(value, &items); err != nil {
				return nil, fmt.Errorf("schema %s: %s must be an array", path, field)
			}
			if field == "enum" {
				result.Enum = items
			} else {
				result.Examples = items
			}
		}
	}
	if value, ok := raw["required"]; ok {
		var required []string
		if err := json.Unmarshal(value, &required); err != nil {
			return nil, fmt.Errorf("schema %s: required must be an array of strings", path)
		}
		for _, name := range required {
			result.Required[name] = true
		}
	}
	if value, ok := raw["properties"]; ok {
		var properties map[string]json.RawMessage
		if err := json.Unmarshal(value, &properties); err != nil || properties == nil {
			return nil, fmt.Errorf("schema %s: properties must be an object", path)
		}
		for name, property := range properties {
			var object map[string]json.RawMessage
			if err := json.Unmarshal(property, &object); err != nil || object == nil {
				return nil, fmt.Errorf("schema %s.%s: property must be an object", path, name)
			}
			child, err := parseObject(object, path+"."+name)
			if err != nil {
				return nil, err
			}
			result.Properties[name] = child
		}
	}
	if value, ok := raw["items"]; ok {
		var object map[string]json.RawMessage
		if err := json.Unmarshal(value, &object); err != nil || object == nil {
			return nil, fmt.Errorf("schema %s: items must be an object", path)
		}
		items, err := parseObject(object, path+"[]")
		if err != nil {
			return nil, err
		}
		result.Items = items
	}
	for _, key := range []string{"minLength", "maxLength", "pattern", "minimum", "maximum", "exclusiveMinimum", "exclusiveMaximum", "multipleOf", "minItems", "maxItems", "uniqueItems", "format"} {
		if value, ok := raw[key]; ok {
			var decoded any
			if err := json.Unmarshal(value, &decoded); err != nil {
				return nil, fmt.Errorf("schema %s: invalid %s: %w", path, key, err)
			}
			result.Constraints[key] = decoded
		}
	}
	for key := range raw {
		if isUnsupported(key) {
			result.Unsupported = append(result.Unsupported, key)
		}
	}
	sort.Strings(result.Unsupported)
	return result, nil
}

func isUnsupported(key string) bool {
	switch key {
	case "$schema", "$id", "$comment", "title", "deprecated", "readOnly", "writeOnly":
		// Recognized JSON Schema metadata that does not affect the rendered API.
		return false
	case "type", "description", "default", "enum", "examples", "required", "properties", "items",
		"minLength", "maxLength", "pattern", "minimum", "maximum", "exclusiveMinimum", "exclusiveMaximum",
		"multipleOf", "minItems", "maxItems", "uniqueItems", "format":
		return false
	default:
		// Do not silently discard a keyword that this deliberately small model
		// does not understand. The renderer will report it at its schema path.
		return true
	}
}
