// Package jsonschema translates JSON Schema documents into the domain model.
package jsonschema

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/smedijan/helm-values-docs/internal/domain"
)

// Parse decodes a JSON Schema document.
func Parse(data []byte) (*domain.Schema, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("decode JSON schema: %w", err)
	}
	definitions := make(map[string]json.RawMessage)
	if value, ok := raw["definitions"]; ok {
		if err := json.Unmarshal(value, &definitions); err != nil {
			return nil, fmt.Errorf("schema $: definitions must be an object")
		}
	}
	return parseObject(raw, "$", definitions)
}

func parseObject(raw map[string]json.RawMessage, path string, definitions map[string]json.RawMessage) (*domain.Schema, error) {
	result := &domain.Schema{
		Properties:  make(map[string]*domain.Schema),
		Required:    make(map[string]bool),
		Constraints: make(map[string]any),
	}
	if value, ok := raw["$ref"]; ok {
		if err := json.Unmarshal(value, &result.Reference); err != nil {
			return nil, fmt.Errorf("schema %s: $ref must be a string", path)
		}
		if strings.HasPrefix(result.Reference, "#/definitions/") {
			name := strings.TrimPrefix(result.Reference, "#/definitions/")
			definition, ok := definitions[name]
			if !ok {
				return nil, fmt.Errorf("schema %s: unresolved reference %q", path, result.Reference)
			}
			var target map[string]json.RawMessage
			if err := json.Unmarshal(definition, &target); err != nil || target == nil {
				return nil, fmt.Errorf("schema %s: definition %q must be an object", path, name)
			}
			resolved, err := parseObject(target, path, definitions)
			if err != nil {
				return nil, err
			}
			if description, ok := raw["description"]; ok {
				if err := json.Unmarshal(description, &resolved.Description); err != nil {
					return nil, fmt.Errorf("schema %s: description must be a string", path)
				}
			}
			return resolved, nil
		}
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
	if value, ok := raw["title"]; ok {
		if err := json.Unmarshal(value, &result.Title); err != nil {
			return nil, fmt.Errorf("schema %s: title must be a string", path)
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
			child, err := parseObject(object, path+"."+name, definitions)
			if err != nil {
				return nil, err
			}
			result.Properties[name] = child
		}
	}
	if value, ok := raw["oneOf"]; ok {
		var branches []json.RawMessage
		if err := json.Unmarshal(value, &branches); err != nil || len(branches) == 0 {
			return nil, fmt.Errorf("schema %s: oneOf must be a non-empty array of schemas", path)
		}
		for index, branch := range branches {
			var object map[string]json.RawMessage
			if err := json.Unmarshal(branch, &object); err != nil || object == nil {
				return nil, fmt.Errorf("schema %s: oneOf[%d] must be an object", path, index)
			}
			parsed, err := parseObject(object, fmt.Sprintf("%s.oneOf[%d]", path, index), definitions)
			if err != nil {
				return nil, err
			}
			result.OneOf = append(result.OneOf, parsed)
		}
	}
	if value, ok := raw["not"]; ok {
		var object map[string]json.RawMessage
		if err := json.Unmarshal(value, &object); err != nil || object == nil {
			return nil, fmt.Errorf("schema %s: not must be a schema object", path)
		}
		parsed, err := parseObject(object, path+".not", definitions)
		if err != nil {
			return nil, err
		}
		result.Not = parsed
	}
	if value, ok := raw["items"]; ok {
		var object map[string]json.RawMessage
		if err := json.Unmarshal(value, &object); err != nil || object == nil {
			return nil, fmt.Errorf("schema %s: items must be an object", path)
		}
		items, err := parseObject(object, path+"[]", definitions)
		if err != nil {
			return nil, err
		}
		result.Items = items
	}
	if value, ok := raw["additionalProperties"]; ok {
		var allowed bool
		if err := json.Unmarshal(value, &allowed); err == nil {
			result.AdditionalProperties = &allowed
		} else {
			var object map[string]json.RawMessage
			if err := json.Unmarshal(value, &object); err != nil || object == nil {
				return nil, fmt.Errorf("schema %s: additionalProperties must be a boolean or schema object", path)
			}
			additional, err := parseObject(object, path+".*", definitions)
			if err != nil {
				return nil, err
			}
			result.AdditionalSchema = additional
		}
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
	case "$schema", "$id", "$comment", "title", "deprecated", "readOnly", "writeOnly", "definitions", "$ref", "additionalProperties":
		// Supported keywords are rendered or resolved by the schema adapter.
		return false
	case "type", "description", "default", "enum", "examples", "required", "properties", "items", "oneOf", "not",
		"minLength", "maxLength", "pattern", "minimum", "maximum", "exclusiveMinimum", "exclusiveMaximum",
		"multipleOf", "minItems", "maxItems", "uniqueItems", "format":
		return false
	default:
		// Do not silently discard a keyword that this deliberately small model
		// does not understand. The renderer will report it at its schema path.
		return true
	}
}
