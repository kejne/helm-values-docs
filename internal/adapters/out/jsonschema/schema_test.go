package jsonschema_test

import (
	"testing"

	"github.com/kejne/helm-values-docs/internal/adapters/out/jsonschema"
)

func TestParseSchema(t *testing.T) {
	model, err := jsonschema.Parse([]byte(`{
		"type":"object",
		"required":["replicas"],
		"properties":{
			"replicas":{"type":"integer","default":1,"minimum":1},
			"image":{"type":"object","properties":{"tag":{"type":"string"}}},
			"choice":{"oneOf":[{"type":"string"}]}
		}
	}`))
	if err != nil {
		t.Fatal(err)
	}
	if !model.Required["replicas"] || model.Properties["replicas"].Default != float64(1) {
		t.Fatalf("schema was not decoded: %#v", model)
	}
	if got := model.PropertyNames(); len(got) != 3 || got[0] != "choice" || got[1] != "image" || got[2] != "replicas" {
		t.Fatalf("property order = %v", got)
	}
	if len(model.Properties["choice"].Unsupported) != 1 || model.Properties["choice"].Unsupported[0] != "oneOf" {
		t.Fatalf("unsupported constructs = %v", model.Properties["choice"].Unsupported)
	}
}

func TestParseRejectsInvalidSchema(t *testing.T) {
	if _, err := jsonschema.Parse([]byte(`{"type":`)); err == nil {
		t.Fatal("Parse accepted malformed JSON")
	}
}

func TestParsePreservesOptionalSchemaDetails(t *testing.T) {
	model, err := jsonschema.Parse([]byte(`{
		"description":"A chart setting",
		"enum":["small","large"],
		"examples":["small"],
		"items":{"type":"string","description":"An entry"}
	}`))
	if err != nil {
		t.Fatal(err)
	}
	if model.Description != "A chart setting" || len(model.Enum) != 2 || model.Enum[0] != "small" || len(model.Examples) != 1 || model.Examples[0] != "small" {
		t.Fatalf("optional schema details were not preserved: %#v", model)
	}
	if model.Items == nil || len(model.Items.Types) != 1 || model.Items.Types[0] != "string" || model.Items.Description != "An entry" {
		t.Fatalf("item schema was not preserved: %#v", model.Items)
	}
}

func TestParseRejectsInvalidOptionalSchemaDetails(t *testing.T) {
	for _, tc := range []struct {
		name string
		data string
	}{
		{name: "description must be a string", data: `{"description":false}`},
		{name: "enum must be an array", data: `{"enum":"small"}`},
		{name: "examples must be an array", data: `{"examples":{}}`},
		{name: "items must be an object", data: `{"items":[]}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := jsonschema.Parse([]byte(tc.data)); err == nil {
				t.Fatalf("Parse(%s) accepted invalid input", tc.data)
			}
		})
	}
}

func TestParseSupportsAdditionalPropertiesAndLocalReferences(t *testing.T) {
	model, err := jsonschema.Parse([]byte(`{
		"properties": {
			"mounts": {"type":"object", "additionalProperties":{"$ref":"#/definitions/mount"}},
			"labels": {"type":"object", "additionalProperties":true}
		},
		"definitions": {"mount":{"title":"Mount", "type":"object", "properties":{"path":{"type":"string"}}}}
	}`))
	if err != nil {
		t.Fatal(err)
	}
	mounts := model.Properties["mounts"]
	if mounts.AdditionalSchema == nil || mounts.AdditionalSchema.Title != "Mount" || mounts.AdditionalSchema.Properties["path"] == nil {
		t.Fatalf("local ref not resolved in additionalProperties: %#v", mounts)
	}
	labels := model.Properties["labels"]
	if labels.AdditionalProperties == nil || !*labels.AdditionalProperties {
		t.Fatalf("additionalProperties=true not parsed: %#v", labels)
	}
	if len(model.Unsupported) != 0 || len(mounts.Unsupported) != 0 {
		t.Fatalf("supported constructs reported unsupported: root=%v mounts=%v", model.Unsupported, mounts.Unsupported)
	}
}

func TestParseReportsEveryUnsupportedKeyword(t *testing.T) {
	model, err := jsonschema.Parse([]byte(`{
		"additionalProperties":false,
		"patternProperties":{"^x-":{"type":"string"}},
		"minProperties":1,
		"const":"fixed",
		"$defs":{"thing":{"type":"string"}},
		"dependentRequired":{"name":["value"]}
	}`))
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"$defs", "const", "dependentRequired", "minProperties", "patternProperties"}
	if len(model.Unsupported) != len(want) {
		t.Fatalf("unsupported constructs = %v, want %v", model.Unsupported, want)
	}
	for i := range want {
		if model.Unsupported[i] != want[i] {
			t.Fatalf("unsupported constructs = %v, want %v", model.Unsupported, want)
		}
	}
}
