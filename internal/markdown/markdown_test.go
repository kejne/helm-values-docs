package markdown

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/kejne/helm-values-docs/internal/schema"
)

func TestRenderGolden(t *testing.T) {
	model, err := schema.Parse([]byte(`{
		"type":"object",
		"required":["replicas"],
		"properties":{
			"replicas":{"type":"integer","default":1,"minimum":1,"description":"Number of pods"},
			"image":{"type":"object","properties":{"tag":{"type":"string","description":"Image | tag"}}},
			"missing":{"type":"string"}
		}
	}`))
	if err != nil {
		t.Fatal(err)
	}
	got, err := Render(model, map[string]any{
		"image":    map[string]any{"tag": "stable"},
		"replicas": 3,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertApproval(t, "standalone.md", got)
}

func TestRenderLooksUpPropertyNamesContainingDots(t *testing.T) {
	model, err := schema.Parse([]byte(`{"properties":{"a.b":{"type":"string"}}}`))
	if err != nil {
		t.Fatal(err)
	}
	got, err := Render(model, map[string]any{"a.b": "value"})
	if err != nil {
		t.Fatal(err)
	}
	if !containsLine(string(got), `| ["a.b"] | string | no | "value" |  |  |`) {
		t.Fatalf("dotted property value missing from:\n%s", got)
	}
}

func TestRenderUsesRootDefaultAndValuesPrecedence(t *testing.T) {
	model, err := schema.Parse([]byte(`{"type":"object","default":{"enabled":true,"nested":{"count":1}},"properties":{"enabled":{"type":"boolean"},"nested":{"type":"object","properties":{"count":{"type":"integer"}}}}}`))
	if err != nil {
		t.Fatal(err)
	}
	got, err := Render(model, map[string]any{"nested": map[string]any{"count": 2}})
	if err != nil {
		t.Fatal(err)
	}
	if !containsLine(string(got), `| ["enabled"] | boolean | no | true |  |  |`) || !containsLine(string(got), `| ["nested"]["count"] | integer | no | 2 |  |  |`) {
		t.Fatalf("root default or values precedence missing from:\n%s", got)
	}
}

func TestRenderUsesRootDefaultWhenValuesAreEmptyWithoutMutation(t *testing.T) {
	model, err := schema.Parse([]byte(`{"type":"object","default":{"enabled":true,"nested":{"count":1}},"properties":{"enabled":{"type":"boolean"},"nested":{"type":"object","properties":{"count":{"type":"integer"}}}}}`))
	if err != nil {
		t.Fatal(err)
	}
	values := map[string]any{}
	before := map[string]any{}
	got, err := Render(model, values)
	if err != nil {
		t.Fatal(err)
	}
	if !containsLine(string(got), `| ["enabled"] | boolean | no | true |  |  |`) || !containsLine(string(got), `| ["nested"]["count"] | integer | no | 1 |  |  |`) {
		t.Fatalf("root default missing from:\n%s", got)
	}
	if !reflect.DeepEqual(values, before) {
		t.Fatalf("render mutated values: got %#v, want %#v", values, before)
	}
}

func TestRenderEscapesAllPropertyPathSegments(t *testing.T) {
	model, err := schema.Parse([]byte(`{"properties":{"parent":{"properties":{"a.b|c\"d":{"type":"string"}}}}}`))
	if err != nil {
		t.Fatal(err)
	}
	got, err := Render(model, map[string]any{"parent": map[string]any{"a.b|c\"d": "value"}})
	if err != nil {
		t.Fatal(err)
	}
	if !containsLine(string(got), `| ["parent"]["a.b\|c\\"d"] | string | no | "value" |  |  |`) {
		t.Fatalf("escaped property path missing from:\n%s", got)
	}
}

func TestRenderReportsUnsupportedKeywordsAtTheirPaths(t *testing.T) {
	model, err := schema.Parse([]byte(`{"additionalProperties":false,"properties":{"setting":{"const":"fixed","dependentRequired":{"x":["y"]}}}}`))
	if err != nil {
		t.Fatal(err)
	}
	got, err := Render(model, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	document := string(got)
	for _, note := range []string{"- `$`: `additionalProperties`", "- `[\"setting\"]`: `const`", "- `[\"setting\"]`: `dependentRequired`"} {
		if !containsLine(document, note) {
			t.Fatalf("unsupported note %q missing from:\n%s", note, document)
		}
	}
}

func TestRenderNonObjectRootDefaultRemainsMissing(t *testing.T) {
	model, err := schema.Parse([]byte(`{"type":"object","default":false,"properties":{"enabled":{"type":"boolean"}}}`))
	if err != nil {
		t.Fatal(err)
	}
	got, err := Render(model, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if !containsLine(string(got), `| ["enabled"] | boolean | no | — |  |  |`) {
		t.Fatalf("non-object root default changed missing-value behavior:\n%s", got)
	}
}

func assertApproval(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", name)
	if os.Getenv("UPDATE_APPROVALS") == "1" {
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatalf("read approval directory: %v", err)
	}
	for _, entry := range entries {
		if entry.Name() != name {
			t.Fatalf("unexpected approval artifact %s", filepath.Join(filepath.Dir(path), entry.Name()))
		}
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read approval %s: %v", path, err)
	}
	if string(got) != string(want) {
		t.Errorf("approval %s differs:\n got:\n%s\nwant:\n%s", path, got, want)
	}
}

func TestRenderReportsUnsupportedRootConstructs(t *testing.T) {
	model, err := schema.Parse([]byte(`{"oneOf":[{"type":"string"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	got, err := Render(model, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if !containsLine(string(got), "- `$`: `oneOf`") {
		t.Fatalf("unsupported root note missing from:\n%s", got)
	}
}

func TestRenderEscapesCellsAndUsesDefaults(t *testing.T) {
	model, err := schema.Parse([]byte(`{"properties":{"a|b":{"type":"string","default":"x|y"}}}`))
	if err != nil {
		t.Fatal(err)
	}
	got, err := Render(model, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	if want := `| ["a\|b"] | string | no | "x\|y" |  |  |`; !containsLine(string(got), want) {
		t.Fatalf("escaped/default row missing from:\n%s", got)
	}
}

func containsLine(document, line string) bool {
	for _, candidate := range splitLines(document) {
		if candidate == line {
			return true
		}
	}
	return false
}

func splitLines(value string) []string {
	result := []string{}
	for len(value) > 0 {
		index := 0
		for index < len(value) && value[index] != '\n' {
			index++
		}
		result = append(result, value[:index])
		if index == len(value) {
			break
		}
		value = value[index+1:]
	}
	return result
}
