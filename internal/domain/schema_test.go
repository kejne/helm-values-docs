package domain_test

import (
	"testing"

	"github.com/kejne/helm-values-docs/internal/domain"
)

func TestPropertyNamesAreCanonical(t *testing.T) {
	schema := &domain.Schema{Properties: map[string]*domain.Schema{
		"zeta":   {},
		"alpha":  {},
		"middle": {},
	}}

	got := schema.PropertyNames()
	want := []string{"alpha", "middle", "zeta"}
	if len(got) != len(want) {
		t.Fatalf("PropertyNames() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("PropertyNames() = %v, want %v", got, want)
		}
	}
}

func TestJSONIsCompactAndDeterministic(t *testing.T) {
	value := map[string]any{"zeta": []any{true, "x"}, "alpha": 2}
	if got, want := domain.JSON(value), `{"alpha":2,"zeta":[true,"x"]}`; got != want {
		t.Fatalf("JSON() = %q, want %q", got, want)
	}
}

func TestJSONUsesStableRepresentationForUnencodableValues(t *testing.T) {
	if got, want := domain.JSON(make(chan int)), "chan int"; got != want {
		t.Fatalf("JSON() = %q, want stable type name %q", got, want)
	}
}
