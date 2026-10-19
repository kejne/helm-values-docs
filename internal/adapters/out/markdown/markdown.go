// Package markdown adapts the domain model to canonical standalone Markdown.
package markdown

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kejne/helm-values-docs/internal/domain"
)

// Render produces deterministic Markdown for a schema and effective values.
func Render(root *domain.Schema, values map[string]any) ([]byte, error) {
	if root == nil {
		return nil, fmt.Errorf("schema is nil")
	}
	var b strings.Builder
	b.WriteString("# Helm values API\n\n")
	if len(root.Properties) > 0 {
		b.WriteString("## Contents\n\n")
		for _, name := range root.PropertyNames() {
			renderTOC(&b, root.Properties[name], name, name, 0)
		}
		b.WriteByte('\n')
	}
	b.WriteString("| **Schema** | **Details** |\n| --- | --- |\n")
	b.WriteString("| Type | `")
	b.WriteString(cell(typeName(root)))
	b.WriteString("` |\n")
	renderAdditionalProperties(&b, root)
	b.WriteByte('\n')
	notes := make([]string, 0)
	for _, unsupported := range root.Unsupported {
		notes = append(notes, fmt.Sprintf("`$`: `%s`", unsupported))
	}
	for _, name := range root.PropertyNames() {
		renderProperty(&b, root.Properties[name], name, name, values, root.Required[name], 2, &notes)
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

func renderTOC(b *strings.Builder, current *domain.Schema, name, path string, depth int) {
	b.WriteString(strings.Repeat("  ", depth))
	b.WriteString("- [")
	b.WriteString(name)
	b.WriteString("](#")
	b.WriteString(anchor(path))
	b.WriteString(")\n")
	children := current.PropertyNames()
	if current.AdditionalSchema != nil && len(current.AdditionalSchema.Properties) > 0 {
		children = append(children, additionalName(name, current.AdditionalSchema))
	}
	for _, child := range children {
		childModel, childName := childSchema(current, child)
		renderTOC(b, childModel, childName, childPath(path, childName), depth+1)
	}
}

func renderProperty(b *strings.Builder, current *domain.Schema, name, path string, parent any, required bool, depth int, notes *[]string) {
	value, present := propertyValue(parent, name)
	b.WriteString("\n<a id=\"")
	b.WriteString(anchor(path))
	b.WriteString("\"></a>\n")
	if depth > 6 {
		depth = 6
	}
	b.WriteString(strings.Repeat("#", depth))
	b.WriteByte(' ')
	b.WriteString("`")
	b.WriteString(name)
	b.WriteString("`\n\n")
	b.WriteString("| **Property** | **Details** |\n| --- | --- |\n")
	b.WriteString("| Type | `")
	b.WriteString(cell(typeName(current)))
	b.WriteString("` |\n| Required | ")
	if required {
		b.WriteString("Yes")
	} else {
		b.WriteString("No")
	}
	b.WriteString(" |\n")
	if len(current.Properties) > 0 {
		b.WriteString("| Default | See nested values |\n")
	} else {
		b.WriteString("| Default | ")
		if present {
			b.WriteString("`")
			b.WriteString(cell(domain.JSON(value)))
			b.WriteString("`")
		} else {
			b.WriteString("—")
		}
		b.WriteString(" |\n")
	}
	renderAdditionalProperties(b, current)
	if current.Reference != "" {
		b.WriteString("| Schema reference | [Referenced schema](")
		b.WriteString(current.Reference)
		b.WriteString(") |\n")
	}
	if details := constraints(current); details != "" {
		b.WriteString("| Constraints | ")
		b.WriteString(cell(details))
		b.WriteString(" |\n")
	}
	if current.Description != "" {
		b.WriteString("| Description | ")
		b.WriteString(cell(current.Description))
		b.WriteString(" |\n")
	}
	for _, unsupported := range current.Unsupported {
		*notes = append(*notes, fmt.Sprintf("`%s`: `%s`", path, unsupported))
	}
	children := current.PropertyNames()
	if current.AdditionalSchema != nil && len(current.AdditionalSchema.Properties) > 0 {
		children = append(children, additionalName(name, current.AdditionalSchema))
	}
	if len(children) > 0 {
		b.WriteString("\n| Property | Type | Required | Default | Description |\n| --- | --- | --- | --- | --- |\n")
		for _, child := range children {
			childSchema, childName := childSchema(current, child)
			childPath := childPath(path, childName)
			b.WriteString("| [")
			b.WriteString(cell(childName))
			b.WriteString("](#")
			b.WriteString(anchor(childPath))
			b.WriteString(") | `")
			b.WriteString(cell(typeName(childSchema)))
			b.WriteString("` | ")
			if childSchema.Required[childName] {
				b.WriteString("Yes")
			} else {
				b.WriteString("No")
			}
			b.WriteString(" | ")
			childValue, childPresent := propertyValue(value, childName)
			if len(childSchema.Properties) > 0 {
				b.WriteString("See nested values")
			} else if childPresent {
				b.WriteString("`" + cell(domain.JSON(childValue)) + "`")
			} else {
				b.WriteString("—")
			}
			b.WriteString(" | ")
			if childSchema.Description != "" {
				b.WriteString(cell(childSchema.Description))
			} else {
				b.WriteString("—")
			}
			b.WriteString(" |\n")
			renderProperty(b, childSchema, childName, childPath, value, current.Required[childName], depth+1, notes)
		}
	}
}

func renderAdditionalProperties(b *strings.Builder, schema *domain.Schema) {
	if schema.AdditionalProperties != nil {
		b.WriteString("| Additional properties | ")
		if *schema.AdditionalProperties {
			b.WriteString("Any type allowed")
		} else {
			b.WriteString("Not allowed")
		}
		b.WriteString(" |\n")
	} else if schema.AdditionalSchema != nil {
		b.WriteString("| Additional properties | Each additional property must conform to the schema")
		if len(schema.AdditionalSchema.Types) > 0 {
			b.WriteString(" (`")
			b.WriteString(cell(typeName(schema.AdditionalSchema)))
			b.WriteString("`)")
		}
		b.WriteString(" |\n")
	}
}

func childSchema(schema *domain.Schema, name string) (*domain.Schema, string) {
	if child, ok := schema.Properties[name]; ok {
		return child, name
	}
	if schema.AdditionalSchema != nil {
		return schema.AdditionalSchema, name
	}
	return nil, name
}

func additionalName(parent string, schema *domain.Schema) string {
	if title := strings.TrimSpace(schema.Title); title != "" {
		return title
	}
	if strings.HasSuffix(parent, "s") && len(parent) > 1 {
		return strings.TrimSuffix(parent, "s")
	}
	return parent + "-object"
}

func childPath(path, name string) string {
	return path + "-" + name
}

func anchor(path string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(path) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
		} else if b.Len() > 0 && !strings.HasSuffix(b.String(), "_") {
			b.WriteByte('_')
		}
	}
	return strings.Trim(b.String(), "_")
}

func propertyValue(parent any, name string) (any, bool) {
	mapping, ok := parent.(map[string]any)
	if !ok {
		return nil, false
	}
	value, ok := mapping[name]
	return value, ok
}

func typeName(current *domain.Schema) string {
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

func constraints(current *domain.Schema) string {
	parts := make([]string, 0, len(current.Constraints)+2)
	keys := make([]string, 0, len(current.Constraints))
	for key := range current.Constraints {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		parts = append(parts, key+"="+domain.JSON(current.Constraints[key]))
	}
	if len(current.Enum) > 0 {
		parts = append(parts, "enum="+domain.JSON(current.Enum))
	}
	if len(current.Examples) > 0 {
		parts = append(parts, "examples="+domain.JSON(current.Examples))
	}
	return strings.Join(parts, ", ")
}

func cell(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "|", "&#124;")
	value = strings.ReplaceAll(value, "\r\n", "<br>")
	value = strings.ReplaceAll(value, "\n", "<br>")
	value = strings.ReplaceAll(value, "\r", "<br>")
	return value
}
