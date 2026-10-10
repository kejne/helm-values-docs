// Package input loads and combines Helm values files.
package input

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadValues reads YAML files in order. Mapping values are merged recursively;
// scalars and sequences from a later file replace earlier values.
func LoadValues(paths []string) (map[string]any, error) {
	result := make(map[string]any)
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read values %q: %w", path, err)
		}
		value, err := parseYAML(data)
		if err != nil {
			return nil, fmt.Errorf("parse values %q: %w", path, err)
		}
		mapping, ok := value.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("parse values %q: root must be a mapping", path)
		}
		mergeMap(result, mapping)
	}
	return result, nil
}

func parseYAML(data []byte) (any, error) {
	var node yaml.Node
	if err := yaml.Unmarshal(data, &node); err != nil {
		return nil, err
	}
	if len(node.Content) == 0 {
		return map[string]any{}, nil
	}
	return nodeValue(node.Content[0])
}

func nodeValue(node *yaml.Node) (any, error) {
	switch node.Kind {
	case yaml.MappingNode:
		result := make(map[string]any, len(node.Content)/2)
		for i := 0; i < len(node.Content); i += 2 {
			key := node.Content[i]
			if key.Tag != "!!str" {
				return nil, fmt.Errorf("mapping key %q is not a string", key.Value)
			}
			if _, exists := result[key.Value]; exists {
				return nil, fmt.Errorf("duplicate mapping key %q", key.Value)
			}
			value, err := nodeValue(node.Content[i+1])
			if err != nil {
				return nil, err
			}
			result[key.Value] = value
		}
		return result, nil
	case yaml.SequenceNode:
		result := make([]any, len(node.Content))
		for i, child := range node.Content {
			value, err := nodeValue(child)
			if err != nil {
				return nil, err
			}
			result[i] = value
		}
		return result, nil
	case yaml.ScalarNode:
		var value any
		if err := node.Decode(&value); err != nil {
			return nil, err
		}
		return value, nil
	case yaml.AliasNode:
		return nodeValue(node.Alias)
	default:
		return nil, fmt.Errorf("unsupported YAML node kind %d", node.Kind)
	}
}

func mergeMap(dst, src map[string]any) {
	for key, value := range src {
		if srcMap, ok := value.(map[string]any); ok {
			if dstMap, ok := dst[key].(map[string]any); ok {
				mergeMap(dstMap, srcMap)
				continue
			}
		}
		dst[key] = value
	}
}
