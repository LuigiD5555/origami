package capacity

import (
	"strings"
	"testing"

	"github.com/LuigiD5555/origami/internal/lab/canonicaldocument"
	"github.com/LuigiD5555/origami/internal/lab/superindex"
)

func TestMeasureKeepsRepresentationScopesSeparate(t *testing.T) {
	document := canonicaldocument.Document{SourceSHA256: strings.Repeat("a", 64), Pages: []canonicaldocument.Page{{Number: 1}}}
	index := superindex.Index{SourceSHA256: document.SourceSHA256, Pages: map[int][]string{1: {"heading"}}, Entities: map[string][]int{}, Relations: map[string][]string{}, Terms: map[string][]int{"swarm": {1}}}
	measurement, err := Measure(document, index, 100, 1600, 2000)
	if err != nil {
		t.Fatal(err)
	}
	if len(measurement.Candidates) != 3 {
		t.Fatalf("candidates=%d", len(measurement.Candidates))
	}
	if measurement.Candidates[0].Scope != "SELF_CONTAINED" || measurement.Candidates[2].Scope != "REFERENCED" {
		t.Fatalf("representation scopes collapsed: %+v", measurement.Candidates)
	}
	if measurement.Candidates[2].Dependency == "NONE" {
		t.Fatal("referenced candidate hid its dependency")
	}
}
