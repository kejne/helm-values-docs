// Package markdown renders the canonical standalone API document.
package markdown

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/kejne/helm-values-docs/internal/schema"
)

// Render produces deterministic Markdown for a schema and effective values.
func Render(root *schema.Schema, values map[string]any) ([]byte, error) {
	if root == nil {
		return nil, fmt.Errorf("schema is nil")
	}
	var b strings.Builder
	b.WriteString("# Helm values API\n\n")
	b.WriteString("| Path | Type | Required | Value | Description | Constraints |\n")
	b.WriteString("| --- | --- | --- | --- | --- | --- |\n")
	notes := make([]string, 0)
	for _, unsupported := range root.Unsupported {
		notes = append(notes, fmt.Sprintf("`$`: `%s`", unsupported))
	}
	effective := effectiveValues(root, values)
	for _, name := range root.PropertyNames() {
		renderProperty(&b, root.Properties[name], pathForProperty(name), name, effective, root.Required[name], &notes)
	}
	if len(notes) > 0 {
		b.WriteString("\n## Unsupported schema constructs\n\n")
		sort.Strings(notes)
		for _, note := range notes {
			b.WriteString("- ")
			b.WriteString(note)
			b.WriteByte('\n')
		}
	}
	return []byte(b.String()), nil
}

func renderProperty(b *strings.Builder, current *schema.Schema, path, name string, parent any, required bool, notes *[]string) {
	value, present := propertyValue(parent, name)
	if !present && current.HasDefault {
		value, present = current.Default, true
	}
	b.WriteString("| ")
	b.WriteString(cell(path))
	b.WriteString(" | ")
	b.WriteString(cell(typeName(current)))
	b.WriteString(" | ")
	if required {
		b.WriteString("yes")
	} else {
		b.WriteString("no")
	}
	b.WriteString(" | ")
	if present {
		b.WriteString(cell(schema.JSON(value)))
	} else {
		b.WriteString("—")
	}
	b.WriteString(" | ")
	b.WriteString(cell(current.Description))
	b.WriteString(" | ")
	b.WriteString(cell(constraints(current)))
	b.WriteString(" |\n")

	for _, unsupported := range current.Unsupported {
		*notes = append(*notes, fmt.Sprintf("`%s`: `%s`", path, unsupported))
	}
	childValues := value
	for _, childName := range current.PropertyNames() {
		renderProperty(b, current.Properties[childName], childPath(path, childName), childName, childValues, current.Required[childName], notes)
	}
}

func propertyValue(parent any, name string) (any, bool) {
	mapping, ok := parent.(map[string]any)
	if !ok {
		return nil, false
	}
	value, ok := mapping[name]
	return value, ok
}

func effectiveValues(root *schema.Schema, values map[string]any) map[string]any {
	result := make(map[string]any)
	if defaults, ok := root.Default.(map[string]any); root.HasDefault && ok {
		result = cloneMap(defaults)
	}
	mergeMaps(result, values)
	return result
}

func cloneMap(source map[string]any) map[string]any {
	result := make(map[string]any, len(source))
	for key, value := range source {
		result[key] = cloneValue(value)
	}
	return result
}

func cloneValue(value any) any {
	switch value := value.(type) {
	case map[string]any:
		return cloneMap(value)
	case []any:
		result := make([]any, len(value))
		for index, item := range value {
			result[index] = cloneValue(item)
		}
		return result
	default:
		return value
	}
}

func mergeMaps(dst, src map[string]any) {
	for key, value := range src {
		if srcMap, ok := value.(map[string]any); ok {
			if dstMap, ok := dst[key].(map[string]any); ok {
				mergeMaps(dstMap, srcMap)
				continue
			}
		}
		dst[key] = value
	}
}

func pathForProperty(name string) string {
	return pathSegment(name)
}

func childPath(path, name string) string {
	return path + pathSegment(name)
}

func pathSegment(name string) string {
	return "[" + strconv.Quote(name) + "]"
}

func typeName(current *schema.Schema) string {
	if len(current.Types) > 0 {
		return strings.Join(current.Types, " | ")
	}
	if len(current.Properties) > 0 {
		return "object"
	}
	if current.Items != nil {
		return "array"
	}
	return "any"
}

func constraints(current *schema.Schema) string {
	parts := make([]string, 0, len(current.Constraints)+2)
	keys := make([]string, 0, len(current.Constraints))
	for key := range current.Constraints {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		parts = append(parts, key+"="+schema.JSON(current.Constraints[key]))
	}
	if len(current.Enum) > 0 {
		parts = append(parts, "enum="+schema.JSON(current.Enum))
	}
	if len(current.Examples) > 0 {
		parts = append(parts, "examples="+schema.JSON(current.Examples))
	}
	return strings.Join(parts, ", ")
}

func cell(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "|", "\\|")
	value = strings.ReplaceAll(value, "\r\n", "<br>")
	value = strings.ReplaceAll(value, "\n", "<br>")
	value = strings.ReplaceAll(value, "\r", "<br>")
	return value
}
