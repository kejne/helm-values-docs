package schema

import (
	"testing"
)

func TestParseSchema(t *testing.T) {
	model, err := Parse([]byte(`{
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
	if _, err := Parse([]byte(`{"type":`)); err == nil {
		t.Fatal("Parse accepted malformed JSON")
	}
}

func TestParseReportsEveryUnsupportedKeyword(t *testing.T) {
	model, err := Parse([]byte(`{
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
	want := []string{"$defs", "additionalProperties", "const", "dependentRequired", "minProperties", "patternProperties"}
	if len(model.Unsupported) != len(want) {
		t.Fatalf("unsupported constructs = %v, want %v", model.Unsupported, want)
	}
	for i := range want {
		if model.Unsupported[i] != want[i] {
			t.Fatalf("unsupported constructs = %v, want %v", model.Unsupported, want)
		}
	}
}
